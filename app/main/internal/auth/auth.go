package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	s "gin/internal/structs"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func CheckSession(session_id string, redis_db *redis.Client, timeout int, ctx context.Context) (error, int) {
	newcontext, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	var value s.SessionValue
	baitiki, err := redis_db.Get(newcontext, session_id).Result()
	if err != nil {
		return err, 0
	}
	fmt.Println("value in string:", baitiki)

	err = json.Unmarshal([]byte(baitiki), &value)
	if err != nil {
		return err, 0
	}
	fmt.Println("session from cookie:", session_id)
	fmt.Println("stored user id:", value.UserId)
	fmt.Println("stored exp session:", value.Exp)

	if err != nil {
		return err, value.UserId
	}
	if time.Now().After(value.Exp) {
		return fmt.Errorf("session exp"), value.UserId
	}
	daun, err := json.Marshal(value)
	if err != nil {
		return err, 0
	}
	err = redis_db.Set(newcontext, session_id, string(daun), time.Minute*30).Err()
	if err != nil {
		return err, value.UserId
	}
	return nil, value.UserId
}
func CreateCSRF(secret string, userId int) (string, error) {
	val := s.CSRFvalue{
		UserId: userId,
		Exp:    time.Now().Add(30 * time.Minute),
	}

	data, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write(data)
	signature := h.Sum(nil)

	// кодируем обе части в base64
	payloadEnc := base64.URLEncoding.EncodeToString(data)
	signatureEnc := base64.URLEncoding.EncodeToString(signature)

	// склеиваем payload + сигнатура
	token := fmt.Sprintf("%s.%s", payloadEnc, signatureEnc)

	return token, nil
}
func VerifyCSRF(secret string, token string) (int, string, error) {
	var payloadEnc, signatureEnc string
	parts := strings.SplitN(token, ".", 2)
	var id int
	if len(parts) != 2 {
		return id, "", fmt.Errorf("invalid token format")
	}
	payloadEnc, signatureEnc = parts[0], parts[1]

	// декодируем
	data, err := base64.URLEncoding.DecodeString(payloadEnc)
	if err != nil {
		return id, token, fmt.Errorf("err while decoding csrf data:", err.Error())
	}
	signature, err := base64.URLEncoding.DecodeString(signatureEnc)
	if err != nil {
		return id, token, fmt.Errorf("err while decoding csrf signature:", err.Error())
	}
	fmt.Println("unmurshalled data:", data)
	fmt.Println("unmurshalled signature:", signature)

	// пересчёт HMAC
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(data)
	expected := h.Sum(nil)
	fmt.Println("expected sign:", expected)

	if !hmac.Equal(signature, expected) {
		return id, token, fmt.Errorf("wrong token")
	}

	// распарсим payload
	var val s.CSRFvalue
	if err := json.Unmarshal(data, &val); err != nil {
		return id, token, fmt.Errorf("error while unurshall csrf data", err.Error())
	}

	// проверим срок жизни
	if time.Now().After(val.Exp) {
		new_scrf, err := CreateCSRF(secret, val.UserId)
		if err != nil {
			return id, new_scrf, fmt.Errorf("err while CreateCSRF", err.Error())
		}
		return val.UserId, new_scrf, nil
	}

	return val.UserId, token, nil
}

func SetSessionToCash() {

}

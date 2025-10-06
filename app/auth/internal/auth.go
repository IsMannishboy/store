package internal

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(password string, hash []byte) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return err
	}
	return nil
}
func MakePassword(passord string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passord), bcrypt.DefaultCost)
	if err != nil {
		return []byte("error"), err
	}
	return hash, nil
}
func CreateSession(id string, redis_db *redis.Client) (Session, error) {
	var session Session
	sesion_id := make([]byte, 32)
	_, err := rand.Read(sesion_id)
	if err != nil {
		return session, err
	}
	session.Id = hex.EncodeToString(sesion_id)
	session.UserId = id
	session.Exp = time.Now().Add(24 * time.Hour)
	return session, nil
}
func AddSessionToCash(ctx context.Context, session Session, redis_db *redis.Client, timeout int) error {
	newcontext, cancel := context.WithTimeout(ctx, time.Second*time.Duration(timeout))
	defer cancel()
	var value SessionValue
	value.Exp = session.Exp
	value.UserId = session.UserId
	err := redis_db.Set(newcontext, session.Id, value, 30*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}
func GetSession(redis_db *redis.Client, sessionID string, ctx context.Context, timeout int) (SessionValue, error) {
	newcontext, cancel := context.WithTimeout(ctx, time.Second*time.Duration(timeout))
	defer cancel()
	var value SessionValue
	err := redis_db.Get(newcontext, sessionID).Scan(&value)
	return value, err
}
func CreateCSRF(secret string, userId string) (string, error) {
	val := CSRFvalue{
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
func VerifyCSRF(secret string, token string) (string, error) {
	var payloadEnc, signatureEnc string
	n, err := fmt.Sscanf(token, "%s.%s", &payloadEnc, &signatureEnc)
	if err != nil || n != 2 {
		return token, err
	}

	// декодируем
	data, err := base64.URLEncoding.DecodeString(payloadEnc)
	if err != nil {
		return token, err
	}
	signature, err := base64.URLEncoding.DecodeString(signatureEnc)
	if err != nil {
		return token, err
	}

	// пересчёт HMAC
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(data)
	expected := h.Sum(nil)

	if !hmac.Equal(signature, expected) {
		return token, fmt.Errorf("wrong token")
	}

	// распарсим payload
	var val CSRFvalue
	if err := json.Unmarshal(data, &val); err != nil {
		return token, err
	}

	// проверим срок жизни
	if time.Now().After(val.Exp) {
		new_scrf, err := CreateCSRF(secret, val.UserId)
		if err != nil {
			return new_scrf, nil
		}
	}

	return token, nil
}

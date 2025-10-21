package internal

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func CheckSession(ctx context.Context, sessionid string, redis_db *redis.Client, RwTimeout time.Duration) error {
	newcontext, cancel := context.WithTimeout(ctx, RwTimeout*time.Second)
	defer cancel()
	var value SessionValue
	var strvalue string
	err := redis_db.Get(newcontext, sessionid).Scan(&strvalue)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(strvalue), &value)
	if err != nil {
		return err
	}
	if value.Role != "admin" {
		return fmt.Errorf("wrong role")
	}
	if time.Now().After(value.Exp) {
		return fmt.Errorf("session exp")
	}
	err = redis_db.Set(newcontext, sessionid, strvalue, 30*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}
func AddSessionToCash(ctx context.Context, session Session, RwTimeout time.Duration, redis_db *redis.Client) error {
	var value = SessionValue{UserId: session.UserId, Exp: session.Exp}
	baitiki, err := json.Marshal(value)
	if err != nil {
		return err
	}
	newctx, cancel := context.WithTimeout(ctx, RwTimeout*time.Second)
	defer cancel()
	err = redis_db.Set(newctx, session.Id, baitiki, time.Minute*30).Err()
	if err != nil {
		return err
	}
	return nil
}
func CheckCSRF(csrf string, secret []byte) (int, string, error) {
	var csrfdata, signature string
	parts := strings.SplitN(csrf, ".", 2)
	if len(parts) != 2 {
		return 0, csrf, fmt.Errorf("wrong csrf")
	}
	csrfdata = parts[0]
	signature = parts[1]
	data, err := base64.StdEncoding.DecodeString(csrfdata)
	if err != nil {
		return 0, csrf, err
	}
	sign, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return 0, csrf, err
	}
	h := hmac.New(sha256.New, secret)
	h.Write(data)
	expected := h.Sum(nil)
	if !hmac.Equal(expected, sign) {
		return 0, csrf, fmt.Errorf("signature doesnt much")
	}
	var csrfvalue CSRFvalue
	err = json.Unmarshal(data, &csrfvalue)
	if err != nil {
		return 0, csrf, err

	}
	if time.Now().After(csrfvalue.Exp) {
		newcsrf, err := CreateCSRF(secret, csrfvalue.UserId)
		if err != nil {
			return 0, csrf, err
		}
		return csrfvalue.UserId, newcsrf, nil

	}
	return csrfvalue.UserId, csrf, nil

}
func CreateCSRF(secret []byte, id int) (string, error) {
	var csrfvalue = CSRFvalue{UserId: id, Exp: time.Now().Add(time.Minute * 30)}
	baitiki, err := json.Marshal(csrfvalue)
	if err != nil {
		return "", err
	}
	h := hmac.New(sha256.New, secret)
	h.Write(baitiki)
	sign := h.Sum(nil)
	csrfdata := base64.URLEncoding.EncodeToString(baitiki)
	csrfsign := base64.URLEncoding.EncodeToString(sign)
	newcsrf := csrfdata + "." + csrfsign
	return newcsrf, nil

}

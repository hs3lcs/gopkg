package crypto

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JwtPack struct {
	ISS  string `json:"iss"`
	IAT  int64  `json:"iat"`
	EXP  int64  `json:"exp"`
	AGCY uint32 `json:"agcy"`
	ACCT uint32 `json:"acct"`
	Type uint8  `json:"type"`
	Role uint8  `json:"role"`
}

func JwtEncode(jwtpack JwtPack) (string, error) {
	exp := time.Duration(Config.JWT_EXP) * time.Minute
	claims := jwt.MapClaims{
		"iss":  "iamsvc",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(exp).Unix(),
		"agcy": jwtpack.AGCY,
		"acct": jwtpack.ACCT,
		"type": jwtpack.Type,
		"role": jwtpack.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(Config.JWT_KEY))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func JwtDecode(token string) (JwtPack, error) {
	var jwtpack JwtPack
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(Config.JWT_KEY), nil
	})
	if err != nil {
		return jwtpack, err
	}
	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		str, err := json.Marshal(claims)
		if err != nil {
			return jwtpack, err
		}
		err = json.Unmarshal(str, &jwtpack)
		return jwtpack, err
	}
	return jwtpack, errors.New("invalid token")
}

func JwtParse(token string) JwtPack {
	var jwtpack JwtPack
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return jwtpack
	}
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return jwtpack
	}
	json.Unmarshal(decoded, &jwtpack)
	return jwtpack
}

func HashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	return string(bytes), err
}

func ComparePassword(hash, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

const CHAR_HASH = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func HashString(n uint8) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = CHAR_HASH[rand.IntN(len(CHAR_HASH))]
	}
	return string(b)
}

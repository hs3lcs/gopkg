package crypto

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func JwtEncode(jwtClaims *JwtClaims) (string, error) {
	exp := time.Duration(Config.JWT_EXP) * time.Second
	claims := jwt.MapClaims{
		"iss":  jwtClaims.ISS,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(exp).Unix(),
		"uid":  jwtClaims.UID,
		"org":  jwtClaims.ORG,
		"type": jwtClaims.Type,
		"role": jwtClaims.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := MD5Hash(Config.JWT_KEY)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func JwtDecode(token string) (*JwtClaims, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		secretKey := MD5Hash(Config.JWT_KEY)
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		str, err := json.Marshal(claims)
		if err != nil {
			return nil, err
		}
		jwtpack := new(JwtClaims)
		err = json.Unmarshal(str, jwtpack)
		return jwtpack, err
	}
	return nil, errors.New("invalid token")
}

func JwtParse(token string) *JwtClaims {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil
	}
	jwtpack := new(JwtClaims)
	json.Unmarshal(decoded, jwtpack)
	return jwtpack
}

func HashPassword(pass string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), cost)
	return string(bytes), err
}

func ComparePassword(hash string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

func StringHash(n uint8) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = CHAR_HASH[rand.IntN(len(CHAR_HASH))]
	}
	return string(b)
}

func SHA256Hash(str string) string {
	hash := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", hash)
}

func MD5Hash(str string) string {
	hash := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", hash)
}

package middleware

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/mergermarket/go-pkcs7"

	"github.com/lestrrat-go/jwx/jwt"
)

//TokenAuth will be used by verifier
var TokenAuth *jwtauth.JWTAuth

// CtxKey will be used to fetch the jwt claims in db connection
type CtxKey struct {
	AccountID   string `json:"account_id"`
	Role        string `json:"role"`
	AccessToken string `json:"access_token"`
}

//ZenAPIVersion key to api version in context
type ZenAPIVersion struct {
}

//InitAuth the jwt token auth
func InitAuth(secret string) {
	TokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

//Authenticator custom auth middleware
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			JSON(w, http.StatusUnauthorized, Map{"status": "error", "data": NewError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))})
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			JSON(w, http.StatusUnauthorized, Map{"status": "error", "data": NewError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))})
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

// GetJwtClaims extracts the jwt claims from Authorization header
// this is passed into the DB context for postgres rls
func GetJwtClaims(r *http.Request, defaults ...string) (context.Context, *CtxKey) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	claimsCtx := &CtxKey{}
	if claims != nil && len(claims) > 0 {
		claimsCtx.AccountID = getClaimValue(claims, "account_id", "")
		claimsCtx.Role = getClaimValue(claims, "role", "zeno_anon")
		claimsCtx.AccessToken = jwtauth.TokenFromHeader(r)
	} else {
		if len(defaults) > 1 {
			claimsCtx.AccountID = defaults[0]
			claimsCtx.Role = defaults[1]
		} else if len(defaults) == 1 {
			claimsCtx.AccountID = defaults[0]
		}
	}
	zCtx := context.WithValue(r.Context(), CtxKey{}, claimsCtx)
	return zCtx, claimsCtx
}

func getClaimValue(claims map[string]interface{}, field string, defaultValue string) string {
	if val, ok := claims[field]; ok {
		return val.(string)
	}
	return defaultValue
}

// Encrypt encrypts plain text string into cipher text string
func Encrypt(unencrypted string, cipherKey string) (string, error) {
	key := []byte(cipherKey)
	plainText := []byte(unencrypted)
	plainText, err := pkcs7.Pad(plainText, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf(`plainText: "%s" has error`, plainText)
	}
	if len(plainText)%aes.BlockSize != 0 {
		err := fmt.Errorf(`plainText: "%s" has the wrong block size`, plainText)
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return fmt.Sprintf("%x", cipherText), nil
}

// Decrypt decrypts cipher text string into plain text string
func Decrypt(encrypted string, cipherKey string) (string, error) {
	key := []byte(cipherKey)
	cipherText, _ := hex.DecodeString(encrypted)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", NewError(http.StatusInternalServerError, "Filed to create cipher")
	}

	if len(cipherText) < aes.BlockSize {
		return "", NewError(http.StatusInternalServerError, "cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		return "", NewError(http.StatusInternalServerError, "cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, _ = pkcs7.Unpad(cipherText, aes.BlockSize)
	return fmt.Sprintf("%s", cipherText), nil
}

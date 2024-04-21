package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	corev1 "k8s.io/client-go/listers/core/v1"
)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

type JWT struct {
	secret []byte
}

func NewJWT(secret []byte) *JWT {
	return &JWT{secret: secret}
}

func (j *JWT) GenerateToken(username, password string) (string, error) {
	claims := Claims{
		Username: username,
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "openapp",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func JWTAuth(cmLister corev1.ConfigMapLister) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		systemCfg, err := cmLister.ConfigMaps(SystemNamespace).Get(SystemConfigMap)
		if err != nil {
			ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
			ctx.Abort()
			return
		}
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ReturnFormattedData(ctx, http.StatusUnauthorized, "Authorization token is required", nil)
			ctx.Abort()
			return
		}

		jwt := NewJWT([]byte(systemCfg.Data["password"]))
		if _, err := jwt.ParseToken(token); err != nil {
			ReturnFormattedData(ctx, http.StatusUnauthorized, err.Error(), nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

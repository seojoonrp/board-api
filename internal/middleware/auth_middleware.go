package middleware

import (
	"fmt"
	"net/http"
	"seojoonrp/board-api/internal/apperror"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return apperror.NewUnauthorized("missing or invalid authorization header")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			keyFunc := func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(secret), nil
			}

			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, keyFunc)
			if err != nil || !token.Valid {
				return apperror.NewUnauthorized("invalid or expired token")
			}

			claims, ok := token.Claims.(*JWTClaims)
			if !ok {
				return apperror.NewUnauthorized("invalid token claims")
			}

			c.Set("userID", claims.UserID)
			return next(c)
		}
	}
}

func UserIDFrom(c echo.Context) (string, error) {
	uid, ok := c.Get("userID").(string)
	if !ok || uid == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized)
	}
	return uid, nil
}

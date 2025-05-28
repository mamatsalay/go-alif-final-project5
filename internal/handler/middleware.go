package handler

import (
	"go.uber.org/dig"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"strings"
	"workout-tracker/internal/erorrs"
	"workout-tracker/internal/model/user"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MiddlewareParams struct {
	dig.In

	Log     *zap.SugaredLogger
	Service AuthService
}

type Middleware struct {
	Log     *zap.SugaredLogger
	Service AuthService
	Secret  string
}

func NewMiddleware(params MiddlewareParams) *Middleware {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
	return &Middleware{
		Log:     params.Log,
		Service: params.Service,
		Secret:  secret,
	}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			m.Log.Errorw("missing or invalid token ", "header", auth)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{erorrs.ErrorKey: "missing or invalid token"})
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				m.Log.Errorf("Unexpected signing method: %v", token.Header["alg"])
				return nil, jwt.ErrTokenSignatureInvalid
			}

			return []byte(m.Secret), nil
		})

		if err != nil {
			m.Log.Errorw("Unauthorized", "header", auth)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{erorrs.ErrorKey: "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			m.Log.Errorw("Unauthorized", "header", auth)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{erorrs.ErrorKey: "invalid claims"})
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			m.Log.Errorw("Unauthorized", "header", auth)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{erorrs.ErrorKey: "invalid user id in token"})
			return
		}

		userID := int(userIDFloat)

		roleStr, ok := claims["role"].(string)
		if !ok {
			m.Log.Errorw("Unauthorized", "header", auth)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{erorrs.ErrorKey: "invalid role in token"})
			return
		}

		role := user.Role(roleStr)

		versionFloat, ok := claims["version"].(float64)
		if !ok {
			m.Log.Errorw("missing version in token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{erorrs.ErrorKey: "invalid token version"})
			return
		}

		tokenVersion := int(versionFloat)

		u, err := m.Service.GetUserByUserID(c.Request.Context(), userID)
		if err != nil {
			m.Log.Errorw("failed to fetch user", "userID", userID, "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{erorrs.ErrorKey: "user not found"})
			return
		}
		if u.TokenVersion != tokenVersion {
			m.Log.Warnw("token version mismatch", "tokenVersion", tokenVersion, "dbVersion", u.TokenVersion)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{erorrs.ErrorKey: "token has been invalidated"})
			return
		}

		c.Set("userID", userID)
		c.Set("role", role)

		c.Next()
	}
}

func (m *Middleware) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		role, ok := roleValue.(user.Role)
		if !exists || !ok {
			m.Log.Errorw("Unauthorized", "reason", "missing or invalid role in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{erorrs.ErrorKey: "invalid role"})
			return
		}

		if role != user.AdminRole {
			m.Log.Errorw("Unauthorized", "reason", "user is not admin")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{erorrs.ErrorKey: "access denied"})
			return
		}

		c.Next()
	}
}

package utility

import (
	"redditBack/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("this-is-a-definitly-safe-key")

type Claims struct {
	UserID string
	jwt.RegisteredClaims
}
type UtilityFunctions struct {
	CacheRepo repository.CacheRepository
}

func NewUtility(cacheRepo repository.CacheRepository) UtilityFunctions {
	return UtilityFunctions{CacheRepo: cacheRepo}
}
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(1000 * time.Minute)

	claims := &Claims{
		UserID: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func (u *UtilityFunctions) JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		exist, _ := u.CacheRepo.IsTokenInvalid(c.Request.Context(), tokenString)
		if exist {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("this-is-a-definitly-safe-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["UserID"])
		c.Next()
	}
}

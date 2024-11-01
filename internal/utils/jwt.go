package utils

import (
	"errors"
	"time"

	"ddup-apis/internal/config"
	"ddup-apis/internal/db"
	"ddup-apis/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token 并保存到数据库
func GenerateToken(userID uint, username string) (string, time.Time, int64, time.Time, error) {
	cfg := config.GetConfig()
	now := time.Now()
	expiresIn := int64(cfg.JWT.ExpiresIn.Seconds())
	expiredAt := now.Add(cfg.JWT.ExpiresIn)
	createdAt := now

	// 生成新的 Token
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			IssuedAt:  jwt.NewNumericDate(createdAt),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", time.Time{}, 0, time.Time{}, err
	}

	// 将用户现有的 token 标记为无效并软删除
	if err := db.DB.Model(&model.Session{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Updates(map[string]interface{}{
			"is_valid":   false,
			"deleted_at": now,
		}).Error; err != nil {
		return "", time.Time{}, 0, time.Time{}, err
	}

	// 保存新 Token 到数据库
	newSession := model.Session{
		UserID:    userID,
		Token:     token,
		IsValid:   true,
		ExpiredAt: expiredAt,
	}

	if err = db.DB.Create(&newSession).Error; err != nil {
		return "", time.Time{}, 0, time.Time{}, err
	}

	return token, createdAt, expiresIn, expiredAt, nil
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig()
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

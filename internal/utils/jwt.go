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
func GenerateToken(userID uint, username string) (token string, createdAt time.Time, expiresIn int64, expiredAt time.Time, err error) {
	cfg := config.GetConfig()
	now := time.Now()
	expiresIn = int64(cfg.JWT.ExpiresIn.Seconds())
	expiredAt = now.Add(cfg.JWT.ExpiresIn)
	createdAt = now

	// 查找该用户所有有效的 Token
	var userSessions []model.UserSession
	err = db.DB.Where("user_id = ? AND is_valid = ? AND expired_at > ?", userID, true, now).
		Order("created_at DESC").
		Find(&userSessions).Error
	if err != nil {
		return "", time.Time{}, 0, time.Time{}, err
	}

	if len(userSessions) > 0 {
		// 保留最新的一个 Token，将其他的设置为无效
		latestSession := userSessions[0]
		if len(userSessions) > 1 {
			// 获取除最新 Token 外的所有 Token ID
			var oldSessionIDs []int64
			for i := 1; i < len(userSessions); i++ {
				oldSessionIDs = append(oldSessionIDs, userSessions[i].ID)
			}

			// 将旧 Token 标记为无效
			if err = db.DB.Model(&model.UserSession{}).
				Where("id IN ?", oldSessionIDs).
				Update("is_valid", false).Error; err != nil {
				return "", time.Time{}, 0, time.Time{}, err
			}
		}

		// 延长最新 Token 的过期时间
		expiredAt = now.Add(cfg.JWT.ExpiresIn)
		latestSession.ExpiredAt = expiredAt
		latestSession.UpdatedAt = now

		if err = db.DB.Save(&latestSession).Error; err != nil {
			return "", time.Time{}, 0, time.Time{}, err
		}

		return latestSession.Token, latestSession.CreatedAt, int64(cfg.JWT.ExpiresIn.Seconds()), expiredAt, nil
	}

	// 将该用户所有旧 Token 标记为无效（包括已过期的）
	if err = db.DB.Model(&model.UserSession{}).
		Where("user_id = ?", userID).
		Update("is_valid", false).Error; err != nil {
		return "", time.Time{}, 0, time.Time{}, err
	}

	// 生成新的 Token
	expiredAt = createdAt.Add(cfg.JWT.ExpiresIn)
	expiresIn = int64(cfg.JWT.ExpiresIn.Seconds())

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			IssuedAt:  jwt.NewNumericDate(createdAt),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", time.Time{}, 0, time.Time{}, err
	}

	// 保存新 Token 到数据库
	newUserSession := model.UserSession{
		UserID:    int64(userID),
		Token:     token,
		IsValid:   true,
		ExpiredAt: expiredAt,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	if err = db.DB.Create(&newUserSession).Error; err != nil {
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

// ShouldRefreshToken 检查是否需要刷新 Token
func ShouldRefreshToken(claims *Claims) bool {
	cfg := config.GetConfig()
	expiresAt := claims.ExpiresAt
	if expiresAt == nil {
		return false
	}

	timeUntilExpiry := time.Until(expiresAt.Time)
	return timeUntilExpiry < cfg.JWT.RefreshGracePeriod
}

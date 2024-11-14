package service

import (
	"context"
	"ddup-apis/internal/config"
	"ddup-apis/internal/db"
	"ddup-apis/internal/dto"
	"ddup-apis/internal/errors"
	"ddup-apis/internal/model"
	"ddup-apis/internal/repository"
	"ddup-apis/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type IUserService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) error
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetUserByID(ctx context.Context, id uint) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id uint, req *dto.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id uint) error
	ChangePassword(ctx context.Context, id uint, req *dto.ChangePasswordRequest) error
	ValidateToken(token string) (*TokenValidationResult, error)
	Logout(ctx context.Context, token string) error
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GenerateToken(ctx context.Context, userID uint, username string) (string, time.Time, int64, time.Time, error)
}

type UserService struct {
	userRepo    IUserRepository
	sessionRepo ISessionRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		userRepo:    repository.NewUserRepository(db),
		sessionRepo: repository.NewSessionRepository(db),
	}
}

// 1. 添加仓储层接口定义
type IUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, id uint, updates map[string]interface{}) error
	UpdatePassword(ctx context.Context, id uint, hashedPassword string) error
	UpdateLastLogin(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}

type ISessionRepository interface {
	CreateSession(ctx context.Context, session *model.Session) error
	GetSessionByToken(ctx context.Context, token string) (*model.Session, error)
	InvalidateSession(ctx context.Context, token string) error
	InvalidateUserSessions(ctx context.Context, userID uint) error
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 2. 将 TokenValidationResult 移到更合适的位置（比如 dto 包）
type TokenValidationResult struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Valid    bool   `json:"valid"`
}

// 3. 统一错误处理
func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) error {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return errors.Wrap(err, "查询用户失败")
	}
	if user != nil {
		return errors.New(400, "用户名已存在", nil)
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.Wrap(err, "密码加密失败")
	}

	// 创建用户
	user = &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Username,
	}

	return s.userRepo.Create(ctx, user)
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New(401, "用户名或密码错误", nil)
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New(401, "用户名或密码错误", nil)
	}

	// 生成token
	token, createdAt, expiresIn, expiredAt, err := s.GenerateToken(ctx, user.ID, user.Username)
	if err != nil {
		return nil, errors.Wrap(err, "生成token失败")
	}

	// 更新最后登录时间
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, errors.Wrap(err, "更新登录时间失败")
	}

	return &dto.LoginResponse{
		Token:     token,
		CreatedAt: createdAt,
		ExpiresIn: expiresIn,
		ExpiredAt: expiredAt,
		User: dto.UserResponse{
			Username:  user.Username,
			Email:     user.Email,
			Mobile:    user.Mobile,
			Location:  user.Location,
			Nickname:  user.Nickname,
			Bio:       user.Bio,
			Gender:    user.Gender,
			Birthday:  user.Birthday,
			Avatar:    user.Avatar,
			LastLogin: user.LastLogin,
		},
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New(404, "用户不存在", err)
	}

	return &dto.UserResponse{
		Username:  user.Username,
		Email:     user.Email,
		Mobile:    user.Mobile,
		Location:  user.Location,
		Nickname:  user.Nickname,
		Bio:       user.Bio,
		Gender:    user.Gender,
		Birthday:  user.Birthday,
		Avatar:    user.Avatar,
		LastLogin: user.LastLogin,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, req *dto.UpdateUserRequest) error {
	updates := make(map[string]interface{})

	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Mobile != "" {
		updates["mobile"] = req.Mobile
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Language != "" {
		updates["language"] = req.Language
	}

	return s.userRepo.Update(ctx, id, updates)
}

func (s *UserService) ChangePassword(ctx context.Context, id uint, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New(404, "用户不存在", err)
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return errors.New(400, "旧密码错误", nil)
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.Wrap(err, "密码加密失败")
	}

	return s.userRepo.UpdatePassword(ctx, id, hashedPassword)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) Logout(ctx context.Context, token string) error {
	return s.sessionRepo.InvalidateSession(ctx, token)
}

// 4. 优化 ValidateToken 方法
func (s *UserService) ValidateToken(token string) (*TokenValidationResult, error) {
	if token == "" {
		return &TokenValidationResult{Valid: false}, nil
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		return &TokenValidationResult{Valid: false}, errors.Wrap(err, "无效的token")
	}

	ctx := context.Background()
	session, err := s.sessionRepo.GetSessionByToken(ctx, token)
	if err != nil {
		return &TokenValidationResult{Valid: false}, errors.Wrap(err, "会话不存在")
	}

	if !session.IsValid || time.Now().After(session.ExpiredAt) {
		return &TokenValidationResult{Valid: false}, nil
	}

	return &TokenValidationResult{
		UserID:   claims.UserID,
		Username: claims.Username,
		Valid:    true,
	}, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

// GenerateToken 生成 JWT Token 并保存到数据库
func (s *UserService) GenerateToken(ctx context.Context, userID uint, username string) (string, time.Time, int64, time.Time, error) {
	cfg := config.GetConfig()
	now := time.Now()
	expiresIn := int64(cfg.JWT.ExpiresIn.Seconds())
	expiredAt := now.Add(cfg.JWT.ExpiresIn)
	createdAt := now

	// 生成新 Token
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
	if err := s.sessionRepo.InvalidateUserSessions(ctx, userID); err != nil {
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

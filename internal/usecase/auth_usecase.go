package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginResult struct {
	Token string
	User  *entity.AppUser
}

type jwtClaims struct {
	UserID        string `json:"user_id"`
	CooperativeID string `json:"cooperative_id"`
	Role          string `json:"role"`
	MemberID      string `json:"member_id"`
	jwt.RegisteredClaims
}

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (*LoginResult, error)
}

type authUsecase struct {
	userRepo      repository.UserRepository
	jwtSecret     string
	jwtExpiryHours int
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtSecret string, jwtExpiryHours int) AuthUsecase {
	return &authUsecase{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
	}
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	memberID := ""
	if user.MemberID != nil {
		memberID = user.MemberID.String()
	}

	claims := jwtClaims{
		UserID:        user.ID.String(),
		CooperativeID: user.CooperativeID.String(),
		Role:          user.Role,
		MemberID:      memberID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(u.jwtExpiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &LoginResult{Token: tokenStr, User: user}, nil
}

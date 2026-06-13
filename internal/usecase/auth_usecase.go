package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrEmailAlreadyUsed      = errors.New("email sudah terdaftar")
	ErrNotMemberOfCooperative = errors.New("user bukan anggota koperasi ini")
)

type LoginResult struct {
	Token       string
	User        *entity.AppUser
	Memberships []repository.MembershipInfo // non-empty = anggota punya >1 koperasi, perlu pilih
}

type jwtClaims struct {
	UserID        string `json:"user_id"`
	CooperativeID string `json:"cooperative_id"`
	Role          string `json:"role"`
	MemberID      string `json:"member_id"`
	jwt.RegisteredClaims
}

type AuthUsecase interface {
	Login(ctx context.Context, identifier, password string) (*LoginResult, error)
	SelectCooperative(ctx context.Context, identifier, password, cooperativeID string) (*LoginResult, error)
	RegisterPengurus(ctx context.Context, req *model.RegisterPengurusRequest) (*LoginResult, error)
	GetUserCooperatives(ctx context.Context, userID string) ([]repository.MembershipInfo, error)
}

type authUsecase struct {
	userRepo       repository.UserRepository
	jwtSecret      string
	jwtExpiryHours int
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtSecret string, jwtExpiryHours int) AuthUsecase {
	return &authUsecase{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
	}
}

func (u *authUsecase) Login(ctx context.Context, identifier, password string) (*LoginResult, error) {
	user, err := u.findAndVerify(ctx, identifier, password)
	if err != nil {
		return nil, err
	}

	memberships, err := u.userRepo.FindMembershipsByUserID(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	// Anggota dengan lebih dari satu koperasi: tetap issue token (scoped ke koperasi pertama),
	// dan sertakan daftar koperasi agar FE bisa prompt user untuk switch.
	if len(memberships) > 1 {
		m := memberships[0]
		coopID, _ := uuid.Parse(m.CooperativeID)
		user.CooperativeID = coopID
		if m.MemberID != "" {
			mid, _ := uuid.Parse(m.MemberID)
			user.MemberID = &mid
		}
		result, err := u.buildLoginResult(user)
		if err != nil {
			return nil, err
		}
		result.Memberships = memberships
		return result, nil
	}

	// Single membership: gunakan data dari tabel membership jika ada.
	if len(memberships) == 1 {
		m := memberships[0]
		coopID, _ := uuid.Parse(m.CooperativeID)
		user.CooperativeID = coopID
		if m.MemberID != "" {
			mid, _ := uuid.Parse(m.MemberID)
			user.MemberID = &mid
		}
	}

	return u.buildLoginResult(user)
}

func (u *authUsecase) SelectCooperative(ctx context.Context, identifier, password, cooperativeID string) (*LoginResult, error) {
	user, err := u.findAndVerify(ctx, identifier, password)
	if err != nil {
		return nil, err
	}

	membership, err := u.userRepo.FindMembership(ctx, user.ID.String(), cooperativeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotMemberOfCooperative
		}
		return nil, err
	}

	user.CooperativeID = membership.CooperativeID
	user.MemberID = membership.MemberID

	return u.buildLoginResult(user)
}

func (u *authUsecase) RegisterPengurus(ctx context.Context, req *model.RegisterPengurusRequest) (*LoginResult, error) {
	if _, err := uuid.Parse(req.CooperativeID); err != nil {
		return nil, err
	}

	if _, err := u.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, ErrEmailAlreadyUsed
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	coopID, err := uuid.Parse(req.CooperativeID)
	if err != nil {
		return nil, err
	}

	user := &entity.AppUser{
		CooperativeID: coopID,
		Email:         req.Email,
		PasswordHash:  string(hashedPassword),
		Role:          "pengurus",
		FullName:      req.FullName,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	membership := &entity.UserCooperativeMembership{
		UserID:        user.ID,
		CooperativeID: coopID,
		MemberID:      nil,
	}
	_ = u.userRepo.CreateMembership(ctx, membership) // best-effort; user sudah terbuat

	return u.buildLoginResult(user)
}

func (u *authUsecase) GetUserCooperatives(ctx context.Context, userID string) ([]repository.MembershipInfo, error) {
	return u.userRepo.FindMembershipsByUserID(ctx, userID)
}

func (u *authUsecase) findAndVerify(ctx context.Context, identifier, password string) (*entity.AppUser, error) {
	user, err := u.userRepo.FindByEmail(ctx, identifier)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (u *authUsecase) buildLoginResult(user *entity.AppUser) (*LoginResult, error) {
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

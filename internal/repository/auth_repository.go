package repository

import (
	"context"
	"errors"

	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

// MembershipInfo is a flat view of user_cooperative_membership joined with cooperative name.
type MembershipInfo struct {
	CooperativeID   string
	CooperativeName string
	MemberID        string // empty string for pengurus
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.AppUser, error)
	Create(ctx context.Context, user *entity.AppUser) error
	FindMembershipsByUserID(ctx context.Context, userID string) ([]MembershipInfo, error)
	FindMembership(ctx context.Context, userID, cooperativeID string) (*entity.UserCooperativeMembership, error)
	CreateMembership(ctx context.Context, m *entity.UserCooperativeMembership) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.AppUser, error) {
	var user entity.AppUser
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.AppUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindMembershipsByUserID(ctx context.Context, userID string) ([]MembershipInfo, error) {
	type row struct {
		CooperativeID   string `gorm:"column:cooperative_id"`
		CooperativeName string `gorm:"column:cooperative_name"`
		MemberID        string `gorm:"column:member_id"`
	}
	var rows []row
	err := r.db.WithContext(ctx).
		Table("user_cooperative_membership ucm").
		Select("ucm.cooperative_id::text AS cooperative_id, c.name AS cooperative_name, COALESCE(ucm.member_id::text, '') AS member_id").
		Joins("JOIN cooperative c ON c.id = ucm.cooperative_id").
		Where("ucm.user_id = ?", userID).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make([]MembershipInfo, len(rows))
	for i, r := range rows {
		result[i] = MembershipInfo{
			CooperativeID:   r.CooperativeID,
			CooperativeName: r.CooperativeName,
			MemberID:        r.MemberID,
		}
	}
	return result, nil
}

func (r *userRepository) FindMembership(ctx context.Context, userID, cooperativeID string) (*entity.UserCooperativeMembership, error) {
	var m entity.UserCooperativeMembership
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND cooperative_id = ?", userID, cooperativeID).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *userRepository) CreateMembership(ctx context.Context, m *entity.UserCooperativeMembership) error {
	return r.db.WithContext(ctx).
		Where(entity.UserCooperativeMembership{UserID: m.UserID, CooperativeID: m.CooperativeID}).
		FirstOrCreate(m).Error
}

package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"gorm.io/gorm"
)

var (
	ErrMemberNotFound = errors.New("anggota tidak ditemukan")
	ErrDuplicateNIK   = errors.New("nik sudah terdaftar di koperasi ini")
)

type MemberRepository interface {
	Create(ctx context.Context, member *entity.Member) error
	FindByID(ctx context.Context, cooperativeID, memberID string) (*entity.Member, error)
	FindByNIK(ctx context.Context, cooperativeID, nik string) (*entity.Member, error)
	FindAll(ctx context.Context, cooperativeID, search, status string) ([]*entity.Member, error)
	Update(ctx context.Context, member *entity.Member) error
	GetSavingsSummary(ctx context.Context, memberID string) (*model.SavingsSummary, error)
}

type memberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{db: db}
}

func (r *memberRepository) Create(ctx context.Context, member *entity.Member) error {
	err := r.db.WithContext(ctx).Create(member).Error
	if err != nil {
		if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key") {
			return ErrDuplicateNIK
		}
		return err
	}
	return nil
}

func (r *memberRepository) FindByID(ctx context.Context, cooperativeID, memberID string) (*entity.Member, error) {
	var member entity.Member
	err := r.db.WithContext(ctx).
		Where("id = ? AND cooperative_id = ?", memberID, cooperativeID).
		First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) FindByNIK(ctx context.Context, cooperativeID, nik string) (*entity.Member, error) {
	var member entity.Member
	err := r.db.WithContext(ctx).
		Where("nik = ? AND cooperative_id = ?", nik, cooperativeID).
		First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) FindAll(ctx context.Context, cooperativeID, search, status string) ([]*entity.Member, error) {
	var members []*entity.Member
	query := r.db.WithContext(ctx).Where("cooperative_id = ?", cooperativeID)
	if search != "" {
		q := "%" + search + "%"
		query = query.Where("(full_name ILIKE ? OR nik ILIKE ?)", q, q)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("joined_at DESC").Find(&members).Error
	return members, err
}

func (r *memberRepository) Update(ctx context.Context, member *entity.Member) error {
	return r.db.WithContext(ctx).Save(member).Error
}

type savingsSummaryRow struct {
	Pokok    int64 `gorm:"column:pokok"`
	Wajib    int64 `gorm:"column:wajib"`
	Sukarela int64 `gorm:"column:sukarela"`
}

func (r *memberRepository) GetSavingsSummary(ctx context.Context, memberID string) (*model.SavingsSummary, error) {
	// Validate memberID is a valid UUID before querying
	if _, err := uuid.Parse(memberID); err != nil {
		return &model.SavingsSummary{}, nil
	}

	var row savingsSummaryRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(SUM(amount) FILTER (WHERE savings_type='pokok'    AND direction='setor'), 0)
		  - COALESCE(SUM(amount) FILTER (WHERE savings_type='pokok'    AND direction='tarik'), 0) AS pokok,
			COALESCE(SUM(amount) FILTER (WHERE savings_type='wajib'    AND direction='setor'), 0)
		  - COALESCE(SUM(amount) FILTER (WHERE savings_type='wajib'    AND direction='tarik'), 0) AS wajib,
			COALESCE(SUM(amount) FILTER (WHERE savings_type='sukarela' AND direction='setor'), 0)
		  - COALESCE(SUM(amount) FILTER (WHERE savings_type='sukarela' AND direction='tarik'), 0) AS sukarela
		FROM savings_transaction WHERE member_id = ?
	`, memberID).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	return &model.SavingsSummary{
		Pokok:    row.Pokok,
		Wajib:    row.Wajib,
		Sukarela: row.Sukarela,
	}, nil
}

// parseBirthDate converts "YYYY-MM-DD" string to *time.Time for entity.
func ParseBirthDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}

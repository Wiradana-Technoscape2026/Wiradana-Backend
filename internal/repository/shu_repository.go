package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

var (
	ErrShuPeriodNotFound = errors.New("periode SHU tidak ditemukan")
)

type ShuRepository interface {
	CreatePeriod(ctx context.Context, p *entity.ShuPeriod) error
	FindPeriods(ctx context.Context, coopID string) ([]entity.ShuPeriod, error)
	FindPeriodByID(ctx context.Context, coopID, periodID string) (*entity.ShuPeriod, error)
	UpdatePeriodStatus(ctx context.Context, periodID, status string) error
	CreateDistributions(ctx context.Context, dists []entity.ShuDistribution) error
	FindDistributionsByPeriod(ctx context.Context, periodID string) ([]ShuDistributionWithName, error)
	FindDistributionsByMember(ctx context.Context, memberID string) ([]ShuDistributionWithPeriod, error)
	GetTotalSimpananAktif(ctx context.Context, coopID string) (int64, error)
	GetTotalJasaPinjaman(ctx context.Context, coopID string, periodYear int) (int64, error)
	GetSimpananMember(ctx context.Context, memberID string) (int64, error)
	GetJasaPinjamanMember(ctx context.Context, memberID string, periodYear int) (int64, error)
	GetActiveMembers(ctx context.Context, coopID string) ([]entity.Member, error)
}

type ShuDistributionWithName struct {
	entity.ShuDistribution
	MemberName string `gorm:"column:member_name"`
}

type ShuDistributionWithPeriod struct {
	entity.ShuDistribution
	MemberName string `gorm:"column:member_name"`
	PeriodYear int    `gorm:"column:period_year"`
}

type shuRepository struct {
	db *gorm.DB
}

func NewShuRepository(db *gorm.DB) ShuRepository {
	return &shuRepository{db: db}
}

func (r *shuRepository) CreatePeriod(ctx context.Context, p *entity.ShuPeriod) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *shuRepository) FindPeriods(ctx context.Context, coopID string) ([]entity.ShuPeriod, error) {
	var periods []entity.ShuPeriod
	err := r.db.WithContext(ctx).
		Where("cooperative_id = ?", coopID).
		Order("year DESC").
		Find(&periods).Error
	return periods, err
}

func (r *shuRepository) FindPeriodByID(ctx context.Context, coopID, periodID string) (*entity.ShuPeriod, error) {
	var p entity.ShuPeriod
	err := r.db.WithContext(ctx).
		Where("id = ? AND cooperative_id = ?", periodID, coopID).
		First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrShuPeriodNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *shuRepository) UpdatePeriodStatus(ctx context.Context, periodID, status string) error {
	return r.db.WithContext(ctx).Model(&entity.ShuPeriod{}).
		Where("id = ?", periodID).
		Update("status", status).Error
}

func (r *shuRepository) CreateDistributions(ctx context.Context, dists []entity.ShuDistribution) error {
	if len(dists) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&dists).Error
}

func (r *shuRepository) FindDistributionsByPeriod(ctx context.Context, periodID string) ([]ShuDistributionWithName, error) {
	var result []ShuDistributionWithName
	err := r.db.WithContext(ctx).
		Table("shu_distribution").
		Select("shu_distribution.*, member.full_name AS member_name").
		Joins("JOIN member ON member.id = shu_distribution.member_id").
		Where("shu_distribution.shu_period_id = ?", periodID).
		Order("member.full_name ASC").
		Find(&result).Error
	return result, err
}

func (r *shuRepository) FindDistributionsByMember(ctx context.Context, memberID string) ([]ShuDistributionWithPeriod, error) {
	var result []ShuDistributionWithPeriod
	err := r.db.WithContext(ctx).
		Table("shu_distribution").
		Select("shu_distribution.*, member.full_name AS member_name, shu_period.year AS period_year").
		Joins("JOIN member ON member.id = shu_distribution.member_id").
		Joins("JOIN shu_period ON shu_period.id = shu_distribution.shu_period_id").
		Where("shu_distribution.member_id = ?", memberID).
		Order("shu_period.year DESC").
		Find(&result).Error
	return result, err
}

func (r *shuRepository) GetTotalSimpananAktif(ctx context.Context, coopID string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(saldo), 0) FROM (
			SELECT member_id,
				COALESCE(SUM(amount) FILTER (WHERE savings_type='pokok' AND direction='setor'), 0)
			  + COALESCE(SUM(amount) FILTER (WHERE savings_type='wajib' AND direction='setor'), 0)
			  - COALESCE(SUM(amount) FILTER (WHERE savings_type='wajib' AND direction='tarik'), 0) AS saldo
			FROM savings_transaction
			WHERE cooperative_id = ?
			GROUP BY member_id
		) sub
		JOIN member m ON m.id = sub.member_id AND m.status = 'aktif'
	`, coopID).Scan(&total).Error
	return total, err
}

func (r *shuRepository) GetTotalJasaPinjaman(ctx context.Context, coopID string, periodYear int) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(p.amount), 0)
		FROM payment p
		JOIN installment_schedule i ON i.id = p.schedule_id
		JOIN loan l ON l.id = i.loan_id
		WHERE l.cooperative_id = ?
		  AND i.interest_due > 0
		  AND EXTRACT(YEAR FROM p.paid_at) = ?
	`, coopID, periodYear).Scan(&total).Error
	return total, err
}

func (r *shuRepository) GetSimpananMember(ctx context.Context, memberID string) (int64, error) {
	mid, err := uuid.Parse(memberID)
	if err != nil {
		return 0, nil
	}
	var total int64
	err = r.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(SUM(amount) FILTER (WHERE savings_type='pokok' AND direction='setor'), 0)
		  + COALESCE(SUM(amount) FILTER (WHERE savings_type='wajib' AND direction='setor'), 0)
		  - COALESCE(SUM(amount) FILTER (WHERE savings_type='wajib' AND direction='tarik'), 0)
		FROM savings_transaction
		WHERE member_id = ?
	`, mid).Scan(&total).Error
	return total, err
}

func (r *shuRepository) GetJasaPinjamanMember(ctx context.Context, memberID string, periodYear int) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(p.amount), 0)
		FROM payment p
		JOIN installment_schedule i ON i.id = p.schedule_id
		JOIN loan l ON l.id = i.loan_id
		WHERE l.member_id = ?
		  AND i.interest_due > 0
		  AND EXTRACT(YEAR FROM p.paid_at) = ?
	`, memberID, periodYear).Scan(&total).Error
	return total, err
}

func (r *shuRepository) GetActiveMembers(ctx context.Context, coopID string) ([]entity.Member, error) {
	var members []entity.Member
	err := r.db.WithContext(ctx).
		Where("cooperative_id = ? AND status = ?", coopID, "aktif").
		Order("full_name ASC").
		Find(&members).Error
	return members, err
}

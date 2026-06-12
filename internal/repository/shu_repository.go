package repository

import (
	"context"
	"errors"

	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

var (
	ErrShuPeriodNotFound = errors.New("periode SHU tidak ditemukan")
	ErrShuPeriodNotDraft = errors.New("periode SHU bukan dalam status draft")
)

type ShuRepository interface {
	CreatePeriod(ctx context.Context, sp *entity.ShuPeriod) error
	FindPeriods(ctx context.Context, coopID string) ([]entity.ShuPeriod, error)
	FindPeriodByID(ctx context.Context, periodID string) (*entity.ShuPeriod, error)
	UpdatePeriodStatus(ctx context.Context, periodID, status string) error
	CreateDistributions(ctx context.Context, dists []entity.ShuDistribution) error
	FindDistributionsByPeriod(ctx context.Context, periodID string) ([]ShuDistributionWithMember, error)
	FindDistributionsByMember(ctx context.Context, memberID string) ([]ShuDistributionWithMember, error)
	// SHU calculation queries
	GetTotalSimpananAktif(ctx context.Context, coopID string) (int64, error)
	GetTotalJasaPinjaman(ctx context.Context, coopID string) (int64, error)
	GetSimpananPerMember(ctx context.Context, coopID string) ([]MemberSimpananRow, error)
	GetJasaPinjamanPerMember(ctx context.Context, coopID string) ([]MemberJasaPinjamanRow, error)
	GetActiveMemberIDs(ctx context.Context, coopID string) ([]string, error)
}

type ShuDistributionWithMember struct {
	entity.ShuDistribution
	MemberName string `gorm:"column:member_name"`
}

type MemberSimpananRow struct {
	MemberID       string `gorm:"column:member_id"`
	TotalSimpanan  int64  `gorm:"column:total_simpanan"`
}

type MemberJasaPinjamanRow struct {
	MemberID         string `gorm:"column:member_id"`
	TotalJasaPinjaman int64 `gorm:"column:total_jasa_pinjaman"`
}

type shuRepository struct {
	db *gorm.DB
}

func NewShuRepository(db *gorm.DB) ShuRepository {
	return &shuRepository{db: db}
}

func (r *shuRepository) CreatePeriod(ctx context.Context, sp *entity.ShuPeriod) error {
	return r.db.WithContext(ctx).Create(sp).Error
}

func (r *shuRepository) FindPeriods(ctx context.Context, coopID string) ([]entity.ShuPeriod, error) {
	var periods []entity.ShuPeriod
	err := r.db.WithContext(ctx).
		Where("cooperative_id = ?", coopID).
		Order("year DESC").
		Find(&periods).Error
	return periods, err
}

func (r *shuRepository) FindPeriodByID(ctx context.Context, periodID string) (*entity.ShuPeriod, error) {
	var sp entity.ShuPeriod
	err := r.db.WithContext(ctx).Where("id = ?", periodID).First(&sp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrShuPeriodNotFound
	}
	return &sp, err
}

func (r *shuRepository) UpdatePeriodStatus(ctx context.Context, periodID, status string) error {
	return r.db.WithContext(ctx).Model(&entity.ShuPeriod{}).
		Where("id = ?", periodID).Update("status", status).Error
}

func (r *shuRepository) CreateDistributions(ctx context.Context, dists []entity.ShuDistribution) error {
	return r.db.WithContext(ctx).Create(&dists).Error
}

func (r *shuRepository) FindDistributionsByPeriod(ctx context.Context, periodID string) ([]ShuDistributionWithMember, error) {
	var result []ShuDistributionWithMember
	err := r.db.WithContext(ctx).
		Table("shu_distribution").
		Select("shu_distribution.*, member.full_name as member_name").
		Joins("JOIN member ON member.id = shu_distribution.member_id").
		Where("shu_distribution.shu_period_id = ?", periodID).
		Scan(&result).Error
	return result, err
}

func (r *shuRepository) FindDistributionsByMember(ctx context.Context, memberID string) ([]ShuDistributionWithMember, error) {
	var result []ShuDistributionWithMember
	err := r.db.WithContext(ctx).
		Table("shu_distribution").
		Select("shu_distribution.*, member.full_name as member_name").
		Joins("JOIN member ON member.id = shu_distribution.member_id").
		Where("shu_distribution.member_id = ?", memberID).
		Scan(&result).Error
	return result, err
}

// GetTotalSimpananAktif returns total simpanan (pokok+wajib) semua anggota aktif
func (r *shuRepository) GetTotalSimpananAktif(ctx context.Context, coopID string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(
			CASE WHEN direction='setor' THEN amount ELSE -amount END
		), 0)
		FROM savings_transaction st
		JOIN member m ON m.id = st.member_id
		WHERE m.cooperative_id = ? AND m.status = 'aktif'
		  AND st.savings_type IN ('pokok', 'wajib')`, coopID).Scan(&total).Error
	return total, err
}

// GetTotalJasaPinjaman returns total interest payments across all loans in cooperative
func (r *shuRepository) GetTotalJasaPinjaman(ctx context.Context, coopID string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(i.interest_due), 0)
		FROM installment_schedule i
		JOIN loan l ON l.id = i.loan_id
		WHERE l.cooperative_id = ? AND i.status = 'lunas'`, coopID).Scan(&total).Error
	return total, err
}

// GetSimpananPerMember returns simpanan (pokok+wajib) per anggota aktif
func (r *shuRepository) GetSimpananPerMember(ctx context.Context, coopID string) ([]MemberSimpananRow, error) {
	var rows []MemberSimpananRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT st.member_id,
			COALESCE(SUM(CASE WHEN st.direction='setor' THEN st.amount ELSE -st.amount END), 0) AS total_simpanan
		FROM savings_transaction st
		JOIN member m ON m.id = st.member_id
		WHERE m.cooperative_id = ? AND m.status = 'aktif'
		  AND st.savings_type IN ('pokok', 'wajib')
		GROUP BY st.member_id`, coopID).Scan(&rows).Error
	return rows, err
}

// GetJasaPinjamanPerMember returns interest paid per member
func (r *shuRepository) GetJasaPinjamanPerMember(ctx context.Context, coopID string) ([]MemberJasaPinjamanRow, error) {
	var rows []MemberJasaPinjamanRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT l.member_id,
			COALESCE(SUM(i.interest_due), 0) AS total_jasa_pinjaman
		FROM installment_schedule i
		JOIN loan l ON l.id = i.loan_id
		WHERE l.cooperative_id = ? AND i.status = 'lunas'
		GROUP BY l.member_id`, coopID).Scan(&rows).Error
	return rows, err
}

// GetActiveMemberIDs returns all active member IDs for a cooperative
func (r *shuRepository) GetActiveMemberIDs(ctx context.Context, coopID string) ([]string, error) {
	var ids []string
	err := r.db.WithContext(ctx).Model(&entity.Member{}).
		Where("cooperative_id = ? AND status = 'aktif'", coopID).
		Pluck("id", &ids).Error
	return ids, err
}

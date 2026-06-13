package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/model/converter"
	"github.com/wiradana/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

var (
	ErrMemberNotFound = errors.New("anggota tidak ditemukan")
	ErrDuplicateNIK   = errors.New("nik sudah terdaftar di koperasi ini")
	ErrAccountExists  = errors.New("akun login untuk NIK ini sudah ada")
)

type MemberUsecase interface {
	Create(ctx context.Context, cooperativeID string, req *model.CreateMemberRequest) (*model.MemberResponse, error)
	FindByID(ctx context.Context, cooperativeID, memberID string) (*model.MemberResponse, error)
	FindAll(ctx context.Context, cooperativeID, search, status string) ([]model.MemberResponse, error)
	Update(ctx context.Context, cooperativeID, memberID string, req *model.UpdateMemberRequest) (*model.MemberResponse, error)
}

type memberUsecase struct {
	memberRepo repository.MemberRepository
	userRepo   repository.UserRepository
}

func NewMemberUsecase(memberRepo repository.MemberRepository, userRepo repository.UserRepository) MemberUsecase {
	return &memberUsecase{memberRepo: memberRepo, userRepo: userRepo}
}

func (u *memberUsecase) Create(ctx context.Context, cooperativeID string, req *model.CreateMemberRequest) (*model.MemberResponse, error) {
	coopUUID, err := uuid.Parse(cooperativeID)
	if err != nil {
		return nil, errors.New("cooperative_id tidak valid")
	}

	// NIK sudah terdaftar → kembalikan data existing tanpa error.
	existing, err := u.memberRepo.FindByNIK(ctx, cooperativeID, req.NIK)
	if err == nil {
		summary, _ := u.memberRepo.GetSavingsSummary(ctx, existing.ID.String())
		resp := converter.ToMemberResponse(existing, summary)
		return &resp, nil
	}
	if !errors.Is(err, repository.ErrMemberNotFound) {
		return nil, err
	}

	hasPassword := req.Password != nil && *req.Password != ""

	attrs := req.CustomAttributes
	if attrs == nil {
		attrs = datatypes.JSON("{}")
	}

	member := &entity.Member{
		CooperativeID:    coopUUID,
		NIK:              req.NIK,
		FullName:         req.FullName,
		Status:           "aktif",
		CustomAttributes: attrs,
	}

	if req.Address != "" {
		member.Address = &req.Address
	}
	if req.PhoneNumber != "" {
		member.PhoneNumber = &req.PhoneNumber
	}
	if bd := repository.ParseBirthDate(req.BirthDate); bd != nil {
		member.BirthDate = bd
	}

	if err := u.memberRepo.Create(ctx, member); err != nil {
		// Race condition: NIK terdaftar di antara cek dan insert — kembalikan data existing.
		if errors.Is(err, repository.ErrDuplicateNIK) {
			existing, ferr := u.memberRepo.FindByNIK(ctx, cooperativeID, req.NIK)
			if ferr == nil {
				summary, _ := u.memberRepo.GetSavingsSummary(ctx, existing.ID.String())
				resp := converter.ToMemberResponse(existing, summary)
				return &resp, nil
			}
		}
		return nil, err
	}

	if hasPassword {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		memberID := member.ID

		// NIK disimpan di kolom email sebagai identifier login anggota.
		// Satu NIK → satu AppUser global; multi-koperasi ditangani via membership.
		existingUser, err := u.userRepo.FindByEmail(ctx, member.NIK)
		if err == nil {
			// AppUser sudah ada (anggota terdaftar di koperasi lain) — cukup tambah membership.
			membership := &entity.UserCooperativeMembership{
				UserID:        existingUser.ID,
				CooperativeID: coopUUID,
				MemberID:      &memberID,
			}
			_ = u.userRepo.CreateMembership(ctx, membership)
		} else {
			// AppUser belum ada — buat baru + membership.
			appUser := &entity.AppUser{
				CooperativeID: coopUUID,
				MemberID:      &memberID,
				Email:         member.NIK,
				PasswordHash:  string(hash),
				Role:          "anggota",
			}
			if err := u.userRepo.Create(ctx, appUser); err != nil {
				if isDuplicateEmail(err) {
					return nil, ErrAccountExists
				}
				return nil, err
			}
			membership := &entity.UserCooperativeMembership{
				UserID:        appUser.ID,
				CooperativeID: coopUUID,
				MemberID:      &memberID,
			}
			_ = u.userRepo.CreateMembership(ctx, membership)
		}
	}

	summary, _ := u.memberRepo.GetSavingsSummary(ctx, member.ID.String())
	resp := converter.ToMemberResponse(member, summary)
	return &resp, nil
}

func isDuplicateEmail(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "23505") || strings.Contains(msg, "duplicate key")
}

func (u *memberUsecase) FindByID(ctx context.Context, cooperativeID, memberID string) (*model.MemberResponse, error) {
	member, err := u.memberRepo.FindByID(ctx, cooperativeID, memberID)
	if err != nil {
		if errors.Is(err, repository.ErrMemberNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	summary, _ := u.memberRepo.GetSavingsSummary(ctx, memberID)
	resp := converter.ToMemberResponse(member, summary)
	return &resp, nil
}

func (u *memberUsecase) FindAll(ctx context.Context, cooperativeID, search, status string) ([]model.MemberResponse, error) {
	members, err := u.memberRepo.FindAll(ctx, cooperativeID, search, status)
	if err != nil {
		return nil, err
	}

	responses := make([]model.MemberResponse, 0, len(members))
	for _, m := range members {
		summary, _ := u.memberRepo.GetSavingsSummary(ctx, m.ID.String())
		responses = append(responses, converter.ToMemberResponse(m, summary))
	}
	return responses, nil
}

func (u *memberUsecase) Update(ctx context.Context, cooperativeID, memberID string, req *model.UpdateMemberRequest) (*model.MemberResponse, error) {
	member, err := u.memberRepo.FindByID(ctx, cooperativeID, memberID)
	if err != nil {
		if errors.Is(err, repository.ErrMemberNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	if req.FullName != nil {
		member.FullName = *req.FullName
	}
	if req.Address != nil {
		member.Address = req.Address
	}
	if req.PhoneNumber != nil {
		member.PhoneNumber = req.PhoneNumber
	}
	if req.Status != nil {
		member.Status = *req.Status
	}
	if req.CustomAttributes != nil {
		member.CustomAttributes = req.CustomAttributes
	}

	if err := u.memberRepo.Update(ctx, member); err != nil {
		return nil, err
	}

	summary, _ := u.memberRepo.GetSavingsSummary(ctx, memberID)
	resp := converter.ToMemberResponse(member, summary)
	return &resp, nil
}

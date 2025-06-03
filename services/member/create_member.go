package member

import (
	"context"
	"database/sql"
	"errors"
	"kithli-api/models"
	"kithli-api/repositories"
)

// Learning note: this is basically our class with dependency injection
type MemberService struct {
	DB   *sql.DB
	RepoFactory func(tx *sql.Tx) repositories.IMemberRepository
}

// Learning note: this is a constructor
func NewMemberService(db *sql.DB) *MemberService {
	return &MemberService{
		DB: db,
		RepoFactory: func(tx *sql.Tx) repositories.IMemberRepository {
			return repositories.NewMemberRepository(tx)
		},
	}
}


func (memberService *MemberService) CreateMember(ctx context.Context, req models.MemberRequest) (int, error) {
	if req.UID == "" || req.MyHeadline == "" {
		return 0, errors.New("UID and myHeadline are required")
	}

	tx, err := memberService.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	repo := memberService.RepoFactory(tx)

	address := &models.Address{
		Street:     req.StreetAddress,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		AptNumber:  req.AptNumber,
	}

	addressID, err := repo.InsertAddress(ctx, address)
	if err != nil {
		return 0, err
	}

	member := &models.Member{
		MyHeadline:            req.MyHeadline,
		AboutMe:               req.AboutMe,
		AdditionalInformation: req.AdditionalInformation,
	}

	memberID, err := repo.InsertMember(ctx, member)
	if err != nil {
		return 0, err
	}

	err = repo.LinkAddress(ctx, memberID, addressID)
	if err != nil {
		return 0, err
	}

	err = repo.UpdateUserMember(ctx, req.UID, memberID)
	if err != nil {
		return 0, err
	}

	return memberID, nil
}

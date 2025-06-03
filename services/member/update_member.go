package member

import (
	"context"
	"errors"
	"kithli-api/models"
)

func (memberService *MemberService) UpdateMember(ctx context.Context, req models.MemberRequest) (int, error) {
	if req.UID == "" || req.MyHeadline == "" {
		return 0, errors.New("UID and myHeadline are required")
	}

	if req.MemberID == nil || req.AddressID == nil {
		return 0, errors.New("Member ID and Address ID are required")
	}

	tx, txErr := memberService.DB.BeginTx(ctx, nil)
	if txErr != nil {
		return 0, txErr
	}
	defer func() {
		if txErr != nil {
			tx.Rollback()
		}
	}()

	repo := memberService.RepoFactory(tx)

	addressID := *req.AddressID
	memberID := *req.MemberID

	address := &models.Address{
		Street:     req.StreetAddress,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		AptNumber:  req.AptNumber,
	}

	member := &models.Member{
		MyHeadline:            req.MyHeadline,
		AboutMe:               req.AboutMe,
		AdditionalInformation: req.AdditionalInformation,
	}

	memberID, err := repo.UpdateMember(ctx, member, memberID)
	if err != nil {
		return 0, err
	}

	addrErr := repo.UpdateAddress(ctx, address, addressID)
	if addrErr != nil {
		return 0, addrErr
	}

	return memberID, nil
}

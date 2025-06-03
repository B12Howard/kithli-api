package user

import (
	"context"
	"kithli-api/firebase"
	"kithli-api/models"
	"kithli-api/repositories"
	"strings"
)

func ExtractUserIDFromToken(authClient *firebase.FirebaseClient, tokenString string) (string, error) {
	tokenString = strings.TrimSpace(strings.Replace(tokenString, "Bearer ", "", 1))

	token, err := authClient.FirebaseAuth.VerifyIDToken(context.Background(), tokenString)
	if err != nil {
		return "", err
	}

	uid, ok := token.Claims["user_id"].(string)
	if !ok {
		return "", err
	}
	return uid, nil
}


type UserService struct {
	Repo repositories.IUserRepository
}

func NewUserService(repo repositories.IUserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (userService *UserService) GetFullUserData(ctx context.Context, externalID string) (*models.UserDataResponse, error) {
	user, err := userService.Repo.GetUserByExternalID(ctx, externalID)
	if err != nil {
		return nil, err
	}

	userResp := &models.UserDataResponse{
		ID:                 user.ID,
		Email:              user.Email,
		Phone:              user.Phone,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		EmailConfirmed:     user.EmailConfirmed,
		ActiveSubscription: user.ActiveSubscription,
		ExternalID:         user.ExternalID,
		MemberID:           user.MemberID,
		KithID:             user.KithID,
		CreatedAt:          user.CreatedAt,
	}

	if user.MemberID != nil {
		memberDetails, err := userService.Repo.GetMemberDetailsByID(ctx, *user.MemberID)
		if err != nil && err.Error() != "sql: no rows in result set" {
			return nil, err
		}

		if memberDetails != nil {
			userResp.MyHeadline = memberDetails.MyHeadline
			userResp.AboutMe = memberDetails.AboutMe
			userResp.AdditionalInfo = memberDetails.AdditionalInfo
			userResp.PostalCode = memberDetails.PostalCode
			userResp.StreetAddress = memberDetails.StreetAddress
			userResp.City = memberDetails.City
			userResp.State = memberDetails.State
			userResp.AddressID = memberDetails.AddressID
		}
	}

	return userResp, nil
}
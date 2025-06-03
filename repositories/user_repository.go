package repositories

import (
	"context"
	"database/sql"
	"kithli-api/models"
)

type IUserRepository interface {
	GetUserByExternalID(ctx context.Context, externalID string) (*models.User, error)
	GetMemberDetailsByID(ctx context.Context, memberID int) (*models.MemberDetails, error)
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) GetUserByExternalID(ctx context.Context, externalID string) (*models.User, error) {
	var u models.User
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, email, phone, first_name, last_name, email_confirmed, active_subscription,
		       external_id, member, kith, created_at
		FROM users
		WHERE external_id = $1
	`, externalID).Scan(
		&u.ID, &u.Email, &u.Phone, &u.FirstName, &u.LastName,
		&u.EmailConfirmed, &u.ActiveSubscription, &u.ExternalID,
		&u.MemberID, &u.KithID, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetMemberDetailsByID(ctx context.Context, memberID int) (*models.MemberDetails, error) {
	var m models.MemberDetails
	err := r.DB.QueryRowContext(ctx, `
		SELECT m.my_headline, m.about_me, m.additional_information,
		       a.postal_code, a.street, a.city, a.state, a.id
		FROM members m
		LEFT JOIN member_addresses ma ON ma.user_id = m.id AND ma.is_primary = TRUE
		LEFT JOIN addresses a ON a.id = ma.address
		WHERE m.id = $1
	`, memberID).Scan(
		&m.MyHeadline, &m.AboutMe, &m.AdditionalInfo,
		&m.PostalCode, &m.StreetAddress, &m.City, &m.State, &m.AddressID,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

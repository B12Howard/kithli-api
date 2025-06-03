package repositories

import (
	"context"
	"database/sql"
	"kithli-api/models"
)

type IMemberRepository interface {
	InsertAddress(ctx context.Context, addr *models.Address) (int, error)
	InsertMember(ctx context.Context, m *models.Member) (int, error)
	LinkAddress(ctx context.Context, memberID, addressID int) error
	UpdateUserMember(ctx context.Context, uid string, memberID int) error
	UpdateAddress(ctx context.Context, addr *models.Address, addressID int) error
	UpdateMember(ctx context.Context, m *models.Member, memberID int) (int, error)
}

type memberRepository struct {
	Tx *sql.Tx
}

func NewMemberRepository(tx *sql.Tx) IMemberRepository {
	return &memberRepository{Tx: tx}
}

func (r *memberRepository) InsertAddress(ctx context.Context, addr *models.Address) (int, error) {
	var id int
	err := r.Tx.QueryRowContext(ctx, `
		INSERT INTO addresses (postal_code, street, apt_number, city, state)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, addr.PostalCode, addr.Street, addr.AptNumber, addr.City, addr.State).Scan(&id)
	return id, err
}

func (r *memberRepository) InsertMember(ctx context.Context, m *models.Member) (int, error) {
	var id int
	err := r.Tx.QueryRowContext(ctx, `
		INSERT INTO members (my_headline, about_me, additional_information)
		VALUES ($1, $2, $3) RETURNING id
	`, m.MyHeadline, m.AboutMe, m.AdditionalInformation).Scan(&id)
	return id, err
}

func (r *memberRepository) LinkAddress(ctx context.Context, memberID, addressID int) error {
	_, err := r.Tx.ExecContext(ctx, `
		INSERT INTO member_addresses (user_id, address) VALUES ($1, $2)
	`, memberID, addressID)
	return err
}

func (r *memberRepository) UpdateUserMember(ctx context.Context, uid string, memberID int) error {
	_, err := r.Tx.ExecContext(ctx, `
		UPDATE users SET member = $1 WHERE external_id = $2
	`, memberID, uid)
	return err
}

func (r *memberRepository) UpdateMember(ctx context.Context, m *models.Member, memberID int) (int, error) {
	var id int
	err := r.Tx.QueryRowContext(ctx, `
		UPDATE members SET my_headline = $1, about_me = $2, additional_information = $3 WHERE id = $4
		RETURNING id
	`, m.MyHeadline, m.AboutMe, m.AdditionalInformation, memberID).Scan(&id)
	return id, err
}

func (r *memberRepository) UpdateAddress(ctx context.Context, addr *models.Address, addressID int) error {
	_, err := r.Tx.ExecContext(ctx, `
		UPDATE addresses SET postal_code = $1, street = $2, apt_number = $3, city = $4, state = $5 WHERE id = $6
	`, addr.PostalCode, addr.Street, addr.AptNumber, addr.City, addr.State, addressID)
	return err
}
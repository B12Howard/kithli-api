package models

type Member struct {
	MyHeadline            string
	AboutMe               *string
	AdditionalInformation *string
}
type MemberDetails struct {
	MyHeadline     *string
	AboutMe        *string
	AdditionalInfo *string
	PostalCode     *string
	StreetAddress  *string
	City           *string
	State          *string
	AddressID      *int
}

type UpdateMemberResponse struct {
	MemberID int    `json:"memberId"`
	Message  string `json:"message"`
}
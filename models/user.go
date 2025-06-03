package models

type User struct {
	ID                 int
	Email              string
	Phone              *string
	FirstName          *string
	LastName           *string
	EmailConfirmed     bool
	ActiveSubscription bool
	ExternalID         string
	MemberID           *int
	KithID             *int
	CreatedAt          string
}

type UserDataResponse struct {
	ID                 int     `json:"id"`
	Email              string  `json:"email"`
	Phone              *string `json:"phone,omitempty"`
	FirstName          *string `json:"firstName,omitempty"`
	LastName           *string `json:"lastName,omitempty"`
	EmailConfirmed     bool    `json:"emailConfirmed"`
	ActiveSubscription bool    `json:"activeSubscription"`
	ExternalID         string  `json:"externalId"`
	MemberID           *int    `json:"memberId,omitempty"`
	KithID             *int    `json:"kithId,omitempty"`
	CreatedAt          string  `json:"createdAt"`
	// Member details (if exists)
	MyHeadline     *string `json:"myHeadline,omitempty"`
	AboutMe        *string `json:"aboutMe,omitempty"`
	PostalCode     *string `json:"postalCode,omitempty"`
	StreetAddress  *string `json:"streetAddress,omitempty"`
	AptNumber      *string `json:"aptNumber,omitempty"`
	City           *string `json:"city,omitempty"`
	State          *string `json:"state,omitempty"`
	AdditionalInfo *string `json:"additionalInformation,omitempty"`
	AddressID	   *int `json:"addressId,omitempty"`
}

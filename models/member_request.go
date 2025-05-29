package models

// MemberRequest represents the incoming request
type MemberRequest struct {
	UID                   string  `json:"uid"`
	MyHeadline            string  `json:"myHeadline"`
	AboutMe               *string `json:"aboutMe,omitempty"`
	PostalCode            *string `json:"postalCode,omitempty"`
	StreetAddress         *string `json:"streetAddress,omitempty"`
	AptNumber             *string `json:"aptNumber,omitempty"`
	City                  *string `json:"city,omitempty"`
	State                 *string `json:"state,omitempty"`
	AdditionalInformation *string `json:"additionalInformation,omitempty"`
}
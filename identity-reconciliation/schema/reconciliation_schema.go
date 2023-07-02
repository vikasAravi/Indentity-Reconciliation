package schema

type PingResponse struct {
	Success string `json:"success"`
}

type ContactSchema struct {
	PrimaryContactId    *int64    `json:"primary_contact_id"`
	Emails              []*string `json:"emails"`
	PhoneNumbers        []*string `json:"phone_numbers"`
	SecondaryContactIds []*int64  `json:"secondary_contact_ids"`
}

type IdentityRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type IdentityResponse struct {
	Contact ContactSchema `json:"contact"`
}

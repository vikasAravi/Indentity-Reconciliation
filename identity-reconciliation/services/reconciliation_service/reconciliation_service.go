package reconciliation_service

import (
	"bitespeed/identity-reconciliation/models"
	"bitespeed/identity-reconciliation/persistance"
	"bitespeed/identity-reconciliation/schema"
	"context"
)

func GetIdentityResponse(ctx context.Context, req schema.IdentityRequest) (*schema.IdentityResponse, error) {
	identity := persistance.ReconciliationDBClient()
	contactDetails, err := identity.GetContactDetails(ctx, req)
	if err != nil {
		return nil, err
	}
	var primaryContactDetails models.Contact
	secondaryContactDetails := make([]models.Contact, 0)
	for _, contactDetail := range contactDetails {
		if *contactDetail.LinkPrecedence == "primary" {
			primaryContactDetails = *contactDetail
		} else {
			secondaryContactDetails = append(secondaryContactDetails, *contactDetail)
		}
	}

	var contact schema.ContactSchema
	contact.PrimaryContactId = primaryContactDetails.Id
	contact.Emails = getEmails(primaryContactDetails, secondaryContactDetails)
	contact.PhoneNumbers = getPhoneNumbers(primaryContactDetails, secondaryContactDetails)
	contact.SecondaryContactIds = getSecondaryContactIds(secondaryContactDetails)

	return &schema.IdentityResponse{Contact: contact}, nil
}

func getEmails(primaryContactDetail models.Contact, secondaryContactDetails []models.Contact) []*string {
	emails := make([]*string, 0)
	uniqueEmailsMap := make(map[string]bool)
	emails = append(emails, primaryContactDetail.Email)
	uniqueEmailsMap[*primaryContactDetail.Email] = true
	for _, secondaryContactDetail := range secondaryContactDetails {
		if !uniqueEmailsMap[*secondaryContactDetail.Email] {
			uniqueEmailsMap[*secondaryContactDetail.Email] = true
			emails = append(emails, secondaryContactDetail.Email)
		}
	}
	return emails
}

func getPhoneNumbers(primaryContactDetail models.Contact, secondaryContactDetails []models.Contact) []*string {
	phoneNumbers := make([]*string, 0)
	uniquePhoneMap := make(map[string]bool)
	phoneNumbers = append(phoneNumbers, primaryContactDetail.PhoneNumber)
	uniquePhoneMap[*primaryContactDetail.PhoneNumber] = true
	for _, secondaryContactDetail := range secondaryContactDetails {
		if !uniquePhoneMap[*secondaryContactDetail.PhoneNumber] {
			uniquePhoneMap[*secondaryContactDetail.PhoneNumber] = true
			phoneNumbers = append(phoneNumbers, secondaryContactDetail.PhoneNumber)
		}
	}
	return phoneNumbers
}

func getSecondaryContactIds(secondaryContactDetails []models.Contact) []*int64 {
	secondaryContactIds := make([]*int64, 0)
	for _, secondaryContactDetail := range secondaryContactDetails {
		secondaryContactIds = append(secondaryContactIds, secondaryContactDetail.Id)
	}
	return secondaryContactIds
}

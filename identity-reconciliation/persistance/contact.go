package persistance

import (
	DB "bitespeed/identity-reconciliation/common/database/postgres"
	"bitespeed/identity-reconciliation/models"
	"bitespeed/identity-reconciliation/query"
	"bitespeed/identity-reconciliation/schema"
	"context"
	"errors"
	"gorm.io/gorm"
	"strconv"
)

var reconciliationDBClient *ReconciliationDB

type IReconciliation interface {
	GetContactDetails(ctx context.Context, identityRequest schema.IdentityRequest) ([]*models.Contact, error)
	SaveContactDetails(ctx context.Context, contact *models.Contact) (*models.Contact, error)
	UpdatePrimaryContactUpdatedByTS(ctx context.Context, primaryContact *models.Contact, request schema.IdentityRequest) (*models.Contact, error)
	GetSecondaryContactDetailsById(ctx context.Context, primaryContactId *int64) ([]*models.Contact, error)
	GetPrimaryContactDetailsWithEmailAndPhoneNumber(ctx context.Context, email string, phoneNumber string) ([]*models.Contact, error)
	UpdatePrimaryContactDetailsToSecondary(ctx context.Context, primaryContactId *int64, secondaryContactId *int64) (*models.Contact, error)
	UpdateSecondaryContactLinkedIdWhenPrimaryContactMigration(ctx context.Context, primaryContactLinkedId *int64, secondaryContactLinkedId *int64) ([]*models.Contact, error)
}

type ReconciliationDB struct {
	DB *gorm.DB
}

func ReconciliationDBClient() *ReconciliationDB {
	if reconciliationDBClient != nil {
		return reconciliationDBClient
	}
	reconciliationDBClient = &ReconciliationDB{
		DB.GetDB(),
	}
	return reconciliationDBClient
}

// [DEF] - GET THE CONTACT DETAILS FOR THE GIVEN REQUEST

func (r *ReconciliationDB) GetContactDetails(ctx context.Context, identityRequest schema.IdentityRequest) ([]*models.Contact, error) {
	contactDetails := make([]*models.Contact, 0)
	primaryContactDetails, _ := r.GetPrimaryContactDetailsWithEmailAndPhoneNumber(ctx, identityRequest.Email, identityRequest.PhoneNumber)
	rowsFetched := len(primaryContactDetails)
	switch rowsFetched {
	case 0:
		// [NOTES]
		// 1. THIS CASE HAPPENS IF THE CUSTOMER COMES FOR THE FIRST TIME
		// 2. SAVING THE CONTACT TO THE DB CONSIDERING AS PRIMARY PREFERENCE
		// 3. NO SECONDARY CONTACTS
		linkPreference := "primary"
		primaryContactDetails, _ := r.SaveContactDetails(ctx, &models.Contact{Email: &identityRequest.Email, PhoneNumber: &identityRequest.PhoneNumber,
			LinkPrecedence: &linkPreference})
		contactDetails = append(contactDetails, primaryContactDetails)
		return contactDetails, nil
	case 1:
		// [NOTES]
		// 1. THIS CASE HAPPENS IF THE CUSTOMER COMES AGAIN ( EITHER WITH SAME ACCOUNT DETAILS / DIFFERENT ACCOUNT DETAILS )
		// 2. IF SAME ACCOUNT DETAILS => UPDATE THE `UPDATED_BY` FIELD AND GET THE SECONDARY CONTACTS USING THE PRIMARY CONTACT ID
		// 3. IF DIFFERENT ACCOUNT BUT MATCHING PARTIALLY - MAKE THE REQUEST ENTRY AS SECONDARY WITH PRIMARY ID AS LINK ID
		// 4. GET ALL THE SECONDARY CONTACTS AND MERGE
		if isSameAccountDetails(identityRequest, *primaryContactDetails[0]) {
			primaryContactDetails, _ := r.UpdatePrimaryContactUpdatedByTS(ctx, primaryContactDetails[0], identityRequest)
			contactDetails = append(contactDetails, primaryContactDetails)
		} else {
			contactDetails = append(contactDetails, primaryContactDetails[0])
			linkPreference := "secondary"
			r.SaveContactDetails(ctx, &models.Contact{Email: &identityRequest.Email,
				PhoneNumber: &identityRequest.PhoneNumber, LinkPrecedence: &linkPreference, LinkedId: primaryContactDetails[0].Id})
		}
		secondaryContactDetails, _ := r.GetSecondaryContactDetailsById(ctx, primaryContactDetails[0].Id)
		for _, secondaryContactDetail := range secondaryContactDetails {
			contactDetails = append(contactDetails, secondaryContactDetail)
		}
		return contactDetails, nil
	case 2:
		// [NOTES]
		// 1. THIS CASE HAPPENS IF THERE ARE TWO PRIMARY ACCOUNTS AND BOTH ARE RELATED TO SAME USER
		// 2. UPDATE THE LATEST ONE AS SECONDARY AND LINK ID SHOULD BE THE OLDER CONTACT
		// 3. UPDATE ALL THE LINKED IDS OF THIS PRIMARY CONTACT ( GOING TO BE MIGRATED )  TO THE CURRENT PRIMARY CONTACT
		// 3. FETCH ALL THE SECONDARY CONTACTS
		r.UpdatePrimaryContactDetailsToSecondary(ctx, primaryContactDetails[0].Id, primaryContactDetails[1].Id)
		r.UpdateSecondaryContactLinkedIdWhenPrimaryContactMigration(ctx, primaryContactDetails[0].Id, primaryContactDetails[1].Id)
		contactDetails = append(contactDetails, primaryContactDetails[0])
		secondaryContactDetails, _ := r.GetSecondaryContactDetailsById(ctx, primaryContactDetails[0].Id)
		for _, secondaryContactDetail := range secondaryContactDetails {
			contactDetails = append(contactDetails, secondaryContactDetail)
		}
		return contactDetails, nil
	default:
		return nil, errors.New("something wrong with the data")
	}
}

// [DEF] - SAVE THE CONTACT DETAILS

func (r *ReconciliationDB) SaveContactDetails(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	var savedContactDetails models.Contact
	r.DB.Raw(query.SaveContactDetails, contact.PhoneNumber, contact.Email, contact.LinkedId, contact.LinkPrecedence).Scan(&savedContactDetails)
	return &savedContactDetails, nil
}

// [DEF] - UPDATE THE PRIMARY CONTACT
// 1. THIS WILL GET TRIGGERED ONLY IF THE USER TRIED DOING SOMETHING WITH THE IDENTICAL ACCOUNT DETAILS
// EXAMPLE -
// CURRENT STATE -  {"email": "hi@gmail.com", phone_number:null}
// NOW THE NEW REQUEST IS {"email": "hi@gmail.com", "phone_number": "0123"} => THIS IS ALREADY EXISTING ACCOUNT BUT CUSTOMER IS PROVIDING ADDITIONAL INFO
// IN THE ABOVE CASE THIS FUNCTION WILL TRIGGER

func (r *ReconciliationDB) UpdatePrimaryContactUpdatedByTS(ctx context.Context, primaryContact *models.Contact, request schema.IdentityRequest) (*models.Contact, error) {
	var updatedPrimaryContact models.Contact
	if *primaryContact.Email == "" && request.Email != "" {
		r.DB.Raw(query.UpdatePrimaryEmailDetailsWithUpdatedTs, request.Email, primaryContact.Id).Scan(&updatedPrimaryContact)
	} else if *primaryContact.PhoneNumber == "" && request.PhoneNumber != "" {
		r.DB.Raw(query.UpdatePrimaryPhoneNumberDetailsWithUpdatedTs, request.PhoneNumber, primaryContact.Id).Scan(&updatedPrimaryContact)
	} else {
		r.DB.Raw(query.UpdatePrimaryUpdatedTS, primaryContact.Id).Scan(&updatedPrimaryContact)
	}
	return &updatedPrimaryContact, nil
}

// [DEF] - GET ALL THE LINKED CONTACTS FOR THE GIVEN PRIMARY CONTACT

func (r *ReconciliationDB) GetSecondaryContactDetailsById(ctx context.Context, primaryContactId *int64) ([]*models.Contact, error) {
	secondaryContactDetails := make([]*models.Contact, 0)
	r.DB.Raw(prepareGetContactDetailsQueryWithId(primaryContactId)).Scan(&secondaryContactDetails)
	return secondaryContactDetails, nil
}

// [DEF] - GET THE PRIMARY CONTACT DETAILS WITH THE GIVEN REQUEST

func (r *ReconciliationDB) GetPrimaryContactDetailsWithEmailAndPhoneNumber(ctx context.Context, email string, phoneNumber string) ([]*models.Contact, error) {
	primaryContactDetails := make([]*models.Contact, 0)
	r.DB.Raw(prepareGetContactDetailsQueryWithPhoneNumberAndEmail(phoneNumber, email)).Scan(&primaryContactDetails)
	return primaryContactDetails, nil
}

// [DEF] - UPDATE PRIMARY CONTACT TO SECONDARY CONTACT

func (r *ReconciliationDB) UpdatePrimaryContactDetailsToSecondary(ctx context.Context, primaryContactId *int64, secondaryContactId *int64) (*models.Contact, error) {
	var updatedContactDetail models.Contact
	r.DB.Raw(query.UpdatePrimaryContactDetails, primaryContactId, "secondary", secondaryContactId).Scan(&updatedContactDetail)
	return &updatedContactDetail, nil
}

func (r *ReconciliationDB) UpdateSecondaryContactLinkedIdWhenPrimaryContactMigration(ctx context.Context, primaryContactLinkedId *int64, secondaryContactLinkedId *int64) ([]*models.Contact, error) {
	var updatedContactDetail []*models.Contact
	r.DB.Raw(query.UpdateLinkedIdWhenPrimaryContactIsMigratedToSecondary, primaryContactLinkedId, secondaryContactLinkedId).Scan(updatedContactDetail)
	return updatedContactDetail, nil
}

func prepareGetContactDetailsQueryWithId(id *int64) string {
	return query.GetPrimaryContactDetails + "linked_id = '" + strconv.FormatInt(*id, 10) + "'"
}

func prepareGetContactDetailsQueryWithPhoneNumberAndEmail(phoneNumber string, email string) string {
	q := query.GetPrimaryContactDetails
	filterQuery := " "
	if phoneNumber != "" {
		filterQuery += "  phone_number = '" + phoneNumber + "' "
	}
	if email != "" {
		if filterQuery != " " {
			filterQuery += " OR "
		}
		filterQuery += "  email = '" + email + "' "
	}
	filterQuery = "(" + filterQuery + ")"
	//filterQuery += " AND link_precedence = '" + "primary" + "' " // THIS SHOULD NOT BE THE CASE
	return q + filterQuery + " ORDER BY created_at"
}

// [DEF]
// CONSIDER THE BELOW EXAMPLE
// IF THE CURRENT STATE OF THE DB IS {"email": "hi@gmail.com", phone_number:null}
// NOW THE NEW REQUEST IS {"email": "hi@gmail.com", "phone_number": "0123"} => THIS IS ALREADY EXISTING ACCOUNT SO TRUE
// NEED TO UPDATE phone_number FOR THE PRIMARY CONTACT
func isSameAccountDetails(request schema.IdentityRequest, primaryContactDetails models.Contact) bool {
	var resp bool
	if request.PhoneNumber != "" && request.Email != "" {
		resp = *primaryContactDetails.PhoneNumber == "" || *primaryContactDetails.Email == ""
	}
	return resp || (request.PhoneNumber == "" || *primaryContactDetails.PhoneNumber == request.PhoneNumber) &&
		(request.Email == "" || *primaryContactDetails.Email == request.Email)
}

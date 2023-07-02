package query

const (
	GetPrimaryContactDetails                              = "SELECT id, phone_number, email, linked_id, link_precedence, created_at, updated_at, deleted_at FROM CONTACT WHERE "
	SaveContactDetails                                    = "INSERT INTO CONTACT(phone_number, email, linked_id, link_precedence) values(?,?,?,?) returning id, phone_number, email, linked_id, link_precedence;"
	UpdatePrimaryContactDetails                           = "UPDATE CONTACT SET linked_id = ?, link_precedence = ? WHERE id = ? returning id, phone_number, email, linked_id, link_precedence;"
	UpdatePrimaryEmailDetailsWithUpdatedTs                = "UPDATE CONTACT SET updated_at = CURRENT_TIMESTAMP, email = ? WHERE id = ? returning id, phone_number, email, linked_id, link_precedence;"
	UpdatePrimaryPhoneNumberDetailsWithUpdatedTs          = "UPDATE CONTACT SET updated_at = CURRENT_TIMESTAMP, phone_number = ? WHERE id = ? returning id, phone_number, email, linked_id, link_precedence;"
	UpdatePrimaryUpdatedTS                                = "UPDATE CONTACT SET updated_at = CURRENT_TIMESTAMP WHERE id = ? returning id, phone_number, email, linked_id, link_precedence;"
	UpdateLinkedIdWhenPrimaryContactIsMigratedToSecondary = "UPDATE CONTACT SET updated_at = CURRENT_TIMESTAMP, linked_id = ? WHERE linked_id = ?"
)

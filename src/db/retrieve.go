package db

import (
	"backend/src/common"
	"gorm.io/gorm"
)

func RetrieveUser(db *gorm.DB, email string) (*common.User, error) {
	var user common.User
	result := db.Where("email = ?", email).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// RetrieveKeystoneById retrieves a single Keystone based on its ID.
func RetrieveKeystoneById(db *gorm.DB, keystoneID uint) (*common.Keystone, error) {
	var keystone common.Keystone
	result := db.Where("id = ?", keystoneID).Preload("User").Preload("Reflections").Find(&keystone)
	if result.Error != nil {
		return nil, result.Error
	}
	return &keystone, nil
}

// RetrieveKeystoneByIdFull retrieves the full scope of a single Keystone based on its ID.
func RetrieveKeystoneByIdFull(db *gorm.DB, keystoneID uint) (*common.Keystone, error) {
	var keystone common.Keystone
	result := db.Where("id = ?", keystoneID).Preload("User").Preload("Reflections").Preload("MailingDetails").Find(&keystone)
	if result.Error != nil {
		return nil, result.Error
	}
	return &keystone, nil
}

// RetrieveAllKeystones retrieves all Keystones from the database.
func RetrieveAllKeystones(db *gorm.DB) ([]*common.Keystone, error) {
	var keystones []*common.Keystone
	result := db.Preload("Reflections").Preload("User").Find(&keystones)
	if result.Error != nil {
		return nil, result.Error
	}
	return keystones, nil
}

// RetrieveReflectionsByKeystoneID retrieves all reflections associated with a given keystoneID.
func RetrieveReflectionsByKeystoneID(db *gorm.DB, keystoneID uint) ([]*common.Reflection, error) {
	var reflections []*common.Reflection

	// Perform the query: find all reflections where the KeystoneID matches the provided keystoneID.
	result := db.Where("keystone_id = ?", keystoneID).Preload("User").Find(&reflections)

	if result.Error != nil {
		// Handle any errors (e.g., no records found, DB connection issues)
		return nil, result.Error
	}

	return reflections, nil
}

// RetrieveAllUsers retrieves all Users from the database.
func RetrieveAllUsers(db *gorm.DB) ([]*common.User, error) {
	var users []*common.User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// RetrieveLastMailingDetailsByKeystoneID retrieves the last mailing details for a given keystoneID.
func RetrieveLastMailingDetailsByKeystoneID(db *gorm.DB, keystoneID uint) (*common.MailingDetails, error) {
	var reflections *common.MailingDetails

	// Perform the query: find all reflections where the KeystoneID matches the provided keystoneID.
	result := db.Where("keystone_id = ?", keystoneID).Last(&reflections)

	if result.Error != nil {
		// Handle any errors (e.g., no records found, DB connection issues)
		return nil, result.Error
	}

	return reflections, nil
}

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

// RetrieveKeystone retrieves a single Keystone based on its ID.
func RetrieveKeystone(db *gorm.DB, keystoneID uint) (*common.Keystone, error) {
	var keystone common.Keystone
	result := db.Where("id = ?", keystoneID).Preload("User").Preload("Reflections").Find(&keystone)
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

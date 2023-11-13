package db

import (
	"backend/src/common"
	"gorm.io/gorm"
)

func InsertUser(db *gorm.DB, user *common.User) error {
	result := db.Create(user)

	if result.Error != nil {
		return result.Error
	}

	common.DebugLogger.Println("persisted user to db")
	return nil
}

func InsertKeystone(db *gorm.DB, user *common.Keystone) error {
	result := db.Create(user)

	if result.Error != nil {
		return result.Error
	}

	common.DebugLogger.Println("persisted keystone to db")
	return nil
}

func InsertReflection(db *gorm.DB, user *common.Reflection) error {
	result := db.Create(user)

	if result.Error != nil {
		return result.Error
	}

	common.DebugLogger.Println("persisted reflection to db")
	return nil
}

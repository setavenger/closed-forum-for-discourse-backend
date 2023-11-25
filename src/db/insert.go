package db

import (
	"backend/src/common"
	"gorm.io/gorm"
)

//func InsertUser(db *gorm.DB, user *common.User) error {
//	result := db.Create(user)
//
//	if result.Error != nil {
//		return result.Error
//	}
//
//	common.DebugLogger.Println("persisted user to db")
//	return nil
//}

func InsertKeystone(db *gorm.DB, keystone *common.Keystone) error {
	result := db.Create(keystone)

	if result.Error != nil {
		return result.Error
	}

	common.DebugLogger.Println("persisted keystone to db")
	return nil
}

func InsertReflection(db *gorm.DB, reflection *common.Reflection) error {
	result := db.Create(reflection)

	if result.Error != nil {
		return result.Error
	}

	common.DebugLogger.Println("persisted reflection to db")
	return nil
}

func InsertMailDetails(db *gorm.DB, info *common.MailingDetails) error {
	result := db.Create(info)

	if result.Error != nil {
		return result.Error
	}

	common.DebugLogger.Println("persisted reflection to db")
	return nil
}

package common

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uint           `gorm:"column:id"` // Unique identifier for the user
	EMail        string         `gorm:"column:email"`
	PasswordHash string         `gorm:"column:password_hash"`
	Nickname     string         `gorm:"column:nickname"`
	Joined       time.Time      `gorm:"column:joined"`
	Keystones    []Keystone     // One-to-many relationship with Keystones
	Reflections  []Reflection   // One-to-many relationship with Reflections
	CreatedAt    time.Time      `json:"created_at,omitempty"`
	UpdatedAt    time.Time      `json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type Keystone struct {
	ID          uint           `gorm:"column:id" json:"id"`
	Timestamp   time.Time      `gorm:"column:timestamp" json:"timestamp"`
	UserID      uint           `gorm:"column:user_id" json:"user_id"` // Foreign key for User
	Title       string         `gorm:"column:title" json:"title"`
	Content     string         `gorm:"column:content" json:"content"`
	Reflections []Reflection   `json:"reflections"`
	User        User           // This allows you to access user information from the Keystone
	CreatedAt   time.Time      `json:"created_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type KeystoneTransfer struct {
	ID              uint      `json:"id"`
	Timestamp       time.Time `json:"timestamp"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	ReflectionCount uint      `json:"reflection_count"`
	Author          string    `json:"author"`
}

type Reflection struct {
	ID         uint           `gorm:"column:id" json:"id"`
	Timestamp  time.Time      `gorm:"column:timestamp" json:"timestamp"`
	KeystoneID *uint          `gorm:"column:keystone_id" json:"keystone_id"` // Pointer to allow nil
	ParentID   *uint          `gorm:"column:parent_id" json:"parent_id"`     // Pointer to allow nil
	UserID     uint           `gorm:"column:user_id" json:"user_id"`         // Foreign key for User
	Content    string         `gorm:"column:content" json:"content"`
	User       User           // This allows you to access user information from the Reflection
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type ReflectionTransfer struct {
	ID         uint      `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	KeystoneID *uint     `json:"keystone_id"` // Pointer to allow nil
	ParentID   *uint     `json:"parent_id"`   // Pointer to allow nil
	Content    string    `json:"content"`
	Author     string    `json:"author"`
}

func (k *Keystone) ToTransfer() *KeystoneTransfer {
	return &KeystoneTransfer{
		ID:              k.ID,
		Timestamp:       k.Timestamp,
		Title:           k.Title,
		Content:         k.Content,
		ReflectionCount: uint(len(k.Reflections)),
		Author:          k.User.Nickname,
	}
}

func (r *Reflection) ToTransfer() *ReflectionTransfer {
	return &ReflectionTransfer{
		ID:         r.ID,
		Timestamp:  r.Timestamp,
		Content:    r.Content,
		Author:     r.User.Nickname,
		KeystoneID: r.KeystoneID,
		ParentID:   r.ParentID,
	}
}

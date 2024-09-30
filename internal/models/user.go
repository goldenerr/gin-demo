package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

// TableName 指定模型的表名
func (User) TableName() string {
	return "users"
}

func (u *User) Create(db *gorm.DB) error {
	return db.Create(u).Error
}

func (u *User) Update(db *gorm.DB) error {
	return db.Save(u).Error
}

func (u *User) Delete(db *gorm.DB) error {
	return db.Delete(u).Error
}

func GetUserByID(db *gorm.DB, id string) (*User, error) {
	var user User
	err := db.First(&user, id).Error
	return &user, err
}

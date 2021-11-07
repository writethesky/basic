package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50)" json:"username"`
	Password string `gorm:"type:char(32)" json:"password"`
	Salt     string `gorm:"type:char(6)" json:"salt"`
}

package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`

	Followings []*User    `gorm:"many2many:user_follows;joinForeignKey:follower_id;joinReferences:following_id"`
	Favorites  []*Article `gorm:"many2many:user_favorites;"`
}

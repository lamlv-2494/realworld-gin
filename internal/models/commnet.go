package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Body string `json:"body" gorm:"type:text;not null"`

	ArticleID uint
	Article   Article `gorm:"foreignKey:ArticleID"`

	AuthorID uint
	Author   User `gorm:"foreignKey:AuthorID"`
}

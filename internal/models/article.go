package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Slug        string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Title       string `gorm:"not null"`
	Description string
	Body        string
	AuthorID    uint
	Author      User   `gorm:"foreignKey:AuthorID"`     // Quan hệ Article thuộc về User
	Tags        []*Tag `gorm:"many2many:article_tags;"` // Quan hệ Nhiều-Nhiều
}

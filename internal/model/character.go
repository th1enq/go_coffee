package model

import (
	"time"

	"gorm.io/gorm"
)

type Character struct {
	gorm.Model
	Name          string
	Rarity        int
	Region        string
	Vision        string
	WeaponType    string
	Constellation string
	Birthday      time.Time
	Affilliation  string
	ReleaseDate   time.Time
}


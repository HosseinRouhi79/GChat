package models

type Country struct {
	BaseModel
	Name   string  `gorm:"not null"`
	Cities *[]City `gorm:"not null"`
}

type City struct {
	BaseModel
	Name      string `gorm:"not null;"`
	CountryID int
	Country   Country `json:"-" gorm:"foreignKey:CountryID;"`
}

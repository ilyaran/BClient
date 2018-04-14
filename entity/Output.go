package entity

import "github.com/jinzhu/gorm"

type Output struct {
	gorm.Model

	PrevOut
	AddrTagLink string  `json:"addrTagLink"`
	AddrTag string		`json:"addrTag"`
	Spent bool			`json:"spent"`
}

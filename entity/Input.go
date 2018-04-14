package entity

import "github.com/jinzhu/gorm"

type InputBase struct {
	Sequence int64		`json:"sequence"`
	PrevOut PrevOut	`json:"prev_out"`
	Script string		`json:"script"`
}

type Input struct {
	gorm.Model

	InputBase
	Witness string		`json:"witness"`
}

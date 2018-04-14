package entity


type AddressTransaction struct {
	Address string  			`json:"address_id" gorm:"column:address_id;primary_key" sql:"type:VARCHAR(512) REFERENCES addresses(address_id) ON DELETE cascade;DEFAULT ''"`
	Transaction string 		`json:"transaction_id" gorm:"column:transaction_id;primary_key" sql:"type:VARCHAR(512) REFERENCES transactions(hash_id) ON DELETE cascade;DEFAULT ''"`
}

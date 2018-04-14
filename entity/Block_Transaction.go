package entity

type BlockTransaction struct {

	BlockId string  			`json:"block_id" gorm:"column:block_id;primary_key" sql:"type:VARCHAR(512) REFERENCES blocks(block_id) ON DELETE cascade;DEFAULT ''"`
	TransactionId string 		`json:"transaction_id" gorm:"column:block_transaction_id;primary_key" sql:"type:VARCHAR(512) REFERENCES transactions(hash_id) ON DELETE cascade;DEFAULT ''"`

}


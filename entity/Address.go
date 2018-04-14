package entity


type Address struct {

	Address string			`json:"address" gorm:"column:address_id;primary_key;unique;type:varchar(512)"`
	Hash160 string			`json:"hash160"`
	NTx int64				`json:"n_tx"`
	TotalReceived int64		`json:"total_received"`
	TotalSent int64			`json:"total_sent"`
	FinalBalance int64		`json:"final_balance"`
	Txs []Transaction		`json:"txs"`

}
type TransformedAddress struct {
	Transactions []*TransformedTransaction	`json:"transactions"`
}
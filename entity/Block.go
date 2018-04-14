package entity


type Block struct {

	Hash string			`json:"hash" gorm:"column:block_id;primary_key;unique;type:varchar(512)"`
	Ver int64			`json:"ver"`
	PrevBlock string	`json:"prev_block"`
	MrklRoot string		`json:"mrkl_root"`
	Time int64				`json:"time" gorm:"column:timestamp"`
	Bits int64			`json:"bits"`
	Fee int64			`json:"fee"`
	Nonce int64			`json:"nonce"`
	NTx int64			`json:"n_tx"`
	Size int64			`json:"size"`
	BlockIndex int64	`json:"block_index"`
	MainChain bool		`json:"main_chain"`
	Height int64		`json:"height"`
	Row string 			`json:"-"`
	Tx []Transaction	`json:"tx"`

}

type TransformedBlock struct {

	Raw string			`json:"raw"`
	Height int64		`json:"height"`
	Time int64			`json:"time"`
}

func somefunc(){



}





















package entity



type PrevOut struct {

	Spent bool			`json:"spent"`
	TxIndex int64		`json:"tx_index"`
	Type int64			`json:"type"`
	Addr string			`json:"addr"`
	Value int64			`json:"value"`
	N int64				`json:"n"`
	Script string		`json:"script"`
}


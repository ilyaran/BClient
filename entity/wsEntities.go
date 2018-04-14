package entity


type WsBase struct {
	Op string		`json:"op"`
}
type WsResp struct {
	WsBase
	X TransactionBase 	`json:"x"`
}

type WsReq struct {
	WsBase
	Addr string		`json:"addr"`
}


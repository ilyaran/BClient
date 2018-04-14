package entity



type TransactionBase struct {
	Hash string				`json:"hash" gorm:"column:hash_id;primary_key;unique;type:varchar(512)"`

	LockTime int64			`json:"lock_time"`
	Ver int					`json:"ver"`
	Size int64				`json:"size"`
	Inputs []Input			`json:"inputs"`

	Time int64				`json:"time" gorm:"column:timestamp"`
	TxIndex int64			`json:"tx_index"`
	VinSz int64				`json:"vin_sz"`

	VoutSz int64			`json:"vout_sz"`
	RelayedBy string		`json:"relayed_by"`
	Outputs []PrevOut		`json:"out"`
}


type Transaction struct {

	TransactionBase

	Weight int64			`json:"weight"`
	BlockHeight int64		`json:"block_height"`
	Result int64			`json:"result"`
}



type TransformedTransaction struct {

	Raw string				`json:"raw"`
	Block *TransformedBlock	`json:"block"`
}

/*


var baseUrl="https://blockchain.info/ru/"
type Blocks struct {
	Blocks []Block `json:"blocks"`
}
func TTestRequest2(db *gorm.DB){
	address:=requestAddr("1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F")
	tcount:=len(address.Txs)
	//db := Database()
	signalChan :=make(chan *TransformedTransaction, tcount)

	for _, v := range address.Txs {
		go func() {

			block := requestBlock(int64(v.BlockHeight))
			db.Save(&v)
			if block!=nil{
				ttt:=TransformedTransaction{v.Hash,TransformedBlock{block.Hash, block.Height, block.Time}}
				signalChan <- &ttt
			}else {
				signalChan <- nil
			}

		}()
	}
	addrTransactions:=TransformedAddress{make([]*TransformedTransaction,0,tcount)}
	for i:=0; i< tcount; i++ {
		addrTransactions.Transactions=append(addrTransactions.Transactions,<-signalChan)
		fmt.Printf("%v === %v \n", i,addrTransactions.Transactions[i])


	}
}

func testRequest(){
	address:=requestAddr("1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F")
	tcount:=len(address.Txs)

	signalStringChan :=make(chan *string, tcount)
	for _, v := range address.Txs {
		go func() {

			block := requestBlock(int64(v.BlockHeight))

			if block!=nil{
				rrr:=fmt.Sprintf(`,{"raw": "%s", "block": {"hash": "%s", "height": %d, "time": %d }`, v.Hash, block.Hash, block.Height, block.Time)
				signalStringChan <- &rrr
			}else {
				rr:=fmt.Sprintf(`,{"raw": "%s", "block": null`, v.Hash)
				signalStringChan <- &rr
			}
		}()
	}
	var out string
	var res *string
	for i:=0; i< tcount; i++ {
		res=<-signalStringChan
		out += *res
		fmt.Printf("%v === %v \n", i,*res)
	}
	if len(out)>0 {
		fmt.Println(`{"transactions": [`+out[1:]+`]}`)
	}
}




func request(url string)(*[]byte,error){
	spaceClient := http.Client{
		Timeout: time.Second * 10, // Maximum of 2 secs
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		//panic(err)
		return nil,err
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		//panic(getErr)
		return nil,err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//panic(err)
		return nil,err
	}
	//fmt.Printf("%T \n",body)

	return &body,nil
}

func requestAddr(addr string)*Address{

	url := baseUrl+"rawaddr/"+addr

	res,err:=request(url)
	if err!=nil {return nil}

	address := Address{}
	jsonErr := json.Unmarshal(*res, &address)
	if jsonErr != nil {
		//panic(jsonErr)
		return nil
	}


	return &address

}

func requestBlock(height int64)*Block{

	url := baseUrl+"block-height/"+strconv.FormatInt(height,10)+"?format=json"

	res,err:=request(url)
	if(err!=nil) {return nil}

	blocks := Blocks{}
	jsonErr := json.Unmarshal(*res, &blocks)
	if jsonErr != nil {
		//panic(jsonErr)
		return nil
	}

	*/
/*for _, v := range blocks.Blocks {
		fmt.Println( v.Hash)
	}*//*

	return &blocks.Blocks[0]

}









*/

package main

import (
	"github.com/ilyaran/Blockchain/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/net/websocket"
	"os"

	"database/sql"
)
var Config = struct {
	BaseUrl string
	WsUrl string
	DbUrl string
	DB_user string
	DB_password string
	DB_host string
	DB_name string

}{
	"https://blockchain.info/ru/",
	"ws.blockchain.info/inv",
	"",
	"postgres",
	"postgres",
	"localhost:5432",
	"blockchain",

}
func Database() *gorm.DB {
	//open a db connection
	db, err := gorm.Open("postgres", "postgres://" + Config.DB_user+ ":" + Config.DB_password+ "@" + Config.DB_host+ "/" + Config.DB_name+ "?sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
type Blocks struct {
	Blocks []entity.Block `json:"blocks"`
}
type HttpResponse struct {
	url      string
	response *http.Response
	err      error
	transactionHash string
}

func main() {

	//Migrate the schema
	db := Database()
	db.AutoMigrate(&entity.Address{})
	db.AutoMigrate(&entity.Block{})
	db.AutoMigrate(&entity.Transaction{})
	db.AutoMigrate(&entity.AddressTransaction{})
	db.AutoMigrate(&entity.BlockTransaction{})

	/*err :=db.Save(&entity.AddressTransaction{"1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F","ce1048bff11524cbce99c06d36cb9a50dc9dc1a52f2b57479f946f9bc3c6d25f"}).Error
	if err != nil {
		fmt.Println("Error insert", err)
	}*/

	router := gin.Default()
	router.GET("/", startPage)
	router.LoadHTMLGlob("views/*")
	v1 := router.Group("/api/v1/")
	{
		v1.GET("address", getAddress)
	}
	router.Run()


}
func startPage(c *gin.Context)  {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Welcome",
	})
}
func getAddress(c *gin.Context)  {
	key := c.DefaultQuery("key", "")
	match, _ := regexp.MatchString(`^[0-9a-zA-Z\=]{0,512}$`, key)
	if !match {
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable, "result": "invalid key"})
		return
	}
	address:=requestAddr(key)
	if address==nil{
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "result": "address not found"})
		return
	}

	db := Database()
	var _transactionsAtAddress []*entity.TransformedTransaction
	var addressFromDb = entity.Address{Address:address.Address}

	if db.First(&addressFromDb).RecordNotFound() {
		db.Save(address)
		_transactionsAtAddress = startGetTransactions(address)

		//start for listen changes on address
		go startWebsocketClient(address.Address)

		//response out
		c.JSON(http.StatusOK, gin.H{"transactions": _transactionsAtAddress})
		return
	}else {
		rows, err := db.Raw(`
			SELECT address_transactions.transaction_id,blocks.row,blocks.height,blocks.timestamp
			FROM address_transactions 
			inner join block_transactions on address_transactions.transaction_id = block_transactions.block_transaction_id
			inner join blocks on blocks.block_id = block_transactions.block_id
			WHERE address_transactions.address_id = ? `, address.Address).Rows()

		defer rows.Close()
		if err!=nil{
			panic(err)
		}
		//transactionsFromDb:=[]entity.Transaction{}
		for rows.Next() {
			t:=entity.TransformedTransaction{Block:&entity.TransformedBlock{}}

			var tBlockRow sql.NullString
			var tBlockHeight,tBlockTime sql.NullInt64


			err = rows.Scan(&t.Raw,&tBlockRow,&tBlockHeight,&tBlockTime)
			if err!=nil{
				panic(err)
			}
			if tBlockRow.Valid || tBlockHeight.Valid || tBlockTime.Valid {
				t.Block.Raw=tBlockRow.String
				t.Block.Height=tBlockHeight.Int64
				t.Block.Time=tBlockTime.Int64
			}else {
				t.Block=nil
			}
			_transactionsAtAddress=append(_transactionsAtAddress,&t)
		}
		c.JSON(http.StatusOK, gin.H{"transactions": _transactionsAtAddress})
	}

}

func startWebsocketClient(addr string) {
	fmt.Println("Starting Client")
	ws, err := websocket.Dial(fmt.Sprintf("wss://%s", Config.WsUrl), "", fmt.Sprintf("http://%s/", Config.WsUrl))
	if err != nil {
		fmt.Printf("Dial failed: %s\n", err.Error())
		os.Exit(1)
	}

	incomingMessages := make(chan string)

	go readClientMessages(ws, incomingMessages)

	response := new(entity.WsReq)
	response.Op = "addr_sub"
	response.Addr = addr

	err = websocket.JSON.Send(ws, response)
	if err != nil {
		fmt.Printf("Send failed: %s\n", err.Error())
		os.Exit(1)
	}

	for {
		message := <-incomingMessages
		fmt.Println(`Message Received:`,message)
		newTransaction := entity.Transaction{}
		jsonErr := json.Unmarshal([]byte(message), &newTransaction)
		if jsonErr == nil {
			//save new transaction in to db
			db := Database()
			db.Save(&newTransaction)
			newAddressTransaction:=entity.AddressTransaction{addr,newTransaction.Hash}
			db.Save(&newAddressTransaction)
		}
	}
}

func readClientMessages(ws *websocket.Conn, incomingMessages chan string) {
	for {
		var message string
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			fmt.Printf("Error::: %s\n", err.Error())
			return
		}
		incomingMessages <- message
	}
}

func startGetTransactions(address *entity.Address) []*entity.TransformedTransaction{
	db := Database()
	_transactions :=make([]*entity.TransformedTransaction,0,len(address.Txs))
	httpResponses :=asyncHttpGets(address)
	for i, v := range httpResponses {

		var _transaction *entity.TransformedTransaction
		if v!=nil {
			res, err := ioutil.ReadAll(v.response.Body)
			if err == nil {
				blocks := Blocks{}
				jsonErr := json.Unmarshal(res, &blocks)
				if jsonErr == nil && len(blocks.Blocks) > 0 {
					var block *entity.Block
					OuterLoop:
					for _,vv:=range blocks.Blocks{
					InnerLoop:
						for _,vvv:=range vv.Tx {
							if vvv.Hash==v.transactionHash{
								block=&vv
								break OuterLoop
								break InnerLoop
							}
						}
					}

					if block != nil {
						//get json from block object
						blockJson, errEncode := json.Marshal(block)
						if errEncode==nil{
							//set hash from json
							hasher := md5.New()
							hasher.Write([]byte(blockJson))

							//convert hash to hex
							hh:=hex.EncodeToString(hasher.Sum(nil))

							block.Row=hh

							db.Save(block)
							db.Save(&entity.BlockTransaction{block.Hash,v.transactionHash})

							_transaction = &entity.TransformedTransaction{Raw: v.transactionHash, Block: &entity.TransformedBlock{Raw: block.Row, Height: block.Height, Time: block.Time}}
						}
					}
				}else {
					//fmt.Println(jsonErr)
				}
			}else {
				//fmt.Println(err)
			}
			fmt.Printf("%v === %v \n", i, _transaction)
			_transactions =append(_transactions, _transaction)
		}


	}
	return _transactions
}

func asyncHttpGets(address *entity.Address) []*HttpResponse {
	transactions := address.Txs
	db := Database()
	tcount:=len(transactions)
	ch := make(chan *HttpResponse, tcount) // buffered
	responses := make([]*HttpResponse,0,tcount)
	for i:=0; i<tcount;i++ {
		go func(i int) {
			//fmt.Printf("Fetching %s \n", url)
			transaction :=  transactions[i]
			db.Save(&transaction)
			db.Save(&entity.AddressTransaction{address.Address,transaction.Hash})

			url:=Config.BaseUrl+"block-height/"+strconv.FormatInt(transaction.BlockHeight,10)+"?format=json"
			resp, err := request(url)
			if err != nil {
				//fmt.Println(err)
				//panic(getErr)
				ch <- nil
				return
			}
			//resp, err := http.Get(url)
			ch <- &HttpResponse{url, resp, err,transaction.Hash}
		}(i)
	}

	for {
		select {
		case r := <-ch:
			//fmt.Printf("%s was fetched\n", r.url)
			responses = append(responses, r)
			if len(responses) == tcount {
				return responses
			}
		default:
			fmt.Print("*")
			time.Sleep(5e7)
		}
	}
	return responses

}
func request(url string)(*http.Response,error){
	spaceClient := http.Client{
		Timeout: time.Second * 10, // Maximum of 2 secs
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		//panic(err)
		return nil,err
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	resp, err := spaceClient.Do(req)
	if err != nil {
		//panic(err)
		return nil,err
	}
	return resp,nil
}
func requestAddr(addr string)*entity.Address{

	url := Config.BaseUrl+"rawaddr/"+addr

	resp,err:=request(url)
	if err!=nil {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//panic(err)
		return nil
	}
	//fmt.Println(string(body))
	address := entity.Address{}
	jsonErr := json.Unmarshal(body, &address)
	if jsonErr != nil {
		//panic(jsonErr)
		return nil
	}
	//fmt.Println(string(*res))

	return &address

}


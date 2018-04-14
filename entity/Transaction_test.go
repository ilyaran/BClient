package entity


import (
	"testing"
	"fmt"
	"time"
	"io/ioutil"
	"strconv"
	"encoding/json"
	"net/http"
	"github.com/jinzhu/gorm"


)
func TestTestRequest(t *testing.T) {
	TTestRequest1()

}

var Config = struct {
	BaseUrl string
	DbUrl string
	DB_user string
	DB_password string
	DB_host string
	DB_name string
}{
	"https://blockchain.info/ru/",
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

var baseUrl="https://blockchain.info/ru/"
type Blocks struct {
	Blocks []Block `json:"blocks"`
}

/*

func TTestRequest1(){
	address:=requestAddr1("1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F")

	tcount:=len(address.Txs)

	transactions:=make([]*TransformedTransaction,0,tcount)

	//var wg sync.WaitGroup
	//wg.Add(tcount)

	for i, v := range address.Txs {
		//fmt.Println(v)
		//go func() {
			//fmt.Println(i)
			//block := requestBlock1(int64(v.BlockHeight))
			var block *Block
			url := baseUrl+"block-height/"+strconv.FormatInt(v.BlockHeight,10)+"?format=json"

			res,err:=request1(url)
			if(err!=nil) { }

			blocks := Blocks{}
			jsonErr := json.Unmarshal(*res, &blocks)
			if jsonErr != nil {
				//panic(jsonErr)

			}

			*/
/*for _, v := range blocks.Blocks {
				fmt.Println( v.Hash)
			}*//*

			block= &blocks.Blocks[0]

			var ttt *TransformedTransaction
			if block!=nil{
				ttt=&TransformedTransaction{v.Hash,TransformedBlock{block.Hash, block.Height, block.Time}}
			}
			fmt.Printf("%v === %v \n", i,ttt)
			transactions=append(transactions,ttt)

			//wg.Done()
		//}()
	}
	//wg.Wait()



}

*/

func requestBlock1(height int64)*Block{

	url := baseUrl+"block-height/"+strconv.FormatInt(height,10)+"?format=json"

	res,err:=request1(url)
	if(err!=nil) {return nil}

	blocks := Blocks{}
	jsonErr := json.Unmarshal(*res, &blocks)
	if jsonErr != nil {
		//panic(jsonErr)
		return nil
	}

	/*for _, v := range blocks.Blocks {
		fmt.Println( v.Hash)
	}*/
	return &blocks.Blocks[0]

}
func request1(url string)(*[]byte,error){
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


func requestAddr1(addr string)*Address{

	url := baseUrl+"rawaddr/"+addr

	res,err:=request1(url)
	if err!=nil {return nil}

	address := Address{}
	jsonErr := json.Unmarshal(*res, &address)
	if jsonErr != nil {
		//panic(jsonErr)
		return nil
	}
	//fmt.Println(string(*res))

	return &address

}
type HttpResponse struct {
	url      string
	response *http.Response
	err      error
	transactionHash string
}

func TTestRequest1(){
	address:=requestAddr1("1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F")

	tcount:=len(address.Txs)

	transactions:=make([]*TransformedTransaction,0,tcount)

	mr:=asyncHttpGets(&address.Txs,tcount)
	for i, v := range mr {

		res, err := ioutil.ReadAll(v.response.Body)
		if err != nil {
			fmt.Printf("%v === %v \n", i,nil)
			transactions=append(transactions,nil)
			return
		}

		blocks := Blocks{}
		jsonErr := json.Unmarshal(res, &blocks)
		if jsonErr != nil {
			fmt.Printf("%v === %v \n", i,nil)
			transactions=append(transactions,nil)
		}else if len(blocks.Blocks)>0{
			block:= &blocks.Blocks[0]

			var ttt *TransformedTransaction
			if block!=nil{
				ttt=&TransformedTransaction{v.transactionHash,TransformedBlock{block.Hash, block.Height, block.Time}}
			}
			fmt.Printf("%v === %v \n", i,ttt)
			transactions=append(transactions,ttt)
		}else {
			fmt.Printf("%v === %v \n", i,nil)
			transactions=append(transactions,nil)
		}

	}
}


func asyncHttpGets(urls *[]Transaction,tcount int) []*HttpResponse {

	ch := make(chan *HttpResponse, tcount) // buffered
	responses := make([]*HttpResponse,0,tcount)
	for _, url := range *urls {
		go func(url string,hash string) {
			//fmt.Printf("Fetching %s \n", url)
			resp, err := http.Get(url)
			ch <- &HttpResponse{url, resp, err,hash}
		}(baseUrl+"block-height/"+strconv.FormatInt(url.BlockHeight,10)+"?format=json",url.Hash)
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
			fmt.Printf(".")
			time.Sleep(5e7)
		}
	}
	return responses

}

/*

func requestBlock1(height int64,i int,hash string,transactions *[]*TransformedTransaction,wg sync.WaitGroup){

	url := baseUrl+"block-height/"+strconv.FormatInt(height,10)+"?format=json"

	res,err:=request1(url)
	if(err!=nil) {
		wg.Done()
		return
	}

	blocks := Blocks{}
	jsonErr := json.Unmarshal(*res, &blocks)
	if jsonErr != nil {
		//panic(jsonErr)
		wg.Done()
		return
	}

	*/
/*for _, v := range blocks.Blocks {
		fmt.Println( v.Hash)
	}*//*

	block:=&blocks.Blocks[0]
	var ttt TransformedTransaction
	if block!=nil{
		ttt=TransformedTransaction{hash,TransformedBlock{block.Hash, block.Height, block.Time}}
		*transactions=append(*transactions,&ttt)
	}else {
		*transactions=append(*transactions,nil)
	}
	fmt.Printf("%v === %v \n", i,ttt)

	wg.Done()
	//return &blocks.Blocks[0]

}

*/


/*


func testRequest1(){
	address:=requestAddr1("1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F")
	tcount:=len(address.Txs)

	signalStringChan :=make(chan *string, tcount)
	for _, v := range address.Txs {
		fmt.Println(v)
		go func() {

			block := requestBlock1(int64(v.BlockHeight))

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
*/

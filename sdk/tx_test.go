package sdk

import (
	"errors"
	"fmt"
	"github.com/imroc/req"
	"github.com/iost-official/go-iost/account"
	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/crypto"
	"github.com/iost-official/go-iost/rpc/pb"
	"github.com/tidwall/gjson"
	"net/http"
	"strconv"
	"testing"
)

// NewIOSTDevSDK creatimg an SDK with reasonable params
func NewIOSTDevSDK1() *IOSTDevSDK {
	keypair,_ := account.NewKeyPair(common.Base58Decode("2yquS3ySrGWPEKywCPzX4RTJugqRh7kJSo5aehsLYPEWkUxBWA39oMrZ7ZxuM4fgyXYs2cPwh5n8aNNpH5x2VyK1"), crypto.Ed25519)

	return &IOSTDevSDK{
		server:              "localhost:30002",
		checkResult:         false,
		checkResultDelay:    3,
		checkResultMaxRetry: 20,
		signAlgo:            "ed25519",
		gasLimit:            1000000,
		gasRatio:            1.0,
		amountLimit:         []*rpcpb.AmountLimit{{Token: "*", Value: "unlimited"}},
		expiration:          90,
		chainID:             uint32(1024),
		keyPair:			 keypair,
		accountName:		 "admin",
		verbose:  			 false,
	}
}

// A Client is a Tron RPC client. It performs RPCs over HTTP using JSON
// request and responses. A Client must be configured with a secret token
// to authenticate with other Cores on the network.
type Client struct {
	BaseURL string
	// AccessToken string
	Debug  bool
	client *req.Req
}

// NewClient create new client to connect
func NewClient(url, token string, debug bool) *Client {
	c := Client{
		BaseURL: url,
		// AccessToken: token,
		Debug: debug,
	}

	api := req.New()
	//trans, _ := api.Client().Transport.(*http.Transport)
	//trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c.client = api

	return &c
}


// Call calls a remote procedure on another node, specified by the path.
func (c *Client) Call(path string, param interface{}) (*gjson.Result, error) {

	if c == nil || c.client == nil {
		return nil, errors.New("API url is not setup. ")
	}

	url := c.BaseURL + path
	authHeader := req.Header{"Accept": "application/json"}

//	r, err := req.Post(url, req.BodyJSON(&param), authHeader)
	r, err := req.Post(url, param, authHeader)
	if err != nil {
	//	log.Error("Failed: %+v >\n", err)
		return nil, err
	}
	// log.Std.Info("%+v", r)

	if r.Response().StatusCode != http.StatusOK {
		message := gjson.ParseBytes(r.Bytes()).String()
		message = fmt.Sprintf("[%s]%s", r.Response().Status, message)
	//	log.Error(message)
		return nil, errors.New(message)
	}

	res := gjson.ParseBytes(r.Bytes())
	return &res, nil
}

func (c * Client) SendTranscation(sign string)(txhash map[string]string, err error) {

	r, err := c.Call("/sendTx", sign)
	if err != nil {
		return nil, err
	}

	txhash_hex := r.Get("hash").String()

	txhash = map[string]string{
		"txhash":    txhash_hex,
	}
	return txhash, nil
}

func SignTX() string{

	client := NewIOSTDevSDK1()

	data := "[\"iost\", \"admin\", \"abcd1\", \"10\", \"\"]"
	action := []*rpcpb.Action{NewAction("token.iost", "transfer", string(data))}

	trx,err := client.CreateTxFromActions(action)
	if err != nil {
		fmt.Println("trx create failed")
	}
	//	fmt.Println("trx create successful")
	signedTx, err := client.SignTx(trx, client.signAlgo)
	if err != nil {

	}
	return MarshalTextString(signedTx)
}

func TestSingTX(t *testing.T) {
	sign := SignTX()

	fmt.Println("sign=",sign)
}


func SendMutilTx(){
	c := NewClient("http://127.0.0.1:30001","",true)

	for i := 1; i <= 600; i++ {
		sign := SignTX()

		//	fmt.Println("sign=",sign)

		r,err := c.SendTranscation(sign)
		if err != nil {
			//t.Errorf("EasyTransferByPrivate failed: %v\n", err)
		} else {
			fmt.Printf("txhash return: %+v\n", r["txhash"])
		}
	}
}

func TestSendMutilTx(t *testing.T) {

	SendMutilTx()

}

var quit chan int // 只开一个信道
func TestSent(t *testing.T) {

		count := 500
		quit = make(chan int) // 无缓冲

		for i := 0; i < count; i++ {
			go SendMutilTx()

		}
		for i := 0; i < count; i++ {
			<- quit
		}
}

//查询GasPrice
func GetBlockByNumber(blockcount int ) (int ,error){

	url := "http://127.0.0.1:30001/getBlockByNumber/" + strconv.Itoa(blockcount) + "/false"
	param := make(req.QueryParam, 0)

	r, err := req.Get(url, param)
	if err != nil {
	//	log.Info(err)
		return 0,err
	}

/*	if c.Debug {
		log.Info("Request API Completed")
	}

	if c.Debug {
		log.Std.Info("%+v", r)
	}
*/
	if err != nil {
		return 0,err
	}

	resp := gjson.ParseBytes(r.Bytes())
/*	err = isError(&resp)
	if err != nil {
		return ""
	}*/

	result := resp.Get("block.tx_count")
	return strconv.Atoi(result.Str)
}
func TestGetBlockByNumber(t *testing.T) {
	start_block :=  219074
	end_block := 219280
	sum_block := end_block - start_block
	all := true
	var sum_txs = 0

	for i := start_block; i <= end_block; i++ {

		txs,_ := GetBlockByNumber(i)
		if all {

			fmt.Printf(" GetBlockByNumber ,height=%d, contain transaction num= %+v ,TPS=%d \n", i, txs ,txs * 2)
		}else{

		}
		sum_txs += txs
	}

	fmt.Printf(" GetBlockByNumber over,sum_time=%d s ,average_txs=%d, average_TPS=%d \n",sum_block/2 , sum_txs/sum_block, sum_txs/sum_block*2)
}

package sdk

import (
	"fmt"
	"github.com/iost-official/go-iost/account"
	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/crypto"
	"github.com/iost-official/go-iost/rpc/pb"
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

func SendTX(){

	client := NewIOSTDevSDK1()

	data := "[\"iost\", \"admin\", \"abcd1\", \"10\", \"\"]"
	action := []*rpcpb.Action{NewAction("token.iost", "transfer", string(data))}

	for i := 0; i < 100 ; i++ {
		trx,err := client.CreateTxFromActions(action)
		if err != nil {
			fmt.Println("trx create failed")
		}
		//	fmt.Println("trx create successful")
		txhash, err := client.SendTx(trx)
		if err != nil {
			fmt.Println("txhash=", txhash)
		}
	}
}

var quit chan int // 只开一个信道
func TestSendTx(t *testing.T) {

	count := 500
	quit = make(chan int) // 无缓冲

	for i := 0; i < count; i++ {
		go SendTX()

	}
	for i := 0; i < count; i++ {
		<- quit
	}
}
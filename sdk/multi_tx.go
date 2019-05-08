package sdk

import (
	"fmt"
	"github.com/iost-official/go-iost/account"
	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/crypto"
	"github.com/iost-official/go-iost/rpc/pb"
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

/*
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
*/
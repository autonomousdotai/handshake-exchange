package bitcoin_service

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/levigross/grequests"
	"github.com/shopspring/decimal"
	"math/big"
	"os"
	// "github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/blockcypher/gobcy"
)

type BitcoinService struct {
}

type Transaction struct {
	TxId               string `json:"txid"`
	SourceAddress      string `json:"source_address"`
	DestinationAddress string `json:"destination_address"`
	Amount             int64  `json:"amount"`
	UnsignedTx         string `json:"unsignedtx"`
	SignedTx           string `json:"signedtx"`
}

var BTC_IN_SATOSHI = decimal.NewFromBigInt(big.NewInt(100000000), 0)

func (c BitcoinService) SendTransaction(address string, amount decimal.Decimal) (string, error) {
	////note the change to BlockCypher Testnet
	bcy := gobcy.API{"b80ed3816ca1414d87e0f7b994f27b16", "btc", "main"}

	//Post New TXSkeleton
	sendAmount := int(amount.Mul(BTC_IN_SATOSHI).IntPart())
	fmt.Println(sendAmount)
	skel, err := bcy.NewTX(gobcy.TempNewTX(address, "1BbPvEanf23t5XkFs3QWpTBUxsSka7f5Y4", sendAmount), false)

	if err != nil {
		fmt.Println(err)
	}
	//Sign it locally
	err = skel.Sign([]string{"c224eaa3609ad43a6226a95e4ea7448f81b240740dfde8c3d4193790b5a6c764"})
	if err != nil {
		fmt.Println(err)
	}
	//Send TXSkeleton
	skel, err = bcy.SendTX(skel)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", skel)

	////note the change to BlockCypher Testnet
	//bcy := gobcy.API{"b80ed3816ca1414d87e0f7b994f27b16", "bcy", "test"}
	////generate two addresses
	//addr1, err := bcy.GenAddrKeychain()
	//addr2, err := bcy.GenAddrKeychain()
	////use faucet to fund first
	//_, err = bcy.Faucet(addr1, 3e5)
	//if err != nil {
	//	fmt.Println(err)
	//}
	////Post New TXSkeleton
	//fmt.Println(addr1.Address)
	//fmt.Println(addr1.OriginalAddress)
	//fmt.Println(addr1.Private)
	//fmt.Println(addr1.Wif)
	//skel, err := bcy.NewTX(gobcy.TempNewTX(addr1.Address, addr2.Address, 2e5), false)
	////Sign it locally
	//err = skel.Sign([]string{addr1.Private})
	//if err != nil {
	//	fmt.Println(err)
	//}
	////Send TXSkeleton
	//skel, err = bcy.SendTX(skel)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("%+v\n", skel)

	return skel.Trans.Hash, err
}

func (c *BitcoinService) SendTransactionRaw(address string, amount decimal.Decimal) (Transaction, error) {
	var transaction Transaction
	wif, err := btcutil.DecodeWIF(os.Getenv("BTC_KEY"))
	nilTransaction := Transaction{}

	if err != nil {
		return nilTransaction, err
	}

	txHash := "ac590a51dac53726db95243d08ae426cce293c71bc6c434414ff0fa6f9c752ea"
	fmt.Println(txHash)
	sendAmount := amount.Mul(BTC_IN_SATOSHI).IntPart()
	addresspubkey, _ := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.MainNetParams)
	sourceTx := wire.NewMsgTx(wire.TxVersion)
	sourceUtxoHash, _ := chainhash.NewHashFromStr(txHash)
	sourceUtxo := wire.NewOutPoint(sourceUtxoHash, 0)
	sourceTxIn := wire.NewTxIn(sourceUtxo, nil, nil)
	destinationAddress, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	sourceAddress, err := btcutil.DecodeAddress(addresspubkey.EncodeAddress(), &chaincfg.MainNetParams)
	if err != nil {
		return nilTransaction, err
	}
	destinationPkScript, _ := txscript.PayToAddrScript(destinationAddress)
	sourcePkScript, _ := txscript.PayToAddrScript(sourceAddress)
	sourceTxOut := wire.NewTxOut(sendAmount, sourcePkScript)
	sourceTx.AddTxIn(sourceTxIn)
	sourceTx.AddTxOut(sourceTxOut)
	sourceTxHash := sourceTx.TxHash()
	redeemTx := wire.NewMsgTx(wire.TxVersion)
	prevOut := wire.NewOutPoint(&sourceTxHash, 0)
	redeemTxIn := wire.NewTxIn(prevOut, nil, nil)
	redeemTx.AddTxIn(redeemTxIn)
	redeemTxOut := wire.NewTxOut(sendAmount, destinationPkScript)
	redeemTx.AddTxOut(redeemTxOut)
	sigScript, err := txscript.SignatureScript(redeemTx, 0, sourceTx.TxOut[0].PkScript, txscript.SigHashAll, wif.PrivKey, false)
	if err != nil {
		return nilTransaction, err
	}
	redeemTx.TxIn[0].SignatureScript = sigScript
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sourceTx.TxOut[0].PkScript, redeemTx, 0, flags, nil, nil, sendAmount)
	if err != nil {
		return nilTransaction, err
	}
	if err := vm.Execute(); err != nil {
		return nilTransaction, err
	}
	var unsignedTx bytes.Buffer
	var signedTx bytes.Buffer
	sourceTx.Serialize(&unsignedTx)
	redeemTx.Serialize(&signedTx)
	transaction.TxId = sourceTxHash.String()
	transaction.UnsignedTx = hex.EncodeToString(unsignedTx.Bytes())
	transaction.Amount = sendAmount
	transaction.SignedTx = hex.EncodeToString(signedTx.Bytes())
	transaction.SourceAddress = sourceAddress.EncodeAddress()
	transaction.DestinationAddress = destinationAddress.EncodeAddress()

	fmt.Println(transaction.SourceAddress)
	url := fmt.Sprintf("https://insight.bitpay.com/api/tx/send")
	body := fmt.Sprintf("rawtx: %s", transaction.SignedTx)
	fmt.Println(body)
	r := bytes.NewReader([]byte(body))
	headers := map[string]string{}
	ro := &grequests.RequestOptions{Headers: headers, RequestBody: r}
	resp, err := grequests.Post(url, ro)
	if resp.Ok != true {
		// return nilTransaction, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return transaction, nil
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/itering/substrate-api-rpc/keyring"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/vedhavyas/go-subkey"
	"github.com/vedhavyas/go-subkey/sr25519"
	"io"
	"log"
	"math/big"
	"net/http"
	"time"
)

const Network = "paseo"
const ApiKey = ""
const TransfersCount = 5

const SourceWalletMnemonic = ""

const DestinationWalletAddress = ""
const DestinationWalletPubKey = ""

var lastTransferId = 639977000046

func main() {
	httpClient := &http.Client{}

	scheme := sr25519.Scheme{}
	kr, err := subkey.DeriveKeyPair(scheme, SourceWalletMnemonic)
	if err != nil {
		log.Fatalf("subkey.DeriveKeyPair(scheme, SourceWalletMnemonic): %v", err)
	}

	sourceWalletSeed := types.HexEncodeToString(kr.Seed())
	sourceWalletAddress := kr.SS58Address(0)

	rpcClient := &rpc.Client{}
	rpcClient.SetKeyRing(keyring.New(keyring.Sr25519Type, sourceWalletSeed))

	setLastTransferId(httpClient, sourceWalletAddress)
	log.Printf("last transfer id: %d\n", lastTransferId)

	transfers := make(chan Transfer, TransfersCount)
	go subscribeWalletTransfers(sourceWalletAddress, httpClient, transfers)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case transfer := <-transfers:
			if transfer.Amount == nil {
				log.Fatal("transfer.Amount is nil")
			}

			amount := new(big.Float)
			amount.SetString(*transfer.Amount)
			log.Printf("transfer.Amount %s", amount.String())
			send(rpcClient, amount)

		case <-ticker.C:
			log.Println("nothing happens")
		}
	}
}

func send(cli *rpc.Client, amount *big.Float) {
	signedTransaction, err := cli.SignTransaction(
		"Balances",
		"transfer",
		map[string]interface{}{"Id": DestinationWalletPubKey},
		amount,
	)

	if err != nil {
		log.Fatalf("Failed to sign transaction with %s signature: %v", signedTransaction, err)
	}

	transactionHash, err := cli.SendAuthorSubmitExtrinsic(signedTransaction)
	if err != nil {
		log.Fatalf("Failed to send transaction with %s txHash: %v", transactionHash, err)
	}
}

func subscribeWalletTransfers(address string, cli *http.Client, transfers chan Transfer) {
	ticker := time.NewTicker(200 * time.Millisecond)

	defer ticker.Stop()

	for range ticker.C {
		ts := getWalletLastTransfers(address, cli)
		for _, t := range ts {
			transfers <- t
		}
	}
}

func setLastTransferId(cli *http.Client, address string) {
	transfers := getWalletLastTransfers(address, cli)
	if len(transfers) == 0 {
		log.Fatal("setLastTransferId.len(transfers) == 0")
	}

	if transfers[len(transfers)-1].TransferID == nil {
		log.Fatalf("transfers[len(transfers)-1].TransferID is nil")
	}

	lastTransferId = *transfers[len(transfers)-1].TransferID
}

func getWalletLastTransfers(address string, cli *http.Client) []Transfer {
	apiUrl := fmt.Sprintf("https://%s.api.subscan.io/api/v2/scan/transfers", Network)

	ticker := time.Tick(200 * time.Millisecond)

	request := TransfersListRequest{
		Address: StringPtr(address),
		Page:    IntPtr(0),
		Row:     IntPtr(100),
		AfterID: []*int{IntPtr(lastTransferId)},
	}

	transfers := make([]Transfer, 0, TransfersCount)

	for len(transfers) < TransfersCount {
		<-ticker

		payload, err := json.Marshal(request)
		if err != nil {
			log.Fatalf("json.Marshal: %v", err)
		}

		req, err := http.NewRequest("POST", apiUrl, bytes.NewReader(payload))
		if err != nil {
			log.Fatalf("http.NewRequest: %v", err)
		}

		req.Header.Add("x-api-key", ApiKey)
		req.Header.Add("Content-Type", "application/json")

		res, err := cli.Do(req)
		if err != nil {
			log.Fatalf("client.Do: %v", err)
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("io.ReadAll: %v", err)
		}

		response := TransfersListResponse{}
		if err = json.Unmarshal(body, &response); err != nil {
			fmt.Println(string(body))
			log.Fatalf("json.Unmarshal: %v", err)
		}

		if response.Code != 0 {
			log.Fatalf("code %d: %s", response.Code, response.Message)
		}

		if len(response.Data.Transfers) == 0 {
			break
		}

		transfers = append(transfers, response.Data.Transfers...)

		if response.Data.Transfers[len(response.Data.Transfers)-1].TransferID == nil {
			log.Fatal("response.Data.Transfers[len(response.Data.Transfers)-1].TransferID is nil")
		}

		if len(response.Data.Transfers) == 0 {
			log.Fatal("getWalletLastTransfers.len(response.Data.Transfers) == 0")
		}

		lastTransferId = *response.Data.Transfers[len(response.Data.Transfers)-1].TransferID

		if *request.Page < 99 {
			*request.Page += 1
		}

		request.AfterID = []*int{&lastTransferId}
	}

	return transfers
}

func getSymbolTokenList(cli *http.Client) SymbolTokenListResponse {
	apiURL := fmt.Sprintf("https://%s.api.subscan.io/api/scan/token", Network)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatal("creating request:", err)
	}

	req.Header.Add("x-api-key", ApiKey)
	res, err := cli.Do(req)
	if err != nil {
		log.Fatal("making request:", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("reading response:", err)
	}

	var response SymbolTokenListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal("decoding JSON:", err)
	}

	return response
}

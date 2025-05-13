package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/itering/substrate-api-rpc/keyring"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/vedhavyas/go-subkey"
	"github.com/vedhavyas/go-subkey/sr25519"
	"io"
	"log"
	"net/http"
	"time"
)

const Network = "paseo"
const ApiKey = ""

const SourceWalletMnemonic = ""

const DestinationWalletAddress = "1569vLJdqPR4eHCcWtp1Z72GbK23XCsiPLWNzzppYknF4X3r"
const DestinationWalletPubKey = "0xb4df2fc883f41c27cdbb67443aad34d5c1eb3ef6e97037121ae3fb95a10e5679"

var transferIds = map[int]struct{}{}

func main() {
	httpClient := &http.Client{}
	rpcClient := &rpc.Client{}

	websocket.SetEndpoint("ws://127.0.0.1:9944")

	scheme := sr25519.Scheme{}
	kr, err := subkey.DeriveKeyPair(scheme, SourceWalletMnemonic)
	if err != nil {
		log.Fatalf("subkey.DeriveKeyPair(scheme, SourceWalletMnemonic): %v", err)
	}

	sourceWalletSeed := types.HexEncodeToString(kr.Seed())
	sourceWalletAddress := kr.SS58Address(0)

	getWalletLastTransfers(sourceWalletAddress, httpClient)
	rpcClient.SetKeyRing(keyring.New(keyring.Sr25519Type, sourceWalletSeed))

	raw, err := rpc.GetMetadataByHash(nil)
	if err != nil {
		log.Fatalf("GetMetadataByHash: %v", err)
	}

	rpcClient.SetMetadata(metadata.RegNewMetadataType(92, raw))

	send(rpcClient, "1")

	//transfers := make(chan Transfer, 10)
	//go subscribeWalletTransfers(sourceWalletAddress, httpClient, transfers)
	//
	//ticker := time.NewTicker(5 * time.Second)
	//defer ticker.Stop()
	//
	//for transfer := range transfers {
	//	if transfer.Amount == nil {
	//		log.Fatal("transfer.Amount is nil")
	//	}
	//
	//	log.Printf("transfer.Amount %s %s", *transfer.Amount, *transfer.AssetUniqueID)
	//	send(rpcClient, *transfer.Amount)
	//}
}

func send(cli *rpc.Client, amount string) {
	signedTransaction, err := cli.SignTransaction(
		"balances",
		"transfer",
		map[string]interface{}{"Id": DestinationWalletAddress},
		amount,
	)

	if err != nil {
		log.Fatalf("Failed to sign transaction with %s signature: %v", signedTransaction, err)
	}

	transactionHash, err := cli.SendAuthorSubmitAndWatchExtrinsic(signedTransaction)
	if err != nil {
		log.Fatalf("Failed to send transaction with %s txHash: %v", transactionHash, err)
	}
}

func subscribeWalletTransfers(address string, cli *http.Client, transfers chan Transfer) {
	for {
		ts := getWalletLastTransfers(address, cli)
		for _, t := range ts {
			transfers <- t
		}
	}
}

func getWalletLastTransfers(address string, cli *http.Client) []Transfer {
	apiUrl := fmt.Sprintf("https://%s.api.subscan.io/api/v2/scan/transfers", Network)

	ticker := time.Tick(200 * time.Millisecond)

	request := TransfersListRequest{
		Address: StringPtr(address),
		Page:    IntPtr(0),
		Row:     IntPtr(100),
		Order:   StringPtr("asc"),
	}

	transfers := make([]Transfer, 0, 100)

	for {
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

		if len(response.Data.Transfers) == 0 {
			log.Fatal("getWalletLastTransfers.len(response.Data.Transfers) == 0")
		}

		if *request.Page < 99 {
			*request.Page += 1
		}

		for _, transfer := range response.Data.Transfers {
			if transfer.To == nil {
				log.Fatal("transfer.To is nil")
			}

			if transfer.Amount == nil {
				log.Fatal("transfer.Amount is nil")
			}

			if transfer.AssetSymbol == nil {
				log.Fatal("transfer.AssetSymbol is nil")
			}

			if _, ok := transferIds[*transfer.TransferID]; ok {
				continue
			}

			//log.Printf("transfer.Amount %s %s", *transfer.Amount, *transfer.AssetUniqueID)
			transferIds[*transfer.TransferID] = struct{}{}

			if *transfer.To == address {
				//log.Printf("received %s %s", *transfer.Amount, *transfer.AssetSymbol)
				transfers = append(transfers, transfer)
			}
		}
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

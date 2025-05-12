package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const Network = "acala"
const ApiKey = ""
const TransfersCount = 5

func main() {
	client := &http.Client{}
	apiUrl := fmt.Sprintf("https://%s.api.subscan.io/api/v2/scan/transfers", Network)

	symbolTokens := getSymbolTokenList(client)

	ticker := time.Tick(200 * time.Millisecond)

	request := TransfersListRequest{
		Address: StringPtr("23M5ttkmR6Kco7bReRDve6bQUSAcwqebatp3fWGJYb4hDSDJ"),
		Page:    IntPtr(0),
		Row:     IntPtr(TransfersCount),
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

		res, err := client.Do(req)
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

		lastTransferID := response.Data.Transfers[len(response.Data.Transfers)-1].TransferID

		if *request.Page < 99 {
			*request.Page += 1
		}

		request.AfterID = []*int{lastTransferID}
		fmt.Println(*lastTransferID)
	}

	for _, transfer := range transfers {
		if transfer.From == nil || transfer.To == nil {
			log.Fatal("transfer.From or transfer.To is empty")
		}

		if transfer.Amount == nil || transfer.AssetSymbol == nil {
			log.Fatal("transfer.AssetSymbol is empty")
		}

		if symbolTokens.Data.Detail[*transfer.AssetSymbol] == nil {
			log.Fatal("symbolTokens.Data.Detail[*transfer.AssetSymbol] is empty")
		}

		if symbolTokens.Data.Detail[*transfer.AssetSymbol].Price == nil {
			log.Fatal("symbolTokens.Data.Detail[*transfer.AssetSymbol].Price is empty")
		}

		amount, err := strconv.ParseFloat(*transfer.Amount, 64)
		if err != nil {
			log.Fatalf("1. strconv.ParseFloat: %v", err)
		}

		symbUsdPrice, err := strconv.ParseFloat(*symbolTokens.Data.Detail[*transfer.AssetSymbol].Price, 64)
		if err != nil {
			log.Fatalf("2. strconv.ParseFloat: %v", err)
		}

		dotUsdPrice, err := strconv.ParseFloat(*symbolTokens.Data.Detail["DOT"].Price, 64)
		if err != nil {
			log.Fatalf("2. strconv.ParseFloat: %v", err)
		}

		symbUsd := amount * symbUsdPrice
		dots := symbUsd / dotUsdPrice

		fmt.Printf("\n%s ---> %s\n", *transfer.From, *transfer.To)
		fmt.Printf("\t%s %s\n", *transfer.Amount, *transfer.AssetSymbol)
		fmt.Printf("\t%f DOT\n", dots)
		fmt.Printf("\t%f USD\n", symbUsd)
	}
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

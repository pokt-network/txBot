package main

import (
	"context"
	"fmt"

	Spec "github.com/pokt-network/txbot/spec"
)


func QueryHeight(config Config, client *Spec.ClientWithResponses, clientCtx context.Context) {
	var body interface{}
	res, err := client.PostQueryHeightWithResponse(clientCtx, body)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Printf("Height: %d\n", *res.JSON200.Height)
}

func RelayEth(config Config, client *Spec.ClientWithResponses, clientCtx context.Context) {

	// Prepare metadata.
	blockHeight := int64(0)
	meta := Spec.RelayMetadata{
		BlockHeight: &blockHeight,
	}
	// Prepare payload.
	data := "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\":[\"0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8\",\"latest\"],\"id\":67}"
	headers := Spec.RelayHeader{
		AdditionalProperties: make(map[string]string),
	}
	method := "POST"
	path := "/v1/" + "config.AAT"
	payload := Spec.RelayPayload{
		Data: &data,
		Headers: &headers,
		Method: &method,
		Path: &path,
	}

	// Prepare AAT.
	appPubKey := ""
	clientPubKey := ""
	signature := ""
	version := ""
	aat := Spec.AAT{
		AppPubKey: &appPubKey,
		ClientPubKey: &clientPubKey,
		Signature: &signature,
		Version: &version,
	}

	// Prepare proof.
	blockchain := "eth"
	entropy := int64(0)
	requestHash := ""
	servicerPubKey := ""
	sessionBlockHeight := int64(0)
	signature = ""
	proof := Spec.RelayProof{
		Aat: &aat,
		Blockchain: &blockchain,
		Entropy: &entropy,
		RequestHash: &requestHash,
		ServicerPubKey: &servicerPubKey,
		SessionBlockHeight: &sessionBlockHeight,
		Signature: &signature,
	}

	// Prepare and make request.
	body := Spec.PostClientRelayJSONRequestBody {
		Meta: &meta,
		Payload: &payload,
		Proof: &proof,
	}
	res, err := client.PostClientRelayWithResponse(clientCtx, body)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Printf("Eth relay response payload: %v\n", *res.JSON200.Payload)
}
package rpc

import (
	sha "crypto"

	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	config "github.com/pokt-network/txbot/config"
	spec "github.com/pokt-network/txbot/spec"
	// cryptoTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
)

const (
	HARMONY string = "0040" // Harmony mainnet shard 0.
	ETHEREUM string = "0021" // Ethereum mainnet.
	IPFS string = "1111" // IPFS.
)

var Hasher = sha.SHA3_256

func Hash(b []byte) []byte {
	hasher := Hasher.New()
	hasher.Write(b) //nolint:golint,errcheck
	return hasher.Sum(nil)
}

func QueryHeight(config config.Config, client *spec.ClientWithResponses, clientCtx context.Context) {
	var body interface{}
	res, err := client.PostQueryHeightWithResponse(clientCtx, body)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	fmt.Printf("Height: %d\n", *res.JSON200.Height)
}

func RelayHmy(config config.Config, client *spec.ClientWithResponses, clientCtx context.Context) {
	data := `{"jsonrpc":"2.0", "method":"hmyv2_latestHeader", "params":[], "id":1}`
	blockchain := HARMONY
	Relay(blockchain, data, config, client, clientCtx)
}

func RelayEth(config config.Config, client *spec.ClientWithResponses, clientCtx context.Context) {
	data := `{"jsonrpc":"2.0", "method":"eth_getBalance", "params":["0xe7a24E61b2ec77d3663ec785d1110688d2A32ecc", "latest"], "id":1}`
	blockchain := ETHEREUM
	Relay(blockchain, data, config, client, clientCtx)
}

func Relay(blockchain string, data string, config config.Config, client *spec.ClientWithResponses, clientCtx context.Context) {
	// Get client keys.
	clientPrivKey := config.GetRandomPrivateKey()
	clientPubKey := clientPrivKey.PublicKey().String()

	// Get the blockchain service node.
	serviceNode := getBlockchainServiceNode(blockchain, client, clientCtx)

	// Prepare metadata.
	blockHeight := int64(1) // Does this need to match the current pocket height?
	meta := spec.RelayMetadata{
		BlockHeight: &blockHeight,
	}

	// Prepare payload.
	headers := spec.RelayHeader{
		AdditionalProperties: make(map[string]string),
	}
	method := "POST"
	path := "" // TODO: What shold this be?
	payload := spec.RelayPayload{
		Data: &data,
		Method: &method,
		Path: &path,
		Headers: &headers,
	}

	// Prepare request.
	request := struct {
		Payload spec.RelayPayload  `json:"payload"`
		Meta    spec.RelayMetadata `json:"meta"`
	}{payload, meta}

	// Compute request hash.
	requestBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	requestHash := Hash(requestBytes)
	requestHashString := hex.EncodeToString(requestHash)

	indented, _ := json.MarshalIndent(request, "", "  ")
	fmt.Printf("%s \n\n %v \n\n %s \n\n %v \n\n", string(indented), requestBytes, requestHashString, requestHash)

	// Prepare AAT.
	aatVersion := "0.0.1"
	appPubKey := clientPubKey // Is this okay?
	appPrivKey := clientPrivKey
	aat := spec.AAT{
		AppPubKey: &appPubKey,
		ClientPubKey: &clientPubKey,
		Signature: new(string),
		Version: &aatVersion,
	}

	// Sign AAT.
	aatBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	aatHash := Hash(aatBytes)
	appSig, err := appPrivKey.Sign(aatHash)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	appSigString := hex.EncodeToString(appSig)
	aat.Signature = &appSigString

	// Prepare proof.
	entropy := int64(1)
	sessionBlockHeight := int64(1)
	servicerPubKey := serviceNode.PublicKey
	if servicerPubKey == nil {
		return
	}
	proof := spec.RelayProof{
		Aat: &aat,
		Blockchain: &blockchain,
		Entropy: &entropy,
		RequestHash: &requestHashString,
		ServicerPubKey: servicerPubKey,
		SessionBlockHeight: &sessionBlockHeight,
		Signature: new(string),
	}

	// Sign proof.
	proofBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	proofHash := Hash(proofBytes)
	clientSig, err := clientPrivKey.Sign(proofHash)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	clientSigString := hex.EncodeToString(clientSig)
	proof.Signature = &clientSigString

	// Prepare relay request.
	body := spec.PostClientRelayJSONRequestBody {
		Meta: &meta,
		Payload: &payload,
		Proof: &proof,
	}

	// Do relay.
	res, err := client.PostClientRelayWithResponse(clientCtx, body)
	if err != nil {
		fmt.Println("PostClientRelayWithResponse Error: ", err.Error())
		return
	}
	fmt.Println(res)
}

func getBlockchainServiceNode(blockchain string, client *spec.ClientWithResponses, clientCtx context.Context) (val spec.Node) {
	body := spec.PostQueryNodesJSONRequestBody {
		Blockchain: &blockchain,
	}
	res, err := client.PostQueryNodesWithResponse(clientCtx, body)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	serviceNodes := *res.JSON200.Result
	if len(serviceNodes) == 0 {
		fmt.Println("Error: No service nodes found for blockchain: ", blockchain)
		return
	}
	val = serviceNodes[0]
	return
}
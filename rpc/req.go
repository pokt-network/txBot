package rpc

import (
	sha "crypto"

	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"

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

var globalSessionBlockHeight = int64(-1)

func Hash(b []byte) []byte {
	hasher := Hasher.New()
	hasher.Write(b) //nolint:golint,errcheck
	return hasher.Sum(nil)
}

func QueryHeight(config config.Config, client *spec.ClientWithResponses, clientCtx context.Context) int64 {
	var body interface{}
	res, err := client.PostQueryHeightWithResponse(clientCtx, body)
	if err != nil {
		fmt.Println("Error", err)
		return int64(-1) // TODO: return error
	}
	height := *res.JSON200.Height
	fmt.Printf("Height: %d\n", height)
	return height
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
	// fmt.Printf("Client private key: %s\n", clientPrivKey)
	clientPubKey := clientPrivKey.PublicKey().RawString()
	// fmt.Printf("Client pub key: %s\n", clientPubKey)

	// Get the blockchain service node.
	// serviceNode := getBlockchainServiceNode(blockchain, client, clientCtx)

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

	// indented, _ := json.MarshalIndent(request, "", "  ")
	// fmt.Printf("%s \n\n %v \n\n %s \n\n %v \n\n", string(indented), requestBytes, requestHashString, requestHash)

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
	aatBytes, err := json.Marshal(aat)
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
	// fmt.Printf("%v \n\n %v \n\n", appSigString, hex.EncodeToString(aatHash))
	// fmt.Println("---", string(aatBytes), hex.EncodeToString(aatHash))
	aat.Signature = &appSigString

	// fmt.Println(appPrivKey.PublicKey().RawString(),hex.EncodeToString(aatHash), appSigString)
	// t_sig, _ := hex.DecodeString(appSigString)
	// t_pk, _ := crypto.NewPublicKey(appPrivKey.PublicKey().RawString())
	// t_msg, _ := hex.DecodeString(hex.EncodeToString(aatHash))
	// fmt.Println(t_sig, t_pk, t_msg)
	// if ok := t_pk.VerifyBytes(t_msg, t_sig); !ok {
	// 	fmt.Println("FAILED")
	// } else {
	// 	fmt.Println("OK")
	// }

	// if ok := appPrivKey.PublicKey().VerifyBytes(aatHash, appSig); !ok {
	// 	fmt.Println("FAILED")
	// } else {
	// 	fmt.Println("OK")
	// }

	// Prepare proof.
	entropy := int64(rand.Uint32())
	// sessionBlockHeight := QueryHeight(config, client, clientCtx)// - 1
	sessionBlockHeight := globalSessionBlockHeight

	servicerPubKey := "eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743"
	// servicerPubKey := serviceNode.PublicKey
	// if servicerPubKey == nil {
	// 	return
	// }
	proof := spec.RelayProof{
		Aat: &aat,
		Blockchain: &blockchain,
		Entropy: &entropy,
		RequestHash: &requestHashString,
		ServicerPubKey: &servicerPubKey,
		SessionBlockHeight: &sessionBlockHeight,
		Signature: new(string),
	}

	proofForSig := struct {
		Entropy            int64  `json:"entropy"`
		SessionBlockHeight int64  `json:"session_block_height"`
		ServicerPubKey     string `json:"servicer_pub_key"`
		Blockchain         string `json:"blockchain"`
		Signature          string `json:"signature"`
		Token              string `json:"token"`
		RequestHash        string `json:"request_hash"`
	}{entropy, sessionBlockHeight, servicerPubKey, blockchain, "", hex.EncodeToString(aatHash), requestHashString}

	fmt.Println()

	// Sign proof.
	proofBytes, err := json.Marshal(proofForSig)
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

	// t_sig, _ := hex.DecodeString(clientSigString)
	// t_pk, _ := crypto.NewPublicKey(clientPrivKey.PublicKey().RawString())
	// t_msg, _ := hex.DecodeString(hex.EncodeToString(proofHash))
	// fmt.Println(t_sig, t_pk, t_msg)
	// if ok := t_pk.VerifyBytes(t_msg, t_sig); !ok {
	// 	fmt.Println("FAILED")
	// } else {
	// 	fmt.Println("OK")
	// }

	// if ok := clientPrivKey.PublicKey().VerifyBytes(proofHash, clientSig); !ok {
	// 	fmt.Println("FAILED")
	// } else {
	// 	fmt.Println("OK")
	// }


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

	if res.StatusCode() == 400 {
		if res.JSON400.Dispatch != nil {
			globalSessionBlockHeight = *res.JSON400.Dispatch.Session.Header.SessionHeight
		}
		// res, _ := json.MarshalIndent(res.JSON400, " ", "")
		// fmt.Println("Error 400: ", string(res))
		// fmt.Println(res.JSON400.Dispatch.BlockHeight, )
		return
	}

	if res.StatusCode() == 200 {
		fmt.Println("200", *res.JSON200.Signature, res.JSON200.Payload, string(res.Body))
	}
	// fmt.Println(res)
	// return res
	// if res.StatusCode != 200 {
	// 	res400, _ := json.Unmarshal(*res.JSON400)
	// 	fmt.Println(res400)
	// }
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

// func getSessionBlockHeight(blockchain string, client *spec.ClientWithResponses, clientCtx context.Context) (val int64) {
// 	body := spec.PostQuerySessionBlockHeightJSONRequestBody {
// 	}
// 	client.PostSessionBlockHeightWithResponse(clientCtx, blockchain)
// }

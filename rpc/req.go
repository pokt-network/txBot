package rpc

import (
	sha "crypto"

	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/pokt-network/pocket-core/crypto"
	config "github.com/pokt-network/txbot/config"
	spec "github.com/pokt-network/txbot/spec"
)

const (
	Polygon        string = "0009" // Polygon mainnet
	Ethereum       string = "0021" // Ethereum mainnet.
	EthereumGoerli string = "0026" // Ethereum Goerli testnet

	PocketMainNet string = "0001" // Pocket TestNet
	PocketTestNet string = "0002" // Pocket TestNet

	Hasher = sha.SHA3_256
)

type RpcContext struct {
	Config  config.Config
	Client  *spec.ClientWithResponses // This client can be used for general RPC calls such as retrieving the height
	Context context.Context

	// Application specific fields.
	AppPrivKey crypto.PrivateKey
	AppPubKey  string

	// Session specific fields.
	Session       *spec.Session
	Servicer      *spec.Node                // The servicer node that will be used for relaying.
	SessionClient *spec.ClientWithResponses // A client specific to the servicer being used in the session
}

func NewRpcContext(config config.Config) *RpcContext {
	client, clientErr := spec.NewClientWithResponses(fmt.Sprintf("%s/v1", config.PocketEndpoint))
	if clientErr != nil {
		panic("Could not initialize RPC client.")
	}
	return &RpcContext{
		Config:  config,
		Client:  client,
		Context: context.TODO(), // Not important at the moment.
		Session: nil,
	}
}

func QueryHeight(config config.Config, rpcCtx *RpcContext) int64 {
	var body interface{}
	res, err := rpcCtx.Client.PostQueryHeightWithResponse(rpcCtx.Context, body)
	if err != nil {
		fmt.Println(err)
	}
	if res == nil {
		fmt.Println("ERROR: Please check your RPC endpoint.")
	}
	if res.JSON200 != nil {
		return *res.JSON200.Height
	} else {
		fmt.Printf("Error querying height: %v\n", *res.HTTPResponse)
	}
	return int64(0)
}

func RelayPolyHeight(config config.Config, rpcCtx *RpcContext) {
	data := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":"v0_localnet"}`
	relay(Polygon, data, config, rpcCtx)
}

func RelayEthHeight(config config.Config, rpcCtx *RpcContext) {
	data := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":"v0_localnet"}`
	relay(Ethereum, data, config, rpcCtx)
}

func RelayEthGoerliHeight(config config.Config, rpcCtx *RpcContext) {
	data := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":"v0_localnet"}`
	relay(EthereumGoerli, data, config, rpcCtx)
}

func RelayPocketHeight(config config.Config, rpcCtx *RpcContext) {
	data := `{"jsonrpc":"2.0","method":"height","params":[],"id":"v0_localnet"}`
	relay(PocketTestNet, data, config, rpcCtx)
}

func relay(blockchain string, data string, config config.Config, rpcCtx *RpcContext) {
	if rpcCtx.Session == nil {
		fmt.Println("No session found...")
		createSession(blockchain, config, rpcCtx)
	}

	// Get the blockchain service node.

	servicerPubKey := rpcCtx.Servicer.PublicKey

	// Prepare metadata.
	meta := spec.RelayMetadata{
		BlockHeight: rpcCtx.Session.Header.SessionHeight,
	}

	// Prepare payload.
	headers := spec.RelayHeader{
		AdditionalProperties: make(map[string]string),
	}
	method := "POST"
	path := "" // TODO: What should this be?
	payload := spec.RelayPayload{
		Data:    &data,
		Method:  &method,
		Path:    &path,
		Headers: &headers,
	}

	// Prepare request. NOTE: request serialization need to be customized.
	request := struct {
		Payload spec.RelayPayload  `json:"payload"`
		Meta    spec.RelayMetadata `json:"meta"`
	}{payload, meta}

	// Compute request hash.
	requestBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	requestHash := hash(requestBytes)
	requestHashString := hex.EncodeToString(requestHash)

	// Prepare AAT.
	aatVersion := "0.0.1"
	aat := spec.AAT{
		AppPubKey:    &rpcCtx.AppPubKey,
		ClientPubKey: &rpcCtx.AppPubKey,
		Signature:    new(string),
		Version:      &aatVersion,
	}

	// Sign AAT.
	aatBytes, err := json.Marshal(aat)
	if err != nil {
		fmt.Println(err)
		return
	}
	aatHash := hash(aatBytes)
	appSig, err := rpcCtx.AppPrivKey.Sign(aatHash)
	if err != nil {
		fmt.Println(err)
		return
	}
	appSigString := hex.EncodeToString(appSig)
	aat.Signature = &appSigString

	// Prepare proof.
	entropy := int64(rand.Uint32())
	proof := spec.RelayProof{
		Aat:                &aat,
		Blockchain:         &blockchain,
		Entropy:            &entropy,
		RequestHash:        &requestHashString,
		ServicerPubKey:     servicerPubKey,
		SessionBlockHeight: rpcCtx.Session.Header.SessionHeight,
		Signature:          new(string),
	}

	// NOTE: proof serialization need to be customized.
	proofForSig := struct {
		Entropy            int64  `json:"entropy"`
		SessionBlockHeight int64  `json:"session_block_height"`
		ServicerPubKey     string `json:"servicer_pub_key"`
		Blockchain         string `json:"blockchain"`
		Signature          string `json:"signature"`
		Token              string `json:"token"`
		RequestHash        string `json:"request_hash"`
	}{entropy, *rpcCtx.Session.Header.SessionHeight, *servicerPubKey, blockchain, "", hex.EncodeToString(aatHash), requestHashString}

	// Sign proof.
	proofBytes, err := json.Marshal(proofForSig)
	if err != nil {
		fmt.Println(err)
		return
	}
	proofHash := hash(proofBytes)
	appSigProof, err := rpcCtx.AppPrivKey.Sign(proofHash)
	if err != nil {
		fmt.Println(err)
		return
	}
	appSigProofString := hex.EncodeToString(appSigProof)
	proof.Signature = &appSigProofString

	// Prepare relay request.
	body := spec.PostClientRelayJSONRequestBody{
		Meta:    &meta,
		Payload: &payload,
		Proof:   &proof,
	}

	// Do relay.
	res, err := rpcCtx.SessionClient.PostClientRelayWithResponse(rpcCtx.Context, body)
	if err != nil {
		fmt.Println("PostClientRelayWithResponse error: ", err.Error())
		return
	}

	if res == nil {
		fmt.Println("ERROR: Please check your RPC endpoint.")
		return
	}

	switch res.StatusCode() {
	case 200:
		{
			fmt.Println(string(res.Body))
		}
	case 400:
		{
			if res.JSON400.Error != nil {
				fmt.Printf("Error sending relay (height %d): %s", *rpcCtx.Session.Header.SessionHeight, *res.JSON400.Error.Message)
			}
			if res.JSON400.Dispatch != nil {
				panic(fmt.Sprintf("Could not send the relay due to %+v\n", *res.JSON400))
			}
			return
		}
	default:
		fmt.Println("Unexpected status code: ", res.StatusCode())
	}
}

// WARNING: There is no business logic to expiring sessions. The client needs to be restarted.
func createSession(blockchain string, config config.Config, rpcCtx *RpcContext) {
	rpcCtx.AppPrivKey = config.GetRandomAppPrivateKey()
	rpcCtx.AppPubKey = rpcCtx.AppPrivKey.PublicKey().RawString()

	height := QueryHeight(config, rpcCtx)
	body := spec.PostClientDispatchJSONRequestBody{
		AppPublicKey:  &rpcCtx.AppPubKey,
		Chain:         &blockchain,
		SessionHeight: &height,
	}
	fmt.Printf("Creating a new session for App with Pub Key %s for Chain %s at height %d\n", *body.AppPublicKey, *body.Chain, *body.SessionHeight)

	res, err := rpcCtx.Client.PostClientDispatchWithResponse(rpcCtx.Context, body)
	if err != nil {
		fmt.Println("PostClientRelayWithResponse error: ", err.Error())
		return
	}

	if res.JSON200 == nil {
		fmt.Printf("Error dispatching new session: %+v\n", *res.HTTPResponse)
		panic("Need to be able to dispatch a session to continue")
	} else {
	}

	// fmt.Printf("Dispatched a new session. \n\t Body: %+v\n\t Response: %+v\n", res.Body, res.HTTPResponse)
	fmt.Printf("Dispatched a new session at height: %d \n", *res.JSON200.Session.Header.SessionHeight)
	rpcCtx.Session = res.JSON200.Session
	rpcCtx.Servicer = &(*rpcCtx.Session.Nodes)[0]

	mappedUrl, ok := rpcCtx.Config.UrlMapping[*rpcCtx.Servicer.ServiceUrl]
	var client *spec.ClientWithResponses
	var clientErr error
	if ok {
		fmt.Printf("Using mapped url from %s to %s\n.", *rpcCtx.Servicer.ServiceUrl, mappedUrl)
		client, clientErr = spec.NewClientWithResponses(fmt.Sprintf("%s/v1", mappedUrl))
	} else {
		fmt.Println("Using original url: ", *rpcCtx.Servicer.ServiceUrl)
		client, clientErr = spec.NewClientWithResponses(fmt.Sprintf("%s/v1", *rpcCtx.Servicer.ServiceUrl))
	}
	if clientErr != nil {
		panic("Could not initialize RPC client for servicer.")
	}
	rpcCtx.SessionClient = client
}

func hash(b []byte) []byte {
	hasher := Hasher.New()
	hasher.Write(b) //nolint:golint,errcheck
	return hasher.Sum(nil)
}

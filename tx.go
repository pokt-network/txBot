package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
	"github.com/pokt-network/pocket-core/codec"
	types2 "github.com/pokt-network/pocket-core/codec/types"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/types/module"
	apps "github.com/pokt-network/pocket-core/x/apps"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/auth"
	authTypes "github.com/pokt-network/pocket-core/x/auth/types"
	"github.com/pokt-network/pocket-core/x/gov"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocket "github.com/pokt-network/pocket-core/x/pocketcore"
	"github.com/tjarratt/babble"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)

var memCDC *codec.Codec

func memCodec() *codec.Codec {
	if memCDC == nil {
		memCDC = codec.NewCodec(types2.NewInterfaceRegistry())
		module.NewBasicManager(
			apps.AppModuleBasic{},
			auth.AppModuleBasic{},
			gov.AppModuleBasic{},
			nodes.AppModuleBasic{},
			pocket.AppModuleBasic{},
		).RegisterCodec(memCDC)
		types.RegisterCodec(memCDC)
		crypto.RegisterAmino(memCDC.AminoCodec().Amino)
	}
	return memCDC
}

func StakeNodeTransaction(config Config) {
	signer := config.GetRandomPrivateKey()
	pk := signer.PublicKey()
	msg := nodesTypes.MsgStake{
		PublicKey:  pk,
		Chains:     GetRandomChains(),
		Value:      GetRandomStake(),
		ServiceUrl: GetRandomDomain(),
	}
	node := GetCurrentNode(types.Address(pk.Address()), config)
	if node.PublicKey != nil {
		msg = RandomizeNodeStakeMsg(node, msg)
	}
	SendRawTx(&msg, config, signer)
}

func UnstakeNodeTransaction(config Config) {
	signer := config.GetRandomPrivateKey()
	pk := signer.PublicKey()
	msg := nodesTypes.MsgBeginUnstake{
		Address: types.Address(pk.Address()),
	}
	SendRawTx(&msg, config, signer)
}

func UnjailNodeTransaction(config Config) {
	signer := config.GetRandomPrivateKey()
	pk := signer.PublicKey()
	msg := nodesTypes.MsgUnjail{
		ValidatorAddr: types.Address(pk.Address()),
	}
	SendRawTx(&msg, config, signer)
}

func SendTx(config Config) {
	signer := config.GetRandomPrivateKey()
	pk := signer.PublicKey()

	signer2 := config.GetRandomPrivateKey()
	pk2 := signer2.PublicKey()
	msg := nodesTypes.MsgSend{
		FromAddress: types.Address(pk.Address()),
		ToAddress:   types.Address(pk2.Address()),
		Amount:      GetRandomAmount(),
	}
	SendRawTx(&msg, config, signer)
}

func StakeAppTransaction(config Config) {
	signer := config.GetRandomPrivateKey()
	pk := signer.PublicKey()
	msg := appsTypes.MsgStake{
		PubKey: pk,
		Chains: GetRandomChains(),
		Value:  GetRandomStake(),
	}
	app := GetCurrentApp(types.Address(pk.Address()), config)
	if app.PublicKey != nil {
		msg = RandomizeAppStakeMsg(app, msg)
	}
	SendRawTx(&msg, config, signer)
}

func UnstakeAppTransaction(config Config) {
	signer := config.GetRandomPrivateKey()
	pk := signer.PublicKey()
	msg := appsTypes.MsgBeginUnstake{
		Address: types.Address(pk.Address()),
	}
	SendRawTx(&msg, config, signer)
}

func SendRawTx(msg types.ProtoMsg, config Config, signer crypto.PrivateKey) {
	b := babble.NewBabbler()
	txBz, err := newTxBz(memCodec(), msg, config.ChainID, signer, int64(10000), b.Babble(), GetLegacyCodec(config))
	if err != nil {
		fmt.Println(err)
		return
	}
	pk := signer.PublicKey()
	res := rpc.SendRawTxParams{
		Addr:        types.Address(pk.Address()).String(),
		RawHexBytes: hex.EncodeToString(txBz),
	}
	j, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonMsg, _ := memCodec().MarshalJSON(msg)
	fmt.Println(string(jsonMsg))
	resp, err := QueryRPC(config, "/v1/client/rawtx", j)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}

func QueryRPC(config Config, path string, jsonArgs []byte) (string, error) {
	//cliURL := app.GlobalConfig.PocketConfig.RemoteCLIURL + ":" + app.GlobalConfig.PocketConfig.RPCPort + path
	cliURL := config.PocketEndpoint + path
	req, err := http.NewRequest("POST", cliURL, bytes.NewBuffer(jsonArgs))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	res, err := strconv.Unquote(string(bz))
	if err == nil {
		bz = []byte(res)
	}
	if resp.StatusCode == http.StatusOK {
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, bz, "", "    ")
		if err == nil {
			return prettyJSON.String(), nil
		}
		return string(bz), nil
	}
	return "", fmt.Errorf("the http status code was not okay: %d, and the status was: %s, with a response of %v", resp.StatusCode, resp.Status, string(bz))
}

func RandomizeAppStakeMsg(app appsTypes.Application, msg appsTypes.MsgStake) appsTypes.MsgStake {
	if rand.Intn(4) < 2 {
		msg.PubKey = app.PublicKey
	}
	if rand.Intn(5) < 3 {
		msg.Value = app.StakedTokens
	}
	if rand.Intn(3) < 1 {
		msg.Chains = app.Chains
	}
	return msg
}

func RandomizeNodeStakeMsg(node nodesTypes.Validator, msg nodesTypes.MsgStake) nodesTypes.MsgStake {
	if rand.Intn(4) < 2 {
		msg.PublicKey = node.PublicKey
	}
	if rand.Intn(5) < 3 {
		msg.Value = node.StakedTokens
	}
	if rand.Intn(3) < 1 {
		msg.Chains = node.Chains
	}
	if rand.Intn(3) < 1 {
		msg.ServiceUrl = node.ServiceURL
	}
	return msg
}

func GetLegacyCodec(c Config) bool {
	if c.LegacyCodecMode == 0 {
		return false
	} else if c.LegacyCodecMode == 1 {
		return true
	} else {
		return 0 == rand.Intn(2)
	}
}

func GetCurrentNode(addr types.Address, config Config) (val nodesTypes.Validator) {
	url := config.PocketEndpoint + "/v1/query/node"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"address":"` + addr.String() + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error creating request for node @ ", url)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error requesting a node @ ", url, err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error gettinr resp body for a node ", url, err.Error())
		return
	}
	fmt.Println(string(body))
	err = val.UnmarshalJSON(body)
	//err = memCodec().UnmarshalJSON(body, &val)
	if err != nil {
		fmt.Println("Error unmarshalling a node ", url, err.Error())
		return
	}
	return
}

func GetCurrentApp(addr types.Address, config Config) (app appsTypes.Application) {
	url := config.PocketEndpoint + "/v1/query/app"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"address":"` + addr.String() + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error creating request for app @ ", url)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error requesting a app @ ", url, err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error getting resp body for an app ", url, err.Error())
		return
	}
	err = app.UnmarshalJSON(body)
	//err = application.UnmarshalJSON(body, &app)
	if err != nil {
		fmt.Println("Error unmarshalling a app ", url, err.Error())
		return
	}
	return
}

func GetRandomChains() []string {
	var chains []string
	for i := 0; i < rand.Intn(15); i++ {
		chain := fmt.Sprintf("%04d", rand.Intn(99))
		chains = append(chains, chain)
	}
	return chains
}

func GetRandomStake() types.BigInt {
	b := types.NewInt(rand.Int63())
	return b
}

func GetRandomAmount() types.BigInt {
	b := types.NewInt(int64(rand.Int31()))
	return b
}

func GetRandomDomain() string {
	prefix := "https://"
	suffix := ":8081"
	babbler := babble.NewBabbler()
	return prefix + babbler.Babble() + suffix
}

func newTxBz(cdc *codec.Codec, msg types.ProtoMsg, chainID string, pk crypto.PrivateKey, fee int64, memo string, legacyCodec bool) (transactionBz []byte, err error) {
	// fees
	fees := types.NewCoins(types.NewCoin(types.DefaultStakeDenom, types.NewInt(fee)))
	// entroyp
	entropy := rand.Int63()
	signBytes, err := auth.StdSignBytes(chainID, entropy, fees, msg, memo)
	if err != nil {
		return nil, err
	}
	sig, err := pk.Sign(signBytes)
	if err != nil {
		return nil, err
	}
	s := authTypes.StdSignature{PublicKey: pk.PublicKey(), Signature: sig}
	tx := authTypes.NewTx(msg, fees, s, memo, entropy)
	if legacyCodec {
		return auth.DefaultTxEncoder(cdc)(tx, 0)
	}
	return auth.DefaultTxEncoder(cdc)(tx, -1)
}

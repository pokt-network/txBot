package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/pokt-network/pocket-core/crypto"
	"io/ioutil"
	"math/rand"
	"os"
)

type Config struct {
	PocketEndpoint   string   `json:"pocket_endpoint"`
	LegacyCodecMode  int      `json:"legacy_codec_mode"`
	TransactionTypes []string `json:"transaction_types"`
	PrivateKeys      []string `json:"private_keys"`
	ChainID          string   `json:"chain_id"`
}

func GetConfigFromFile() Config {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		panic("Error opening config file: " + err.Error())
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		panic("Error unmarshalling config file: " + err.Error())
	}
	return config
}

func (c *Config) GetPrivateKeys() (pks []crypto.PrivateKey) {
	for _, k := range c.PrivateKeys {
		pk, err := hex.DecodeString(k)
		if err != nil {
			panic("Error in parsing private key to hex: " + k + " err: " + err.Error())
		}
		privKey, err := crypto.NewPrivateKeyBz(pk)
		if err != nil {
			panic("Error in parsing private key to hex: " + k + " err: " + err.Error())
		}
		pks = append(pks, privKey)
	}
	return
}

func (c *Config) GetRandomPrivateKey() crypto.PrivateKey {
	pks := c.GetPrivateKeys()
	return pks[rand.Intn(len(pks))]
}

func (c *Config) GetRandomTransactionType() string {
	return c.TransactionTypes[rand.Intn(len(c.TransactionTypes))]
}

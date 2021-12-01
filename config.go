package main

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/pokt-network/pocket-core/crypto"
)

func (c *Config) GetRandomPrivateKey() crypto.PrivateKey {
	return c.PrivateKeys[rand.Intn(len(c.PrivateKeys))].Key
}

func (c *Config) GetRandomTransactionType() string {
	return c.TxReqTypes[rand.Intn(len(c.TxReqTypes))]
}

func (mode* RequestMode) UnmarshalJSON(data []byte) error {
	i, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	*mode = RequestMode(i)
	return nil
}

func (pk* PrivateKey) UnmarshalJSON(data []byte) error {
	stringPk := strings.Trim(string(data), "\"")
	decodedPk, err := hex.DecodeString(stringPk)
	if err != nil {
		return err
	}
	privKey, err := crypto.NewPrivateKeyBz(decodedPk)
	if err != nil {
		return err
	}
	*pk = PrivateKey{Key: privKey}
	return nil
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
package config

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/pokt-network/pocket-core/crypto"
)

func (c *Config) GetRandomNodePrivateKey() crypto.PrivateKey {
	return c.NodePrivateKeys[rand.Intn(len(c.NodePrivateKeys))].Key
}

func (c *Config) GetRandomAppPrivateKey() crypto.PrivateKey {
	return c.AppPrivateKeys[rand.Intn(len(c.AppPrivateKeys))].Key
}

func (c *Config) GetRandomTransactionType() string {
	return c.TxReqTypes[rand.Intn(len(c.TxReqTypes))]
}

func (mode *RequestMode) UnmarshalJSON(data []byte) error {
	// Check if request mode is specified as a int.
	if i, err := strconv.Atoi(string(data)); err == nil {
		requestMode := RequestMode(i)
		if _, ok := RequestModeToString[requestMode]; ok {
			*mode = requestMode
			return nil
		}
	}

	// Check if request mode is specified as a string.
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if val, ok := RequestModeToId[s]; ok {
		*mode = val
		return nil
	}

	// Could not identify request mode.
	return errors.New("Invalid request mode: " + s)
}

func (pk *PrivateKey) UnmarshalJSON(data []byte) error {
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

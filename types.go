package main

import "github.com/pokt-network/pocket-core/crypto"

type PrivateKey struct {
	Key crypto.PrivateKey
}

type Config struct {
	Mode 		     RequestMode  `json:"mode"`
	TimerDuration    int          `json:"timer_duration_in_s"`
	PocketEndpoint   string       `json:"pocket_endpoint"`
	LegacyCodecMode  int          `json:"legacy_codec_mode"`
	TransactionTypes []string     `json:"transaction_types"`
	PrivateKeys      []PrivateKey `json:"private_keys"`
	ChainID          string       `json:"chain_id"`
}

type RequestMode int

const (
	TimerMode RequestMode = iota  // 0
	RandomMode	// 1
	ManualMode  // 2
	BurstMode   // 3
)
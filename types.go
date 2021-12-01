package main

import "github.com/pokt-network/pocket-core/crypto"

type TimerModeConfig struct {
	IntervalMs uint64 `json:"interval_ms"`
}

type RandomModeConfig struct {
}

type ManualModeConfig struct {
}

type BurstModeConfig struct {
	IntervalMs uint64 `json:"interval_ms"`
	NumRequests uint64 `json:"num_requests"`
	TxReqName string `json:"tx_req_name"`
}

type ModeConfig struct {
	TimerModeConfig TimerModeConfig `json:"timer_mode_config"`
	RandomModeConfig RandomModeConfig `json:"random_mode_config"`
	ManualModeConfig ManualModeConfig `json:"manual_mode_config"`
	BurstModeConfig BurstModeConfig  `json:"burst_mode_config"`
}

type PrivateKey struct {
	Key crypto.PrivateKey
}


type RequestMode int

const (
	TimerMode RequestMode = iota  // 0
	RandomMode	// 1
	ManualMode  // 2
	BurstMode   // 3
)

type Config struct {
	ChainID          string       `json:"chain_id"`
	PocketEndpoint   string       `json:"pocket_endpoint"`
	LegacyCodecMode  int          `json:"legacy_codec_mode"`
	Mode 		     RequestMode  `json:"mode"`
	ModeConfigs      ModeConfig   `json:"mode_configs"`
	TxReqTypes       []string     `json:"tx_req_types"`
	PrivateKeys      []PrivateKey `json:"private_keys"`
}
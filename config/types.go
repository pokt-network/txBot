package config

import "github.com/pokt-network/pocket-core/crypto"

type TimerModeConfig struct {
	IntervalMs uint64 `json:"interval_ms"`
}

type RandomModeConfig struct {
}

type SelectModeConfig struct {
}

type ManualModeConfig struct {
	TxReqName string `json:"tx_req_name"`
}

type BurstModeConfig struct {
	IntervalMs  uint64 `json:"interval_ms"`
	NumRequests uint64 `json:"num_requests"`
	TxReqName   string `json:"tx_req_name"`
}

type ModeConfig struct {
	TimerModeConfig  TimerModeConfig  `json:"timer_mode_config"`
	RandomModeConfig RandomModeConfig `json:"random_mode_config"`
	SelectModeConfig ManualModeConfig `json:"select_mode_config"`
	ManualModeConfig ManualModeConfig `json:"manual_mode_config"`
	BurstModeConfig  BurstModeConfig  `json:"burst_mode_config"`
}

type PrivateKey struct {
	Key crypto.PrivateKey
}

type RequestMode int

const (
	TimerMode  RequestMode = iota // 0
	RandomMode                    // 1
	SelectMode                    // 2
	ManualMode                    // 3
	BurstMode                     // 4
)

var RequestModeToString = map[RequestMode]string{
	TimerMode:  "timerMode",
	RandomMode: "randomMode",
	SelectMode: "selectMode",
	ManualMode: "manualMode",
	BurstMode:  "burstMode",
}

var RequestModeToId = map[string]RequestMode{
	"timerMode":  TimerMode,
	"randomMode": RandomMode,
	"selectMode": SelectMode,
	"manualMode": ManualMode,
	"burstMode":  BurstMode,
}

type AppConfig struct {
	AppAddress string `json:"app_address"`
	SecretKey  string `json:"secret_key"`
	PortalId   string `json:"portal_id"`
}

type Config struct {
	ChainID            string       `json:"chain_id"`
	PocketEndpoint     string       `json:"pocket_endpoint"`
	LegacyCodecMode    int          `json:"legacy_codec_mode"`
	Mode               RequestMode  `json:"mode"`
	ModeConfigs        ModeConfig   `json:"mode_configs"`
	TxReqTypes         []string     `json:"tx_req_types"`
	NodePrivateKeys    []PrivateKey `json:"node_private_keys"`
	AppPrivateKeys     []PrivateKey `json:"app_private_keys"`
	ServicerPrivateKey PrivateKey   `json:"servicer_private_key"`
}

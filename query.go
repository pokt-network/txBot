package main

import (
	"encoding/json"
	"fmt"

	cli "github.com/pokt-network/pocket-core/app/cmd/cli"
	rpc "github.com/pokt-network/pocket-core/app/cmd/rpc"
	types "github.com/pokt-network/pocket-core/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
)

type (
	Client struct {
		config  *Config
		queries map[string]func(...interface{}) (string, interface{})
	}
)

func GetClient() *Client {
	queries := make(map[string]func(...interface{}) (string, interface{}))
	client := Client{queries: queries}
	return &client
}

func (c *Client) Init(config *Config) {
	c.config = config
	c.queries["GetTx"] = func(params ...interface{}) (string, interface{}) {
		return cli.GetTxPath, rpc.HashAndProveParams{Hash: params[0].(string), Prove: false}
	}

	c.queries["GetHeight"] = func(params ...interface{}) (string, interface{}) {
		return cli.GetHeightPath, []byte{}
	}

	c.queries["GetBalance"] = func(params ...interface{}) (string, interface{}) {
		return cli.GetBalancePath, rpc.HeightAndAddrParams{
			Height:  params[0].(int64),
			Address: params[1].(string),
		}
	}

	c.queries["GetAccount"] = func(params ...interface{}) (string, interface{}) {
		return cli.GetAccountPath, rpc.HeightAndAddrParams{Height: params[0].(int64), Address: params[1].(string)}
	}

	c.queries["GetNodes"] = func(params ...interface{}) (string, interface{}) {
		return cli.GetNodesPath, rpc.HeightAndValidatorOptsParams{
			Height: params[0].(int64),
			Opts: nodeTypes.QueryValidatorsParams{
				Blockchain:    params[1].(string),
				Page:          params[2].(int),
				Limit:         params[3].(int),
				StakingStatus: params[4].(types.StakeStatus),
				JailedStatus:  params[5].(int),
			},
		}
	}

	c.queries["GetApp"] = func(params ...interface{}) (string, interface{}) {
		return cli.GetAppPath, rpc.HeightAndAddrParams{
			Height:  params[0].(int64),
			Address: params[1].(string),
		}
	}
}

func (c *Client) Call(queryName string, paramValues ...interface{}) (interface{}, error) {
	path, params := c.queries[queryName](paramValues...)

	j, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := QueryRPC(*c.config, path, j)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

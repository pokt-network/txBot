package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	config "github.com/pokt-network/txbot/config"
	rpc "github.com/pokt-network/txbot/rpc"
	spec "github.com/pokt-network/txbot/spec"
)

func main() {
	config := config.GetConfigFromFile()
	client, clientErr := spec.NewClientWithResponses(fmt.Sprintf("%s/v1",config.PocketEndpoint))
	if clientErr != nil {
		fmt.Println(clientErr)
		os.Exit(1)
	}
	StartSendingTxOrReqs(config, client)
}

func printAllTxOrRequests(config config.Config) {
	fmt.Printf("\n%s\nEnter the number associated with the request or transaction you want to trigger.\n%s\n", strings.Repeat("~", 50), strings.Repeat("~", 50))
	for idx, txOrReq := range config.TxReqTypes {
		fmt.Printf(" %d: %s\n", idx + 1, txOrReq)
	}
	fmt.Println("")
}

func getRandomTxOrReq(c config.Config) string {
	return c.TxReqTypes[rand.Intn(len(c.TxReqTypes))]
}

func execTransOrReq(c config.Config, client *spec.ClientWithResponses, txOrReq string) {
	clientCtx := context.TODO()
	switch txOrReq {
		case "appStake":
			rpc.AppStakeTransaction(c)
		case "appUnstake":
			rpc.AppUnstakeTransaction(c)
		case "nodeSend":
			rpc.NodeSendTx(c)
		case "nodeStake":
			rpc.NodeStakeTransaction(c)
		case "nodeUnstake":
			rpc.NodeUnstakeTransaction(c)
		case "nodeUnjail":
			rpc.NodeUnjailTransaction(c)
		case "queryHeight":
			rpc.QueryHeight(c, client, clientCtx)
		case "relayEth":
			// TODO: https://infura.io/docs/ethereum/json-rpc/eth-blockNumber
			rpc.RelayEth(c, client, clientCtx)
		case "relayHmy":
			// TODO: https://docs.harmony.one/home/developers/api/methods/account-methods/hmy_getbalancebyblocknumber
			rpc.RelayHmy(c, client, clientCtx)
		default:
			break
	}
}

func StartSendingTxOrReqs(c config.Config, client *spec.ClientWithResponses) {
	for {
		switch c.Mode {
			case config.TimerMode: {
				interval := c.ModeConfigs.TimerModeConfig.IntervalMs
				fmt.Printf("%s\nSleep for %d seconds \n%s\n", strings.Repeat("~", 50), interval, strings.Repeat("~", 50))
				time.Sleep(time.Duration(interval) * time.Millisecond)
				execTransOrReq(c, client, getRandomTxOrReq(c))
			}
			case config.RandomMode: {
				fmt.Printf("%s\nPress any key to submit a random transaction or request...\n%s\n", strings.Repeat("~", 50), strings.Repeat("~", 50))
				bufio.NewScanner(os.Stdin).Scan()
				execTransOrReq(c, client, c.TxReqTypes[rand.Intn(len(c.TxReqTypes))])
			}
			case config.SelectMode: {
				printAllTxOrRequests(c)
				input := bufio.NewScanner(os.Stdin)
				input.Scan()
				idx, err := strconv.Atoi(input.Text())
				if err != nil {
					fmt.Println("Couldn't parse input...")
					continue
				}
				idx--
				if idx < 0 || idx >= len(c.TxReqTypes) {
					fmt.Println("Invalid index")
					continue
				}
				txOrReq := c.TxReqTypes[idx]
				fmt.Printf("Selection was %s. Output is below:\n%s\n", txOrReq, strings.Repeat("-", 50))
				execTransOrReq(c, client, txOrReq)
			}
			case config.ManualMode: {
				txOrReq := c.ModeConfigs.ManualModeConfig.TxReqName
				fmt.Printf("%s\nPress any key to submit %s transaction or request...\n%s\n", strings.Repeat("~", 50), txOrReq, strings.Repeat("~", 50))
				bufio.NewScanner(os.Stdin).Scan()
				execTransOrReq(c, client, txOrReq)

			}
			case config.BurstMode: {
				interval := c.ModeConfigs.BurstModeConfig.IntervalMs
				num_req := c.ModeConfigs.BurstModeConfig.NumRequests
				txOrReq := c.ModeConfigs.BurstModeConfig.TxReqName
				fmt.Printf(
					"%s\nAbout to send %d %s requests separated by %f seconds\n%s\n",
					strings.Repeat("~", 50),
					num_req,
					txOrReq,
					float64(interval) / 1e3,
					strings.Repeat("~", 50))
				for i := uint64(0); i < num_req; i++ {
					time.Sleep(time.Duration(interval) * time.Millisecond)
					execTransOrReq(c, client, txOrReq)
				}
				fmt.Println("Press any key to continue...")
				bufio.NewScanner(os.Stdin).Scan()
			}
			default:
				fmt.Printf("Mode not supported yet: %d\n", c.Mode)
				os.Exit(1)
		}
	}
}

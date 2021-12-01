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

	Spec "github.com/pokt-network/txbot/spec"
)

func main() {
	config := GetConfigFromFile()
	client, clientErr := Spec.NewClientWithResponses(fmt.Sprintf("%s/v1",config.PocketEndpoint))
	if clientErr != nil {
		fmt.Println(clientErr)
		os.Exit(1)
	}
	StartSendingTxOrReqs(config, client)
}

func printAllTxOrRequests(config Config) {
	fmt.Printf("%s\nEnter the number associated with the request or transaction you want to trigger.\n%s\n", strings.Repeat("~", 50), strings.Repeat("~", 50))
	for idx, txOrReq := range config.TxReqTypes {
		fmt.Printf(" %d: %s\n", idx + 1, txOrReq)
	}
}

func getRandomTxOrReq(c Config) string {
	return c.TxReqTypes[rand.Intn(len(c.TxReqTypes))]
}

func execTransOrReq(c Config, client *Spec.ClientWithResponses, txOrReq string) {
	clientCtx := context.TODO()
	switch txOrReq {
		case "appStake":
			AppStakeTransaction(c)
		case "appUnstake":
			AppStakeTransaction(c)
		case "nodeSend":
			NodeSendTx(c)
		case "nodeStake":
			NodeStakeTransaction(c)
		case "nodeUnstake":
			NodeUnstakeTransaction(c)
		case "nodeUnjail":
			NodeUnjailTransaction(c)
		case "queryHeight":
			QueryHeight(c, client, clientCtx)
		case "relayEth":
			// TODO: https://infura.io/docs/ethereum/json-rpc/eth-blockNumber
			// RelayEth(c, client, clientCtx)
			break
		case "relayHmy":
			// TODO: https://docs.harmony.one/home/developers/api/methods/account-methods/hmy_getbalancebyblocknumber
			break
		default:
			break
	}
}

func StartSendingTxOrReqs(c Config, client *Spec.ClientWithResponses) {
	for {
		switch c.Mode {
			case TimerMode: {
				interval := c.ModeConfigs.TimerModeConfig.IntervalMs
				fmt.Printf("%s\nSleep for %d seconds \n%s\n", strings.Repeat("~", 50), interval, strings.Repeat("~", 50))
				time.Sleep(time.Duration(interval) * time.Millisecond)
				execTransOrReq(c, client, getRandomTxOrReq(c))
			}
			case RandomMode: {
				fmt.Printf("%s\nPress any key to submit a random transaction or request...\n%s\n", strings.Repeat("~", 50), strings.Repeat("~", 50))
				bufio.NewScanner(os.Stdin).Scan()
				execTransOrReq(c, client, c.TxReqTypes[rand.Intn(len(c.TxReqTypes))])
			}
			case ManualMode: {
				printAllTxOrRequests(c)
				input := bufio.NewScanner(os.Stdin)
				input.Scan()
				idx, err := strconv.Atoi(input.Text())
				idx = idx - 1
				if err != nil || idx < 0 || idx > len(c.TxReqTypes) {
					fmt.Println("Invalid index")
					os.Exit(1)
				}
				execTransOrReq(c, client, c.TxReqTypes[idx])
			}
			case BurstMode: {
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

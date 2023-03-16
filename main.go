package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	config "github.com/pokt-network/txbot/config"
	rpc "github.com/pokt-network/txbot/rpc"
)

func main() {
	config := config.GetConfigFromFile()
	startSendingTxOrReqs(config, rpc.NewRpcContext(config))
}

func printAllTxOrRequestOptions(config config.Config) {
	fmt.Printf(
		"\n%s\nEnter the number associated with the request or transaction you want to trigger.\n%s\n",
		strings.Repeat("~", 50),
		strings.Repeat("~", 50))
	for idx, txOrReq := range config.TxReqTypes {
		fmt.Printf(" %d: %s\n", idx+1, txOrReq)
	}
	fmt.Println("")
}

func getRandomTxOrReq(c config.Config) string {
	return c.TxReqTypes[rand.Intn(len(c.TxReqTypes))]
}

func execTransOrReq(c config.Config, rpcCtx *rpc.RpcContext, txOrReq string) {
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
		height := rpc.QueryHeight(c, rpcCtx)
		fmt.Println("Current height: ", height)
	case "relayEthHeight":
		rpc.RelayEthHeight(c, rpcCtx)
	case "relayPolyHeight":
		rpc.RelayPolyHeight(c, rpcCtx)
	default:
		break
	}
}

func startSendingTxOrReqs(c config.Config, rpcCtx *rpc.RpcContext) {
	for {
		switch c.Mode {
		case config.TimerMode:
			{
				interval := c.ModeConfigs.TimerModeConfig.IntervalMs
				fmt.Printf(
					"%s\nSleep for %d seconds \n%s\n",
					strings.Repeat("~", 50),
					interval/1000,
					strings.Repeat("~", 50))
				time.Sleep(time.Duration(interval) * time.Millisecond)
				execTransOrReq(c, rpcCtx, getRandomTxOrReq(c))
			}
		case config.RandomMode:
			{
				fmt.Printf(
					"%s\nPress any key to submit a random transaction or request...\n%s\n",
					strings.Repeat("~", 50),
					strings.Repeat("~", 50))
				bufio.NewScanner(os.Stdin).Scan()
				execTransOrReq(c, rpcCtx, getRandomTxOrReq(c))
			}
		case config.SelectMode:
			{
				printAllTxOrRequestOptions(c)
				input := bufio.NewScanner(os.Stdin)
				input.Scan()
				idx, err := strconv.Atoi(input.Text())
				if err != nil {
					fmt.Println("Couldn't parse input...")
					continue
				}
				idx--
				if idx < 0 || idx >= len(c.TxReqTypes) {
					fmt.Println("Invalid index: " + strconv.Itoa(idx))
					continue
				}
				txOrReq := c.TxReqTypes[idx]
				fmt.Printf(
					"Selection was %s. Output is below:\n%s\n",
					txOrReq,
					strings.Repeat("-", 50))
				execTransOrReq(c, rpcCtx, txOrReq)
			}
		case config.ManualMode:
			{
				txOrReq := c.ModeConfigs.ManualModeConfig.TxReqName
				fmt.Printf(
					"%s\nPress any key to submit %s transaction or request...\n%s\n",
					strings.Repeat("~", 50),
					txOrReq,
					strings.Repeat("~", 50))
				bufio.NewScanner(os.Stdin).Scan()
				execTransOrReq(c, rpcCtx, txOrReq)

			}
		case config.BurstMode:
			{
				interval := c.ModeConfigs.BurstModeConfig.IntervalMs
				num_req := c.ModeConfigs.BurstModeConfig.NumRequests
				txOrReq := c.ModeConfigs.BurstModeConfig.TxReqName
				fmt.Printf(
					"%s\nAbout to send %d %s requests separated by %f seconds\n%s\n. Press any key to continue...\n",
					strings.Repeat("~", 50),
					num_req,
					txOrReq,
					float64(interval)/1e3,
					strings.Repeat("~", 50))
				bufio.NewScanner(os.Stdin).Scan()
				for i := uint64(0); i < num_req; i++ {
					time.Sleep(time.Duration(interval) * time.Millisecond)
					execTransOrReq(c, rpcCtx, txOrReq)
				}

			}
		default:
			panic("Unsupported execution mode: " + config.RequestModeToString[c.Mode])
		}
	}
}

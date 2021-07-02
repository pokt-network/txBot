package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	c := GetConfigFromFile()
	// SendRandomTx(c)
	client := GetClient()

	client.Init(&c)

	var height int64 = 0
	// hash := "FCF719CA739DCCBC281B12BC0D671AAA7A015848"
	// client.Call("GetTx", hash)

	// client.Call("GetHeight", nil)

	// address := "fcf719ca739dccbc281b12bc0d671aaa7a015848"
	// client.Call("GetBalance", height, address)

	// appAddress := "fcf719ca739dccbc281b12bc0d671aaa7a015848"
	// client.Call("GetAccount", height, appAddress)

	//client.Call("GetNodes", height, blockchain, page, limit, stakingStatus, jailingStatus)
	//client.Call("GetApp", height, appAddress)
}

func SendRandomTx(c Config) {
	for {
		switch c.TransactionTypes[rand.Intn(len(c.TransactionTypes))] {
		case "nodeStake":
			StakeNodeTransaction(c)
		case "appStake":
			StakeAppTransaction(c)
		case "nodeUnstake":
			UnstakeNodeTransaction(c)
		case "appUnstake":
			UnstakeAppTransaction(c)
		case "nodeUnajil":
			UnjailNodeTransaction(c)
		case "send":
			SendTx(c)
		default:
			continue
		}
		if c.TimerMode {
			t := time.NewTicker(time.Duration(c.TimerDuration) * time.Second)
			select {
			case <-t.C:
				continue
			}
		}
		fmt.Println("Press return to submit another tx")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
	}
}

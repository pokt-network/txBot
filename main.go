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
	SendRandomTx(c)
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

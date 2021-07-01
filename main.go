package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
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
		fmt.Println("Press return to submit another tx")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
	}
}

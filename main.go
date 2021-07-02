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
			case "nodeUnjail":
				UnjailNodeTransaction(c)
			case "send":
				SendTx(c)
			case "height":
				GetHeight(c)
			default:
				continue
		}

		switch c.Mode {
			case TimerMode: {
				time.Sleep(time.Duration(c.TimerDuration) * time.Second)
				continue
			}
			case RandomMode: {
				fmt.Println("Press return to run another random submit another tx")
				input := bufio.NewScanner(os.Stdin)
				input.Scan()
			}
			case ManualMode: fallthrough
			case BurstMode: fallthrough
			default:
				fmt.Printf("Mode not supported yet: %d\n", c.Mode)
				os.Exit(1)
		}
	}
}

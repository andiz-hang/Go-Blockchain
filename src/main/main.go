package main

import (
	"fmt"
	"blockchain"
)

func main() {
	b := blockchain.Initial(16)
	b.Mine(10)

	fmt.Println(b.ValidHash())
}
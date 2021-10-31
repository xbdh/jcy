package main

import (
	"fmt"
	"github.com/xbdh/jcy/database"
	"os"
	"time"
)



func main() {
	state, err := database.NewStateFromDisk()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer state.Close()

	block0 := database.NewBlock(
		database.Hash{},
		uint64(time.Now().Unix()),
		[]database.TX{
			database.NewTX("andrej", "andrej", 3, ""),
			database.NewTX("andrej", "andrej", 700, "reward"),
		},
	)

	state.AddBlock(block0)
	block0hash, _ := state.Persiet()

	block1 := database.NewBlock(
		block0hash,
		uint64(time.Now().Unix()),
		[]database.TX{
			database.NewTX("andrej", "babayaga", 2000, ""),
			database.NewTX("andrej", "andrej", 100, "reward"),
			database.NewTX("babayaga", "andrej", 1, ""),
			database.NewTX("babayaga", "caesar", 1000, ""),
			database.NewTX("babayaga", "andrej", 50, ""),
			database.NewTX("andrej", "andrej", 600, "reward"),
		},
	)

	state.AddBlock(block1)
	state.Persiet()
}
package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xbdh/jcy/database"
	"os"
	"time"
)

var migrateCmd = func() *cobra.Command {
	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrates the blockchain database according to new business rules.",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, _ := cmd.Flags().GetString(flagDataDir)

			state, err := database.NewStateFromDisk(dataDir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			block0 := database.NewBlock(
				database.Hash{},
				0,
				uint64(time.Now().Unix()),
				[]database.Tx{
					database.NewTx("andrej", "andrej", 3, ""),
					database.NewTx("andrej", "andrej", 700, "reward"),
				},
			)

			block0hash,err:=state.AddBlock(block0)
			//block0hash, err := state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			block1 := database.NewBlock(
				block0hash,
				1,
				uint64(time.Now().Unix()),
				[]database.Tx{
					database.NewTx("andrej", "babayaga", 2000, ""),
					database.NewTx("andrej", "andrej", 100, "reward"),
					database.NewTx("babayaga", "andrej", 1, ""),
					database.NewTx("babayaga", "caesar", 1000, ""),
					database.NewTx("babayaga", "andrej", 50, ""),
					database.NewTx("andrej", "andrej", 600, "reward"),
				},
			)

			block1hash, err :=state.AddBlock(block1)
			//block1hash, err := state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			block2 := database.NewBlock(
				block1hash,
				2,
				uint64(time.Now().Unix()),
				[]database.Tx{
					database.NewTx("andrej", "andrej", 24700, "reward"),
				},
			)

			block2hash, err :=state.AddBlock(block2)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Println("block 2 hash:",block2hash.Hex())
		},
	}

	addDefaultRequiredFlags(migrateCmd)

	return migrateCmd
}
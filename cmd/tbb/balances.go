package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xbdh/jcy/database"

	"os"
)

func balancesListCmd() *cobra.Command{
	var balancesListCmd =&cobra.Command{
		Use:   "list",
		Short: "Lists all balances..",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir,_:=cmd.Flags().GetString(flagDataDir)

			state, err := database.NewStateFromDisk(dataDir)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			fmt.Println("Account balances:")
			fmt.Println(".................")

			for account, balance := range state.Balances {
				fmt.Printf("%s: %d\n", account, balance)
			}
		},
	}
	addDefaultRequireFlags(balancesListCmd)
	return balancesListCmd
}



func balancesCmd()*cobra.Command  {
	var balanceCmd =&cobra.Command{
		Use: "balances",
		Short: "interact with balances (list...)",
		PreRunE: func(cmd *cobra.Command, args []string) error{

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

		},

	}

	balanceCmd.AddCommand(balancesListCmd())
	return balanceCmd
}




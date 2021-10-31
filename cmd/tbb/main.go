package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const flagDataDir = "datadir"

func main() {
	var tbbCmd =cobra.Command{
		Use: "tbb",
		Short: "the blockchain jcy CLI",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	tbbCmd.AddCommand(versionCmd)
	tbbCmd.AddCommand(balancesCmd())
	//tbbCmd.AddCommand(txCmd())
	tbbCmd.AddCommand(runCmd())

	err:=tbbCmd.Execute()
	if err != nil {
		fmt.Fprint(os.Stderr,err)
		os.Exit(1)
	}
}

func addDefaultRequireFlags(cmd *cobra.Command)  {
	cmd.Flags().String(flagDataDir,"","absoulte path to the node data dir where the db will stored")
	cmd.MarkFlagRequired(flagDataDir)
}

func incorrectUsageErr()error  {
	return fmt.Errorf("incorrect usage")
}
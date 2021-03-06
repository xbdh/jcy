package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const flagDataDir = "datadir"
const flagPort = "port"
const flagIp = "ip"

func main() {
	var tbbCmd =cobra.Command{
		Use: "tbb",
		Short: "the blockchain jcy CLI",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	tbbCmd.AddCommand(versionCmd)
	tbbCmd.AddCommand(balancesCmd())
	tbbCmd.AddCommand(runCmd())
	tbbCmd.AddCommand(migrateCmd())

	err:=tbbCmd.Execute()
	if err != nil {
		fmt.Fprint(os.Stderr,err)
		os.Exit(1)
	}
}

func addDefaultRequiredFlags(cmd *cobra.Command)  {
	cmd.Flags().String(flagDataDir,"","absoulte path to the node data dir where the db will stored")
	cmd.MarkFlagRequired(flagDataDir)
}

func incorrectUsageErr()error  {
	return fmt.Errorf("incorrect usage")
}
package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xbdh/jcy/node"
	"os"
)

func runCmd()*cobra.Command  {
	var runCmd =&cobra.Command{
		Use: "run",
		Short: "Launches the tbb node and its http api",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir,_:=cmd.Flags().GetString(flagDataDir)

			fmt.Println("lanuching tbb node and its http api ....")

			err:=node.Run(dataDir)
			if err != nil {
				//fmt.Println("fuck")
				fmt.Fprintln(os.Stderr,err)
				os.Exit(1)
			}
		},
	}
	addDefaultRequireFlags(runCmd)

	return runCmd
}

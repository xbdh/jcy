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
			port,_ := cmd.Flags().GetUint64(flagPort)

			fmt.Println("lanuching tbb node and its http api ....")

			bootstrap:= node.NewPeerNode(
				"1234",
				8080,
				true,
				true,
				)

			n:=node.New(dataDir,port,bootstrap)
			err:=n.Run()
			if err != nil {
				//fmt.Println("fuck")
				fmt.Fprintln(os.Stderr,err)
				os.Exit(1)
			}
		},
	}

	addDefaultRequiredFlags(runCmd)
	runCmd.Flags().Uint64(flagPort,node.DefaultHTTPPort,"exposed HTTP port for communication with peers")

	return runCmd
}

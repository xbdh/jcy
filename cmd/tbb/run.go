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
			ip,_:=cmd.Flags().GetString(flagIp)
			fmt.Println("lanuching tbb node and its http api ....")

			bootstrap:= node.NewPeerNode(
				"127.0.0.1",
				8080,
				true,
				false,
				)

			n:=node.New(dataDir,ip, port,bootstrap)
			err:=n.Run()
			if err != nil {
				//fmt.Println("fuck")
				fmt.Fprintln(os.Stderr,err)
				os.Exit(1)
			}
		},
	}

	addDefaultRequiredFlags(runCmd)
	runCmd.Flags().String(flagIp,node.DefaultIP,"exposed ip  for communication with peers")
	runCmd.Flags().Uint64(flagPort,node.DefaultHTTPPort,"exposed http port for communication with peers")

	return runCmd
}

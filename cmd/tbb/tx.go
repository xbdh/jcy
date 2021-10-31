package main

//import (
//	"fmt"
//	"github.com/spf13/cobra"
//	"github.com/xbdh/jcy/database"
//	"os"
//)
//
//const flagFrom = "from"
//const flagTo = "to"
//const flagValue ="value"
//const flagData ="data"
//
//
//func txAddcmd() *cobra.Command{
//	var txadd =&cobra.Command{
//		Use: "add",
//		Short: "Add new tx to database",
//		Run: func(cmd *cobra.Command, args []string) {
//			from,_:=cmd.Flags().GetString(flagFrom)
//			to ,_ :=cmd.Flags().GetString(flagTo)
//			value,_:=cmd.Flags().GetUint(flagValue)
//			data,_:=cmd.Flags().GetString(flagData)
//
//			accountFrom :=database.NewAccount(from)
//			accountTo :=database.NewAccount(to)
//			tx:=database.NewTX(accountFrom,accountTo,value,data)
//
//			state,err:=database.NewStateFromDisk()
//			fmt.Println(state.LatestBlockHash())
//			if err != nil {
//				fmt.Fprintln(os.Stderr,err)
//				os.Exit(1)
//			}
//			defer state.Close()
//
//			err=state.Add(tx)
//			if err != nil {
//				fmt.Fprintln(os.Stderr,err)
//				os.Exit(1)
//			}
//
//			_,err=state.Persiet()
//			if err != nil {
//				fmt.Fprintln(os.Stderr,err)
//				os.Exit(1)
//			}
//
//			fmt.Println("New Tx successfully added to ledger")
//		},
//	}
//
//	txadd.Flags().String(flagFrom,"","From what account to send tokens")
//	txadd.MarkFlagRequired(flagFrom)
//
//	txadd.Flags().String(flagTo,"","To what account to send tokens")
//	txadd.MarkFlagRequired(flagTo)
//
//	txadd.Flags().Uint(flagValue,0,"how many tokens to send")
//	txadd.MarkFlagRequired(flagValue)
//
//	txadd.Flags().String(flagData,"","possible values: reward")
//	return txadd
//
//}
//
//func txCmd()*cobra.Command  {
//
//
//	var txCmd=&cobra.Command{
//		Use: "tx",
//		Short: "Interact with txs (addd...)",
//		PreRunE: func(cmd *cobra.Command, args []string) error {
//			return nil
//		},
//		Run: func(cmd *cobra.Command, args []string) {
//
//		},
//
//	}
//	txCmd.AddCommand(txAddcmd())
//
//	return txCmd
//}
package node

import (
	"github.com/xbdh/jcy/database"
	"net/http"
)

type ErrRes struct {
	Error string `json:"error"`
}

type BalancesRes struct {
	Hash database.Hash     `json:"block_hash"`
	Balances map[database.Account]uint `json:"balances"`
}

type StatusRes struct {
	Hash database.Hash`json:"block_hash"`
	Number uint64     `json:"block_number"`

	KnownPeers map[string]PeerNode `json:"peers_known"`
}
type TxAddReq struct {
	From string `json:"from"`
	To string`json:"to"`
	Value uint `json:"value"`
	Data string `json:"data"`
}

type TxAddRes struct {
	Hash database.Hash `json:"block_hash"`
}
type SyncRes struct {
	Blocks []database.Block `json:"blocks"`
}

func listBalances(writer http.ResponseWriter, request *http.Request,state *database.State)  {
	writeRes(writer,BalancesRes{
		Hash:     state.LatestBlockHash(),
		Balances: state.Balances,
	})
}

func txAdd(writer http.ResponseWriter, request *http.Request,state *database.State)  {
	txreq:=TxAddReq{}
	err:=readReq(request,&txreq)
	if err != nil {
		writeErrRes(writer,err)
		return
	}

	tx:=database.NewTx(database.NewAccount(txreq.From),database.NewAccount(txreq.To),txreq.Value,txreq.Data)

	err =state.AddTx(tx)
	if err != nil {
		writeErrRes(writer,err)
		return
	}

	hash ,err:= state.Persist()
	if err != nil {
		writeErrRes(writer,err)
		return
	}

	writeRes(writer,TxAddRes{Hash: hash})
}

func statusHandler(writer http.ResponseWriter, request *http.Request, n *Node)  {
	res:=StatusRes{
		Hash:   n.state.LatestBlockHash(),
		Number: n.state.LatestBlock().Header.Number,
		KnownPeers: n.knownPeers,
	}
	writeRes(writer,res)
}

func syncaHandler(writer http.ResponseWriter, request *http.Request,dataDir string){
	reqHash := request.URL.Query().Get(endpointSyncQueryKeyFromBlock)

	hash:= database.Hash{}

	err:= hash.UnmarshalText([]byte(reqHash))
	if err != nil {
		writeErrRes(writer,err)
		return
	}

	blocks,err:=database.GetBlocksAfer(hash,dataDir)
	if err != nil {
		writeErrRes(writer,err)
		return
	}

	writeRes(writer,SyncRes{
		Blocks: blocks,
	})
}

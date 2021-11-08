package node

import (
	"github.com/xbdh/jcy/database"
	"net/http"
	"strconv"
	"fmt"
	"time"
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

type AddPeerRes struct {
	Success bool `json:"success"`
	Error string `json:"error"`
}
// 账户余额信息
// 包括最后LatestBlockHash，Balances
func listBalancesHandler(writer http.ResponseWriter, request *http.Request,state *database.State)  {
	writeRes(writer,BalancesRes{
		Hash:     state.LatestBlockHash(),
		Balances: state.Balances,
	})
}

// 新增交易
func txAddHandler(writer http.ResponseWriter, request *http.Request,state *database.State)  {
	// req-> TxAddReq{} -> database.Tx{}
	txreq:=TxAddReq{}
	err:=readReq(request,&txreq)
	if err != nil {
		writeErrRes(writer,err)
		return
	}

	tx:=database.NewTx(database.NewAccount(txreq.From),database.NewAccount(txreq.To),txreq.Value,txreq.Data)

	block:=database.NewBlock(
		state.LatestBlockHash(),
		state.NextBlockNumber(),
		uint64(time.Now().Unix()),

		[]database.Tx{tx},  //这里只有一个交易
		 )
	hash,err:=state.AddBlock(block)
	if err != nil {
		writeErrRes(writer,err)
		return
	}

	writeRes(writer,TxAddRes{Hash: hash})
}

// 节点信息
// 包含 最后块hashLatest :BlockHash， 块的高度：Number，此节点已知的peer信息： KnownPeers
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

func addPeerHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	peerIP := r.URL.Query().Get(endpointAddPeerQueryKeyIP)
	peerPortRaw := r.URL.Query().Get(endpointAddPeerQueryKeyPort)

	peerPort, err := strconv.ParseUint(peerPortRaw, 10, 32)
	if err != nil {
		writeRes(w, AddPeerRes{false, err.Error()})
		return
	}

	peer := NewPeerNode(peerIP, peerPort, false, true)

	node.AddPeer(peer)

	fmt.Printf("Peer '%s' was added into KnownPeers\n", peer.TcpAddress())

	writeRes(w, AddPeerRes{true, ""})
}

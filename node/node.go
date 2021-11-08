package node

import (
	"context"
	"fmt"
	"github.com/xbdh/jcy/database"
	"net/http"
)

const DefaultIP ="127.0.0.1"
const DefaultHTTPPort = 8080
const endpointStatus = "/node/status"
const endpointSync = "/node/sync"
const endpointSyncQueryKeyFromBlock = "fromBlock"

const endpointAddPeer = "/node/peer"
const endpointAddPeerQueryKeyIP = "ip"
const endpointAddPeerQueryKeyPort = "port"

type Node struct {
	dataDir string
	ip string
	port  uint64
	
	state *database.State
	
	knownPeers map[string]PeerNode
}

type PeerNode struct {
	IP string             `json:"ip"`
	Port uint64           `json:"port"`
	IsBootstrap bool      `json:"is_bootstrap"`

	connected bool `json:"connected"`

}

func NewPeerNode(IP string, port uint64, isBootstrap bool, connected bool) PeerNode {
	return PeerNode{IP: IP, Port: port, IsBootstrap: isBootstrap, connected: connected}
}
func (pn PeerNode) TcpAddress()string  {
	return fmt.Sprintf("%s:%d",pn.IP,pn.Port)
}


func New(dataDir string,ip string, port uint64, bootstrap PeerNode) *Node {
	knownPeers:=make(map[string]PeerNode)
	knownPeers[bootstrap.TcpAddress()]=bootstrap
	return &Node{
		dataDir: dataDir,
		ip: ip,
		port:    port,
		knownPeers: knownPeers,
	}
}



func (n* Node)Run() error {
	ctx:=context.Background()

	fmt.Printf("listening on %s:%d\n",n.ip,n.port)

	state,err:=database.NewStateFromDisk(n.dataDir)
	if err != nil {
		//fmt.Println("state wrong")
		return err
	}

	defer state.Close()

	n.state=state

	go n.sync(ctx)

	// 获取余额
	http.HandleFunc("/balances/list",func(writer http.ResponseWriter, request *http.Request) {
		listBalancesHandler(writer,request,state)
	})

	// 增添交易
	http.HandleFunc("/tx/add",func(writer http.ResponseWriter, request *http.Request) {
		txAddHandler(writer,request,state)
	})

	// 节点状态
	http.HandleFunc(endpointStatus,func(writer http.ResponseWriter, request *http.Request) {
		statusHandler(writer,request,n)
	})

	// 同步节点和块信息
	http.HandleFunc(endpointSync,func(writer http.ResponseWriter, request *http.Request) {
		syncaHandler(writer,request,n.dataDir)
	})

	// 添加从peer节点获取的peer节点信息
	http.HandleFunc(endpointAddPeer,func(writer http.ResponseWriter, request *http.Request) {
		addPeerHandler(writer,request,n)
	})


	err= http.ListenAndServe(fmt.Sprintf(":%d" ,n.port),nil)

	return err
}

func (n *Node) AddPeer(peer PeerNode) {
	n.knownPeers[peer.TcpAddress()] = peer
}

func (n *Node) RemovePeer(peer PeerNode) {
	delete(n.knownPeers, peer.TcpAddress())
}

func (n *Node) IsKnownPeer(peer PeerNode) bool {
	if peer.IP == n.ip && peer.Port == n.port {
		return true
	}

	_, isKnownPeer := n.knownPeers[peer.TcpAddress()]

	return isKnownPeer
}

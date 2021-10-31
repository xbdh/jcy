package node

import (
	"fmt"
	"github.com/xbdh/jcy/database"
	"net/http"
)

const DefaultHTTPPort = 8080

type Node struct {
	dataDir string
	port  uint64
	
	state *database.State
	
	knownPeers []PeerNode
}

type PeerNode struct {
	IP string             `json:"ip"`
	Port uint64           `json:"port"`
	IsBootstrap bool      `json:"is_bootstrap"`
	IsActive bool         `json:"is_active"`
}

func NewPeerNode(IP string, port uint64, isBootstrap bool, isActive bool) PeerNode {
	return PeerNode{IP: IP, Port: port, IsBootstrap: isBootstrap, IsActive: isActive}
}



func New(dataDir string, port uint64, bootstrap PeerNode) *Node {
	return &Node{
		dataDir: dataDir,
		port: port,
		knownPeers: []PeerNode{bootstrap}}
}


func (n* Node)Run() error {
	fmt.Printf("listening on http port %d",n.port)

	state,err:=database.NewStateFromDisk(n.dataDir)
	if err != nil {
		//fmt.Println("state wrong")
		return err
	}
	defer state.Close()

	n.state=state
	
	http.HandleFunc("/balances/list",func(writer http.ResponseWriter, request *http.Request) {
		listBalances(writer,request,state)
	})


	http.HandleFunc("/tx/add",func(writer http.ResponseWriter, request *http.Request) {
		txAdd(writer,request,state)
	})

	http.HandleFunc("/node/status",func(writer http.ResponseWriter, request *http.Request) {
		statusHandler(writer,request,n)
	})


	err= http.ListenAndServe(fmt.Sprintf(":%d" ,n.port),nil)

	return err
}


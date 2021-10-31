package node

import (
	"context"
	"fmt"
	"github.com/xbdh/jcy/database"
	"net/http"
	"time"
)

const DefaultHTTPPort = 8080
const endpointStatus = "/node/status"
const endpointSync = "/node/sync"
const endpointSyncQueryKeyFromBlock = "fromBlock"

type Node struct {
	dataDir string
	port  uint64
	
	state *database.State
	
	knownPeers map[string]PeerNode
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
func (pn PeerNode) TcpAddress()string  {
	return fmt.Sprintf("%s:%d",pn.IP,pn.Port)
}


func New(dataDir string, port uint64, bootstrap PeerNode) *Node {
	knownPeers:=make(map[string]PeerNode)
	knownPeers[bootstrap.TcpAddress()]=bootstrap
	return &Node{
		dataDir: dataDir,
		port:    port,
		knownPeers: knownPeers,
	}
}



func (n* Node)Run() error {
	ctx:=context.Background()

	fmt.Printf("listening on http port %d\n",n.port)

	state,err:=database.NewStateFromDisk(n.dataDir)
	if err != nil {
		//fmt.Println("state wrong")
		return err
	}

	go n.sync(ctx)

	defer state.Close()

	n.state=state
	
	http.HandleFunc("/balances/list",func(writer http.ResponseWriter, request *http.Request) {
		listBalances(writer,request,state)
	})


	http.HandleFunc("/tx/add",func(writer http.ResponseWriter, request *http.Request) {
		txAdd(writer,request,state)
	})

	http.HandleFunc(endpointStatus,func(writer http.ResponseWriter, request *http.Request) {
		statusHandler(writer,request,n)
	})

	http.HandleFunc(endpointSync,func(writer http.ResponseWriter, request *http.Request) {
		syncaHandler(writer,request,n.dataDir)
	})


	err= http.ListenAndServe(fmt.Sprintf(":%d" ,n.port),nil)

	return err
}

func (n *Node) sync(ctx context.Context) error {
	ticker:= time.NewTicker(45*time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Searching for new peers and blocks...")

			n.fetchNewBlocksAndPeers()
		case <-ctx.Done():
			ticker.Stop()

		}
	}
}

func (n *Node) fetchNewBlocksAndPeers()  {
	for _,knownPeer:=range n.knownPeers{
		status,err:=queryPeerStatus(knownPeer)
		if err != nil {
			fmt.Println("ERROR :",err)
			continue
		}

		localBlockNumber:=n.state.LatestBlock().Header.Number
		if localBlockNumber < status.Number{
			newBlockCount := status.Number-localBlockNumber

			fmt.Printf("Found %d new block from peer %s\n",newBlockCount,knownPeer.IP)
		}

		for _, mayNewPeer:=range status.KnownPeers{
			_,isKnowPeer:=n.knownPeers[mayNewPeer.TcpAddress()]
			if !isKnowPeer{
				fmt.Sprintf("Found new Peer %s\n",mayNewPeer.TcpAddress())

				n.knownPeers[mayNewPeer.TcpAddress()]=mayNewPeer
			}
		}

	}
}

func queryPeerStatus(peer PeerNode)(StatusRes,error)  {
	url:=fmt.Sprintf("http://%s/%s",peer.TcpAddress(),endpointStatus)
	res,err:=http.Get(url)
	if err != nil {
		return StatusRes{}, err
	}

	statusRes:=StatusRes{}
	err = readRes(res,&statusRes)
	if err != nil {
		return StatusRes{}, err
	}
	return statusRes ,nil
}

// curl -X GET http://127.0.0.1:8081/node/sync?fromBlock=ec31eeac22949972770ad83fd9501849f79813b539bdb343af528d37ff9b426a
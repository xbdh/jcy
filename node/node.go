package node

import (
	"encoding/json"
	"fmt"
	"github.com/xbdh/jcy/database"
	"io/ioutil"
	"net/http"
)

const httpPort = 8080

type ErrRes struct {
	Error string `json:"error"`
}

type BalancesRes struct {
	Hash database.Hash     `json:"block_hash"`
	Balances map[database.Account]uint `json:"balances"`
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

func Run(dataDir string) error {
	fmt.Printf("listening on http port %d",httpPort)

	state,err:=database.NewStateFromDisk(dataDir)
	if err != nil {
		//fmt.Println("state wrong")
		return err
	}
	defer state.Close()


	http.HandleFunc("/balances/list",func(writer http.ResponseWriter, request *http.Request) {
		listBalances(writer,request,state)
	})


	http.HandleFunc("/tx/add",func(writer http.ResponseWriter, request *http.Request) {
		txAdd(writer,request,state)
	})

	err= http.ListenAndServe(fmt.Sprintf(":%d" ,httpPort),nil)

	return err
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

	tx:=database.NewTX(database.NewAccount(txreq.From),database.NewAccount(txreq.To),txreq.Value,txreq.Data)

	err =state.Add(tx)
	if err != nil {
		writeErrRes(writer,err)
		return
	}

	hash ,err:= state.Persiet()
	if err != nil {
		writeErrRes(writer,err)
		return
	}

	writeRes(writer,TxAddRes{Hash: hash})
}

func writeErrRes(w http.ResponseWriter, err error) {
	jsonErrRes, _ := json.Marshal(ErrRes{err.Error()})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonErrRes)
}

func writeRes(w http.ResponseWriter, content interface{}) {
	contentJson, err := json.Marshal(content)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(contentJson)
}

func readReq(r *http.Request,reqBody interface{})error  {
	reqBodyJson,err:= ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	err= json.Unmarshal(reqBodyJson,reqBody)
	if err!=nil{
		return fmt.Errorf("unable to unmarshal request body %s",err.Error())
	}
	return nil
}
package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)


type State struct {
	Balances map[Account]uint
	txMempool []TX

	dbFile *os.File

	latestBlockHash Hash
}

func NewStateFromDisk(dataDir string) (*State ,error){
	//
	//cwd ,err :=os.Getwd()
	//if err!=nil{
	//	return nil,err
	//}

	err:=initDataDirIfNotExists(dataDir)
	if err != nil {
		//fmt.Println("intstat wrong",err)
		return nil, err
	}

	//genesisFilePath:=filepath.Join(cwd,"database","genesis.json")

	gen,err:=loadGenesis(getGenesisJsonFilePath(dataDir))
	if err != nil {
		//fmt.Println("lloadt wrong",err.Error())
		return nil, err
	}
	balances:=make(map[Account]uint)
	for account,balance:=range gen.Balances{
		balances[account]=balance
	}

	//txDbFilePath:=filepath.Join(cwd,"database","tx.db")
	//
	//f,err:=os.OpenFile(txDbFilePath,os.O_APPEND|os.O_RDWR,0600)
	//if err != nil {
	//	return nil, err
	//}

	//blockFilePath:= filepath.Join(cwd,"database","block.db")
	f,err:=os.OpenFile(getBlocksDbFilePath(dataDir),os.O_APPEND|os.O_RDWR,0600)
	if err != nil {
		return nil, err
	}


	scanner:=bufio.NewScanner(f)

	state:=&State{balances,make([]TX,0),f,Hash{}}

	for scanner.Scan(){
		if err:=scanner.Err();err!=nil{
			return nil, err
		}

		blockfsjson:=scanner.Bytes()


		var blockfs BlockFs

		err= json.Unmarshal(blockfsjson,&blockfs)
		if err != nil {
			return nil, err
		}
		if err:=state.applyBlock(blockfs.Value);err!=nil{
			return nil,err
		}
		//fmt.Println(blockfsjson)
		//fmt.Println("====")
		//fmt.Println(blockfs)
		////fmt.Println(state.latestBlockHash)
		//fmt.Println("-----")
		//fmt.Println(scanner.Text())
		state.latestBlockHash =blockfs.Key

	}

	return state,nil
}

func (s *State) LatestBlockHash()  Hash{
	return s.latestBlockHash
}
func (s *State) Add(tx TX)error  {

	if err:=s.apply(tx);err!=nil{
		return err
	}
	s.txMempool=append(s.txMempool,tx)
	return nil
}

func (s *State) AddBlock(b Block) error  {
	   for _,tx:=range b.Txs{
	   	if err:=s.Add(tx);err!=nil{
	   		return err
		}
	   }
	return nil
}

func (s *State) applyBlock(b Block)error  {
	   for _,tx:=range b.Txs{
	   	if err:=s.apply(tx);err!=nil{
	   		       return err
		   }
	   }
	return nil
}
func (s *State) apply(tx TX)error  {
	if tx.IsReward(){
		s.Balances[tx.To]+=tx.Value
		return nil
	}
	if tx.Value>s.Balances[tx.From]{
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From]-=tx.Value
	s.Balances[tx.To] += tx.Value
	return nil
}
func (s *State) Persiet()(Hash,error)  {

	block:=NewBlock(s.latestBlockHash,uint64(time.Now().Unix()),s.txMempool)

	blockhash,err:=block.Hash()
	if err != nil {
		return Hash{},err
	}
	blockfs:=BlockFs{
		Key:   blockhash,
		Value: block,
	}

	blockfsjson,err:=json.Marshal(blockfs)
	if err != nil {
		return Hash{},err
	}
	fmt.Println("Persisting new block to disk ")
	fmt.Printf("%s\n",blockfsjson)

	if _,err:= s.dbFile.Write(append(blockfsjson,'\n'));err!=nil{
		return Hash{},err
	}
	s.txMempool=[]TX{}

	s.latestBlockHash=blockhash

	return s.latestBlockHash,nil

	//mempool:=make([]TX,len(s.txMempool))
	//
	//copy(mempool,s.txMempool)
	//
	//for i:=0;i<len(mempool);i++{
	//	txjson,err:=json.Marshal(mempool[i])
	//
	//	if err!=nil{
	//		return err
	//	}
	//
	//	if _,err:= s.dbFile.Write(append(txjson,'\n'));err!=nil{
	//		return err
	//	}
	//
	//	// 已经存储的交易要删除
	//	s.txMempool=s.txMempool[1:]
	//}
	//return nil
}

func (s *State) Close()  {
	s.dbFile.Close()
}


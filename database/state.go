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
	txMempool []Tx

	dbFile *os.File

	latestBlock Block
	latestBlockHash Hash
}

func NewStateFromDisk(dataDir string) (*State ,error){

	err:=initDataDirIfNotExists(dataDir)
	if err != nil {

		return nil, err
	}


	gen,err:=loadGenesis(getGenesisJsonFilePath(dataDir))
	if err != nil {

		return nil, err
	}
	balances:=make(map[Account]uint)
	for account,balance:=range gen.Balances{
		balances[account]=balance
	}

	f,err:=os.OpenFile(getBlocksDbFilePath(dataDir),os.O_APPEND|os.O_RDWR,0600)
	if err != nil {
		return nil, err
	}


	scanner:=bufio.NewScanner(f)

	state:=&State{balances,make([]Tx,0),f,Block{},Hash{}}

	for scanner.Scan(){
		if err:=scanner.Err();err!=nil{
			return nil, err
		}

		blockfsjson:=scanner.Bytes()

		if len(blockfsjson)==0{
			break
		}

		var blockfs BlockFs

		err= json.Unmarshal(blockfsjson,&blockfs)
		if err != nil {
			return nil, err
		}
		if err:=state.applyBlock(blockfs.Value);err!=nil{
			return nil,err
		}

		state.latestBlock=blockfs.Value
		state.latestBlockHash =blockfs.Key

	}

	return state,nil
}

func (s *State) LatestBlockHash()  Hash{
	return s.latestBlockHash
}
func (s* State)LatestBlock() Block {
	return s.latestBlock
}


func (s *State) AddTx(tx Tx)error  {

	if err:=s.apply(tx);err!=nil{
		return err
	}
	s.txMempool=append(s.txMempool,tx)
	return nil
}

func (s *State) AddBlock(b Block) error  {
	   for _,tx:=range b.Txs{
	   	if err:=s.AddTx(tx);err!=nil{
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
func (s *State) apply(tx Tx)error  {
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
func (s *State) Persist()(Hash,error)  {

	latestBlockHash:=s.latestBlockHash

	block:=NewBlock(latestBlockHash,s.latestBlock.Header.Number+1,uint64(time.Now().Unix()),s.txMempool)

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

	//s.latestBlockHash=latestBlockHash  //? 存疑问
	s.latestBlock=block
	s.latestBlockHash=blockhash
	s.txMempool=[]Tx{}

	return s.latestBlockHash,nil  // 存疑问


}

func (s *State) Close()  {
	s.dbFile.Close()
}


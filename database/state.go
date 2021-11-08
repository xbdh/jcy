package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)


type State struct {
	Balances map[Account]uint
	txMempool []Tx

	dbFile *os.File

	latestBlock Block
	latestBlockHash Hash

	hasGenesisBlock bool
}

// 从Genesis.json 读取的balances信息 和 block.db 读取的tx 交易信息
// 最终形成Balances map[Account]uint 账户和余额信息
// txMempool 此时为空
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

	// 如果没有block.db
	// latestBlock 为 Block{}
	// latestBlockHash 为 Hash{} 0x000000
	// hasGenesisBlock 为 false  没有创世区块，block.db为空
	state:=&State{balances,make([]Tx,0),f,Block{},Hash{},false}

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
		if err:=appyTxs(blockfs.Value.Txs,state);err!=nil{
			return nil,err
		}

		state.latestBlock=blockfs.Value
		state.latestBlockHash =blockfs.Key
		state.hasGenesisBlock=true
	}

	return state,nil
}

func (s *State) LatestBlockHash()  Hash{
	return s.latestBlockHash
}
func (s* State)LatestBlock() Block {
	return s.latestBlock
}

// 返回下一区块的高度，
// 判断 block.db是否为空，为空 为0，不为空+1
func (s *State)NextBlockNumber()uint64  {
	if s.hasGenesisBlock==false{
		return uint64(0)
	}
	return s.latestBlock.Header.Number+1
}



func (s *State) Close()  {
	s.dbFile.Close()
}

// 深度copy ，用来验证是否为恶意交易
func (s *State) copy() State  {
	c:=State{}
	c.latestBlockHash=s.latestBlockHash
	c.latestBlock=s.latestBlock
	c.txMempool=make([]Tx,0)
	c.Balances=make(map[Account]uint)
	c.hasGenesisBlock=s.hasGenesisBlock

	for account ,balance:=range s.Balances{
		c.Balances[account] =balance
	}

	for _,tx :=range s.txMempool{
		c.txMempool=append(c.txMempool,tx)
	}
	return c
}

// 持久化block，更新 state 的 balances 信息
func (s *State) AddBlock(b Block)  (Hash,error){
	pendingState :=s.copy()
	err:=applyBlock(b,pendingState)
	if err != nil {
		return Hash{}, err
	}
	blockhash,err:=b.Hash()
	if err != nil {
		return Hash{}, err

	}

	blockfs:=BlockFs{
		Key:   blockhash,
		Value: b,
	}

	blockfsjson,err:=json.Marshal(blockfs)
	if err != nil {
		return Hash{},err
	}

	fmt.Println("Persisting new block to disk ")
	fmt.Printf("新增块信息为：%s\n",blockfsjson)

	if _,err:= s.dbFile.Write(append(blockfsjson,'\n'));err!=nil{
		return Hash{},err
	}


	s.latestBlock=b
	s.latestBlockHash=blockhash

	// 使用copy后的 state ，计算 balances 后再更新原来的 balances
	s.Balances=pendingState.Balances

	s.hasGenesisBlock=true
	return blockhash,nil

}
// 持久化从peer 获取的 blocks 信息
func (s *State) AddBlocks(blocks []Block)error  {
	for _,b := range blocks{
		_,err:=s.AddBlock(b)
		if err != nil {
			return err
		}
	}
	return nil
}


// 验证block是否合法,并更新 balances 信息
func applyBlock(b Block,s State) error {
	nextExpectedBlock :=s.latestBlock.Header.Number+1

	// block.db不为空 且 块高度不符合
	if s.hasGenesisBlock && b.Header.Number!=nextExpectedBlock  {
		return fmt.Errorf("next expected block must be %d not %d \n ",
			nextExpectedBlock,
			b.Header.Number,
			)
	}
	// ? 重复判断
	// block.db不为空 且 parent hash不符合
	if s.hasGenesisBlock &&  s.latestBlock.Header.Number>0 && !reflect.DeepEqual(b.Header.Parent,s.latestBlockHash) {

		return fmt.Errorf("next block parent must be %x not %x \n",
			s.latestBlockHash,
			b.Header.Parent,
		)
	}
	return appyTxs(b.Txs,&s)
}

// 根据txs 信息 更新 state 的 balances
func appyTxs(txs []Tx,s *State)error  {
	for _,tx:=range txs{
		err:=applyTx(tx,s)
		if err != nil {
			return err
		}
	}
	return nil
}

// 根据tx信息 更新 state 的 balances
func applyTx(tx Tx,s *State) error {
	if tx.IsReward(){
		s.Balances[tx.To]+=tx.Value
		return nil
	}
	if tx.Value>s.Balances[tx.From]{
		return fmt.Errorf("wrong Tx: sender %s balances is %d TBB ,tx cost is %d \n,",
			tx.From,
			s.Balances[tx.From],
			tx.Value,
		)
	}

	s.Balances[tx.From]-=tx.Value
	s.Balances[tx.To] += tx.Value
	return nil
}


//
//func (s *State) Persist()(Hash,error)  {
//
//	latestBlockHash:=s.latestBlockHash
//
//	block:=NewBlock(latestBlockHash,s.latestBlock.Header.Number+1,uint64(time.Now().Unix()),s.txMempool)
//
//	blockhash,err:=block.Hash()
//	if err != nil {
//		return Hash{},err
//	}
//
//	blockfs:=BlockFs{
//		Key:   blockhash,
//		Value: block,
//	}
//
//	blockfsjson,err:=json.Marshal(blockfs)
//	if err != nil {
//		return Hash{},err
//	}
//	fmt.Println("Persisting new block to disk ")
//	fmt.Printf("%s\n",blockfsjson)
//
//	if _,err:= s.dbFile.Write(append(blockfsjson,'\n'));err!=nil{
//		return Hash{},err
//	}
//
//	//s.latestBlockHash=latestBlockHash  //? 存疑问
//	s.latestBlock=block
//	s.latestBlockHash=blockhash
//	s.txMempool=[]Tx{}
//
//	return s.latestBlockHash,nil  // 存疑问
//
//
//}
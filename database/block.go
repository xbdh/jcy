package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Hash [32]byte

//  类型别名的序列化
// 使用 MarshalText支持非字符串作为key的map
func (h Hash) MarshalText()([]byte,error)  {
	return []byte(hex.EncodeToString(h[:])),nil
}
// 一定是指针类型
func (h *Hash) UnmarshalText(data []byte)error  {
	_,err:=hex.Decode(h[:],data)
	return err
}


type Block struct {
	Header BlockHeader `json:"header"`
	Txs []Tx           `json:"payload"`
}

type BlockHeader struct {
	Parent Hash        `json:"parent"`
	Number uint64      `json:"number"`
	Time uint64        `json:"time"`
}

type BlockFs struct {
	Key Hash          `json:"hash"`
	Value Block       `json:"block"`
}

func NewBlock(parent Hash,number uint64,time uint64, txs []Tx) Block {
	return Block{
		Header: BlockHeader{
			Parent: parent,
			Number: number,
			Time:   time,
		},
		Txs:    txs,
	}
}

func (b Block) Hash()(Hash,error)  {
	blockjson,err:=json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(blockjson),nil
}


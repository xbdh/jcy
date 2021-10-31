package database

import (
	"encoding/json"
	"io/ioutil"
)
var genesisJson = `
{
  "genesis_time": "2021-10-28 22:09",
  "chain_id": "the-blockchain-jcy",
  "balances": {
    "andrej": 1000000
  }
}`


type genesis struct {
	Balances map[Account]uint `json:"balances"`
}

func loadGenesis(path string)(genesis,error)  {
	//fmt.Println(path)
	content,err:=ioutil.ReadFile(path)
	if err != nil {
		return genesis{}, err
	}

	var loadedGenesis genesis

	err =json.Unmarshal(content,&loadedGenesis)
	if err != nil {
		return genesis{}, err
	}

	return loadedGenesis,nil
}

func writeGenesisToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(genesisJson), 0644)
}
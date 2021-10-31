package database

import (
	"bufio"
	"encoding/json"
	"os"
)

func GetBlocksAfer(blockHash Hash,dataDir string) ([]Block,error) {
	f, err := os.OpenFile(getBlocksDbFilePath(dataDir), os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}

	blocks :=make([]Block,0)
	shouldstartCollecting :=false

	scanner := bufio.NewScanner(f)

	for scanner.Scan(){
		if err:=scanner.Err();err!=nil{
			return nil, err
		}


		var blockFs BlockFs

		err:=json.Unmarshal(scanner.Bytes(),&blockFs)
		if err != nil {
			return nil, err
		}
		if shouldstartCollecting{
			blocks=append(blocks,blockFs.Value)
			continue
		}
		if blockFs.Key==blockHash{
			shouldstartCollecting=true

		}
	}
	return blocks,nil
}

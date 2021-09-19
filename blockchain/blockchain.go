package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)
type BlockChain struct {
	LastHash []byte
	Database *badger.DB

}

func InitBlockchain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions()
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err := db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh"));err==badger.ErrKeyNotFound{
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis Proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			err = txn.Set([]byte("lh"), genesis.Hash)
			lastHash = genesis.Hash
			return err
		}else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.Value()
			return err
		}
		
	})
	Handle(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}

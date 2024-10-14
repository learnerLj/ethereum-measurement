package measurement

type EthStatus struct {
	MsgMeta `bson:",inline"`

	Version   uint32 `json:"version" bson:"version"`     // protocol version
	NetworkID uint64 `json:"networkId" bson:"networkId"` // blockchain network id
	TD        string `json:"td" bson:"td"`               // total difficulty
	BlockHash string `json:"blockHash" bson:"blockHash"` // hash of the best known block
	Genesis   string `json:"genesis" bson:"genesis"`     // hash of the genesis block
	ForkID    string `json:"forkId" bson:"forkId"`       // EIP-2124 fork identifier
}

type EthTransaction struct {
	MsgMeta `bson:",inline"`

	//	we temporarily ignore the transactions, as they consume huge number of storage.
	//Transactions []Transaction `json:"transactions" bson:"transactions"`
}

type EthNewPooledTransactionHashes struct {
	MsgMeta `bson:",inline"`

	TxHashInfos []TxHashInfo `json:"txHashInfos" bson:"txHashInfos"`
}
type EthGetPooledTransactions struct {
	MsgMeta  `bson:",inline"`
	TxHashes []string `json:"txHashes" bson:"txHashes"`
}
type EthPooledTransactions struct {
	MsgMeta `bson:",inline"`

	//we temporarily ignore the transactions, as they consume huge number of storage.
	//Transactions []Transaction `json:"transactions" bson:"transactions"`
}

type RecordTransactions struct {
	MsgMeta `bson:",inline"`

	//we temporarily ignore the transactions, as they consume huge number of storage.
	//Transactions []Transaction `json:"transactions" bson:"transactions"`
}

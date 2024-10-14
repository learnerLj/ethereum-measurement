# Eth-Net


## Modification
For connection:
```go

// ####### go-ethereum/p2p/server.go ################
defaultDialRatio       = 3 
// -> defaultDialRatio       = 2

defaultMaxPendingPeers = 50
// -> defaultMaxPendingPeers = 2000

// ####### node/defaults.go #########################
MaxPeers:   50,
// -> MaxPeers:   10000,


```

for eth protocol:
- [ ] cancel our response(PooledTransactionsMsg) to request(GetPooledTransactionsMsg)
`eth/protocols/eth/peer.go:201`
```go

// ReplyPooledTransactionsRLP is the response to RequestTxs.
func (p *Peer) ReplyPooledTransactionsRLP(id uint64, hashes []common.Hash, txs []rlp.RawValue) error {
	// Mark all the transactions as known, but ensure we don't overflow our limits
	p.knownTxs.Add(hashes...)
	//++ return nil
	// Not packed into PooledTransactionsResponse to avoid RLP decoding
	return p2p.Send(p.rw, PooledTransactionsMsg, &PooledTransactionsRLPPacket{
		RequestId:                     id,
		PooledTransactionsRLPResponse: txs,
	})
}
```

- [ ] only anounce transaction hash instead. `eth/handler.go:463`
```go
    if broadcast {
        //++ annos[peer] = append(annos[peer], tx.Hash())
        // -- txset[peer] = append(txset[peer], tx.Hash())
    } else {
        annos[peer] = append(annos[peer], tx.Hash())
    }
```


Geth's DHT actually has 17 buckets (instead of 256) with a default size of 16, 
because of even distribution. The occurrence of the most significant n digits 
being the same only happens with a probability of 1/2^n. With n equal to 200, 
this probability is extremely tiny. In fact, Geth only keeps buckets for the 
upper 1/15 of distances. This limits the maximum number of peers, as peers occupy 
all the DHT and a new node only has the opportunity to be accepted when the network fluctuates. 
Therefore, we modify the bucket size to 256 to connect to up to 256*17 = 4352 peers.

Also, we keep the `seedCount` as 1000, the number of node records dump into database,
to speed up peer lookup for re-launch. Making Kademlia concurrency factor 
greater(from original 3 to now 100) improves the lookup ability too.
```go
// p2p/discover/table.go:45
bucketSize      = 16 // Kademlia bucket size
// -> bucketSize      = 256 // Kademlia bucket size

//++
originBucketSize = 16  // Geth original bucket size

seedCount         = 30
// -> 	seedCount         = 1000

//p2p/discover/table.go:44
alpha           = 3   // Kademlia concurrency factor
// -> 	alpha           = 100   // Kademlia concurrency factor
```

avoid extensive search
```go
p2p/discover/v4_udp.go:724
closest := t.tab.findnodeByID(target, originBucketSize, true).entries
```


-----------------
blob tx should not broadcast,
eth/handler.go:530 notice which tx to send
eth/protocols/eth/broadcast.go:73 receive

-----
backup database: `mongodump --host 158.132.255.36 --port 27017 -u jiahao -p luojiahao --authenticationDatabase admin --out ./mongodb_bak`

restore database: `mongorestore ./mongodb_bak`

Or move data: `sudo rsync -apPvh eth.45measurement /media/zihao/disk4/jiahao`

file link: `sudo ln -s /media/zihao/disk5/jiahao/mongodb/txs txs`


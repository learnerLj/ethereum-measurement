

geth控制台筛选节点类型：
`var excludeList = ["Geth", "Nethermind", "erigon", "besu", "reth", "CoreGeth"]; admin.peers.filter(peer => !excludeList.some(exclude => peer.name.startsWith(exclude)));`

查询重复的peer:
```
db.peers.aggregate([
    {
        "$group": {
            "_id": "$id",
            "count": { "$sum": 1 }
        }
    },
    {
        "$match": {
            "count": { "$gte": 2 }
        }
    },
    {
        "$project": {
            "_id": 0,
            "id": "$_id",
            "count": 1
        }
    }
]);
```

neighbors长度统计
```
db.discvNeighbors.aggregate([
    {
        "$match": {
            "neighbors": { "$exists": true, "$not": { "$size": 0 } } // 确保 neighbors 存在且不为空数组
        }
    },
    {
        "$addFields": {
            "neighborsLength": { "$size": "$neighbors" }
        }
    },
    {
        "$match": {
            "neighborsLength": { "$nin": [4, 12, 16] }
        }
    },
    {
        "$group": {
            "_id": "$neighborsLength",
            "count": { "$sum": 1 }
        }
    },
    {
        "$sort": { "_id": 1 }
    },
    {
        "$project": {
            "_id": 0,
            "neighborsLength": "$_id",
            "count": 1
        }
    }
]);

```


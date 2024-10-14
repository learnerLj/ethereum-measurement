package measurement_test

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/measurement"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"sort"
	"strings"
	"testing"
	"time"
)

var peers = []*measurement.Peer{
	{
		Id:          "peer1",
		Type:        measurement.OutboundNodes,
		Caps:        []measurement.Cap{{Name: "eth", Version: 63}, {Name: "snap", Version: 1}},
		NetworkID:   1,
		GenesisHash: "0x123...",
		ForkID:      "0xabc...",
		IsConnected: false,
	},
	{
		Id:          "peer2",
		Type:        measurement.InboundNodes,
		Caps:        []measurement.Cap{{Name: "eth", Version: 63}, {Name: "les", Version: 2}},
		NetworkID:   3,
		GenesisHash: "0x456...",
		ForkID:      "0xdef...",
		IsConnected: true,
	},
	{
		Id:          "peer3",
		Type:        measurement.VisitedNodes,
		Caps:        []measurement.Cap{{Name: "eth", Version: 63}},
		NetworkID:   1,
		GenesisHash: "0x789...",
		ForkID:      "0xghi...",

		IsConnected: true,
	},
}

var messages = []measurement.Msg{
	&measurement.DiscvPing{
		MsgMeta: measurement.MsgMeta{ // 显式初始化嵌套的MsgMeta
			LocalId:  "localId1",
			RemoteId: "remoteId1",
			UUId:     uuid.New().String(),
			Moment:   time.Now(),
			// 其余MsgMeta字段根据需要填充
		},
		// 填充DiscvPing特有的字段
		FromIp:  "192.168.1.1",
		FromUdp: 30303,
		FromTcp: 30303,
		ToIp:    "192.168.1.2",
		ToUdp:   30304,
		// 其余字段根据需要填充
	},
	&measurement.DiscvPong{
		MsgMeta: measurement.MsgMeta{ // 同上，为DiscvPong初始化MsgMeta
			LocalId:  "localId2",
			RemoteId: "remoteId2",
			UUId:     uuid.New().String(),
			Moment:   time.Now(),
		},
		// 填充DiscvPong特有的字段，如果有的话
	},
	// 可以继续添加更多消息实例
}

var envPath = "/media/zihao/disk1/jiahao/workstation/eth-nodes/go-ethereum/measurement/instruct.env"

func testNodeConfig() *node.Config {
	testNodeKey, _ := crypto.GenerateKey()
	return &node.Config{
		Name: "test node",
		P2P:  p2p.Config{PrivateKey: testNodeKey},
	}
}

func measure4Test() *measurement.Measurer {
	m := measurement.NewMeasurement(envPath)
	stack, err := node.New(testNodeConfig())
	if err != nil {
		log.Fatalf("failed to create protocol stack: %v", err)
	}
	s, n := stack.Server(), stack.Server().NodeInfo()
	m.InitNode(s.Logger.New("Measure", ""), s.NodeInfo().ID, n.IP, n.Ports.Listener, n.Ports.Discovery) // add local node interface
	return m
}

func showDbCollections(m *measurement.Measurer) []string {
	collections, err := m.DB().ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	return collections

}
func TestDbConnection(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	m := measure4Test()
	cs := showDbCollections(m)
	sort.Slice(cs, func(i, j int) bool {
		return cs[i][0] < cs[j][0]
	})
	fmt.Println(cs, len(cs))

}

func TestInsertMsg(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	m := measure4Test()

	for _, msg := range messages {
		_, err := m.InsertMsg(msg)
		if err != nil {
			t.Errorf("Failed to insert message %s: %v", msg.UUID(), err)
		}
		fmt.Println(msg)
	}

	for _, msg := range messages {
		var msgFromDb measurement.Msg
		switch msg.(type) {
		case *measurement.DiscvPing:
			msgFromDb = new(measurement.DiscvPing)
		case *measurement.DiscvPong:
			msgFromDb = new(measurement.DiscvPong)
		}

		err := m.GetMsg(msg.UUID(), msgFromDb)
		if err != nil {
			t.Errorf("Failed to retrieve message %s: %v", msg.UUID(), err)
		}
		if msgFromDb.UUID() != msg.UUID() {
			t.Errorf("UUID mismatch: expected %s, got %s", msg.UUID(), msgFromDb.UUID())
		}
		fmt.Println(msg)
	}
}

// TestInsertPeer 测试插入Peer实例到数据库的功能
func TestInsertPeer(t *testing.T) {
	m := measure4Test()

	for _, peer := range peers {
		// 假设InsertPeer是向数据库插入Peer实例的方法
		err := m.InsertPeer(peer)
		if err != nil {
			t.Errorf("Failed to insert peer %s: %v", peer.Id, err)
		}
	}
}

func TestDeleteIndex(t *testing.T) {

	m := measure4Test()

	indexNames := []string{
		"localId_1", // 索引名通常是字段名后跟下划线和排序顺序
		"remoteId_1",
		"uuid_1",
	}

	collections, err := m.DB().ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	for _, collectionName := range collections {
		collection := m.DB().Collection(collectionName)
		indexes := collection.Indexes()

		for _, indexName := range indexNames {
			_, err := indexes.DropOne(context.Background(), indexName)
			if err != nil {
				if strings.Contains(err.Error(), "index not found with name") {
					log.Printf("Index %s does not exist in collection %s\n", indexName, collectionName)
				} else {
					log.Fatal(err)
				}
			}
		}
	}

}

package measurement

import (
	"context"
	"errors"
	"fmt"
	ethlog "github.com/ethereum/go-ethereum/log"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"sync"
	"time"
)

var Measure *Measurer

var (
	ethMsgs   = [5]Msg{&EthStatus{}, &EthTransaction{}, &EthNewPooledTransactionHashes{}, &EthGetPooledTransactions{}, &EthPooledTransactions{}}
	rlpxMsgs  = [4]Msg{&RLPxHello{}, &RLPxDisc{}, &RLPxPing{}, &RLPxPong{}}
	discvMsgs = [6]Msg{&DiscvPing{}, &DiscvPong{}, &DiscvFindNode{}, &DiscvNeighbors{}, &DiscvENRRequest{}, &DiscvENRResponse{}}
)

type Measurer struct {
	db           *mongo.Database
	txdb         *mongo.Database
	msgInstances []Msg
	log          ethlog.Logger
	//myNode       *node.Node

	//localNodeInfo
	node NodeInfo
	//id                   string
	//ip, tcpPort, udpPort string

	initDone     sync.WaitGroup // wait for initialization
	loadSeedOnce sync.Once
}

func (m *Measurer) DoOnce(f func()) {
	m.loadSeedOnce.Do(f)
}

func NewMeasurement(envPath string) *Measurer {
	//config measurement environment
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("Fatal! can't find environment setting for measurement")
	}

	m := Measurer{}

	m.openMongoDb()

	// register all types of messages
	msgInstances := append(append(ethMsgs[:], rlpxMsgs[:]...), discvMsgs[:]...)
	for _, v := range msgInstances {
		m.registerMsg(v)
	}

	//create index to speed up lookup
	m.creatIndex()

	m.initDone.Add(1)
	return &m
}

func (m *Measurer) InitNode(logger ethlog.Logger, id, ip string, tcpPort, udpPort int) {

	//register local Node
	m.log = logger
	m.node.Id = id
	m.node.Ip, m.node.Tcp, m.node.Udp = ip, uint16(tcpPort), uint16(udpPort)

	m.initDone.Done()
}

func (m *Measurer) openMongoDb() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	// check connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	dbName := os.Getenv("MONGODB_NAME")
	if dbName == "" {
		log.Fatal("You must set your 'MONGODB_NAME' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	txDbName := os.Getenv("MONGODB_TXS_DB")
	if txDbName == "" {
		log.Fatal("You must set your 'MONGODB_TXS_DB' environment variable.")
	}

	m.db = client.Database(dbName)
	m.txdb = client.Database(txDbName)
}

func (m *Measurer) DB() *mongo.Database {
	return m.db
}

func (m *Measurer) InsertDialTask(d *DialTask) error {
	m.initDone.Wait()
	nodesColl := m.db.Collection("dialTask")
	_, err := nodesColl.InsertOne(context.TODO(), d)
	if err != nil {
		return err
	}
	return nil
}

// whichMatch: 1 for caps, 2 for eth
func (m *Measurer) UpdateDialTask(id string, whichMatch int) error {
	match := "matchCaps"
	if whichMatch == 2 {
		match = "matchChain"
	}

	nodesColl := m.db.Collection("dialTask")

	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", bson.D{{match, true}}}}

	// 执行更新操作
	_, err := nodesColl.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *Measurer) GetPeer(id string) *Peer {
	nodesColl := m.db.Collection("peers")
	filter := bson.D{{"id", id}}

	// Retrieves the first matching document
	var p Peer
	err := nodesColl.FindOne(context.TODO(), filter).Decode(&p)

	// Prints a message if no documents are matched or if any
	// other errors occur during the operation
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		m.log.Info("db finding", err)
		return nil
	}
	return &p
}

func (m *Measurer) CheckPeerConnected(id string) bool {
	p := m.GetPeer(id)
	if p != nil && p.IsConnected {
		return true
	}
	return false
}

func (m *Measurer) InsertPeer(peer *Peer) error {
	m.initDone.Wait()
	nodesColl := m.db.Collection("peers")
	_, err := nodesColl.InsertOne(context.TODO(), peer)
	if err != nil {
		m.log.Info("db insert", err)
		return err
	}
	return nil
}

func (m *Measurer) SetPeerConnectionStatus(id string, connected bool) error {
	m.initDone.Wait()
	nodesColl := m.db.Collection("peers")

	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", bson.D{{"isConnected", connected}}}}

	// 执行更新操作
	_, err := nodesColl.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		m.log.Info("db update", err)
		return err
	}

	return nil
}

func (m *Measurer) UpdatePeerWithHello(id string, hello *RLPxHello, peerType PeerType) error {
	m.initDone.Wait()
	nodesColl := m.db.Collection("peers")

	filter := bson.D{{"id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"clientId", hello.ClientId},
			{"caps", hello.Caps},
			{"peerType", peerType},
		}},
	}

	result, err := nodesColl.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		m.log.Info("db update", err)
		return err
	}

	// It should never trigger
	if result.MatchedCount == 0 {
		err := fmt.Errorf("no peer found with ID %s", id)
		m.log.Info("db update", err)
		return err
	}

	return nil
}

func (m *Measurer) UpdatePeerWithEth(id string, status *EthStatus) error {
	m.initDone.Wait()
	nodesColl := m.db.Collection("peers")

	filter := bson.D{{"id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"networkID", status.NetworkID},
			{"genesisHash", status.Genesis},
			{"forkID", status.ForkID},
		}},
	}

	result, err := nodesColl.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		m.log.Info("db update", err)
		return err
	}

	// It should never trigger
	if result.MatchedCount == 0 {
		err := fmt.Errorf("no peer found with ID %s", id)
		m.log.Info("db update", err)
		return err
	}

	return nil
}

func (m *Measurer) GetMsg(uid string, msg Msg) error {
	filter := bson.D{{"uuid", uid}}
	collection := m.db.Collection(msg.Name())
	// Retrieves the first matching document
	err := collection.FindOne(context.TODO(), filter).Decode(msg)
	return err
}

func (m *Measurer) InsertMsg(msg Msg) (string, error) {
	m.initDone.Wait() // we need local id
	err := m.insertMsg(msg)
	return msg.UUID(), err
}

func (m *Measurer) insertMsg(msg Msg) error {
	collectionName := msg.Name()
	collection := m.db.Collection(collectionName)
	_, err := collection.InsertOne(context.TODO(), msg)
	return err
}

func (m *Measurer) InsertMsgAddMeta(msg Msg, remoteId string, inbound bool, uid string) (string, error) {
	msg.SetMeta(m, remoteId, inbound, uid)

	// we can derive its id by signature. So insert peer

	//ignore RLPx ping, RLPx pong
	//outbound
	ignore := false

	//	TODO: ENRResponse rest fields analysis
	switch msgWithType := msg.(type) {
	case *DiscvPing:
		//ignore = true
		// add source endpoint if we send ping
		if !inbound {
			msgWithType.FromIp = m.node.Ip
			msgWithType.FromTcp = m.node.Tcp
			msgWithType.FromUdp = m.node.Udp
		}

		//insert node
		if inbound {
			remotePeer := m.GetPeer(remoteId)
			if remotePeer == nil {
				p := &Peer{
					Ip:     msgWithType.FromIp,
					Udp:    msgWithType.FromUdp,
					Tcp:    msgWithType.FromTcp,
					Id:     remoteId,
					Type:   VisitedNodes,
					Moment: time.Now()}
				m.InsertPeer(p)
			}
		}
	case *DiscvPong:
		//ignore = true

	case *DiscvFindNode:
		//ignore = true
	case *DiscvNeighbors:
		//ignore = true
		//remotePeer := m.GetPeer(remoteId)
		// insert node
		if inbound {
			for _, n := range msgWithType.Neighbors {
				//not found peer
				if dbNode := m.GetPeer(n.Id); dbNode == nil {
					p := &Peer{
						Ip:     n.Ip,
						Udp:    n.Udp,
						Tcp:    n.Tcp,
						Id:     n.Id,
						Type:   VisitedNodes,
						Moment: time.Now()}
					m.InsertPeer(p)
				}
			}
		}

	case *DiscvENRRequest:
		ignore = true
	case *DiscvENRResponse:
		ignore = true

	case *RLPxPing:
		ignore = true
	case *RLPxPong:
		ignore = true

	case *RLPxHello:
		remotePeer := m.GetPeer(remoteId)
		// peers send Hello to us
		if inbound {
			t := InboundNodes
			if msgWithType.Dial {
				t = OutboundNodes
			}
			if remotePeer != nil {
				// update capabilities of an existing peer
				m.UpdatePeerWithHello(remoteId, msgWithType, t)
			} else {
				// insert a new peer
				p := &Peer{
					Ip:       msgWithType.Ip,
					Udp:      msgWithType.Udp,
					Tcp:      msgWithType.Tcp,
					Id:       remoteId,
					Type:     t,
					ClientId: msgWithType.ClientId,
					Caps:     msgWithType.Caps,
					Moment:   time.Now(),
				}
				m.InsertPeer(p)
			}
		} else {
			//ignore = true
		}
	case *EthStatus:
		//	we only consider inbound message.
		if inbound {
			m.UpdatePeerWithEth(remoteId, msgWithType)
		} else {
			//ignore = true
		}

	}

	if ignore {
		return msg.UUID(), nil
	}

	return m.InsertMsg(msg)
}

func (m *Measurer) creatIndex() {
	indexModels := []mongo.IndexModel{
		//{
		//	Keys: bson.D{{Key: "localId", Value: 1}}, // 1 indicates the ascending index
		//},
		{
			Keys: bson.D{{Key: "remoteId", Value: 1}},
		},
		//{
		//	Keys: bson.D{{Key: "uuid", Value: 1}},
		//},
	}

	for _, msgInstance := range m.msgInstances {
		collectionName := msgInstance.Name()
		collection := m.db.Collection(collectionName)
		_, err := collection.Indexes().CreateMany(context.Background(), indexModels)
		if err != nil {
			log.Fatal(err)
		}
	}

	peerModel := mongo.IndexModel{
		Keys: bson.D{{Key: "id", Value: 1}},
	}
	//	create peer collection
	_, err := m.db.Collection("peers").Indexes().CreateOne(context.Background(), peerModel)
	m.db.Collection("dialTask").Indexes().CreateOne(context.Background(), peerModel)
	taskModel := mongo.IndexModel{
		Keys: bson.D{{Key: "uuid", Value: 1}},
	}
	m.db.Collection("newDialTask").Indexes().CreateOne(context.Background(), taskModel)
	if err != nil {
		log.Fatal(err)
	}

}

func (m *Measurer) registerMsg(instance Msg) {
	m.msgInstances = append(m.msgInstances, instance)
}

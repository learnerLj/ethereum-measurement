package measurement

import (
	"github.com/google/uuid"
	"time"
)

type PeerType uint

const (
	VisitedNodes PeerType = iota
	InboundNodes
	OutboundNodes
)

type Peer struct {
	Moment time.Time `json:"moment" bson:"moment"`
	//  Discv4 protocol
	// Network information
	Ip  string `json:"ip" bson:"ip"`
	Udp uint16 `json:"udp" bson:"udp"`
	Tcp uint16 `json:"tcp" bson:"tcp"`

	// Id is capable to derive from Pubkey, but still record it for convenient
	Id string `json:"id" bson:"id"`

	// whether We find it or it finds us. The result actually may change but mostly not.
	Type PeerType `json:"peerType" bson:"peerType"`

	// RLPx protocol

	// we will resolve ClientId and get them in the data analysis
	//
	// Client information, we split the RLPxHello message into different parts.
	//ClientName      string `json:"clientName" bson:"clientName"`
	//ClientVersion   string `json:"clientVersion" bson:"clientVersion"`
	//ClientPlatform  string `json:"clientPlatform" bson:"clientPlatform"`
	//CompilerVersion string `json:"compilerVersion" bson:"compilerVersion"`

	ClientId string `json:"clientId" bson:"clientId"`
	Caps     []Cap  `json:"caps" bson:"caps"`

	//reason for disconnect could be found in the RLPxDisc
	//DisconnectReason string `json:"disconnectReason" bson:"disconnectReason"`

	// eth protocol
	NetworkID   int    `json:"networkID" bson:"networkID"`
	GenesisHash string `json:"genesisHash" bson:"genesisHash"`
	ForkID      string `json:"forkID" bson:"forkID"`

	// connected peer
	IsConnected bool `json:"isConnected" bson:"isConnected"`
}

type DialTask struct {
	To         string    `json:"id" bson:"id"`
	Moment     time.Time `json:"moment" bson:"moment"`
	MatchCaps  bool      `json:"matchCaps" bson:"matchCaps"`
	MatchChain bool      `json:"matchChain" bson:"matchChain"`
}

type NewDialTask struct {
	UUId       string    `json:"uuid" bson:"uuid"`
	To         string    `json:"id" bson:"id"`
	Moment     time.Time `json:"moment" bson:"moment"`
	MatchCaps  bool      `json:"matchCaps" bson:"matchCaps"`
	MatchChain bool      `json:"matchChain" bson:"matchChain"`
}

type Msg interface {
	Name() string
	UUID() string
	SetMeta(measure *Measurer, remoteId string, inbound bool, uid string)
}

type MsgMeta struct {
	LocalId  string    `json:"localId" bson:"localId"`
	RemoteId string    `json:"remoteId" bson:"remoteId"`
	Inbound  bool      `json:"inbound" bson:"inbound"`
	Moment   time.Time `json:"moment" bson:"moment"`
	UUId     string    `json:"uuid" bson:"uuid"`
	// only corresponding request and response have the same UUID
}

// Topology is semi-stable at one moment
type Topology struct {
}

const (
	EthStatusMsg = 0x00

	//EthNewBlockHashesMsg  = 0x01

	EthTransactionsMsg = 0x02

	//EthGetBlockHeadersMsg = 0x03
	//EthBlockHeadersMsg    = 0x04
	//EthGetBlockBodiesMsg  = 0x05
	//EthBlockBodiesMsg     = 0x06
	//EthNewBlockMsg        = 0x07

	EthNewPooledTransactionHashesMsg = 0x08
	EthGetPooledTransactionsMsg      = 0x09
	EthPooledTransactionsMsg         = 0x0a

	//EthGetReceiptsMsg                = 0x0f
	//EthReceiptsMsg                   = 0x10

	DiscvPingMsg        = 0x01 // zero is 'reserved'
	DiscvPongMsg        = 0x02
	DiscvFindnodeMsg    = 0x03
	DiscvNeighborsMsg   = 0x04
	DiscvENRRequestMsg  = 0x05
	DiscvENRResponseMsg = 0x06

	RLPxHelloMsg = 0x00
	RLPxDiscMsg  = 0x01
	RLPxPingMsg  = 0x02
	RLPxPongMsg  = 0x03
)

// NodeInfo for Neighbors Packet (0x04)
type NodeInfo struct {
	Ip  string `json:"ip" bson:"ip"`
	Udp uint16 `json:"udp" bson:"udp"`
	Tcp uint16 `json:"tcp" bson:"tcp"`
	Id  string `json:"id" bson:"id"`
}

type Cap struct {
	Name    string `json:"name" bson:"name"`
	Version uint   `json:"version" bson:"version"`
}

type ForkID struct {
	ForkHash string `json:"forkHash" bson:"forkHash"` // fork-hash
	ForkNext string `json:"forkNext" bson:"forkNext"` // fork-next
}

type TxHashInfo struct {
	TxType byte   `json:"txType" bson:"txType"`
	TxSize int64  `json:"txSize" bson:"txSize"`
	TxHash string `json:"txHash" bson:"txHash"`
}

func (m *MsgMeta) UUID() string {
	return m.UUId
}

func (m *MsgMeta) SetMeta(measure *Measurer, remoteId string, inbound bool, uid string) {
	m.LocalId = measure.node.Id
	m.Moment = time.Now()
	if len(uid) != 0 {
		m.UUId = uid
	} else {
		m.UUId = uuid.NewString()
	}

	m.RemoteId = remoteId
	m.Inbound = inbound
}

func (e *EthStatus) Name() string {
	return "ethStatus"
}
func (e *EthTransaction) Name() string {
	return "ethTransaction"
}
func (e *EthNewPooledTransactionHashes) Name() string {
	return "ethNewPooledTransactionHashes"
}
func (e *EthGetPooledTransactions) Name() string {
	return "ethGetPooledTransactions"
}
func (e *EthPooledTransactions) Name() string {
	return "ethPooledTransactions"
}
func (e *RLPxHello) Name() string {
	return "RLPxHello"
}

func (e *RLPxDisc) Name() string {
	return "RLPxDisc"
}

func (e *RLPxPing) Name() string {
	return "RLPxPing"
}

func (e *RLPxPong) Name() string {
	return "RLPxPong"
}

func (e *DiscvPing) Name() string {
	return "discvPing"
}

func (e *DiscvPong) Name() string {
	return "discvPong"
}

func (e *DiscvFindNode) Name() string {
	return "discvFindNode"
}

func (e *DiscvNeighbors) Name() string {
	return "discvNeighbors"
}

func (e *DiscvENRRequest) Name() string {
	return "discvENRRequest"
}

func (e DiscvENRResponse) Name() string {
	return "discvENRResponse"
}

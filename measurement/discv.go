package measurement

// Message types are distinguished by placing them in different collections

// DiscvPing is set from information of outbound ping when insertion
type DiscvPing struct {
	MsgMeta `bson:",inline"`

	FromIp  string `json:"fromIp" bson:"fromIp"`
	FromUdp uint16 `json:"fromUdp" bson:"fromUdp"`
	FromTcp uint16 `json:"fromTcp" bson:"fromTcp"`

	ToIp  string `json:"toIp" bson:"toIp"`
	ToUdp uint16 `json:"toUdp" bson:"toUdp"`
	//ENRSeq           uint64 `json:"seq" bson:"seq"`
}

type DiscvPong struct {
	MsgMeta `bson:",inline"`

	ToIp  string `json:"toIp" bson:"toIp"`
	ToUdp uint16 `json:"toUdp" bson:"toUdp"`
	ToTcp uint16 `json:"toTcp" bson:"toTcp"`
}

// DiscvFindNode for FindNode Packet (0x03)
type DiscvFindNode struct {
	MsgMeta `bson:",inline"`

	Target string `json:"target" bson:"target"`
}

// DiscvNeighbors for Neighbors Packet (0x04)
type DiscvNeighbors struct {
	MsgMeta `bson:",inline"`

	Timeout   bool        `json:"timeout" bson:"timeout"`
	Neighbors []*NodeInfo `json:"neighbors" bson:"neighbors"`
}

// DiscvENRRequest for ENRRequest Packet (0x05)
type DiscvENRRequest struct {
	MsgMeta `bson:",inline"`
}

// DiscvENRResponse for ENRResponse Packet (0x06)
type DiscvENRResponse struct {
	MsgMeta `bson:",inline"`
	//Enr     string `json:"enr" bson:"enr"`
	//ENR without resolve is useless for our measurement. so we record its basic identity information

	Ip  string `json:"ip" bson:"ip"`
	Id  string `json:"id" bson:"id"`
	Udp uint16 `json:"udp" bson:"udp"`
	Tcp uint16 `json:"tcp" bson:"tcp"`
	Seq uint64 `json:"seq" bson:"seq"`
}

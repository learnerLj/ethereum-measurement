package measurement

// Message types are distinguished by placing them in different collections

type RLPxHello struct {
	MsgMeta `bson:",inline"`

	// if an inbound hello is received, we can not really make sure the peer
	// has been recorded with its id and endpoint.
	// only for convenience.
	Ip   string `json:"ip" bson:"ip"`
	Udp  uint16 `json:"udp" bson:"udp"`
	Dial bool   `json:"dial" bson:"dial"`

	ClientId string `json:"clientId" bson:"clientId"`
	Caps     []Cap  `json:"caps" bson:"caps"`
	Tcp      uint16 `json:"tcp" bson:"tcp"` // Ignore if 0
}

type RLPxDisc struct {
	MsgMeta `bson:",inline"`

	DiscReason uint8 `json:"DiscReason" bson:"DiscReason"` // see https://github.com/ethereum/devp2p/blob/master/rlpx.md#disconnect-0x01
}

type RLPxPing struct {
	MsgMeta `bson:",inline"`
}

type RLPxPong struct {
	MsgMeta `bson:",inline"`
}

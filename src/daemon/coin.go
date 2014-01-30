package daemon

// Communication layer for the coin pkg

var (
    GetBlocksPrefix  = gnet.MessagePrefix{'G', 'E', 'T', 'B'}
    GiveBlocksPrefix = gnet.MessagePrefix{'G', 'I', 'V', 'B'}
)

type GetBlocksMessage struct {
    c *gnet.MessageContext `enc:"-"`
}

func NewGetBlocksMessage() *GetBlocksMessage {
    return &GetBlocksMessage{}
}

func (self *GetBlocksMessage) Handle(mc *gnet.MessageContext) error {
    self.c = mc
    return recordMessageEvent(self, mc)
}

type GiveBlocksMessage struct {
    c *gnet.MessageContext `enc:"-"`
}

func NetGiveBlocksMessage() *GiveBlocksMessage {
    return &GiveBlocksMessage{}
}

func (self *GiveBlocksMessage) Handle(mc *gnet.MessageContext) error {
    self.c = mc
    return recordMessageEvent(self, mc)
}

package daemon

import (
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/visor"
)

// A connection's state within the daemon
type Connection struct {
    Id           int    `json:"id"`
    Addr         string `json:"address"`
    LastSent     int64  `json:"last_sent"`
    LastReceived int64  `json:"last_received"`
    // Whether the connection is from us to them (true, outgoing),
    // or from them to us (false, incoming)
    Outgoing bool `json:"outgoing"`
    // Whether the client has identified their version, mirror etc
    Introduced bool   `json:"introduced"`
    Mirror     uint32 `json:"mirror"`
    ListenPort uint16 `json:"listen_port"`
}

// An array of connections
// Arrays must be wrapped in structs to avoid certain javascript exploits
type Connections struct {
    Connections []*Connection `json:"connections"`
}

type SpendResult struct {
    RemainingBalance visor.BalanceResult       `json:"remaining_balance"`
    Transaction      visor.ReadableTransaction `json:"txn"`
    Error            string                    `json:"error"`
}

type BlockchainProgress struct {
    // Our current blockchain length
    Current uint64 `json:"current"`
    // Our best guess at true blockchain length
    Highest uint64 `json:"highest"`
}

type ResendResult struct {
    Sent bool `json:"sent"`
}

type RPC struct{}

func (self RPC) GetConnection(d *Daemon, addr string) *Connection {
    if d.Pool.Pool == nil {
        return nil
    }
    c := d.Pool.Pool.Addresses[addr]
    _, expecting := d.ExpectingIntroductions[addr]
    return &Connection{
        Id:           c.Id,
        Addr:         addr,
        LastSent:     c.LastSent.Unix(),
        LastReceived: c.LastReceived.Unix(),
        Outgoing:     (d.OutgoingConnections[addr] == nil),
        Introduced:   !expecting,
        Mirror:       d.ConnectionMirrors[addr],
        ListenPort:   d.GetListenPort(addr),
    }
}

func (self RPC) GetConnections(d *Daemon) *Connections {
    if d.Pool.Pool == nil {
        return nil
    }
    conns := make([]*Connection, 0, len(d.Pool.Pool.Pool))
    for _, c := range d.Pool.Pool.Pool {
        conns = append(conns, self.GetConnection(d, c.Addr()))
    }
    return &Connections{Connections: conns}
}

func (self RPC) Spend(v *Visor, pool *Pool, vrpc visor.RPC, amt visor.Balance,
    fee uint64, dest coin.Address) *SpendResult {
    if v.Visor == nil {
        return nil
    }
    txn, err := v.Spend(amt, fee, dest, pool)
    errString := ""
    if err != nil {
        errString = err.Error()
        logger.Error("Failed to make a spend: %v", err)
    }
    b := vrpc.GetTotalBalance(v.Visor, true)
    return &SpendResult{
        RemainingBalance: *b,
        Transaction:      visor.NewReadableTransaction(&txn),
        Error:            errString,
    }
}

func (self RPC) GetBlockchainProgress(v *Visor) *BlockchainProgress {
    if v.Visor == nil {
        return nil
    }
    return &BlockchainProgress{
        Current: v.Visor.MostRecentBkSeq(),
        Highest: v.EstimateBlockchainLength(),
    }
}

func (self RPC) ResendTransaction(v *Visor, p *Pool,
    txHash coin.SHA256) *ResendResult {
    if v.Visor == nil {
        return nil
    }
    sent := v.ResendTransaction(txHash, p)
    return &ResendResult{
        Sent: sent,
    }
}

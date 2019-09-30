// +build ignore

// 20160901 - Initial version by user johnstuartmill,
// public key 02fb4acf944c84d48341e3c1cb14d707034a68b7f931d6be6d732bec03597d6ff6
// 20161025 - Code revision by user johnstuartmill.
package main

//
// WARNING: WARNING: WARNING: Do NOT use this code for obtaining any
// research results.  This file is only an illustration. A realistic
// simulation would require to have (i) nonzero latencies for event
// propagation and (ii) an event queue inside the implementation of
// MeshNetworkInterface.
//

import (
	"flag"
	"fmt"
	"log"
	mathrand "math/rand"
	"os"
	"sort"
	"time"

	"github.com/SkycoinProject/skycoin/src/daemon/gnet"
	//
	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/consensus"
)

var Cfg_print_config bool = true
var Cfg_debug_connect_request bool = false
var Cfg_debug_node_final_state bool = false
var Cfg_debug_node_summary bool = false
var Cfg_debug_show_block_maker bool = false

var Cfg_simu_topology_is_random bool = true

var Cfg_simu_num_node int = 10
var Cfg_simu_num_blockmaker int = 2
var Cfg_simu_prob_malicious float64 = 0.0
var Cfg_simu_prob_duplicate float64 = 0.0

var Cfg_simu_num_block_round int = 10
var Cfg_simu_fanout_per_node int = 3

// Will be reset later, based on values of other parameters:
var Cfg_simu_num_iter int = 0

var common_channel uint16 = 3

////////////////////////////////////////////////////////////////////////////////
func pretty_print_flags(prefix string, detail bool) {
	if detail {

		max1 := 0
		max2 := 0

		flag.VisitAll(func(f *flag.Flag) {
			len1 := len(f.Name)
			len2 := len(fmt.Sprintf("%v", f.Value))
			if max1 < len1 {
				max1 = len1
			}
			if max2 < len2 {
				max2 = len2
			}
		})

		format := fmt.Sprintf("    --%%-%ds %%%dv    %%s\n", max1, max2)
		format = "%s" + format

		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf(format, prefix, f.Name, f.Value, f.Usage)
		})

	} else {

		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf("%s--%s=%v\n", prefix, f.Name, f.Value)
		})

	}
}

////////////////////////////////////////////////////////////////////////////////
func cmd_line_args_process() {

	var ip *int = nil
	var qp *uint64 = nil
	var dp *float64 = nil
	var bp *bool = nil

	//
	// Simulation parameters
	//
	ip = &Cfg_simu_num_node
	flag.IntVar(ip, "simu-num-nodes", *ip, "Number of nodes in the network.")

	ip = &Cfg_simu_num_blockmaker
	flag.IntVar(ip, "simu-num-blockmaker", *ip,
		"Number of nodes in the network that make blocks.")

	dp = &Cfg_simu_prob_malicious
	flag.Float64Var(dp, "simu-prob-malicious", *dp,
		"Probability that a node temporarily joins a malicious group that"+
			" publishes same block in order to cause a fork of the blockchain.")

	dp = &Cfg_simu_prob_duplicate
	flag.Float64Var(dp, "simu-prob-duplicate", *dp,
		"Probability that a node sends a duplicate message with same hash but"+
			" different signature. (Duplicate (hash,sig) pairs are easily"+
			" detected and discarded.)")

	ip = &Cfg_simu_num_block_round
	flag.IntVar(ip, "simu-num-rounds", *ip,
		"Number of block rounds. When all them are published and the"+
			" resulting messages propagate, the simulation ends.")

	ip = &Cfg_simu_fanout_per_node
	flag.IntVar(ip, "simu-fanout-per-node", *ip,
		"Number of incoming (and outgoing) connections to (and from) each"+
			" node.")

	bp = &Cfg_debug_connect_request
	flag.BoolVar(bp, "debug-connect-request", *bp, "")

	bp = &Cfg_print_config
	flag.BoolVar(bp, "print-config", *bp, "")

	bp = &Cfg_debug_node_final_state
	flag.BoolVar(bp, "debug-node-final-state", *bp, "")

	bp = &Cfg_debug_node_summary
	flag.BoolVar(bp, "debug-node-summary", *bp, "")

	bp = &Cfg_debug_show_block_maker
	flag.BoolVar(bp, "debug-show-block-maker", *bp, "")

	bp = &Cfg_simu_topology_is_random
	flag.BoolVar(bp, "simu-topology-is-random", *bp,
		"Connect nodes randomly or place them in one circle.")

	//
	// Consensus parameters
	//

	bp = &consensus.Cfg_debug_block_duplicate
	flag.BoolVar(bp, "debug-block-duplicate", *bp, "")

	bp = &consensus.Cfg_debug_block_out_of_sequence
	flag.BoolVar(bp, "debug-block-out-of-sequence", *bp, "")

	bp = &consensus.Cfg_debug_block_accepted
	flag.BoolVar(bp, "debug-block-accepted", *bp, "")

	bp = &consensus.Cfg_debug_HashCandidate
	flag.BoolVar(bp, "debug-hash-candidate", *bp, "")

	ip = &consensus.Cfg_blockchain_tail_length
	flag.IntVar(ip, "blockchain-tail-length", *ip,
		"Blocks held in memory. This limits memory usage.")

	qp = &consensus.Cfg_consensus_candidate_max_seqno_gap
	flag.Uint64Var(qp, "consensus-candidate-max-seqno-gap", *qp,
		"Proposed blocks (or consensus candidates) are ignored if theie seqno"+
			" is too high or too low w.r.t. what is stored. This limits memory"+
			" use and helps prevents some mild attacks.")

	qp = &consensus.Cfg_consensus_waiting_time_as_seqno_diff
	flag.Uint64Var(qp, "consensus-waiting-time-as-seqno-diff", *qp,
		"When to decide on selecting the best hash from BlockStat"+
			" so that it can be moved to blockchain.")

	ip = &consensus.Cfg_consensus_max_candidate_messages
	flag.IntVar(ip, "consensus-max-candidate-messages", *ip,
		"How many (hash,signer_pubkey) pairs to acquire for decision-making."+
			" This also limits forwarded traffic, because the messages in excess"+
			" of this limit are discarded hence not forwarded.")

	//
	//
	//
	show := flag.Bool("show", false, "Show current parameter values and exit.")

	//
	//
	flag.Parse()
	//
	//

	if Cfg_simu_num_node < Cfg_simu_num_blockmaker {
		fmt.Printf("Invalid input: --simu-num-nodes=%d < --simu-num-blockmaker="+
			"%d. Exiting.\n", Cfg_simu_num_node, Cfg_simu_num_blockmaker)
		os.Exit(1)
	}

	if Cfg_simu_prob_malicious < 0. || 1 < Cfg_simu_prob_malicious {
		fmt.Printf("Invalid input: --simu-prob-malicious=%g is outside"+
			" [0 .. 1] range. Exiting.\n", Cfg_simu_prob_malicious)
		os.Exit(1)
	}

	if Cfg_simu_prob_duplicate < 0. || 1 < Cfg_simu_prob_duplicate {
		fmt.Printf("Invalid input: --simu-prob-duplicate=%g is outside"+
			" [0 .. 1] range. Exiting.\n", Cfg_simu_prob_malicious)
		os.Exit(1)
	}

	//
	// Derived parameters
	//

	// Most likely we do not need that many. However, we keep the
	// number high so it would not interfere with message propagation
	// by premature exit from the vent loop. Yet we keep it finite to
	// prevent an infinite run that can be caused by a bug:

	Cfg_simu_num_iter = 10 * // '10' is a heuristic
		Cfg_simu_num_node * Cfg_simu_num_blockmaker *
		Cfg_simu_num_block_round * Cfg_simu_fanout_per_node

	if *show {
		pretty_print_flags("", true)
		os.Exit(1)
	} else {
		if Cfg_print_config {
			pretty_print_flags("FILE_Config.txt|", false)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
//
//
//
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
// The body of this function lends itself to something like
//
//    ConsensusParticipant::BuildAndPropagateNewBlock()
//
// Before doing so, ConsensusParticipant would need to accumulate
// transactions, possibly negotiate with others as to who makes blocks
// etc. FOR NOW, any node can make (and publish) blocks.
//
func propagate_hash_from_node(
	h cipher.SHA256,
	nodePtr *consensus.ConsensusParticipant,
	external_use bool,
	external_seqno uint64) {

	//
	// WARNING: Do NOT use this code for obtaining any research
	// results.  This file is only an illustration. A realistic
	// simulation require to have nonzero latencies for event
	// propagation and to have an event queueu inside the
	// implementation of MeshNetworkInterface.
	//

	o := external_seqno // HACK for DEBUGGING
	if !external_use {
		o = nodePtr.GetNextBlockSeqNo() // So that blocks are ordered.
	}

	b := consensus.BlockBase{}
	b.Init(
		nodePtr.SignatureOf(h), // Signature of hash.
		h,
		o)

	nodePtr.OnBlockHeaderArrived(&b)
}

////////////////////////////////////////////////////////////////////////////////
// MESSAGE
type BlockBaseWrapper struct {
	consensus.BlockBase
}

func (self *BlockBaseWrapper) String() string {
	return fmt.Sprintf("BlockBaseWrapper={Sig=%s,Hash=%s,Seqno=%d}",
		self.Sig.Hex()[:8], self.Hash.Hex()[:8], self.Seqno)
}
func (self *BlockBaseWrapper) Handle(context *gnet.MessageContext, closure interface{}) error {

	if false {
		fmt.Printf("@@@ consensus.BlockBase Handle: context={%v} closure={%v} data={%v}\n", context, closure, self)
	}

	arg := closure.(*PoolOwner)
	arg.DataCallback(context, self)
	return nil
}

//array of messages to register
var messageMap map[string](interface{}) = map[string](interface{}){
	"id01": BlockBaseWrapper{}, //message id, message type
}

////////////////////////////////////////////////////////////////////////////////
//
//
//
////////////////////////////////////////////////////////////////////////////////
type PoolOwner struct {
	pCMan *MinimalConnectionManager

	pDispatcherManager *gnet.DispatcherManager
	pDispatcher        *gnet.Dispatcher
	pConnectionPool    *gnet.ConnectionPool
	//
	isConnSolicited map[string]bool // Addr -> bool
	//
	key_list []string        // To preserve the order.
	key_map  map[string]byte // To avoid duplicates.
	//
	debug_num_id   int
	debug_nickname string
}

func (self *PoolOwner) Shutdown() {
	self.pConnectionPool.Shutdown()
}

func (self *PoolOwner) Init(pCMan *MinimalConnectionManager, listen_port uint16, num_id int, nickname string) {

	config := gnet.NewConfig()
	config.Port = uint16(listen_port)

	cp := gnet.NewConnectionPool(config)

	dm := gnet.NewDispatcherManager()

	self.isConnSolicited = make(map[string]bool)

	cp.Config.MessageCallback = dm.OnMessage
	cp.Config.ConnectCallback = self.OnConnect
	cp.Config.DisconnectCallback = self.OnDisconnect

	d := dm.NewDispatcher(cp, common_channel, self)
	d.RegisterMessages(messageMap)

	self.pCMan = pCMan
	self.pDispatcherManager = dm
	self.pDispatcher = d
	self.key_map = make(map[string]byte)
	self.pConnectionPool = cp

	self.debug_num_id = num_id
	self.debug_nickname = nickname
}

func (self *PoolOwner) Run() {

	cp := self.pConnectionPool // alias

	if cp.Config.Port != 0 {
		if err := cp.StartListen(); err != nil {
			log.Panic(err)
		}
		go cp.AcceptConnections()
	}

	go func() {
		for true {
			// [JSM] 20161025 The 'Sleep' below is a workaround of
			// unresolve concurrency issues in file gnet/pool.go
			// class ConnectionPool.
			time.Sleep(time.Millisecond * 1)
			cp.HandleMessages()
		}
	}()

}
func (self *PoolOwner) DataCallback(context *gnet.MessageContext, xxx *BlockBaseWrapper) {

	if self.isConnSolicited[context.Conn.Addr()] {

		var msg consensus.BlockBase
		msg.Sig = xxx.Sig
		msg.Hash = xxx.Hash
		msg.Seqno = xxx.Seqno

		self.pCMan.GetNode().OnBlockHeaderArrived(&msg)

	} else {
		// Ignoring
	}
}
func (self *PoolOwner) OnConnect(c *gnet.Connection, is_solicited bool) {
	self.isConnSolicited[c.Addr()] = is_solicited

	if false {
		for channel, pDispatcher := range self.pDispatcherManager.Dispatchers {
			// (BTW, pDispatcher.ReceivingObject is '&PoolOwner')

			// Each channel (or subscription subject) might require
			// specific action to be taken on connect, such as request
			// initial values from solicited, authenticate the client
			// from unsolicited one etc.
			_ = channel
			_ = pDispatcher
		}
	}
}
func (self *PoolOwner) OnDisconnect(c *gnet.Connection, reason gnet.DisconnectReason) {

	if false {
		for channel, pDispatcher := range self.pDispatcherManager.Dispatchers {
			// (BTW, pDispatcher.ReceivingObject is '&PoolOwner')

			// Each channel (or subscription subject) might require
			// specific action to be taken on disconnect, such as
			// reconnect, clean up etc.
			_ = channel
			_ = pDispatcher
		}
	}

	delete(self.isConnSolicited, c.Addr())
}

////////////////////////////////////////////////////////////////////////////////
func (self *PoolOwner) RequestConnectionToKeys() {
	for _, key := range self.key_list {
		go func(arg string) {
			_, err := self.BlockingConnectToUrl(arg)
			if err != nil {
				log.Panic(err)
			}
		}(key)
	}
}

////////////////////////////////////////////////////////////////////////////////
func (self *PoolOwner) RegisterKey(key string) bool {

	if _, have := self.key_map[key]; have {
		return false // Let caller handle his issues.
	}

	self.key_list = append(self.key_list, key)
	self.key_map[key] = byte('1')

	return true

}
func (self *PoolOwner) GetListenPort() uint16 {
	return self.pDispatcher.Pool.Config.Port
}
func (self *PoolOwner) BlockingConnectTo(IPAddress string, port uint16) (*gnet.Connection, error) {
	url := fmt.Sprintf("%s:%d", IPAddress, port)
	conn, err := self.pDispatcher.Pool.Connect(url)
	return conn, err
}
func (self *PoolOwner) BlockingConnectToUrl(url string) (*gnet.Connection, error) {
	conn, err := self.pDispatcher.Pool.Connect(url)
	return conn, err
}
func (self *PoolOwner) BroadcastMessage(msg gnet.Message) error {
	return self.pDispatcher.BroadcastMessage(common_channel, msg)
}
func (self *PoolOwner) Print() {
	detail := true

	fmt.Printf("PoolOwner={%d,%s,keys={n=%d",
		self.debug_num_id, self.debug_nickname, len(self.key_list))

	if detail {
		for _, val := range self.key_list {
			fmt.Printf(",%s", val)
		}
	} else {
		fmt.Printf(",...")
	}
	fmt.Printf("}}")
}

////////////////////////////////////////////////////////////////////////////////
func print_stat(X []*MinimalConnectionManager, iter int) {

	n := 0
	for i, _ := range X {
		n += X[i].GetNode().Incoming_block_count
	}

	msg_per_node_per_round :=
		float64(n) / float64(Cfg_simu_num_node*Cfg_simu_num_block_round)

	msg_per_node_per_round_per_link := msg_per_node_per_round /
		float64(Cfg_simu_fanout_per_node)

	msg_per_node_per_round_per_blockmaker := msg_per_node_per_round /
		float64(Cfg_simu_num_blockmaker)

	// Print for viewing:
	fmt.Printf(
		"MSG_STAT iter                     %d\n"+
			"MSG_STAT msg_count_all            %d\n"+
			"MSG_STAT num_node                 %d\n"+
			"MSG_STAT num_blockmaker           %d\n"+
			"MSG_STAT num_block_round          %d\n"+
			"MSG_STAT fanout_per_node          %d\n"+
			"MSG_STAT max_candidate_messages   %d (This limits the effect of"+
			" having large num_blockmaker)\n"+
			"MSG_STAT msg_per_node_per_round                %.3f\n"+
			"MSG_STAT msg_per_node_per_round_per_link       %.3f\n"+
			"MSG_STAT msg_per_node_per_round_per_blockmaker %.3f\n",
		iter, n, Cfg_simu_num_node, Cfg_simu_num_blockmaker,
		Cfg_simu_num_block_round, Cfg_simu_fanout_per_node,
		consensus.Cfg_consensus_max_candidate_messages,
		msg_per_node_per_round,
		msg_per_node_per_round_per_link,
		msg_per_node_per_round_per_blockmaker)

}

////////////////////////////////////////////////////////////////////////////////
func Simulate_compare_node_StateQueue(
	X []*MinimalConnectionManager,
	global_seqno2h map[uint64]cipher.SHA256,
	global_seqno2h_alt map[uint64]cipher.SHA256,
) {
	//
	// Step 1 of 3: for each observed seqno, find the histogram of
	// 'best' hash. The historgam is formed by summing over nodes.
	//
	type QQQ map[uint64]map[cipher.SHA256]int
	xxx := make(QQQ) // Access:       [seqno][hash]=count
	type ZZZ []QQQ   // Access: [node][seqno][hash]=count

	ni := len(X)

	zzz := make(ZZZ, ni)

	for i := 0; i < ni; i++ { // Nodes
		nj := X[i].GetNode().Get_block_stat_queue_Len()

		zzz[i] = make(map[uint64]map[cipher.SHA256]int)

		for j := 0; j < nj; j++ { // Elements in node's BlockStatQueue

			// 'bs' a pointer:
			bs := X[i].GetNode().Get_block_stat_queue_element_at(j)
			seqno := bs.GetSeqno()
			hash, _, _ := bs.GetBestHashPubkeySig()

			if _, have := xxx[seqno]; !have {
				xxx[seqno] = make(map[cipher.SHA256]int)
			}
			xxx[seqno][hash]++

			if _, have := zzz[i][seqno]; !have {
				zzz[i][seqno] = make(map[cipher.SHA256]int)
			}
			zzz[i][seqno][hash]++
		}
	}
	//
	// Step 2 of 3. For each seqno, find the most-frequently observed
	// hash. Also, find the ratio of blocks accepted to blocks
	// published.
	//
	yyy := make(map[uint64]cipher.SHA256) // Access: [seqno]=hash

	var accept_count int = 0
	var total_count int = 0

	for seqno, hash2count := range xxx {
		var best_count int = 0
		var sum_count int = 0
		var best_hash cipher.SHA256 // undef

		initialized := false

		for hash, count := range hash2count {
			if initialized {
				if best_count < count {
					best_count = count
					best_hash = hash
				}
			} else {
				initialized = true

				best_count = count
				best_hash = hash
			}
			sum_count += count
		}

		if initialized {
			yyy[seqno] = best_hash

			// Here all 'seqno' contribute equally:
			accept_count += best_count
			total_count += sum_count
		}
	}

	if true {
		keys := []int{}
		for seqno, _ := range yyy {
			keys = append(keys, int(seqno))
		}

		sort.Ints(keys)

		for _, key := range keys {
			seqno := uint64(key)

			// Most-frequently accepted (across nodes) for the given seqno:
			best_hash := yyy[seqno]

			prescribed := best_hash == global_seqno2h[seqno]
			malicious := best_hash == global_seqno2h_alt[seqno]

			fmt.Printf("CONSENSUS: seqno=%d best_hash=%s prescribed=%t"+
				" malicious=%t\n", seqno, best_hash.Hex()[:8], prescribed,
				malicious)
		}

	}
	fmt.Printf("CONSENSUS: total_count=%d accept_count=%d, accept_ratio=%f\n",
		total_count, accept_count, float32(accept_count)/float32(total_count))

	for i, zzz_i := range zzz {
		join_count := 0       // How many have selected the most popular hash.
		other_count := 0      // How many have selected NOT the most popular.
		prescribed_count := 0 // How many have selected the intended hash.
		malicious_count := 0  // How many have selected the malicious hash.

		for seqno, hash2count := range zzz_i {

			// Most-frequently accepted (across nodes) for the given seqno:
			best_hash := yyy[seqno]
			prescribed := global_seqno2h[seqno]
			malicious := global_seqno2h_alt[seqno]

			for hash, count := range hash2count {
				if hash == best_hash {
					join_count += count
				} else {
					other_count += count
				}

				if hash == prescribed {
					prescribed_count += count
				}
				if hash == malicious {
					malicious_count += count
				}
			}
		}

		fmt.Printf("NODE i=%d join_count=%d other_count=%d prescribed_count="+
			"%d  malicious_count=%d\n",
			i, join_count, other_count, prescribed_count, malicious_count)
	}
}

////////////////////////////////////////////////////////////////////////////////
type MinimalConnectionManager struct {
	theNodePtr *consensus.ConsensusParticipant

	// MinimalConnectionManager solicit a conn to them; receive data
	// from them but do not send anything
	publishers PoolOwner

	// MinimalConnectionManager accept
	// connections. MinimalConnectionManager send data to them; does
	// not receive anything; ignores all incoming data
	subscribers PoolOwner
}

func (self *MinimalConnectionManager) Init(listen_port uint16, num_id int, nickname string) {
	self.publishers.Init(self, 0, num_id, nickname+";recv_from_pub")
	self.subscribers.Init(self, listen_port, num_id, nickname+";send_to_sub")
}
func (self *MinimalConnectionManager) Run() {
	self.publishers.Run()
	self.subscribers.Run()
}
func (self *MinimalConnectionManager) ShutdownPublishing() {
	self.subscribers.Shutdown()
}
func (self *MinimalConnectionManager) ShutdownSubscribing() {
	self.publishers.Shutdown()
}
func (self *MinimalConnectionManager) GetListenPort() uint16 {
	return self.subscribers.pDispatcher.Pool.Config.Port
}

////////////////////////////////////////////////////////////////////////////////
func (self *MinimalConnectionManager) GetNode() *consensus.ConsensusParticipant {
	return self.theNodePtr
}
func (self *MinimalConnectionManager) RegisterPublisher(key string) bool {
	return self.publishers.RegisterKey(key)
}
func (self *MinimalConnectionManager) SendBlockToAllMySubscriber(xxx *consensus.BlockBase) {
	var msg BlockBaseWrapper
	msg.Sig = xxx.Sig
	msg.Hash = xxx.Hash
	msg.Seqno = xxx.Seqno

	self.subscribers.pDispatcher.BroadcastMessage(common_channel, &msg)
}
func (self *MinimalConnectionManager) RequestConnectionToAllMyPublisher() {
	self.publishers.RequestConnectionToKeys()

	// NOTE: The node does not request connection to subscriber;
	// instead, the subscriber request connection to Node.
}
func (self *MinimalConnectionManager) OnSubscriberConnectionRequest(key string) {
	fmt.Printf(">>>>>>>>>>>>>>>> GOT HERE. Please cooment this line out >>>>>>>>>>>>>>>>>>>>>>>>>\n")

	// FOR NOW accept all connection request. TODO: check for black
	// list, latency table etc.
	var acceptable = true
	if acceptable {
		self.subscribers.RegisterKey(key)
	}
}
func (self *MinimalConnectionManager) Print() {
	fmt.Printf("ConnectionManager={Pub={")
	self.publishers.Print()
	fmt.Printf("},Sub={")
	self.subscribers.Print()
	fmt.Printf("}}")
}

////////////////////////////////////////////////////////////////////////////////
//
//
//
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
func get_random_index_subset(N int, S int) []int {
	// N - population size
	// S - subset size

	if N < 0 {
		N = 0
	}
	if S < 0 {
		S = 0
	}
	if S > N {
		S = N
	}

	index_map := make(map[int]int, S)
	if 2*S < N {
		// Include at random
		for i := 0; i < 3*S; i++ { // '3' is a heuristic
			if len(index_map) >= S {
				break
			}
			index_map[mathrand.Intn(N)] = 1
		}
	} else {
		// Fill up
		for i := 0; i < N; i++ {
			index_map[i] = 1
		}
		n := N - S
		// Exclude at random
		for i := 0; i < 3*n; i++ { // '3' is a heuristic
			if len(index_map) <= S {
				break
			}
			delete(index_map, mathrand.Intn(N))
		}
	}

	keys := []int{}

	for k, _ := range index_map {
		keys = append(keys, k)
	}

	return keys
}

////////////////////////////////////////////////////////////////////////////////
//
// main
//
////////////////////////////////////////////////////////////////////////////////
func main() {

	cmd_line_args_process()

	// PERFORMANCE:
	cipher.DebugLevel1 = false
	cipher.DebugLevel2 = false

	var X []*MinimalConnectionManager

	var hack_global_seqno uint64 = 0

	seed := "hdhdhdkjashfy7273"
	_, SecKeyArray :=
		cipher.GenerateDeterministicKeyPairsSeed([]byte(seed), Cfg_simu_num_node)

	for i := 0; i < Cfg_simu_num_node; i++ {
		var cm MinimalConnectionManager
		cm.Init(6060+uint16(i), i, "fox")

		// Reason for mutual registration: (1) when conn man receives
		// messages, it needs to notify the node; (2) when node has
		// processed a mesage, it might need to use conn man to send
		// some data out.
		nodePtr := consensus.NewConsensusParticipantPtr(&cm)
		s := SecKeyArray[i]
		nodePtr.SetPubkeySeckey(cipher.PubKeyFromSecKey(s), s)

		cm.theNodePtr = nodePtr

		X = append(X, &cm)
	}
	if false {
		fmt.Printf("Got %d nodes\n", len(X))
	}

	if Cfg_simu_topology_is_random {

		fmt.Printf("CONFIG Topology: connecting %d nodes randomly with approx"+
			" %d nearest-neighbors in and approx %d nearest-neighbors out.\n",
			Cfg_simu_num_node, Cfg_simu_fanout_per_node,
			Cfg_simu_fanout_per_node)

		n := len(X)

		for i := 0; i < n; i++ {

			cm := X[i]

			indices :=
				get_random_index_subset(Cfg_simu_num_node, Cfg_simu_fanout_per_node)
			for _, j := range indices {
				if i != j {
					cm.RegisterPublisher(fmt.Sprintf("127.0.0.1:%d", X[j].GetListenPort()))
					fmt.Printf("TOPOLOGY: port %d solicits conn to port %d\n", cm.GetListenPort(), X[j].GetListenPort())

				}
			}
		}
	} else {

		fmt.Printf("CONFIG Topology: connecting %d nodes via one (thick)"+
			" circle with approx %d  nearest-neighbors in and approx %d "+
			"nearest-neighbors out.\n",
			Cfg_simu_num_node, Cfg_simu_fanout_per_node,
			Cfg_simu_fanout_per_node)

		n := len(X)

		for i := 0; i < n; i++ {

			cm := X[i]

			c_left := int(Cfg_simu_fanout_per_node / 2)
			c_right := Cfg_simu_fanout_per_node - c_left

			for c := 0; c < c_left; c++ {
				j := (i - 1 - c + n) % n
				cm.RegisterPublisher(fmt.Sprintf("127.0.0.1:%d", X[j].GetListenPort()))
			}

			for c := 0; c < c_right; c++ {
				j := (i + 1 + c) % n
				cm.RegisterPublisher(fmt.Sprintf("127.0.0.1:%d", X[j].GetListenPort()))
			}
		}
	}

	// Start GoRoutines related to connectivity
	for i, _ := range X {
		X[i].Run()
	}

	fmt.Printf("Waiting for start ...\n")
	time.Sleep(time.Millisecond * 10)

	// Connect. PROD: This should request connections. The
	// connections can be accepted, rejected or never answered. Such
	// replies are asynchronous. SIMU: we connect synchronously.
	for i, _ := range X {
		X[i].RequestConnectionToAllMyPublisher()
	}

	fmt.Printf("Waiting for connections ...\n")
	time.Sleep(time.Second * 1)

	global_seqno2h := make(map[uint64]cipher.SHA256)
	global_seqno2h_alt := make(map[uint64]cipher.SHA256)

	iter := 0
	block_round := 0
	done_processing_messages := false
	for ; iter < Cfg_simu_num_iter; iter++ {

		if false {
			fmt.Printf("Iteration %d\n", iter)
		}

		//fmt.Printf("Waiting for messages to propagate ...\n")
		time.Sleep(time.Millisecond * 100)

		if true {
			if block_round < Cfg_simu_num_block_round {

				if true {
					t := time.Now()
					fmt.Printf("wall_time=%02d:%02d:%02d"+
						" block_round=%d\n", t.Hour(), t.Minute(), t.Second(),
						block_round)
				}

				// NOTE: Propagating blocks from here is a
				// simplification/HACK: it implies that we have
				// knowledge of when messaging due to previous
				// activity (blocks and connections) has
				// stopped. Again, we make blocks from here for
				// debugging and testing only.

				//x := secp256k1.RandByte(888) // Random data in SIMU.
				x := make([]byte, 888)
				mathrand.Read(x)

				h := cipher.SumSHA256(x) // Its hash.

				//x_alt := secp256k1.RandByte(888) // Random data in SIMU.
				x_alt := make([]byte, 888)
				mathrand.Read(x)
				h_alt := cipher.SumSHA256(x_alt) // Its hash.

				global_seqno2h[hack_global_seqno] = h
				global_seqno2h_alt[hack_global_seqno] = h_alt

				indices := get_random_index_subset(Cfg_simu_num_node,
					Cfg_simu_num_blockmaker)

				if Cfg_debug_show_block_maker {
					fmt.Printf("block_round=%d, Random indices of block-"+
						"makers: %v\n", block_round, indices)
				}

				n_forkers := int(Cfg_simu_prob_malicious * float64(len(indices)))

				for i := 0; i < len(indices); i++ {
					// TODO: Have many nodes send same block, and a few nodes
					// send a different block. Research the conditions under
					// which the block published by the majority would
					// dominate the other one.
					index := indices[i]
					nodePtr := X[index].GetNode()

					malicious := (i < n_forkers)
					duplicate := (mathrand.Float64() < Cfg_simu_prob_duplicate)

					ph := &h
					if malicious {
						ph = &h_alt
					}

					rep := 1
					if duplicate {
						rep = 2
					}

					//
					// WARNING: In a reslistic simulation, one would
					// need to remove the assumption of knowing global
					// properties such as 'hack_global_seqno'
					//
					if malicious {
						fmt.Printf(">>>>>> NODE (index,pubkey)=(%d,%s) is"+
							" publishing ALTERNATIVE block\n", index,
							nodePtr.Pubkey.Hex()[:8])
					}

					for j := 0; j < rep; j++ {
						// Signing same hash multiple times produces different
						// signatures (for a good reason). We do it
						// here to test if malicious re-publishing is
						// detected properly.
						propagate_hash_from_node(*ph, nodePtr, true,
							hack_global_seqno)
					}
				}

				hack_global_seqno += 1
				block_round += 1
			} else {
				done_processing_messages = true
				break // <<<<<<<<
			}
		}
	}

	fmt.Printf("Waiting to finish ...\n")
	time.Sleep(time.Millisecond * 500)

	zzz := "done"
	if !done_processing_messages {
		zzz = "***NOT done***"
	}

	if false {
		// [JSM] Do not call these FOR NOW: ConnectionPool does not
		// implement shutdown correctly.
		for i, _ := range X {
			X[i].ShutdownPublishing()
		}
		for i, _ := range X {
			X[i].ShutdownSubscribing()
		}
	}

	fmt.Printf("Done (i) making Blocks, %s (ii) processing responses."+
		" See stats on the next few lines. Used iterations=%d, unused"+
		" iterations=%d. Exiting the event loop now.\n",
		zzz, iter, Cfg_simu_num_iter-iter)

	print_stat(X, iter)

	if Cfg_debug_node_final_state {
		for i, _ := range X {
			fmt.Printf("FILE_FinalState.txt|NODE i=%d ", i)
			X[i].GetNode().Print()
			fmt.Printf("\n")
		}
	}

	if Cfg_debug_node_summary {
		Simulate_compare_node_StateQueue(X, global_seqno2h, global_seqno2h_alt)
	}
}

////////////////////////////////////////////////////////////////////////////////

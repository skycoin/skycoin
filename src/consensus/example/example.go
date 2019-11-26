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
	mathrand "math/rand"
	"os"
	"sort"

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

var Cfg_simu_num_node int = 100
var Cfg_simu_num_blockmaker int = 10
var Cfg_simu_prob_malicious float64 = 0.0
var Cfg_simu_prob_duplicate float64 = 0.0

var Cfg_simu_num_block_round int = 10
var Cfg_simu_fanout_per_node int = 3

// Will be reset later, based on values of other parameters:
var Cfg_simu_num_iter int = 0

////////////////////////////////////////////////////////////////////////////////
type MinimalConnectionManager struct {
	theNodePtr *consensus.ConsensusParticipant
	//
	publisher_key_list  []*MinimalConnectionManager
	subscriber_key_list []*MinimalConnectionManager
}

func (self *MinimalConnectionManager) GetNode() *consensus.ConsensusParticipant {
	return self.theNodePtr
}
func (self *MinimalConnectionManager) RegisterPublisher(key *MinimalConnectionManager) bool {

	self.publisher_key_list = append(self.publisher_key_list, key)
	return true
}
func (self *MinimalConnectionManager) SendBlockToAllMySubscriber(blockPtr *consensus.BlockBase) {
	for _, p := range self.subscriber_key_list {
		p.GetNode().OnBlockHeaderArrived(blockPtr)
	}
}
func (self *MinimalConnectionManager) RequestConnectionToAllMyPublisher() {
	for _, p := range self.publisher_key_list {
		p.OnSubscriberConnectionRequest(self)
	}
}
func (self *MinimalConnectionManager) OnSubscriberConnectionRequest(other *MinimalConnectionManager) {
	self.subscriber_key_list = append(self.subscriber_key_list, other)
}
func (self *MinimalConnectionManager) Print() {
	detail := false

	fmt.Printf("ConnectionManager={publisher={n=%d",
		len(self.publisher_key_list))

	if detail {
		for _, val := range self.publisher_key_list {
			fmt.Printf(",%v", val)
		}
	} else {
		fmt.Printf(",...")
	}
	fmt.Printf("}")

	fmt.Printf(",subscriber={n=%d", len(self.subscriber_key_list))
	if detail {
		for _, val := range self.subscriber_key_list {
			fmt.Printf(",%v", val)
		}
	} else {
		fmt.Printf(",...")
	}
	fmt.Printf("}")
}

////////////////////////////////////////////////////////////////////////////////
//
//
//
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
func print_stat(X []*MinimalConnectionManager,
	iter int) {

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
		cm := MinimalConnectionManager{}
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
			" %d  nearest-neighbors in and approx %d nearest-neighbors out.\n",
			Cfg_simu_num_node, Cfg_simu_fanout_per_node,
			Cfg_simu_fanout_per_node)

		for i, _ := range X {
			cm := X[i]
			for g := 0; g < Cfg_simu_fanout_per_node; g++ {
				j := mathrand.Intn(Cfg_simu_num_node)
				if i != j {
					cm.RegisterPublisher(X[j])
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
				cm.RegisterPublisher(X[j])
			}

			for c := 0; c < c_right; c++ {
				j := (i + 1 + c) % n
				cm.RegisterPublisher(X[j])
			}
		}
	}

	// Connect. PROD: This should request connections. The
	// connections can be accepted, rejected or never answered. Such
	// replies are asynchronous. SIMU: we connect synchronously.
	for i, _ := range X {
		X[i].RequestConnectionToAllMyPublisher()
	}

	global_seqno2h := make(map[uint64]cipher.SHA256)
	global_seqno2h_alt := make(map[uint64]cipher.SHA256)

	iter := 0
	block_round := 0
	done_processing_messages := false
	for ; iter < Cfg_simu_num_iter; iter++ {

		if true {
			if block_round < Cfg_simu_num_block_round {

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

	zzz := "done"
	if !done_processing_messages {
		zzz = "***NOT done***"
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

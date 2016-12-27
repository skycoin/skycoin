package nodemanager

import (
	"fmt"
	"math"

	"github.com/skycoin/skycoin/src/cipher"
)

type Graph struct {
	nodes map[cipher.PubKey]map[cipher.PubKey]*DirectRoute
	paths map[cipher.PubKey]*SP
}

type DirectRoute struct {
	from   cipher.PubKey
	to     cipher.PubKey
	weight int
}

func NewGraph() *Graph {
	graph := Graph{}
	graph.Clear()
	return &graph
}

func (s *Graph) ToString() {
	for from, routes := range s.nodes {
		fmt.Println("\nNODE: ", from)
		for _, directRoute := range routes {
			fmt.Println("\tRoute: ", directRoute)
		}
		fmt.Println("======================\n")
	}
}

func (s *Graph) AddDirectRoute(from, to cipher.PubKey, weight int) {
	if from == to || weight < 1 {
		return
	}
	if _, ok := s.nodes[from]; !ok {
		s.nodes[from] = map[cipher.PubKey]*DirectRoute{}
	} else {
		if _, ok = s.nodes[from][to]; ok {
			return
		}
	}
	newDirectRoute := &DirectRoute{from, to, weight}
	s.nodes[from][to] = newDirectRoute
}

/* ---- can be useful in the future, maybe
func (s *Graph) RebuildRoutes() {
	s.clear()
	for node := range(nodes) {
		paths[node] = newSP(s, node)
	}
}
*/
func (s *Graph) FindRoute(from, to cipher.PubKey) ([]cipher.PubKey, bool) {
	sp, found := s.paths[from]
	if !found {
		sp = newSP(s, from)
		s.paths[from] = sp
	}
	route, found := sp.pathTo(to)
	return route, found
}

func (s *Graph) Clear() {
	s.nodes = map[cipher.PubKey]map[cipher.PubKey]*DirectRoute{}
	s.paths = map[cipher.PubKey]*SP{}
}

type SP struct { //ShortestPath
	source cipher.PubKey
	distTo map[cipher.PubKey]int
	edgeTo map[cipher.PubKey]*DirectRoute
	pq     *MinPQ
}

func newSP(graph *Graph, source cipher.PubKey) *SP {

	sp := SP{}

	sp.source = source
	sp.distTo = map[cipher.PubKey]int{}
	sp.edgeTo = map[cipher.PubKey]*DirectRoute{}

	for node := range graph.nodes {
		sp.distTo[node] = math.MaxInt32
	}
	sp.distTo[source] = 0

	sp.pq = newPQ()
	sp.pq.insert(source, 0)
	for !sp.pq.isEmpty() {
		v := sp.pq.delMin()
		for _, directRoute := range graph.nodes[v] {
			sp.relax(directRoute)
		}
	}

	return &sp
}

func (s *SP) pathTo(to cipher.PubKey) ([]cipher.PubKey, bool) { // if the path exists return a path and true, otherwise empty path and false

	path := []cipher.PubKey{to}
	e := s.edgeTo[to]

	for {
		if e == nil {
			return []cipher.PubKey{}, false
		} // no edge, so path doesn't exist
		if e.from == s.source {
			break
		} // we are at the source, work is finished
		path = append(path, e.from)
		e = s.edgeTo[e.from]
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 { // reverse an slice, in the future should apply stack instead of it
		path[i], path[j] = path[j], path[i]
	}

	return path, true
}

func (s *SP) relax(edge *DirectRoute) {
	from := edge.from
	to := edge.to
	newDist := s.distTo[from] + edge.weight
	if s.distTo[to] > newDist {
		s.distTo[to] = newDist
		s.edgeTo[to] = edge
		if s.pq.contains(to) {
			s.pq.decreaseKey(to, s.distTo[to])
		} else {
			s.pq.insert(to, s.distTo[to])
		}
	}
}

type NodeDist struct {
	node cipher.PubKey
	dist int
}

type MinPQ struct {
	keys []*NodeDist
}

func newPQ() *MinPQ {
	pq := MinPQ{}
	zeroND := &NodeDist{}
	pq.keys = []*NodeDist{zeroND}
	return &pq
}

func (pq *MinPQ) isEmpty() bool {
	return len(pq.keys) == 1
}

func (pq *MinPQ) contains(node cipher.PubKey) bool { //needs optimization, it's not effective
	for _, nodeDist := range pq.keys {
		if nodeDist.node == node {
			return true
		}
	}
	return false
}

func (pq *MinPQ) delMin() cipher.PubKey {
	n := len(pq.keys)
	if pq.isEmpty() {
		return pq.keys[0].node
	}
	min := pq.keys[1].node
	n--
	pq.exch(1, n)
	pq.keys = pq.keys[0:n]

	return min
}

func (pq *MinPQ) insert(node cipher.PubKey, dist int) {
	nodeDist := &NodeDist{node, dist}
	pq.keys = append(pq.keys, nodeDist)
	pq.swim(len(pq.keys) - 1)
}

func (pq *MinPQ) decreaseKey(node cipher.PubKey, dist int) { // needs optimizaion, it's not effective
	for _, key := range pq.keys {
		if key.node == node {
			key.dist = dist
			break
		}
	}
}

func (pq *MinPQ) swim(k int) {
	for k > 1 && pq.less(k, k/2) {
		pq.exch(k, k/2)
		k /= 2
	}
}

func (pq *MinPQ) sink(k int) {
	n := len(pq.keys) - 1
	for 2*k < n {
		j := 2 * k
		if j < n && pq.less(j+1, j) {
			j++
		}
		if !pq.less(j, k) {
			break
		}
		pq.exch(k, j)
	}
}

func (pq *MinPQ) less(i, j int) bool {
	return pq.keys[i].dist < pq.keys[j].dist
}

func (pq *MinPQ) exch(i, j int) {
	pq.keys[i], pq.keys[j] = pq.keys[j], pq.keys[i]
}

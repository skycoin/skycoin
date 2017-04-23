package nodemanager

import (
	"fmt"
	"math"
	"sync"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type RouteGraph struct {
	nodes map[cipher.PubKey]map[cipher.PubKey]*DirectRoute
	paths map[cipher.PubKey]*SP
	lock  *sync.Mutex
}

type DirectRoute struct {
	from   cipher.PubKey
	to     cipher.PubKey
	weight int
}

func newGraph() *RouteGraph {
	graph := RouteGraph{}
	graph.lock = &sync.Mutex{}
	graph.clear()
	return &graph
}

func (self *RouteGraph) toString() {
	self.lock.Lock()
	defer self.lock.Unlock()
	for from, routes := range self.nodes {
		fmt.Println("\nNODE: ", from)
		for _, directRoute := range routes {
			fmt.Println("\tRoute: ", directRoute)
		}
		fmt.Println("======================\n")
	}
}

func (self *RouteGraph) addDirectRoute(from, to cipher.PubKey, weight int) {
	if from == to || weight < 1 {
		return
	}

	_, ok := self.getDirectRoutes(from)

	if !ok {
		self.setDirectRoutes(from, map[cipher.PubKey]*DirectRoute{})
	} else {
		if _, ok = self.getDirectRoute(from, to); ok {
			return
		}
	}
	newDirectRoute := &DirectRoute{from, to, weight}
	self.setDirectRoute(from, to, newDirectRoute)
}

/* ---- can be useful in the future, maybe
func (self *RouteGraph) RebuildRoutes() {
	s.clear()
	for node := range(nodes) {
		paths[node] = newSP(s, node)
	}
}
*/
func (self *RouteGraph) findRoute(from, to cipher.PubKey) ([]cipher.PubKey, error) {
	sp, found := self.paths[from]
	if !found {
		sp = newSP(self, from)
		self.paths[from] = sp
	}
	route, err := sp.pathTo(to)
	return route, err
}

func (self *RouteGraph) clear() {
	self.nodes = map[cipher.PubKey]map[cipher.PubKey]*DirectRoute{}
	self.paths = map[cipher.PubKey]*SP{}
}

func (self *RouteGraph) getDirectRoutes(nodeId cipher.PubKey) (map[cipher.PubKey]*DirectRoute, bool) {
	self.lock.Lock()
	defer self.lock.Unlock()

	nodes, ok := self.nodes[nodeId]
	return nodes, ok
}

func (self *RouteGraph) getDirectRoute(nodeFrom, nodeTo cipher.PubKey) (*DirectRoute, bool) {
	self.lock.Lock()
	defer self.lock.Unlock()

	directRoute, ok := self.nodes[nodeFrom][nodeTo]
	return directRoute, ok
}

func (self *RouteGraph) setDirectRoutes(nodeId cipher.PubKey, nodes map[cipher.PubKey]*DirectRoute) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.nodes[nodeId] = nodes
}

func (self *RouteGraph) setDirectRoute(nodeFrom, nodeTo cipher.PubKey, directRoute *DirectRoute) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.nodes[nodeFrom][nodeTo] = directRoute
}

type SP struct { //ShortestPath
	source cipher.PubKey
	distTo map[cipher.PubKey]int
	edgeTo map[cipher.PubKey]*DirectRoute
	pq     *MinPQ
	lock   *sync.Mutex
}

func newSP(graph *RouteGraph, source cipher.PubKey) *SP {

	graph.lock.Lock()
	defer graph.lock.Unlock()

	sp := SP{}
	sp.lock = &sync.Mutex{}

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

func (s *SP) pathTo(to cipher.PubKey) ([]cipher.PubKey, error) { // if the path exists return a path and true, otherwise empty path and false

	s.lock.Lock()
	defer s.lock.Unlock()

	pathStack := newStack()
	pathStack.push(to)
	e := s.edgeTo[to]

	for {
		if e == nil {
			return []cipher.PubKey{}, messages.ERR_NO_ROUTE
		} // no edge, so path doesn't exist

		pathStack.push(e.from)
		if e.from == s.source {
			break
		} // we are at the source, work is finished
		e = s.edgeTo[e.from]
	}

	size := pathStack.size
	path := make([]cipher.PubKey, 0, size)

	for i := 0; i < int(size); i++ {
		path = append(path, pathStack.pop())
	}
	return path, nil
}

func (s *SP) relax(edge *DirectRoute) {
	from := edge.from
	to := edge.to

	s.lock.Lock()
	defer s.lock.Unlock()

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

type NodeStack struct {
	first *StackNode
	size  uint32
}

type StackNode struct {
	item cipher.PubKey
	next *StackNode
}

func newStack() *NodeStack {
	self := &NodeStack{}
	self.first = nil
	self.size = 0
	return self
}

func (self *NodeStack) push(item cipher.PubKey) {
	oldfirst := self.first
	self.first = &StackNode{
		item: item,
		next: oldfirst,
	}
	self.size++
}

func (self *NodeStack) pop() cipher.PubKey {
	if self.size == 0 {
		return cipher.PubKey{}
	}
	item := self.first.item
	self.first = self.first.next
	self.size--
	return item
}

type NodeDist struct {
	node cipher.PubKey
	dist int
}

type MinPQ struct {
	keys      []*NodeDist
	positions map[cipher.PubKey]int
	lock      *sync.Mutex
}

func newPQ() *MinPQ {
	pq := MinPQ{}
	zeroND := &NodeDist{}
	pq.keys = []*NodeDist{zeroND}
	pq.positions = map[cipher.PubKey]int{}
	pq.lock = &sync.Mutex{}
	return &pq
}

func (pq *MinPQ) isEmpty() bool {
	return len(pq.keys) == 1
}

func (pq *MinPQ) contains(node cipher.PubKey) bool {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	_, exists := pq.positions[node]
	return exists
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
	pq.lock.Lock()
	delete(pq.positions, min)
	pq.lock.Unlock()

	return min
}

func (pq *MinPQ) insert(node cipher.PubKey, dist int) {
	nodeDist := &NodeDist{node, dist}
	pq.keys = append(pq.keys, nodeDist)
	position := len(pq.keys) - 1
	pq.lock.Lock()
	pq.positions[node] = position
	pq.lock.Unlock()
	pq.swim(position)
}

func (pq *MinPQ) decreaseKey(node cipher.PubKey, dist int) {
	pq.lock.Lock()
	position, found := pq.positions[node]
	pq.lock.Unlock()
	if found {
		key := pq.keys[position]
		key.dist = dist
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
	pq.lock.Lock()
	pq.positions[pq.keys[i].node], pq.positions[pq.keys[j].node] = j, i
	pq.keys[i], pq.keys[j] = pq.keys[j], pq.keys[i]
	pq.lock.Unlock()
}

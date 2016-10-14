package gui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/mesh/nodemanager"
	"github.com/skycoin/skycoin/src/mesh/transport"
	wh "github.com/skycoin/skycoin/src/util/http"
)

//struct for nodeAddTransportHandler
type ConfigWithID struct {
	NodeID int
	Config nodemanager.TestConfig
}

//struct for nodeRemoveTransportHandler
type TransportWithID struct {
	NodeID    int
	Transport transport.ITransport
}

//struct for nodeGetListNodesHandler
// type ListNodes struct {
// 	PubKey []string
// }

func testHandler(nm *nodemanager.NodeManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		wh.Error400(w, fmt.Sprint("Works!"))

		if addr := r.FormValue("addr"); addr == "" {
			wh.Error404(w)
		} else {
			//wh.SendOr404(w, nm.GetConnection(addr))
			wh.Error404(w)
		}
	}
}

//Handler for /nodemanager/start - add new Node
//mode: GET
//url: /nodemanager/start
func nodeStartHandler(nm *nodemanager.NodeManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Starting Node")
		count := len(nm.NodesList)
		i := nm.AddNode()
		if i != count {
			wh.Error500(w)
		}

	}
}

//Handler for /nodemanager/stop - stop Node
//mode: GET
//url: /nodemanager/stop?id=value
func nodeStopHandler(nm *nodemanager.NodeManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Stoping Node")
		nodeID := r.FormValue("id")
		if nodeID == "" {
			wh.Error400(w, "Missing Node id")
			return
		}
		i, err := strconv.Atoi(nodeID)
		if err != nil {
			wh.Error400(w, "Node id must be integer")
			return
		}

		if len(nm.PubKeyList) < i {
			wh.Error400(w, "Invalid Node id")
			return
		}

		nm.NodesList[nm.PubKeyList[i]].Close()
		delete(nm.NodesList, nm.PubKeyList[i])
		nm.PubKeyList = append(nm.PubKeyList[:i], nm.PubKeyList[i+1:]...)

	}
}

//Handler for /nodemanager/getlistnodes
//mode: GET
//url: /nodemanager/getlistnodes
//return: array of PubKey
func nodeGetListNodesHandler(nm *nodemanager.NodeManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Get list nodes")
		// var list ListNodes
		// for _, PubKey := range nm.PubKeyList {
		// 	list.PubKey = append(list.PubKey, string(PubKey[:]))
		// }
		wh.SendJSON(w, nm.PubKeyList)

	}
}

//Handler for /nodemanager/gettransports
//mode: GET
//url: /nodemanager/gettransports?id=value
func nodeGetTransportsHandler(nm *nodemanager.NodeManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Get transport from Node")
		nodeID := r.FormValue("id")
		if nodeID == "" {
			wh.Error400(w, "Missing Node id")
			return
		}
		i, err := strconv.Atoi(nodeID)
		if err != nil {
			wh.Error400(w, "Node id must be integer")
			return
		}

		if len(nm.PubKeyList) < i {
			wh.Error400(w, "Invalid Node id")
			return
		}

		wh.SendJSON(w, nm.GetTransportsFromNode(i))

	}
}

//Handler for /nodemanager/addtransport
//mode: POST
//url: /nodemanager/addtransport
func nodeAddTransportHandler(nm *nodemanager.NodeManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Add transport to Node")

		var c ConfigWithID
		err := json.NewDecoder(r.Body).Decode(&c)

		if err != nil {
			wh.Error400(w, "Error decoding config for transport")
		}
		if len(nm.PubKeyList) < c.NodeID {
			wh.Error400(w, "Invalid Node id")
			return
		}

		node := nm.GetNodeByIndex(c.NodeID)
		nodemanager.AddTransportToNode(node, c.Config)
	}
}

//Handler for /nodemanager/removetransport
//mode: POST
//url: /nodemanager/removetransport
func nodeRemoveTransportHandler(nm *nodemanager.NodeManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Remove transport from Node")

		var c TransportWithID
		err := json.NewDecoder(r.Body).Decode(&c)

		if err != nil {
			wh.Error400(w, "Error decoding config for transport")
		}
		if len(nm.PubKeyList) < c.NodeID {
			wh.Error400(w, "Invalid Node id")
			return
		}
		logger.Info(strconv.Itoa(c.NodeID))

		nm.RemoveTransportsFromNode(c.NodeID, c.Transport)

	}
}

//RegisterNodeManagerHandlers - create routes for NodeManager
func RegisterNodeManagerHandlers(mux *http.ServeMux, nm *nodemanager.NodeManager) {
	//

	//  Test  Will be assigned name if present.
	mux.HandleFunc("/test", testHandler(nm))

	//Route for start Node
	mux.HandleFunc("/nodemanager/start", nodeStartHandler(nm))

	//Route for stop Node
	mux.HandleFunc("/nodemanager/stop", nodeStopHandler(nm))

	//Route for get transports from Node
	mux.HandleFunc("/nodemanager/gettransports", nodeGetTransportsHandler(nm))

	//Route for add transport to Node
	mux.HandleFunc("/nodemanager/addtransport", nodeAddTransportHandler(nm))

	//Route for remove transport from Node
	mux.HandleFunc("/nodemanager/removetransport", nodeRemoveTransportHandler(nm))

	//Route for get list Nodes
	mux.HandleFunc("/nodemanager/getlistnodes", nodeGetListNodesHandler(nm))

}

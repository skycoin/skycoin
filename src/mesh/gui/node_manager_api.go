package gui

import (
	//"errors"
	"fmt"
	//"log"
	"net/http"
	//"os"
	"encoding/json"

	"strconv"
	//"strings"
	//"github.com/skycoin/skycoin/src/cipher"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers

	"github.com/skycoin/skycoin/src/mesh/nodemanager"
	"github.com/skycoin/skycoin/src/mesh/transport"
)

//struct for nodeAddTransportHandler
type ConfigWithId struct {
	Id     int
	Config nodemanager.TestConfig
}

//struct for nodeRemoveTransportHandler
type TransportWithId struct {
	Id        int
	Transport transport.ITransport
}

//struct for nodeGetListNodesHandler
// type ListNodes struct {
// 	PubKey []string
// }

func testHandler(nm *nodemanager.NodeManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		wh.Error400(w, fmt.Sprintf("Works!"))

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
		id := r.FormValue("id")
		if id == "" {
			wh.Error400(w, "Missing Node id")
			return
		}
		i, err := strconv.Atoi(id)
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
		id := r.FormValue("id")
		if id == "" {
			wh.Error400(w, "Missing Node id")
			return
		}
		i, err := strconv.Atoi(id)
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

		var c ConfigWithId
		err := json.NewDecoder(r.Body).Decode(&c)

		if err != nil {
			wh.Error400(w, "Error decoding config for transport")
		}
		if len(nm.PubKeyList) < c.Id {
			wh.Error400(w, "Invalid Node id")
			return
		}

		nm.AddTransportsToNode(c.Config, c.Id)

	}
}

//Handler for /nodemanager/removetransport
//mode: POST
//url: /nodemanager/removetransport
func nodeRemoveTransportHandler(nm *nodemanager.NodeManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Remove transport from Node")

		var c TransportWithId
		err := json.NewDecoder(r.Body).Decode(&c)

		if err != nil {
			wh.Error400(w, "Error decoding config for transport")
		}
		if len(nm.PubKeyList) < c.Id {
			wh.Error400(w, "Invalid Node id")
			return
		}
		logger.Info(strconv.Itoa(c.Id))

		nm.RemoveTransportsFromNode(c.Id, c.Transport)

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

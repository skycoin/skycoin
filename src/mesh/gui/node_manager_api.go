package gui

import (
	//"errors"
	"fmt"
	//"log"
	"net/http"
	//"os"
	//	"encoding/json"
	"strconv"
	//"strings"

	//"github.com/skycoin/skycoin/src/cipher"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers

	"github.com/skycoin/skycoin/src/mesh/nodemanager"
)

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

//Handler for /nodemanager/gettransport
//mode: GET
//url: /nodemanager/gettransport?id=value
func nodeGetTransportHandler(nm *nodemanager.NodeManager) http.HandlerFunc {
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

//RegisterNodeManagerHandlers - create routes for NodeManager
func RegisterNodeManagerHandlers(mux *http.ServeMux, nm *nodemanager.NodeManager) {
	//

	//  Test  Will be assigned name if present.
	mux.HandleFunc("/test", testHandler(nm))

	//Route for start Node
	mux.HandleFunc("/nodemanager/start", nodeStartHandler(nm))

	//Route for stop Node
	mux.HandleFunc("/nodemanager/stop", nodeStopHandler(nm))

	//Route for get transport from Node
	mux.HandleFunc("/nodemanager/gettransport", nodeGetTransportHandler(nm))

}

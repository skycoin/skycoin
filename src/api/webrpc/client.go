package webrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Do send request to web
func Do(req *Request, rpcAddress string) (*Response, error) {
	d, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	rsp, err := http.Post(fmt.Sprintf("http://%s/webrpc", rpcAddress), "application/json", bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	res := Response{}
	if err := json.NewDecoder(rsp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

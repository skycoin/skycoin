package gui

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"encoding/json"
	"net/http/httptest"

	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
)

func TestHealthCheckHandler(t *testing.T) {
	unspent := uint64(10)
	unconf := uint64(20)

	unconfirmed := []visor.UnconfirmedTxn{{}, {}}
	connections := &daemon.Connections{
		Connections: []*daemon.Connection{
			{},
			{},
			{},
		},
	}
	metadata := &visor.BlockchainMetadata{
		Unspents:    unspent,
		Unconfirmed: unconf,
	}
	version := "1.0.0"
	commit := "abcdef"

	buildInfo := visor.BuildInfo{
		Version: version,
		Commit:  commit,
	}

	gateway := NewGatewayerMock()
	gateway.On("GetAllUnconfirmedTxns", mock.Anything).Return(unconfirmed)
	gateway.On("GetConnections", mock.Anything).Return(connections)
	gateway.On("GetBlockchainMetadata", mock.Anything).Return(metadata)
	gateway.On("GetBuildInfo", mock.Anything).Return(buildInfo)

	endpoint := "/health"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)

	if err != nil {
		t.Error(err)
		return
	}

	rr := httptest.NewRecorder()
	handler := NewServerMux(configuredHost, ".", gateway, &CSRFStore{})
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Wrong response code expected %d actual %d", http.StatusOK, rr.Code)
	}

	resp := &HealthResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), resp)

	if err != nil {
		t.Error(err)
	}

	if len(resp.VersionData.Version) == 0 {
		t.Errorf("Empty version data")
	}

	if resp.UnconfirmedTxCount != len(unconfirmed) {
		t.Errorf("Wrong count of unconfirmed tx expected %d actual %d",
			len(unconfirmed), resp.UnconfirmedTxCount)
	}

	if resp.OpenConnectionCount != len(connections.Connections) {
		t.Errorf("Wrong connection count expected %d actual %d",
			len(connections.Connections), resp.OpenConnectionCount)
	}

	if resp.BlockChainMetadata.Unconfirmed != unconf {
		t.Errorf("Wrong blockchain metadata unconfirmed expected %d actual %d",
			unconf, resp.BlockChainMetadata.Unconfirmed)
	}

	if resp.BlockChainMetadata.Unspents != unspent {
		t.Errorf("Wrong blockchain metadata unspent expected %d actual %d",
			unspent, resp.BlockChainMetadata.Unspents)
	}

	if resp.VersionData.Commit != commit || resp.VersionData.Version != version {
		t.Errorf("Wrong build info expected version %s commit %s actual version %s commit %s",
			version, commit, resp.VersionData.Version, resp.VersionData.Commit)
	}
}

package commands

import (
	"encoding/json"
	"errors"

	"github.com/skycoin/skycoin/cmd/skydex-client/skydexui"
	"github.com/skycoin/skycoin/src/skydex/protocol"
)

// uiMarket is the session as the server rendered client sees it.
type uiMarket struct{ s *session }

func (m uiMarket) Default() string { return m.s.defaultPK }

func (m uiMarket) Status() (bool, string, string) { return m.s.status() }

func (m uiMarket) Connect(ref string) (string, error) {
	_, name, err := m.s.connect(ref)
	return name, err
}

func (m uiMarket) Disconnect() { m.s.disconnect() }

func (m uiMarket) Call(msgType string, req, resp any) error {
	env, err := m.s.do(msgType, req)
	if errors.Is(err, errNotConnected) {
		return skydexui.ErrNotConnected
	}
	if err != nil {
		return err
	}
	if env.IsError() {
		var e protocol.ErrorData
		_ = env.Bind(&e) //nolint:errcheck
		return &skydexui.MarketError{Code: e.Code, Message: e.Message}
	}
	if resp == nil || len(env.Data) == 0 {
		return nil
	}
	return json.Unmarshal(env.Data, resp)
}

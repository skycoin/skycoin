package daemon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPeerBlockchainHeights(t *testing.T) {
	p := newPeerBlockchainHeights()

	addr1 := "127.0.0.1:1234"
	addr2 := "127.0.0.1:5678"
	addr3 := "127.0.0.1:9999"

	require.Empty(t, p.heights)
	p.Remove(addr1)
	require.Empty(t, p.heights)

	e := p.Estimate(1)
	require.Equal(t, uint64(1), e)

	e = p.Estimate(13)
	require.Equal(t, uint64(13), e)

	p.Record(addr1, 10)
	require.Len(t, p.heights, 1)

	records := p.All()
	require.Len(t, records, 1)
	require.Equal(t, PeerBlockchainHeight{
		Address: addr1,
		Height:  10,
	}, records[0])

	p.Record(addr1, 11)
	require.Len(t, p.heights, 1)

	records = p.All()
	require.Len(t, records, 1)
	require.Equal(t, PeerBlockchainHeight{
		Address: addr1,
		Height:  11,
	}, records[0])

	e = p.Estimate(1)
	require.Equal(t, uint64(11), e)

	e = p.Estimate(13)
	require.Equal(t, uint64(13), e)

	p.Record(addr2, 12)
	p.Record(addr3, 12)
	require.Len(t, p.heights, 3)

	records = p.All()
	require.Len(t, records, 3)
	require.Equal(t, []PeerBlockchainHeight{
		{
			Address: addr1,
			Height:  11,
		},
		{
			Address: addr2,
			Height:  12,
		},
		{
			Address: addr3,
			Height:  12,
		},
	}, records)

	e = p.Estimate(1)
	require.Equal(t, uint64(12), e)

	e = p.Estimate(13)
	require.Equal(t, uint64(13), e)

	p.Record(addr3, 24)
	e = p.Estimate(13)
	require.Equal(t, uint64(24), e)
}

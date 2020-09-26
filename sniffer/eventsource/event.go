package eventsource

import (
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
)

// EvtProviderPut should be emitted on every datastore Put() for a peer providing a CID.
type EvtProviderPut struct {
	CID    cid.Cid
	PeerID peer.ID
}

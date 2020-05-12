package filters

import (
	t "github.com/ipfs-search/ipfs-search/types"
	"log"
	"time"
)

type lastSeenFilter struct {
	resources  map[t.Resource]time.Time
	Expiration time.Duration
	PruneLen   int
}

func LastSeenFilter(expiration time.Duration, pruneLen int) *lastSeenFilter {
	// Allocate memory for pruneLen+1
	r := make(map[t.Resource]time.Time, pruneLen+1)

	return &lastSeenFilter{
		Expiration: expiration,
		PruneLen:   pruneLen,
		resources:  r,
	}
}

func (f *lastSeenFilter) prune() {
	if len(f.resources) > f.PruneLen {
		// Delete all expired items
		now := time.Now()
		cnt := 0

		for i, t := range f.resources {
			if now.Sub(t) > f.Expiration {
				delete(f.resources, i)
				cnt++
			}
		}

		log.Printf("Pruned %d resources, len: %d, pruneLen: %d", cnt, len(f.resources), f.PruneLen)
	}
}

func (f *lastSeenFilter) Filter(p t.Provider) (bool, error) {
	f.prune()

	lastSeen, present := f.resources[*(p.Resource)]

	if !present {
		// Not present, add it!
		log.Printf("Adding LastSeen: %v, len: %d", p, len(f.resources))
		f.resources[*(p.Resource)] = p.Date

		// Index it!
		return true, nil
	}

	if p.Date.Sub(lastSeen) > f.Expiration {
		// Last seen longer than expiration ago, update last seen.
		log.Printf("Updating LastSeen: %v, len: %d", p, len(f.resources))
		f.resources[*(p.Resource)] = p.Date

		// Index it!
		return true, nil
	}

	// Too recent, don't index
	log.Printf("Filtering recent %v, LastSeen %s", p, lastSeen)
	return false, nil
}

package gourlhaus

import (
	"context"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	// Check to make sure we can cache all the entries
	cache := NewCache()
	entries, err := cache.GetEntries(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if len(entries) < 1000 {
		// We should have gotten more entries
		t.Errorf("we should have more entries in the cache")
	}

	// Check that we can fetch again in < 1 millisecond (using cache)
	startTime := time.Now()
	entries, err = cache.GetEntries(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if time.Now().Sub(startTime) > time.Millisecond {
		t.Errorf("Fetched too slowly")
	}
}

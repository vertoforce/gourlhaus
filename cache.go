package gourlhaus

import (
	"context"
	"time"
)

// Cache is a cache of the urlhaus database.  It stores all the entires in memory
// Note that this will likely take about 100MB - 300MB of memory.
type Cache struct {
	// How often the cache will update with URLHaus
	// Do not set this less than 5min as the URLHaus endpoints only updates every 5 minutes max
	UpdateInterval time.Duration
	LastUpdated    time.Time
	// FetchHashes Set to true means we will call FillInURLHashDetails on the list stored in memory
	FetchHashes bool

	CachedEntries []URLEntry
}

// GetEntries Gets the urlhaus entries from the cache, updating with the latest results from urlhaus if need to sync again.
// We will sync if time.Now() - cache.LastUpdated > UpdateInterval
func (c *Cache) GetEntries(ctx context.Context) ([]URLEntry, error) {
	if c.CachedEntries == nil || time.Now().Sub(c.LastUpdated) > c.UpdateInterval {
		// We should sync!
		var err error
		c.CachedEntries, err = GetAllURLs(ctx)
		if err != nil {
			return nil, err
		}

		if c.FetchHashes {

		}
		c.LastUpdated = time.Now()
	}

	return c.CachedEntries, nil
}

// NewCache returns a new cache with the default update interval of 20 minutes
func NewCache() *Cache {
	return &Cache{
		UpdateInterval: time.Minute * 20,
	}
}

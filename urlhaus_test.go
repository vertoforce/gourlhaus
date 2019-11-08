package gourlhaus

import (
	"context"
	"testing"
)

func TestGetURLs(t *testing.T) {
	ctx := context.Background()

	// Recent
	recentURLs, err := GetRecentURLs(ctx)
	if err != nil {
		t.Errorf("Error fetching recent URLs: " + err.Error())
	}
	good := false
	for range recentURLs {
		good = true
	}
	if !good {
		t.Errorf("No Recent URL returned")
	}

	// Online
	onlineURLs, err := GetAllOnlineURLs(ctx)
	if err != nil {
		t.Errorf("Error fetching online URLs: " + err.Error())
	}
	good = false
	for range onlineURLs {
		good = true
	}
	if !good {
		t.Errorf("No Recent URL returned")
	}

	// All
	allURLs, err := GetAllURLs(ctx)
	if err != nil {
		t.Errorf("GetAllURLs returned error: " + err.Error())
	}
	good = false
	for range allURLs {
		good = true
	}
	if !good {
		t.Errorf("No Recent URL returned")
	}
}

func TestGetURLsWithHashes(t *testing.T) {
	ctx := context.Background()

	entries, err := GetURLsWithHashes(ctx)
	if err != nil {
		t.Error(err)
	}

	// Check to make sure at least one item has a MD5, SHA256, or filetype
	for entry := range entries {
		if entry.MD5 != "" || entry.Filetype != "" || entry.SHA256 != "" {
			return
		}
	}

	t.Error("No entries got enriched with a hash")
}

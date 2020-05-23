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
	if len(recentURLs) < 1000 {
		t.Errorf("Did not get enough urls")
	}

	// Online
	onlineURLs, err := GetAllOnlineURLs(ctx)
	if err != nil {
		t.Errorf("Error fetching online URLs: " + err.Error())
	}
	if len(onlineURLs) < 1000 {
		t.Errorf("Did not get enough urls")
	}

	// All
	allURLs, err := GetAllURLs(ctx)
	if err != nil {
		t.Errorf("GetAllURLs returned error: " + err.Error())
	}
	if len(allURLs) < 1000 {
		t.Errorf("Did not get enough urls")
	}
}

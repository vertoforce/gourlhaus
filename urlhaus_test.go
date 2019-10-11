package urlhaus

import "testing"

func TestGetRecentURLs(t *testing.T) {
	recentURLs, err := GetRecentURLs()
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(recentURLs) == 0 {
		t.Errorf("No URLs found")
	}
}

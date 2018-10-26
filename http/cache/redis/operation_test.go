package redis

import "testing"

func TestHMSet(t *testing.T) {

	m := map[string]string{
		"title":  "Example2",
		"author": "Steve",
		"body":   "Map",
	}

	HMSet("testkey",m)
}

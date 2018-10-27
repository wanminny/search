package redis

import "testing"

func TestHMSet(t *testing.T) {

	m := map[string]interface{}{
		"condition":"cccc",
		"status":1,
	}

	HMSet("testkey",m)
}

func TestHGetAll(t *testing.T) {

	HGetAll("ebc929a0b119670ff67b233fd07c24e3")
	//HGetAll("testkey")

}
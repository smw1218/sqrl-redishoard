package redishoard

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	ssp "github.com/smw1218/sqrl-ssp"
)

// Hoard implements an ssp.Hoard using redis as
// a backing store
type Hoard struct {
	client redis.UniversalClient
}

// NewHoard creates a redis backed Hoard
func NewHoard(client redis.UniversalClient) *Hoard {
	return &Hoard{
		client: client,
	}
}

// Get implements ssp.Hoard
func (h *Hoard) Get(nut ssp.Nut) (*ssp.HoardCache, error) {
	data, err := h.client.Get(string(nut)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ssp.ErrNotFound
		}
		return nil, fmt.Errorf("redis nut lookup failed: %v", err)
	}
	return h.fromBytes(nut, data)
}

// GetAndDelete implements ssp.Hoard
func (h *Hoard) GetAndDelete(nut ssp.Nut) (*ssp.HoardCache, error) {
	ret, _ := h.client.TxPipelined(func(pipe redis.Pipeliner) error {
		pipe.Get(string(nut))
		pipe.Del(string(nut))
		return nil
	})
	for _, cmd := range ret {
		err := cmd.Err()
		if err != nil {
			if err == redis.Nil {
				return nil, ssp.ErrNotFound
			}
			return nil, fmt.Errorf("redis nut lookup failed: %v", err)
		}
	}
	stringCmd := ret[0].(*redis.StringCmd)
	data, err := stringCmd.Bytes()
	if err != nil {
		return nil, fmt.Errorf("redis HoardCache read failed: %v", err)
	}
	log.Printf("data: %v", string(data))
	return h.fromBytes(nut, data)
}

func (h *Hoard) fromBytes(nut ssp.Nut, data []byte) (*ssp.HoardCache, error) {
	hoardCache := &ssp.HoardCache{}
	err := json.Unmarshal(data, hoardCache)
	if err != nil {
		return nil, fmt.Errorf("can't decode HoardCache object for nut %v: %v", nut, err)
	}
	return hoardCache, nil
}

// Save implements ssp.Hoard
func (h *Hoard) Save(nut ssp.Nut, value *ssp.HoardCache, expiration time.Duration) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed json encoding HoardCache: %v", err)
	}
	return h.client.Set(string(nut), jsonBytes, expiration).Err()
}

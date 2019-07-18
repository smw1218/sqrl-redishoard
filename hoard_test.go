package redishoard

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	ssp "github.com/smw1218/sqrl-ssp"
)

// Redis must be installed locally on the default port to run these tests
func TestSave(t *testing.T) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{})
	h := NewHoard(client)

	var nut ssp.Nut = "blah"
	hoardCache := &ssp.HoardCache{
		State: "boom!",
	}
	err := h.Save(nut, hoardCache, time.Second)
	if err != nil {
		t.Fatalf("Failed saving hoard cache: %v", err)
	}

	val, err := h.Get(nut)
	if err != nil {
		t.Fatalf("Failed getting from Hoard: %v", err)
	}

	if val.State != "boom!" {
		t.Fatalf("Wrong value from Hoard: %#v", val)
	}
}

func TestGetAndDelete(t *testing.T) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{})
	h := NewHoard(client)

	var nut ssp.Nut = "GetAndDelete"
	hoardCache := &ssp.HoardCache{
		State: "boom!",
	}
	err := h.Save(nut, hoardCache, time.Second)
	if err != nil {
		t.Fatalf("Failed saving hoard cache: %v", err)
	}

	val, err := h.Get(nut)
	if err != nil {
		t.Fatalf("Failed Get from Hoard: %v", err)
	}

	if val.State != "boom!" {
		t.Fatalf("Wrong value from Hoard: %#v", val)
	}

	val, err = h.GetAndDelete(nut)
	if err != nil {
		t.Fatalf("Failed GetAndDelete from Hoard: %v", err)
	}

	if val.State != "boom!" {
		t.Fatalf("Wrong value from Hoard: %#v", val)
	}

	val, err = h.GetAndDelete(nut)
	if err != ssp.ErrNotFound {
		t.Fatalf("Should have been deleted but wasn't: %v", err)
	}

	if val != nil {
		t.Fatalf("Wrong value from Hoard: %#v", val)
	}

}

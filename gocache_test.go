package gocache_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/twharmon/gocache"
)

func TestGetSet(t *testing.T) {
	db := gocache.New[string, int]()
	want := 5
	db.Set("foo", want)
	got := db.Get("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestGetExpNoConfig(t *testing.T) {
	db := gocache.New[string, int]()
	want := 0
	db.Set("foo", 5, time.Nanosecond)
	time.Sleep(time.Microsecond)
	got := db.Get("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestDelete(t *testing.T) {
	db := gocache.New[string, int]()
	want := 0
	db.Set("foo", 5)
	db.Delete("foo")
	got := db.Get("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestHasTrue(t *testing.T) {
	db := gocache.New[string, int]()
	want := true
	db.Set("foo", 5)
	got := db.Has("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestHasFalse(t *testing.T) {
	db := gocache.New[string, int]()
	want := false
	db.Set("foo", 5)
	db.Delete("foo")
	got := db.Has("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestSize(t *testing.T) {
	db := gocache.New[string, int]()
	want := 1
	db.Set("foo", 5)
	got := db.Size()
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestKeys(t *testing.T) {
	db := gocache.New[string, int]()
	want := []string{"bar", "foo"}
	db.Set("foo", 5)
	db.Set("bar", 6)
	got := db.Keys()
	sort.Strings(got)
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestClear(t *testing.T) {
	db := gocache.New[string, int]()
	want := 0
	db.Set("foo", 5)
	db.Set("bar", 6)
	db.Clear()
	got := db.Get("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestConfigCapacityHint(t *testing.T) {
	db := gocache.New[string, int](gocache.NewConfig().WithCapacityHint(1))
	want := 5
	db.Set("foo", want)
	got := db.Get("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestConfigTTLHasHit(t *testing.T) {
	db := gocache.New[string, int](gocache.NewConfig().WithDefaultEvictionPolicy(gocache.NewEvictionPolicy(time.Second)))
	want := true
	db.Set("foo", 5)
	got := db.Has("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestConfigTTLHasExp(t *testing.T) {
	db := gocache.New[string, int](gocache.NewConfig().WithDefaultEvictionPolicy(gocache.NewEvictionPolicy(time.Nanosecond)))
	want := false
	db.Set("foo", 5)
	time.Sleep(time.Microsecond)
	got := db.Has("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestConfigTTLGetExp(t *testing.T) {
	db := gocache.New[string, int](gocache.NewConfig().WithDefaultEvictionPolicy(gocache.NewEvictionPolicy(time.Nanosecond)))
	want := 0
	db.Set("foo", 5)
	time.Sleep(time.Microsecond)
	got := db.Get("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestConfigTTLGetHit(t *testing.T) {
	db := gocache.New[string, int](gocache.NewConfig().WithDefaultEvictionPolicy(gocache.NewEvictionPolicy(time.Second)))
	want := 5
	db.Set("foo", want)
	got := db.Get("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestConfigMaxCapacityMiss(t *testing.T) {
	db := gocache.New[string, int](gocache.NewConfig().WithMaxCapacity(1))
	want := 0
	db.Set("foo", 1)
	db.Set("bar", 2)
	got := db.Get("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestConfigTTLUpdateOnGet(t *testing.T) {
	ep := gocache.NewEvictionPolicy(time.Microsecond * 100).UpdateOnGet()
	fmt.Println(ep)
	cfg := gocache.NewConfig().WithDefaultEvictionPolicy(ep)
	db := gocache.New[string, int](cfg)
	want := 5
	db.Set("foo", want)
	start := time.Now()
	for {
		db.Get("foo")
		if time.Since(start).Microseconds() > 10 {
			break
		}
	}
	got := db.Get("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestConfigTTLUpdateOnHas(t *testing.T) {
	ep := gocache.NewEvictionPolicy(time.Microsecond * 100).UpdateOnHas()
	fmt.Println(ep)
	cfg := gocache.NewConfig().WithDefaultEvictionPolicy(ep)
	db := gocache.New[string, int](cfg)
	want := 5
	db.Set("foo", want)
	start := time.Now()
	for {
		db.Has("foo")
		if time.Since(start).Microseconds() > 10 {
			break
		}
	}
	got := db.Get("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestConfigMaxCapacityHit(t *testing.T) {
	db := gocache.New[string, int](gocache.NewConfig().WithMaxCapacity(2))
	want := 1
	db.Set("foo", want)
	db.Set("bar", 2)
	got := db.Get("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

func TestConfigMaxCapacityTTLMiss(t *testing.T) {
	db := gocache.New[string, int](gocache.NewConfig().WithMaxCapacity(2))
	want := 1
	db.Set("foo", want)
	db.Set("bar", 2, time.Second)
	db.Set("baz", 3)
	got := db.Get("foo")
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v; got %v", want, got)
	}
}

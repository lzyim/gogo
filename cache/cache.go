package cache

import (
	"fmt"
	"github.com/soloslee/gogo/hash"
	"sync"
	"time"
)

const (
	NoExpiration int64 = -1
)

type goStruct struct {
	key    string
	value  interface{}
	expire int64
	mux    sync.RWMutex
}

type dLinkList struct {
	data       *goStruct
	prev, next *dLinkList
}

type hashTable struct {
	tableSize int16
	numOfEle  int16
	arHash    [1024]*dLinkList
}

func newDLinkList(key string, value interface{}, expire int64) *dLinkList {
	gs := &goStruct{key: key, value: value, expire: expire}
	return &dLinkList{gs, nil, nil}
}

func (gs *goStruct) getVal() (interface{}, int64) {
	gs.mux.RLock()
	defer gs.mux.RUnlock()
	return gs.value, gs.expire
}

func (gs *goStruct) setVal(value interface{}) {
	gs.mux.Lock()
	defer gs.mux.Unlock()
	gs.value = value
}

func (h *hashTable) getGs(key string) *goStruct {
	index := hash.HashStr(key)
	if h.arHash[index] == nil {
		return nil
	}
	list := h.arHash[index]
	for ; ; list = list.next {
		if list.data.key == key {
			return list.data
		}
		if list.next == nil {
			break
		}
	}
	return nil
}

func (h *hashTable) Incr(key string, n int64) error {
	gs := h.getGs(key)
	if gs == nil {
		return nil
	}
	gs.mux.Lock()
	switch gs.value.(type) {
	case int:
		gs.value = gs.value.(int) + int(n)
	case int8:
		gs.value = gs.value.(int8) + int8(n)
	case int16:
		gs.value = gs.value.(int16) + int16(n)
	case int32:
		gs.value = gs.value.(int32) + int32(n)
	case int64:
		gs.value = gs.value.(int64) + n
	case uint:
		gs.value = gs.value.(uint) + uint(n)
	case uintptr:
		gs.value = gs.value.(uintptr) + uintptr(n)
	case uint8:
		gs.value = gs.value.(uint8) + uint8(n)
	case uint16:
		gs.value = gs.value.(uint16) + uint16(n)
	case uint32:
		gs.value = gs.value.(uint32) + uint32(n)
	case uint64:
		gs.value = gs.value.(uint64) + uint64(n)
	case float32:
		gs.value = gs.value.(float32) + float32(n)
	case float64:
		gs.value = gs.value.(float64) + float64(n)
	default:
		gs.mux.Unlock()
		return fmt.Errorf("The value for %s is not an integer", key)
	}
	gs.mux.Unlock()
	return nil
}

func (h *hashTable) Decr(key string, n int64) error {
	gs := h.getGs(key)
	if gs == nil {
		return nil
	}
	gs.mux.Lock()
	switch gs.value.(type) {
	case int:
		gs.value = gs.value.(int) - int(n)
	case int8:
		gs.value = gs.value.(int8) - int8(n)
	case int16:
		gs.value = gs.value.(int16) - int16(n)
	case int32:
		gs.value = gs.value.(int32) - int32(n)
	case int64:
		gs.value = gs.value.(int64) - n
	case uint:
		gs.value = gs.value.(uint) - uint(n)
	case uintptr:
		gs.value = gs.value.(uintptr) - uintptr(n)
	case uint8:
		gs.value = gs.value.(uint8) - uint8(n)
	case uint16:
		gs.value = gs.value.(uint16) - uint16(n)
	case uint32:
		gs.value = gs.value.(uint32) - uint32(n)
	case uint64:
		gs.value = gs.value.(uint64) - uint64(n)
	case float32:
		gs.value = gs.value.(float32) - float32(n)
	case float64:
		gs.value = gs.value.(float64) - float64(n)
	default:
		gs.mux.Unlock()
		return fmt.Errorf("The value for %s is not an integer", key)
	}
	gs.mux.Unlock()
	return nil
}

func New() *hashTable {
	return &hashTable{1024, 0, [1024]*dLinkList{}}
}

func (h *hashTable) Set(key string, value interface{}, second int64) {
	index := hash.HashStr(key)
	var expire int64
	if second == -1 {
		expire = -1
	} else {
		expire = second + time.Now().Unix()
	}
	if h.arHash[index] == nil {
		h.arHash[index] = newDLinkList(key, value, expire)
		h.numOfEle++
	} else {
		list := h.arHash[index]
		found := false
		for ; ; list = list.next {
			if list.data.key == key {
				list.data.setVal(value)
				found = true
			}
			if list.next == nil {
				break
			}
		}
		if !found {
			newItem := newDLinkList(key, value, expire)
			list.next = newItem
			newItem.prev = list
			h.numOfEle++
		}
	}
}

func (h *hashTable) Get(key string) (interface{}, bool) {
	index := hash.HashStr(key)
	list := h.arHash[index]
	found := false
	if list == nil {
		return nil, false
	} else {
		for ; ; list = list.next {
			if list.data.key == key {
				found = true
				break
			}
			if list.next == nil {
				break
			}
		}
		if found {
			val, exp := list.data.getVal()
			if exp < 0 {
				return val, true
			}
			if time.Now().Unix() > exp {
				h.Del(key)
				return nil, false
			}
			return val, true
		}
		return nil, false
	}
}

func (h *hashTable) Del(key string) {
	index := hash.HashStr(key)
	if h.arHash[index] == nil {
		return
	}
	list := h.arHash[index].next
	if h.arHash[index].data.key == key {
		if h.arHash[index].next != nil {
			h.arHash[index] = h.arHash[index].next
		} else {
			h.arHash[index] = nil
		}
		return
	}
	for ; ; list = list.next {
		if list.data.key == key {
			if list.prev != nil && list.next != nil {
				list.prev.next = list.next
				list.next.prev = list.prev
			} else if list.prev != nil && list.next == nil {
				list.prev.next = nil
			}
			break
		}
		if list.next == nil {
			break
		}
	}
}

func (h *hashTable) Count() int16 {
	return h.numOfEle
}

func (h *hashTable) Flush() {
	*h = hashTable{1024, 0, [1024]*dLinkList{}}
}

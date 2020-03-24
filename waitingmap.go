package waitingmap

import (
	"sync"
	"time"
)

// MAP实现key不存在 get操作等待 直到key存在或者超时，保证并发安全
type (
	Map struct {
		entry map[string]*entry
		rmx   sync.RWMutex
	}
	entry struct {
		ch      chan bool
		value   interface{}
		isExist bool
	}
)

//创建map实例
func NewMap() *Map {
	return &Map{
		entry: make(map[string]*entry),
	}
}

//存入key/val。此方法不会阻塞，时刻都可以执行
func (m *Map) Wt(key string, val interface{}) {
	m.rmx.Lock()
	defer m.rmx.Unlock()
	if e, ok := m.entry[key]; ok {
		//如果key存在，直接赋值，并通知订阅协程
		e.value = val
		e.isExist = true
		if len(e.ch) < 1 {
			e.ch <- true
		}
	} else {
		//如果key不存在，创建entry实例并赋值
		e = &entry{
			ch:      make(chan bool, 1),
			value:   val,
			isExist: true,
		}
		m.entry[key] = e
	}
}

//读取一个key，如果key不存在阻塞，等待key存在或者超时
func (m *Map) Rd(key string, timeout time.Duration) interface{} {
	m.rmx.Lock()
	if e, ok := m.entry[key]; ok && e.isExist {
		//如果key存在且有值，直接返回
		m.rmx.Unlock()
		return e.value

	} else if ok {
		//如果key存在，但无值（Rd创建），等待并返回
		m.rmx.Unlock()
		select {
		case <-e.ch:
			return e.value
		case <-time.After(timeout):
			return nil
		}

	} else {
		//如果key不存在，创建entry实例，等待并返回
		e = &entry{
			ch:      make(chan bool, 1),
			isExist: false,
		}
		m.entry[key] = e
		m.rmx.Unlock()
		select {
		case <-e.ch:
			return e.value
		case <-time.After(timeout):
			return nil
		}
	}
}

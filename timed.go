package genh

import (
	"sync"
	"sync/atomic"
	"time"
)

type tmEle[V any] struct {
	sync.RWMutex
	la atomic.Int64 // last access / read
	v  V
	t  *time.Timer
}

func (e *tmEle[V]) expired(expiry time.Duration) bool {
	if expiry < 1 {
		return false
	}
	return time.Since(time.Unix(0, e.la.Load())) > expiry
}

type TimedMap[K comparable, V any] struct {
	m LMap[K, *tmEle[V]]
}

func (tm *TimedMap[K, V]) Set(k K, v V, timeout time.Duration) {
	ele := &tmEle[V]{v: v}
	ele.la.Store(time.Now().UnixNano())
	if timeout > 0 {
		ele.t = time.AfterFunc(timeout, func() { tm.deleteEle(k, ele) })
	}
	tm.m.Set(k, ele)
}

func (tm *TimedMap[K, V]) SetUpdateFn(k K, vfn func() V, updateEvery time.Duration) {
	tm.SetUpdateExpireFn(k, vfn, updateEvery, -1)
}

func (tm *TimedMap[K, V]) SetUpdateExpireFn(k K, vfn func() V, updateEvery, expireIfNotAccessedFor time.Duration) {
	ele := &tmEle[V]{v: vfn()}
	ele.la.Store(time.Now().UnixNano())
	tm.m.Set(k, ele)
	if updateEvery < time.Millisecond {
		panic("every must be >= time.Millisecond")
	}
	go func() {
		for {
			time.Sleep(updateEvery)
			if tm.m.Get(k) != ele {
				return
			}
			v := vfn()
			ele.Lock()
			if ele.expired(expireIfNotAccessedFor) {
				tm.deleteEle(k, ele)
				ele.Unlock()
				return
			}
			ele.v = v
			ele.Unlock()
		}
	}()
}

func (tm *TimedMap[K, V]) Get(k K) (v V) {
	v, _ = tm.GetOk(k)
	return
}

func (tm *TimedMap[K, V]) GetOk(k K) (v V, ok bool) {
	ele := tm.m.Get(k)
	if ok = ele != nil; ok {
		now := time.Now().UnixNano()
		ele.RLock()
		defer ele.RUnlock()
		v = ele.v
		ele.la.Store(now)
	}
	return
}

func (tm *TimedMap[K, V]) DeleteGet(k K) (v V, ok bool) {
	ele := tm.m.DeleteGet(k)
	if ok = ele != nil; ok {
		ele.RLock()
		defer ele.RUnlock()
		v = ele.v
		if ele.t != nil {
			ele.t.Stop()
		}
	}
	return
}

func (tm *TimedMap[K, V]) Delete(k K) {
	ele := tm.m.DeleteGet(k)
	if ele != nil {
		if ele.t != nil {
			ele.t.Stop()
		}
	}
}

func (tm *TimedMap[K, V]) deleteEle(k K, ele *tmEle[V]) {
	tm.m.Update(func(m map[K]*tmEle[V]) {
		if m[k] == ele {
			delete(m, k)
		}
	})
}

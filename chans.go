package genh

func SliceToChan[T any](s []T, cap int) <-chan T {
	if cap == 0 {
		cap = 1
	}
	ch := make(chan T, cap)
	go func() {
		for _, v := range s {
			ch <- v
		}
		close(ch)
	}()
	return ch
}

func ChanToSlice[T any](s <-chan T) []T {
	out := make([]T, cap(s))
	for v := range s {
		out = append(out, v)
	}
	return Clip(out)
}

func ClosedChan[T any]() chan T {
	ch := make(chan T)
	close(ch)
	return ch
}

func Chan[T any](cap int) (ch <-chan T, pushFn func(T) bool, closeFn func()) {
	var closed AtomicInt32
	done := make(chan struct{})
	rch := make(chan T, cap)
	pushFn = func(v T) bool {
		select {
		case <-done:
			close(rch)
			return false
		case rch <- v:
			return true
		}
	}
	closeFn = func() {
		if closed.Swap(1) == 0 {
			close(done)
		}
	}
	ch = rch
	return
}

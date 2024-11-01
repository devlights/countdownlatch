package countdownlatch

import (
	"sync"
	"sync/atomic"
)

// Latch は、C#のCountdownEventやJavaのCountDownLatchと同様の機能を提供する構造体です.
type Latch struct {
	count atomic.Int32
	mutex sync.Mutex
	cond  *sync.Cond
}

// New は、指定されたカウント数でLatchを初期化します.
func New(initialCount int) *Latch {
	if initialCount < 0 {
		panic("初期カウントは0以上である必要があります")
	}

	var (
		latch Latch
	)
	latch.count.Store(int32(initialCount))
	latch.cond = sync.NewCond(&latch.mutex)

	return &latch
}

// Signal は、カウントを1減らします.
// 戻り値として、カウントダウンが満了したかどうかを返します.
func (me *Latch) Signal() bool {
	return me.SignalCount(1)
}

// SignalCount は、指定された数だけカウントを減らします.
// 戻り値として、カウントダウンが満了したかどうかを返します.
func (me *Latch) SignalCount(count int) bool {
	if count <= 0 {
		return false
	}

	me.mutex.Lock()
	defer me.mutex.Unlock()

	newCount := me.count.Add(-int32(count))
	if newCount <= 0 {
		me.cond.Broadcast()
		return true
	}

	return false
}

// Wait は、カウントが0になるまでブロックします.
func (me *Latch) Wait() {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	for me.count.Load() > 0 {
		me.cond.Wait()
	}
}

// CurrentCount は、現在のカウント値を返します.
func (me *Latch) CurrentCount() int {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	return int(me.count.Load())
}

// Reset は、カウントを指定された値にリセットします.
func (me *Latch) Reset(count int) {
	if count < 0 {
		panic("リセットカウントは0以上である必要があります")
	}

	me.mutex.Lock()
	defer me.mutex.Unlock()

	me.cond.Broadcast()

	me.count.Store(int32(count))
}

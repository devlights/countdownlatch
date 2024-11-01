package countdownlatch_test

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/devlights/countdownlatch"
)

func ExampleLatch() {
	const (
		numLatchs     = 3
		numGoroutines = 5
	)

	log.SetFlags(0)

	var (
		latch = countdownlatch.New(numLatchs)
	)
	for range 2 {
		var (
			wg sync.WaitGroup
		)

		latch.Reset(numLatchs)

		for range numGoroutines {
			wg.Add(1)
			go func() {
				defer wg.Done()

				fmt.Println("待機開始")
				latch.Wait()
				fmt.Println("待機解除")
			}()
		}

		for range numLatchs {
			<-time.After(100 * time.Millisecond)

			fmt.Printf("現在のカウント: %d\n", latch.CurrentCount())
			latch.Signal()
		}

		wg.Wait()
		fmt.Println("----------------")
	}

	// Output:
	// 待機開始
	// 待機開始
	// 待機開始
	// 待機開始
	// 待機開始
	// 現在のカウント: 3
	// 現在のカウント: 2
	// 現在のカウント: 1
	// 待機解除
	// 待機解除
	// 待機解除
	// 待機解除
	// 待機解除
	// ----------------
	// 待機開始
	// 待機開始
	// 待機開始
	// 待機開始
	// 待機開始
	// 現在のカウント: 3
	// 現在のカウント: 2
	// 現在のカウント: 1
	// 待機解除
	// 待機解除
	// 待機解除
	// 待機解除
	// 待機解除
	// ----------------
}

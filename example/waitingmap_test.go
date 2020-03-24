package example

import (
	"fmt"
	wmap "github.com/Pegasus219/waitingmap"
	"sync"
	"testing"
	"time"
)

func TestWaitmap(t *testing.T) {
	var wg sync.WaitGroup
	var wMap = wmap.NewMap()

	wg.Add(1)
	go func() {
		getVal := wMap.Rd("test", time.Millisecond*100)
		fmt.Println("get test value=", getVal)
		wg.Done()
	}()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			//time.Sleep(time.Millisecond * 200) //测试获取超时
			wMap.Wt("test", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

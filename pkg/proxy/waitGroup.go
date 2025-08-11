package proxy

import "sync"

var currWg = &sync.WaitGroup{}

func Wait() {
	currWg.Wait()
}

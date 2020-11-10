// Reference: https://gist.github.com/rushilgupta/228dfdf379121cb9426d5e90d34c5b96
package main

import (
	"container/heap"
	"math"
	"math/rand"
	"time"
)

type Request struct {
	data int
	resp chan float64
}

type Work struct {
	idx int
	wok chan Request
	pending int
}

type Pool []*Work

type Balancer struct {
	pool Pool
	done chan *Work
}

func (w *Work) doWork(done chan *Work) {
	for {
		req := <-w.wok
		req.resp <- math.Sin(float64(req.data))
		done <- w
	}
}

func InitBalancer() *Balancer {
	nWorker := 10
	nRequester := 100
	done := make(chan *Work, nWorker)
	b := &Balancer {
		make(Pool, 0, nWorker),
		done,
	}

	for i := 0; i < nWorker; i++ {
		w := &Work { wok: make(chan Request, nRequester) }
		heap.Push(&b.pool, w)
		go w.doWork(b.done)
	}

	return b
}

func createAndRequest(req chan Request) {
	resp := make(chan float64)

	for {
		time.Sleep(time.Duration(rand.Int63n(int64(time.Millisecond))))
		req <- Request { int(rand.Int31n(90)), resp }

		// resp 채널로 부터 데이터를 읽는다.
		// 채널에서 데이터를 읽기: [data] <- [link]
		// 채널로 데이터를 보내기: [link] <- [data]
		<- resp
	}
}

func main() {

}

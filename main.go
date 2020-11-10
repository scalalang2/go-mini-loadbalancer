// Reference: https://gist.github.com/rushilgupta/228dfdf379121cb9426d5e90d34c5b96
package main

import (
	"container/heap"
	"fmt"
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

// Heap 인터페이스 구현 (Duck Typing)
func (p Pool) Len() int { return len(p) }

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p *Pool) Swap(i, j int) {
	a := *p
	a[i], a[j] = a[j], a[i]
	a[i].idx = i
	a[j].idx = j
}

func (p *Pool) Push(x interface{}) {
	n := len(*p)
	item := x.(*Work) // type casting
	item.idx = n
	*p = append(*p, item)
}

func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	item.idx = -1
	*p = old[0:n-1]
	return item
}

// Mini-Loader Balancer 구현
func (b *Balancer) balance(req chan Request) {
	for {
		select {
			case request := <-req:
				b.dispatch(request)
			case w:= <-b.done:
				b.completed(w)
		}

		b.print()
	}
}

func (b *Balancer) dispatch(req Request) {
	w := heap.Pop(&b.pool).(*Work)
	w.wok <- req
	w.pending++
	heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *Work) {
	w.pending--
	heap.Remove(&b.pool, w.idx)
	heap.Push(&b.pool, w)
}

func (b *Balancer) print() {
	sum := 0
	sumsq := 0

	for _, w := range b.pool {
		fmt.Printf("%d ", w.pending)
		sum += w.pending
		sumsq += w.pending * w.pending
	}

	avg := float64(sum) / float64(len(b.pool))
	variance := float64(sumsq) / float64(len(b.pool)) - avg * avg
	fmt.Printf(" %.2f %.2f \n", avg, variance)
}

// Work 구현
func (w *Work) doWork(done chan *Work) {
	for {
		req := <-w.wok
		req.resp <- math.Sin(float64(req.data))
		done <- w
	}
}

func InitBalancer() *Balancer {
	nWorker := 10
	nRequester := 1000
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
	work := make(chan Request)
	for i:= 0; i < 1000; i++ {
		go createAndRequest(work)
	}
	InitBalancer().balance(work)
}

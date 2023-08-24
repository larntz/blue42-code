package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

const (
	rounds  = 100
	buckets = 100
	balls   = 1_000_000
)

func main() {
	var wg sync.WaitGroup

	rChan := make(chan int, rounds)
	bChan := make(chan int, rounds)
	for i := 0; i < rounds; i++ {
		go RandomFill(&wg, rChan)
		go BestOfTwo(&wg, bChan)
		wg.Add(2)
	}
	wg.Wait()
	close(rChan)
	close(bChan)

	fmt.Printf("Avg spread for RandomFill = %d\n", avg(rChan))
	fmt.Printf("Avg spread for BestOfTwo = %d\n", avg(bChan))
}

func avg(c chan int) int {
	total := 0
	for x := range c {
		total += x
	}
	return total / rounds
}

func BestOfTwo(wg *sync.WaitGroup, c chan int) {
	defer wg.Done()
	var b [buckets]int
	var b1, b2, max int
	min := math.MaxInt
	// fill buckets
	for i := 0; i < balls; i++ {
		for {
			b1 = rand.Intn(buckets)
			b2 = rand.Intn(buckets)
			if b1 != b2 {
				break
			}
		}

		if b[b1] < b[b2] {
			b[b1]++
		} else {
			b[b2]++
		}
	}

	// check buckets
	for i := 0; i < buckets; i++ {
		if b[i] < min {
			min = b[i]
		}
		if b[i] > max {
			max = b[i]
		}
	}
	c <- (max - min)
}

func RandomFill(wg *sync.WaitGroup, c chan int) {
	defer wg.Done()
	var b [buckets]int
	max := 0
	min := math.MaxInt

	// fill buckets
	for i := 0; i < balls; i++ {
		b[rand.Intn(buckets)]++
	}

	// check buckets
	for i := 0; i < buckets; i++ {
		if b[i] < min {
			min = b[i]
		}
		if b[i] > max {
			max = b[i]
		}
	}
	c <- (max - min)
}

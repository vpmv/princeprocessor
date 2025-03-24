package main

import (
	"log"
	"runtime"
	"sync"
)

var (
	Dictionary []string
	Separator  byte = ' '
	MaxWorkers      = runtime.NumCPU() * 2 //  Dynamically adjusts based on CPU
)

type StackItem struct {
	chainBuffer []byte
	pos         int
	depth       int
	used        []bool
}

// Generates chains efficiently (fast, low memory, proper spacing)
func generateChains(elemMin, elemMax int, out *outFile) {
	if len(Dictionary) == 0 {
		log.Fatal("Error: Wordlist is empty")
		return
	}

	var wg sync.WaitGroup
	writeMutex := &sync.Mutex{} // sync to control write buffer in parallel scope
	workQueue := make(chan int, len(Dictionary))

	//  Worker function processing a subset of words
	worker := func() {
		for i := range workQueue {
			chainBuilder(elemMin, elemMax, out, writeMutex, i)
			wg.Done()
		}
	}

	//  Start workers
	for i := 0; i < MaxWorkers; i++ {
		go worker()
	}

	//  Distribute work (each word starts a chain worker - i.e. goroutine)
	for idx := range Dictionary {
		wg.Add(1)
		workQueue <- idx
	}
	close(workQueue)

	wg.Wait()
}

func chainBuilder(elemMin, elemMax int, out *outFile, writeMutex *sync.Mutex, startIdx int) {
	stack := make([]StackItem, 0, 100)

	//  Start from the assigned `startIdx`
	word := Dictionary[startIdx]
	chainBuffer := make([]byte, 512)
	copy(chainBuffer, word)
	used := make([]bool, len(Dictionary))
	used[startIdx] = true
	stack = append(stack, StackItem{chainBuffer, len(word), 1, used})

	for len(stack) > 0 {
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if item.depth >= elemMin {
			output := string(item.chainBuffer[:item.pos]) + "\n"
			writeMutex.Lock()
			out.buf.WriteString(output)
			writeMutex.Unlock()
		}

		if item.depth >= elemMax {
			continue
		}

		for i := 0; i < len(Dictionary); i++ {
			if item.used[i] {
				continue
			}

			newUsed := make([]bool, len(Dictionary))
			copy(newUsed, item.used)
			newUsed[i] = true

			newBuffer := make([]byte, item.pos+len(Dictionary[i])+1)
			copy(newBuffer, item.chainBuffer[:item.pos])

			nextPos := item.pos
			if item.pos > 0 {
				newBuffer[nextPos] = Separator
				nextPos++
			}
			copy(newBuffer[nextPos:], Dictionary[i])

			stack = append(stack, StackItem{newBuffer, nextPos + len(Dictionary[i]), item.depth + 1, newUsed})
		}
	}
}

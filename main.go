package main

import (
	"Task17.7.1/pkg/counter"
	"context"
	"fmt"
	"log"
	"sync"
)

func worker(ctx context.Context, cancel context.CancelFunc, c *counter.Counter, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		c.Add(1, ctx, cancel)
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func main() {
	var wg sync.WaitGroup
	var wgChan sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	amountOfThreads := 0
	maxValue := 0
	fmt.Println("Укажите количество горутин: ")
	_, err := fmt.Scanln(&amountOfThreads)
	if err != nil {
		log.Fatalln("Неверное значение")
	}
	if amountOfThreads < 1 {
		log.Fatalln("Количество горутин не может быть меньше 1")
	}
	fmt.Println("Укажите максимальное значение счётчика: ")
	_, err = fmt.Scanln(&maxValue)
	if err != nil {
		log.Fatalln("Неверное значение")
	}
	if maxValue < 1 {
		log.Fatalln("Максимальное значение счтчика не может быть меньше 1")
	}
	c := counter.NewCounter(maxValue)
	for id := 0; id < amountOfThreads; id++ {
		wg.Add(1)
		go worker(ctx, cancel, c, &wg)
		fmt.Println(id)
	}
	go c.Increment(&wgChan, cancel)
	wgChan.Add(1)
	wg.Wait()
	c.CloseChanel()
	wgChan.Wait()
	fmt.Println("Counter: ", c.Value())
}

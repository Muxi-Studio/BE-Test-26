package main

import (
	"fmt"
	"grandmaster/opponent"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	moves := make(chan opponent.Move)

	wg.Add(5)

	for i := range 5 {
		go func(i int) {
			defer wg.Done()
			var op1 opponent.Move
			op := opponent.New(i)
			for {
				op1 = op.Move()
				moves <- op1
				if op1.Move == "pass" {
					return
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(moves)
	}()

	var i opponent.Move
	for i = range moves {
		if i.Move == "pass" {
			fmt.Printf("Opponent %d: %s, response: pass", i.Id, i.Move)
		}
		if i.Move == "occupy" {
			fmt.Printf("Opponent %d: %s, response: invasion", i.Id, i.Move)
		}
		if i.Move == "hane" {
			fmt.Printf("Opponent %d: %s, response: tsugi", i.Id, i.Move)
		}
		if i.Move == "tsuke" {
			fmt.Printf("Opponent %d: %s, response: hiki", i.Id, i.Move)
		}
		if i.Move == "tsuke-toshi" {
			fmt.Printf("Opponent %d: %s, response: osarru", i.Id, i.Move)
		}
		if i.Move == "kiri" {
			fmt.Printf("Opponent %d: %s, response: nobi", i.Id, i.Move)
		}
	}
}

package main

import (
	"fmt"
	"grandmaster/opponent"
	"sync"
)

var response = map[string]string{
	"pass":        "pass",
	"occupy":      "invasion",
	"hane":        "tsugi",
	"tsuke":       "hiki",
	"tsuke-toshi": "osarru",
	"kiri":        "nobi",
}

func main() {
	moves := make(chan opponent.Move)
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			op := opponent.New(i)
			moves <- op.Move()
		}(i)
	}

	go func() {
		wg.Wait()
		close(moves)
	}()
	for move := range moves {
		fmt.Printf("Opponent %d: %s, response: %s\n", move.Id, move.Move, response[move.Move])
	}
}

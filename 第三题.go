package main

import (
	"fmt"
	"grandmaster/opponent"
	"sync"
)

func main() {
	var response = map[string]string{
		"pass" : "pass",
		"occupy" : "invasion",
		"hane"   : "tsugi",
		"tsuke"  : "hiki",
		"tsuke-toshi" : "osarru",
		"kiri"   : "nibi",
	}
	var wg sync.WaitGroup
	moves := make(chan opponent.Move)

	for i := range 5 {
		go func(id int) {
			wg.Add(1)
			op := opponent.New(id)
			moves <- op.Move()
		if opponent.Move == "pass" {
			return
		}
			defer wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(moves)
	}()

	for i := 0; i < 5;i++ {
		move := <- moves
		fmt.Printf("Opponent %d: %s,response: %s",move.id,move.Move,response[move.Move])
	}
}

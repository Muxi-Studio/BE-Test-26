package main

import (
	"fmt"
	"grandmaster/opponent"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	moves := make(chan opponent.Move, 25)
	var op opponent.Opponent

	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			op = opponent.New(id)
			for op.Move().Move != "pass" {
				moves <- op.Move()
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(moves)
	}()

	wg.Wait()
	for v := range moves {
		var statu string
		if v.Move == "pass" {
			statu = "pass"
		} else if v.Move == "occupy" {
			statu = "invasion"
		} else if v.Move == "hane" {
			statu = "tsugi"
		} else if v.Move == "tsuke" {
			statu = "hiki"
		} else if v.Move == "tsuke-toshi" {
			statu = "osarru"
		} else if v.Move == "kiri" {
			statu = "nobi"
		}
		//我是真没招了朋友

		fmt.Print("Opponent %d: %s, response: %s", v.Id, v.Move, statu)
	}
}

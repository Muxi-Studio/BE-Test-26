package main

import (
	"fmt"
	"grandmaster/opponent"
	"sync"
)

func main() {
	/* A */
	var productWG sync.WaitGroup
	//var consumeWG sync.WaitGroup
	moves := make(chan opponent.Move, 20)

	productWG.Add(5)
	for i := range 5 {
		op := opponent.New(i)
		fmt.Println(op)
		go func(op opponent.Opponent) {
			//time.Sleep(999 * time.Millisecond)
			t := op.Move()
			//fmt.Printf("Opponent %d:%s.,respnse:", t.Id, t.Move)
			moves <- t
			if t.Move == "pass" {
				//fmt.Println(1)
				productWG.Done()
			}

		}(op)
	}
	go func() {
		productWG.Wait()
		close(moves)
	}()

	for id := range moves {
		fmt.Printf("Opponent %d:%s.,respnse:", id.Id, id.Move)
		if id.Move == "pass" {
			fmt.Println("pass")
		} else if id.Move == "occupy" {
			fmt.Println("invasion")
		} else if id.Move == "tsuke" {
			fmt.Println("hiki")
		} else if id.Move == "tsuke-toshi" {
			fmt.Println("osarru")
		} else if id.Move == "hane" {
			fmt.Println("tsugi")
		} else if id.Move == "kiri" {
			fmt.Println("nobi")
		}
		fmt.Println(1)
	}
}

package opponent

import (
	"math/rand"
	"time"
)

type Opponent struct {
	Id     int
	Steps  int
	Target int
}

func New(id int) Opponent {
	return Opponent{
		Id:     id,
		Steps:  0,
		Target: rand.Intn(5) + 10,
	}
}

type Move struct {
	Id   int
	Move string
}

func (op *Opponent) Move() Move {
	op.Steps++

	move := []string{
		"occupy",      // invasion
		"hane",        // tsugi
		"tsuke",       // hiki
		"tsuke-toshi", // osaeru
		"kiri",        // nobi
	}[rand.Intn(5)]

	if op.Steps >= op.Target {
		move = "pass"
	}

	thinkTime := rand.Intn(500)
	time.Sleep(time.Duration(thinkTime) * time.Millisecond)
	return Move{Id: op.Id, Move: move}
}

package main

import (
	"context"
	"errors"
)

type Msg struct {
	Apply string
	Reply string
}

type Answer struct {
	msg Msg
	ID  int
}

// TODO 设计消息和 Actor，替换下面的占位实现。
// 用一个有界缓冲 channel 做共享邮箱，用每次调用独有的 channel 回信。
// 每个应用实例启动一个 Actor，在它的处理循环中初始化并独占 activityStore。
func newActor() executor {
	

	return pendingExecutor{}
}

type pendingExecutor struct {
	Statu   string
	Receive string
	
}

func (pendingExecutor) ask(context.Context, operation) (ActivityResp, error) {
	reply := make(chan Answer, 1)
	rep := newActor()

	mailbox := make(map[int]Answer)
	mailbox[rep.ID] = <- reply
	return ActivityResp{}, errors.New("请实现单 Actor 的投递、处理与回信")
}

func (pendingExecutor) shutdown() {}

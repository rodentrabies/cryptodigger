package main

import "strings"

type Event interface {
	Description() string
	Consequence(coins int) int
}

type simpleEvent struct {
	str string
	fun func(int) int
}

func (event simpleEvent) Description() string {
	return event.str
}

func (event simpleEvent) Consequence(param int) int {
	return event.fun(param)
}

func NewEvent(str string, inc int) Event {
	return &simpleEvent{
		str: str,
		fun: func(param int) int {
			return param + inc
		},
	}
}

func NewRandomEvent() Event {
	str := strings.Join(
		[]string{
			"Buteric Vitalin just deleted you Epherium deposit",
			"and kept going, so you lost 10 coins...",
		}, "\n",
	)
	return NewEvent(str, -10)
}

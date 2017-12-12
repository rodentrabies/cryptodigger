package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

//------------------------------------------------------------------------------
// Event declarations
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
	sentence := []string{}
	for _, parts := range [][]string{subjects, verbs, objects, consequences} {
		sentence = append(sentence, parts[rand.Intn(len(parts))])
	}
	impact, value := randomImpact(impacts)
	fullText := fmt.Sprintf("%s and %s", strings.Join(sentence, " "), impact)
	return NewEvent(reformatText(fullText), value)
}

func reformatText(text string) string {
	words := regexp.MustCompile("\\s").Split(text, -1)
	result, line, spaces := "", "", 0
	for _, word := range words {
		if len(line)+len(word) < LineLength {
			line = line + word + " "
			spaces++
		} else {
			n := (LineLength-len(line))/spaces + 1
			result += strings.Replace(line, " ", strings.Repeat(" ", n), -1) + "\n"
			line, spaces = word+" ", 0
		}
	}
	return result + line + "\n"
}

type impactDescriptor struct {
	desc string
	coef int
}

func randomImpact(impacts []impactDescriptor) (string, int) {
	impact := impacts[rand.Intn(len(impacts))]
	value := rand.Intn(10) + 2
	return fmt.Sprintf(impact.desc, value), impact.coef * value
}

//------------------------------------------------------------------------------
// Event generator data
var subjects = []string{
	"Famous actor Battlefield Counterstrike",
	"The creator of Ephemerium, Bitalic Buttcoin",
	"Montero's CEO @cuddlybeetle",
	"Self-proclaimed Bitcoin creator Sashimi Fakamoto",
}

var verbs = []string{
	"endorsed a",
	"decided to take part in a",
	"started a new cryptocurrency",
}

var objects = []string{
	"Bitcoin ICO",
	"movie about himself",
	"Cryptocoin",
}

var consequences = []string{
	"so the price started to fluctuate",
	"so community started a fundraiser",
	"so market crashed",
	"so hashpower fleed",
}

var impacts = []impactDescriptor{
	{"you lost %d coins...", -1},
	{"you gained %d coins...", 1},
	{"(completely unrelated), someone stole %d coins from your storage...", -1},
}

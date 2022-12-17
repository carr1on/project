package main

import "fmt"

func main() {
	var (
		base               float32 = 18
		timeHourPerTickets         = 60
		percentHundred     float32 = 100
		percentRequired    float32
	)

	someInteger := 32 - base
	percentOne := base / percentHundred

	for someInteger < base {
		base = base - percentOne
		percentRequired++
		if someInteger == base {
			break
		}
	}
	percentRequired = percentHundred - percentRequired
	timeToAccept := float32(timeHourPerTickets) * (percentRequired / 100)

	timeForAccept := fmt.Sprintf("время ожидания до ответа на новый запрос в Support: %d минут", int(timeToAccept))
	fmt.Println(timeForAccept)
}

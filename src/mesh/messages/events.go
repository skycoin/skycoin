package messages

import (
	"fmt"
	"io/ioutil"
	"time"
)

type TimeEvent struct {
	happened int64
	comment  string
}

var events = []TimeEvent{}

var events2 = []TimeEvent{TimeEvent{}, TimeEvent{}}

func RegisterEvent(comment string) {
	return
	eventOld := events2[1]

	eventNew := TimeEvent{
		time.Now().UnixNano(),
		comment,
	}

	events2[0] = eventOld
	events2[1] = eventNew

	delta := eventNew.happened - eventOld.happened
	fmt.Printf("%s - %s\t\t%d\n", eventOld.comment, eventNew.comment, delta)
}

func RegisterEventOld(comment string) {
	return
	event := TimeEvent{
		time.Now().UnixNano(),
		comment,
	}
	events = append(events, event)
	i := len(events) - 1
	if i == 0 {
		return
	}
	delta := events[i].happened - events[i-1].happened
	fmt.Printf("%d\t%s - %s\t\t%d\n", i, events[i-1].comment, events[i].comment, delta)
}

func PrintEvents(treshold float64) {
	n := len(events)
	overall := events[n-1].happened - events[0].happened
	var delta int64
	output := ""
	for i := 1; i < n; i++ {
		delta = events[i].happened - events[i-1].happened
		percent := float64(delta) / float64(overall) * 100
		if percent > treshold {
			fmt.Printf("%d\t%s - %s\t\t%d\t%f\n\n", i, events[i-1].comment, events[i].comment, delta, percent)
			output += fmt.Sprintf("%s - %s,%d,%f\n", events[i-1].comment, events[i].comment, delta, percent)
		}
	}
	outputBytes := []byte(output)
	err := ioutil.WriteFile("output", outputBytes, 0644)
	if err != nil {
		panic(err)
	}
}

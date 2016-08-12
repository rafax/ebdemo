package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/satori/go.uuid"
	vegeta "github.com/tsenart/vegeta/lib"
)

var (
	rps     = 2000
	seconds = 6000
)

func main() {
	rate := uint64(rps) // per second
	duration := time.Duration(seconds) * time.Second

	var metrics vegeta.Metrics
	histogram := vegeta.Histogram{
		Buckets: []time.Duration{
			0,
			1 * time.Millisecond,
			5 * time.Millisecond,
			10 * time.Millisecond,
			100 * time.Millisecond,
			1000 * time.Millisecond,
		},
	}
	attacker := vegeta.NewAttacker()
	for res := range attacker.Attack(buildTargeter(), rate, duration) {
		metrics.Add(res)
		histogram.Add(res)
	}
	metrics.Close()
	vegeta.NewHistogramReporter(&histogram).Report(os.Stdout)

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

func buildTargeter() vegeta.Targeter {
	t := []vegeta.Target{}
	for i := 0; i < rps*seconds; i++ {
		body, _ := json.Marshal(PlayheadUpdate{Mgid: uuid.NewV4().String(), Playhead: strconv.Itoa(rand.Intn(10000))})
		t = append(t, vegeta.Target{
			Method: "POST",
			URL:    "http://localhost:3000/mem/" + uuid.NewV4().String(),
			Body:   body,
		})
	}
	fmt.Println("Built targeter")
	return vegeta.NewStaticTargeter(t...)
}

type PlayheadUpdate struct {
	Mgid     string
	Playhead string
}

package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	traqwriter "github.com/ras0q/traq-writer"
	"github.com/robfig/cron/v3"
)

const BELL_TIMERS = 108

var randSeed *rand.Rand

func main() {
	randSeed = rand.New(rand.NewSource(time.Now().UnixNano()))

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	c := cron.New(cron.WithLocation(jst))

	target := os.Getenv("TARGET")
	switch target {
	case "traq":
		w = traqwriter.NewTraqWebhookWriter(
			os.Getenv("TRAQ_WEBHOOK_ID"),
			os.Getenv("TRAQ_WEBHOOK_SECRET"),
			traqwriter.DefaultHTTPOrigin,
		)
		log.Println("Writing to traQ")
		_, err := w.Write([]byte(":bell: 起動"))
		if err != nil {
			panic(err)
		}
	default:
		w = os.Stdout
		log.Println("No TARGET specified, writing to stdout")
	}

	for i := 0; i < BELL_TIMERS; i++ {
		t := fmt.Sprintf("%d %d 31 12 *", 59-(i%60), 23-i/60)
		log.Println(t)
		_, err := c.AddFunc(t, postMessage(i+1))
		if err != nil {
			panic(err)
		}
	}

	c.Start()

	runtime.Goexit()
}

var w io.Writer

var messages = []string{
	":bell::bell::bell:",
	strings.Repeat(":bell:", 100),
	":no_bell:",
	":bellhop:", //🛎
	":bell.ex-large.wiggle:",
	":joshua_bell:",
	":Weepinbell:",
	":bell_pepper:",
}

func postMessage(count int) func() {
	return func() {
		message := ":bell:"
		id := randSeed.Intn(100)
		if id < len(messages) {
			message = messages[id]
		}
		message = fmt.Sprintf("%s (%d)", message, count)

		_, err := w.Write([]byte(message))
		if err != nil {
			log.Println(err)
		}
	}
}

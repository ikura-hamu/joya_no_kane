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

	nextYear := time.Now().Year() + 1
	newYearTime := time.Date(nextYear, 1, 1, 0, 0, 0, 0, jst)
	cronTime := newYearTime
	for i := 0; i < BELL_TIMERS; i++ {
		// 01/02 03:04:05PM '06 -0700
		t := cronTime.Format("04 03 02 01 *")
		log.Println(t)
		_, err := c.AddFunc(t, postMessage(i))
		if err != nil {
			panic(err)
		}
		cronTime = cronTime.Add(-time.Minute)
	}

	_, err = c.AddFunc("00 22 31 12 *", func() {
		_, err := fmt.Fprintf(w, "%s (1-108)\n:tada:", strings.Repeat(":bell:", BELL_TIMERS))
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		panic(err)
	}

	c.Start()

	runtime.Goexit()
}

var w io.Writer

var messages = []string{
	":bell::bell::bell:",
	strings.Repeat(":bell:", 100),
	":no_bell.large:",
	":bellhop.large:", //🛎
	":bell.ex-large.wiggle:",
	":joshua_bell.large:",
	":Weepinbell.large:",
	":bell_pepper.large:",
	`:null::null::null::null::bell::null::null::null::null:
:null::null::bell::bell::bell::bell::bell::null::null:
:null::bell::bell::bell::bell::bell::bell::bell::null:
:null::bell::bell::bell::bell::bell::bell::bell::null:
:null::bell::bell::bell::bell::bell::bell::bell::null:
:null::bell::bell::bell::bell::bell::bell::bell::null:
:null::bell::bell::bell::bell::bell::bell::bell::null:
:null::bell::bell::bell::bell::bell::bell::bell::null:
:bell::bell::bell::bell::bell::bell::bell::bell::bell:
:null::null::null::null::bell::null::null::null::null:
`,
	":ka-n_zubora.large:",
}

func postMessage(count int) func() {
	return func() {
		message := ":bell.large:"
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

package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ttn-nguyen42/retries"
)

func main() {
	failsUntil := 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	retries.Do(func() error {
		if failsUntil == 0 {
			return nil
		}
		failsUntil -= 1
		return errors.New("random errors")
	},
		retries.Context(ctx),
		retries.Attempts(20),
		retries.OnRetry(func(trial uint, incomingTimeout time.Duration, err error) {
			log.Printf("Retry #%v, until next retry %v", trial, incomingTimeout)
		}),
		retries.Delay(time.Millisecond*200),
		retries.MaxDelay(time.Second*5),
		retries.DelayMethod(retries.ConstantDelay),
	)
	log.Println("Finished")
}

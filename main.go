package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

var fail = flag.Int("fail", 3, "fail at?")
var handle = flag.Bool("handle", false, "handle or not?")

func work(idx int) error {
	fmt.Printf("running: %d\n", idx)
	if idx == *fail {
		return fmt.Errorf("fail at: %d", idx)
	}
	time.Sleep(1 * time.Second)
	return nil
}

func main() {
	flag.Parse()
	fmt.Printf("fail at: %d, handle %v\n", *fail, *handle)

	ctx := context.Background()
	egrp, ctx := errgroup.WithContext(ctx)
	egrp.SetLimit(5)

	for i := 0; i < 10; i++ {
		idx := i
		egrp.Go(func() error {
			select {
			case <-ctx.Done():
				if *handle {
					return ctx.Err()
				}
				return work(idx)
			default:
				return work(idx)
			}

		})
	}

	fmt.Printf("Done: %v\n", egrp.Wait())
}

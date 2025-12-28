package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aileron-projects/go/ztime/zcron"
)

func main() {
	c := &zcron.Config{
		Crontab: "* */3 * * * * *", // every 3s.
		JobFunc: func(ctx context.Context) error {
			fmt.Println("It's", time.Now())
			return nil
		},
	}
	cron, err := zcron.NewCron(c)
	if err != nil {
		panic(err)
	}
	cron.Start()
}

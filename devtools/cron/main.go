package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/aileron-projects/go/ztime/zcron"
)

var (
	cron = flag.String("cron", "* * * * * *", "cron expression")
)

func main() {
	flag.Parse()
	if *cron == "" {
		flag.Usage()
		os.Exit(1)
	}

	count := 0
	c, err := zcron.NewCron(&zcron.Config{
		Crontab: *cron,
		JobFunc: func(ctx context.Context) error {
			count++
			log.Println("job called ", count)
			return nil
		},
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}

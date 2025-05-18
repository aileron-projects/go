package zcron_test

import (
	"context"
	"fmt"
	"time"

	"github.com/aileron-projects/go/ztime/zcron"
)

func ExampleCron() {
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	count := 0
	stop := make(chan struct{})

	job, _ := zcron.NewCron(&zcron.Config{
		Crontab: "* * * * * *",
		JobFunc: func(ctx context.Context) error {
			if count++; count > 3 {
				stop <- struct{}{}
				return nil
			}
			now = now.Add(time.Second)
			fmt.Println(now.Format(time.DateTime), "Hello Go!")
			return nil
		},
	})

	job.WithTimeFunc(func() time.Time { return now })
	go job.Start()
	<-stop

	// Output:
	// 2000-01-01 00:00:01 Hello Go!
	// 2000-01-01 00:00:02 Hello Go!
	// 2000-01-01 00:00:03 Hello Go!
}

func ExampleCrontab_everySeconds() {
	ct, err := zcron.Parse("* * * * * *") // Every seconds.
	if err != nil {
		panic(err)
	}

	// Replace internal clock for testing.
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	ct.WithTimeFunc(func() time.Time { return now })

	for range 6 {
		fmt.Println(now.Format(time.DateTime), "|", ct.Next().Format(time.DateTime))
		now = now.Add(500 * time.Millisecond) // Forward 500ms.
	}
	// Output:
	// 2000-01-01 00:00:00 | 2000-01-01 00:00:01
	// 2000-01-01 00:00:00 | 2000-01-01 00:00:01
	// 2000-01-01 00:00:01 | 2000-01-01 00:00:02
	// 2000-01-01 00:00:01 | 2000-01-01 00:00:02
	// 2000-01-01 00:00:02 | 2000-01-01 00:00:03
	// 2000-01-01 00:00:02 | 2000-01-01 00:00:03
}

func ExampleCrontab() {
	ct, err := zcron.Parse("59 59 23 31 * *") // Schedule at 31st at 23:59:59
	if err != nil {
		panic(err)
	}

	// Replace internal clock for testing.
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	ct.WithTimeFunc(func() time.Time { return now })

	for range 10 {
		fmt.Println(now.Format(time.DateTime), "|", ct.Next().Format(time.DateTime))
		now = now.Add(15 * 24 * time.Hour) // Forward 15 days.
	}

	// Output:
	// 2000-01-01 00:00:00 | 2000-01-31 23:59:59
	// 2000-01-16 00:00:00 | 2000-01-31 23:59:59
	// 2000-01-31 00:00:00 | 2000-01-31 23:59:59
	// 2000-02-15 00:00:00 | 2000-03-31 23:59:59
	// 2000-03-01 00:00:00 | 2000-03-31 23:59:59
	// 2000-03-16 00:00:00 | 2000-03-31 23:59:59
	// 2000-03-31 00:00:00 | 2000-03-31 23:59:59
	// 2000-04-15 00:00:00 | 2000-05-31 23:59:59
	// 2000-04-30 00:00:00 | 2000-05-31 23:59:59
	// 2000-05-15 00:00:00 | 2000-05-31 23:59:59
}

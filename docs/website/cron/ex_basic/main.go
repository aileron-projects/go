package main

import (
	"fmt"
	"time"

	"github.com/aileron-projects/go/ztime/zcron"
)

func main() {
	// Working with the local time zone.
	cron, _ := zcron.Parse("* * * * *")
	fmt.Println("Next:", cron.Next())
	fmt.Println("Next After 1 hour:", cron.NextAfter(time.Now().Add(time.Hour)))
	fmt.Println("Next After 1 day:", cron.NextAfter(time.Now().Add(24*time.Hour)))

	// Working with a specific time zone.
	tzcron, _ := zcron.Parse("TZ=UTC * * * * *")
	fmt.Println("Next:", tzcron.Next())
	fmt.Println("Next After 1 hour:", tzcron.NextAfter(time.Now().UTC().Add(time.Hour)))
	fmt.Println("Next After 1 day:", tzcron.NextAfter(time.Now().UTC().Add(24*time.Hour)))
}

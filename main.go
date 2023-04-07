package main

import (
	// "context"

	"fmt"

	cs "github.com/leftrightleft/gh-ghas-gauge/internal/codescanning"
)

func main() {
	cs_alerts := cs.GetCodeScanningAlerts("burrito-party")

	var closed_count int
	var open_count int
	var dismissed_count int
	var fixed_count int
	for _, alert := range cs_alerts {
		if alert.State == "closed" {
			closed_count++
		}
		if alert.State == "dismissed" {
			dismissed_count++
		}
		if alert.State == "fixed" {
			fixed_count++
		}
		if alert.State == "open" {
			open_count++
		}
	}
	fmt.Println("CODE SCANNING ALERTS")
	fmt.Println("  - total alerts: ", len(cs_alerts))
	fmt.Println("  - open alerts: ", open_count)
	fmt.Println("  - closed alerts: ", closed_count)
	fmt.Println("  - dismissed alerts: ", dismissed_count)
	fmt.Println("  - fixed alerts: ", fixed_count)

}

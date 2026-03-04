package sys

import (
	"fmt"
	"testing"
	"time"
)

func TestAccuracyDrift(t *testing.T) {
	Calibrate()

	durations := []time.Duration{
		100 * time.Microsecond,
		1 * time.Millisecond,
		10 * time.Millisecond,
		100 * time.Millisecond,
	}

	fmt.Println("\n--- ACCURACY DRIFT AUDIT ---")
	fmt.Printf("%-15s | %-15s | %-15s | %-10s\n", "Target", "time.Since", "PreciseDur", "Diff %")
	fmt.Println("----------------------------------------------------------------------")

	for _, d := range durations {
		startTSC := RDTSC()
		startWall := time.Now()

		time.Sleep(d)

		endWall := time.Now()
		endTSC := RDTSC()

		wallDuration := uint64(endWall.Sub(startWall).Nanoseconds())
		nanoDuration := PreciseDuration(startTSC, endTSC)

		diff := float64(int64(nanoDuration-wallDuration)) / float64(wallDuration) * 100

		fmt.Printf("%-15v | %-15d | %-15d | %-10.4f%%\n", d, wallDuration, nanoDuration, diff)
	}
	fmt.Println("----------------------------------------------------------------------")
}

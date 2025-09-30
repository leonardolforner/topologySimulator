package sim

import (
	"fmt"
	"sort"

	"github.com/leonardolforner/topologySimulator/internal/config"
)

type Aggregate struct {
	TimeByQueue map[string]map[int]float64
	Losses      map[string]int
	SumSimTime  float64
	AvgSimTime  float64
	Reps        int
}

func AggregateResults(cfg *config.Config, reps []replicationResult, avg float64) Aggregate {
	out := Aggregate{
		TimeByQueue: map[string]map[int]float64{},
		Losses:      map[string]int{},
		SumSimTime:  0,
		AvgSimTime:  avg,
		Reps:        len(reps),
	}
	for _, r := range reps {
		out.SumSimTime += r.totalTime
		for q, tb := range r.timeByQueue {
			m := out.TimeByQueue[q]
			if m == nil {
				m = map[int]float64{}
				out.TimeByQueue[q] = m
			}
			for n, t := range tb {
				m[n] += t
			}
		}
		for q, L := range r.losses {
			out.Losses[q] += L
		}
	}
	return out
}

func PrintReport(cfg *config.Config, ag Aggregate) {
	fmt.Println()
	fmt.Println("=========================================================")
	fmt.Println("======================    REPORT   ======================")
	fmt.Println("=========================================================")

	names := make([]string, 0, len(cfg.Queues))
	for n := range cfg.Queues {
		names = append(names, n)
	}
	sort.Strings(names)

	for _, name := range names {
		q := cfg.Queues[name]
		if q.Capacity == nil {
			fmt.Printf("*********************************************************\nQueue:   %s (G/G/%d)\n", name, q.Servers)
		} else {
			fmt.Printf("*********************************************************\nQueue:   %s (G/G/%d/%d)\n", name, q.Servers, *q.Capacity)
		}
		if q.MinArrival != nil && q.MaxArrival != nil {
			if _, ok := cfg.Arrivals[name]; ok {
				fmt.Printf("Arrival: %.1f ... %.1f\n", *q.MinArrival, *q.MaxArrival)
			}
		}
		fmt.Printf("Service: %.1f ... %.1f\n", q.MinService, q.MaxService)
		fmt.Println("*********************************************************")
		fmt.Println("   State               Time               Probability")

		tb := ag.TimeByQueue[name]
		var total float64
		for _, t := range tb {
			total += t
		}

		maxN := 0
		for n := range tb {
			if n > maxN {
				maxN = n
			}
		}
		for n := 0; n <= maxN; n++ {
			t := tb[n]
			p := 0.0
			if total > 0 {
				p = 100.0 * t / total
			}
			fmt.Printf("%6d %18.4f %21.2f%%\n", n, t, p)
		}
		fmt.Printf("\nNumber of losses: %d\n\n", ag.Losses[name])
	}

	fmt.Println("=========================================================")
	fmt.Printf("Simulation average time: %.4f\n", ag.AvgSimTime)
	fmt.Println("=========================================================")
}

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/leonardolforner/topologySimulator/internal/config"
	"github.com/leonardolforner/topologySimulator/internal/sim"
)

func main() {
	cfgPath := flag.String("config", "topology.yaml", "YAML contract file")
	maxDraws := flag.Int("max-rands", 100000, "Stop when this many random draws are consumed (per replication)")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	fmt.Println("=========================================================")
	fmt.Println("============   QUEUEING NETWORK SIMULATOR   =============")
	fmt.Println("==================     Leonardo Forner    ===============")
	fmt.Println("=========================================================")

	engine, err := sim.NewSimulator(cfg, sim.Options{MaxRandomDraws: *maxDraws})
	if err != nil {
		log.Fatalf("sim init: %v", err)
	}

	reps, avg, err := engine.Run()
	if err != nil {
		log.Fatalf("run: %v", err)
	}

	fmt.Println("=========================================================")
	fmt.Println("=================    END OF SIMULATION   ================")
	fmt.Println("=========================================================")

	ag := sim.AggregateResults(cfg, reps, avg)
	sim.PrintReport(cfg, ag)
}

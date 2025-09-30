package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/leonardolforner/topologySimulator/internal/config"
)

func main() {
	path := flag.String("config", "topology.yaml", "YAML file")
	flag.Parse()

	cfg, err := config.Load(*path)
	if err != nil {
		log.Fatalf("load: %v", err)
	}

	fmt.Println("Queues loaded:")
	for name, q := range cfg.Queues {
		capStr := "∞"
		if q.Capacity != nil {
			capStr = fmt.Sprint(*q.Capacity)
		}
		fmt.Printf("- %s: servers=%d cap=%s svc=[%.1f..%.1f]\n",
			name, q.Servers, capStr, q.MinService, q.MaxService)
	}
}

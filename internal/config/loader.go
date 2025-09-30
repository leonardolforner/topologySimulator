package config

import (
	"fmt"
	"math"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("yaml unmarshall: %w", err)
	}
	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validate(c *Config) error {
	if len(c.Queues) == 0 {
		return fmt.Errorf("no queues defined")
	}

	for q, t0 := range c.Arrivals {
		if t0 < 0 {
			return fmt.Errorf("arrivals[%s] must be >= 0", q)
		}
		qq, ok := c.Queues[q]
		if !ok {
			return fmt.Errorf("arrivals refers to unknown queue %q", q)
		}
		if qq.MinArrival == nil || qq.MaxArrival == nil {
			return fmt.Errorf("queue %q is an external source but missing minArrival/maxArrival", q)
		}
		if *qq.MinArrival >= *qq.MaxArrival {
			return fmt.Errorf("queue %q: minArrival must be < maxArrival", q)
		}
	}

	for name, q := range c.Queues {
		if q.Servers <= 0 {
			return fmt.Errorf("queue %q: servers must be >= 1", name)
		}
		if q.MinService <= 0 || q.MaxService <= 0 || q.MinService >= q.MaxService {
			return fmt.Errorf("queue %q: invalid service range", name)
		}
		if q.Capacity != nil && *q.Capacity < 0 {
			return fmt.Errorf("queue %q: capacity must be >= 0", name)
		}
	}

	bySrc := map[string]float64{}
	for _, e := range c.Network {
		if _, ok := c.Queues[e.Source]; !ok && e.Source != "" {
			return fmt.Errorf("edge: unknown source queue %q", e.Source)
		}
		if _, ok := c.Queues[e.Target]; !ok {
			return fmt.Errorf("edge: unknown target queue %q", e.Target)
		}
		if e.Probability < 0 || e.Probability > 1 {
			return fmt.Errorf("edge %s->%s: probability out of [0,1]", e.Source, e.Target)
		}
		bySrc[e.Source] += e.Probability
	}
	for src, sum := range bySrc {
		if math.Abs(sum-1.0) > 1e-9 {
			return fmt.Errorf("outgoing probabilities from %q must sum to 1.0 (got %.12f)", src, sum)
		}
	}

	if len(c.Seeds) > 0 && c.RndNumbersPerSeed <= 0 {
		return fmt.Errorf("seeds provided but rndnumbersPerSeed <= 0")
	}
	return nil
}

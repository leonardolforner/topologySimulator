package config

type Config struct {
	Arrivals          map[string]float64
	Queues            map[string]Queue
	Network           []Edge
	RndNumbers        []float64
	RndNumbersPerSeed int
	Seeds             []int64
}

type Queue struct {
	Servers    int
	Capacity   *int // nil == no maximum
	MinArrival *float64
	MaxArrival *float64
	MinService float64
	MaxService float64
}

type Edge struct {
	Source      string
	Target      string
	Probability float64
}

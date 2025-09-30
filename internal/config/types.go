package config

type Config struct {
	Arrivals          map[string]float64 `yaml:"arrivals"`
	Queues            map[string]Queue   `yaml:"queues"`
	Network           []Edge             `yaml:"network"`
	RndNumbers        []float64          `yaml:"rndnumbers"`
	RndNumbersPerSeed int                `yaml:"rndnumbersperseed"`
	Seeds             []int64            `yaml:"seeds"`
}

type Queue struct {
	Servers    int      `yaml:"servers"`
	Capacity   *int     `yaml:"capacity"` // nil == no maximum
	MinArrival *float64 `yaml:"minArrival"`
	MaxArrival *float64 `yaml:"maxArrival"`
	MinService float64  `yaml:"minService"`
	MaxService float64  `yaml:"maxService"`
}

type Edge struct {
	Source      string  `yaml:"source"`
	Target      string  `yaml:"target"`
	Probability float64 `yaml:"probability"`
}

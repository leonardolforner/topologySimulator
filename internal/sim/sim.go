package sim

import (
	"container/heap"
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/leonardolforner/topologySimulator/internal/config"
)

type Options struct{ MaxRandomDraws int }

type Simulator struct {
	cfg  *config.Config
	opts Options
	rngF RNGFactory
}

func NewSimulator(cfg *config.Config, opts Options) (*Simulator, error) {
	if opts.MaxRandomDraws <= 0 {
		return nil, errors.New("MaxRandomDraws must be > 0")
	}
	return &Simulator{cfg: cfg, opts: opts, rngF: BuildRNGFactory(cfg.RndNumbers, cfg.Seeds)}, nil
}

type eventType int

const (
	evExternalArrival eventType = iota
	evServiceEnd
)

type event struct {
	t      float64
	kind   eventType
	q      string
	_order int64
}

type eventPQ []*event

func (p eventPQ) Len() int { return len(p) }
func (p eventPQ) Less(i, j int) bool {
	if p[i].t != p[j].t {
		return p[i].t < p[j].t
	}
	if p[i].kind != p[j].kind {
		return p[i].kind > p[j].kind
	} // service end first
	return p[i]._order < p[j]._order
}
func (p eventPQ) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p *eventPQ) Push(x any)   { *p = append(*p, x.(*event)) }
func (p *eventPQ) Pop() any     { old := *p; n := len(old); x := old[n-1]; *p = old[:n-1]; return x }

type queueState struct {
	name              string
	servers, capacity int
	inSystem, busy    int
	lastT             float64
	timeByN           map[int]float64
	losses            int
}

type replicationResult struct {
	timeByQueue map[string]map[int]float64
	losses      map[string]int
	totalTime   float64
}

func (s *Simulator) Run() ([]replicationResult, float64, error) {
	reps := s.rngF.Replications()
	out := make([]replicationResult, 0, reps)
	var sum float64

	for i := 0; i < reps; i++ {
		r := s.rngF.NewReplication(i)
		res, tot, err := s.runOne(r)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, res)
		sum += tot

		fmt.Printf("Simulation: #%d\n", i+1)
		if len(s.cfg.Seeds) > 0 {
			fmt.Printf("...simulating with random numbers (seed '%d')...\n", s.cfg.Seeds[i])
		}
	}
	return out, sum / float64(reps), nil
}

func (s *Simulator) runOne(rng RNG) (replicationResult, float64, error) {
	// queue states
	qs := map[string]*queueState{}
	for name, q := range s.cfg.Queues {
		capSys := math.MaxInt
		if q.Capacity != nil {
			capSys = *q.Capacity
		}
		qs[name] = &queueState{name: name, servers: q.Servers, capacity: capSys, timeByN: map[int]float64{}}
	}

	// events
	var pq eventPQ
	heap.Init(&pq)
	var order int64
	for qName, t0 := range s.cfg.Arrivals {
		heap.Push(&pq, &event{t: t0, kind: evExternalArrival, q: qName, _order: order})
		order++
	}

	var now float64
	limit := s.opts.MaxRandomDraws

	advance := func(to float64) {
		dt := to - now
		if dt < 0 {
			dt = 0
		}
		for _, st := range qs {
			st.timeByN[st.inSystem] += dt
			st.lastT = to
		}
		now = to
	}
	getQ := func(name string) *queueState { return qs[name] }

	startIfPossible := func(st *queueState) {
		if st.inSystem > st.busy && st.busy < st.servers && rng.Used() < limit {
			svc := Uniform(rng, s.cfg.Queues[st.name].MinService, s.cfg.Queues[st.name].MaxService)
			st.busy++
			heap.Push(&pq, &event{t: now + svc, kind: evServiceEnd, q: st.name, _order: order})
			order++
		}
	}
	admit := func(st *queueState) bool {
		if st.inSystem >= st.capacity {
			st.losses++
			return false
		}
		st.inSystem++
		startIfPossible(st)
		return true
	}
	scheduleNextExternal := func(qName string) {
		qq := s.cfg.Queues[qName]
		if rng.Used() < limit {
			ia := Uniform(rng, *qq.MinArrival, *qq.MaxArrival)
			heap.Push(&pq, &event{t: now + ia, kind: evExternalArrival, q: qName, _order: order})
			order++
		}
	}
	routeFrom := func(source string) (string, bool) {
		edges := s.outgoing(source)
		if len(edges) == 0 {
			return "", false
		}
		if rng.Used() >= limit {
			return "", false
		}
		u := rng.Next()
		acc := 0.0
		for _, e := range edges {
			acc += e.Probability
			if u <= acc {
				return e.Target, true
			}
		}
		return edges[len(edges)-1].Target, true
	}

	for pq.Len() > 0 && rng.Used() < limit {
		ev := heap.Pop(&pq).(*event)
		advance(ev.t)

		switch ev.kind {
		case evExternalArrival:
			if rng.Used() < limit {
				scheduleNextExternal(ev.q)
			}
			admit(getQ(ev.q))

		case evServiceEnd:
			src := getQ(ev.q)
			if src.busy > 0 {
				src.busy--
			}
			if src.inSystem > 0 {
				src.inSystem--
			}
			startIfPossible(src)
			if tgt, ok := routeFrom(src.name); ok {
				admit(getQ(tgt))
			}
		}
	}

	res := replicationResult{
		timeByQueue: map[string]map[int]float64{},
		losses:      map[string]int{},
		totalTime:   now,
	}
	for name, st := range qs {
		res.timeByQueue[name] = st.timeByN
		res.losses[name] = st.losses
	}
	return res, now, nil
}

func (s *Simulator) outgoing(source string) []config.Edge {
	var edges []config.Edge
	for _, e := range s.cfg.Network {
		if e.Source == source {
			edges = append(edges, e)
		}
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i].Probability < edges[j].Probability })
	return edges
}

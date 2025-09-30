package sim

import "math/rand"

type RNG interface {
	Next() float64 // U(0,1)
	Used() int
}

type ListRng struct {
	list []float64
	i    int
}

func NewListRNG(xs []float64) *ListRng {
	return &ListRng{
		list: xs,
	}
}

func (r *ListRng) Next() float64 {
	if r.i >= len(r.list) {
		panic("rndnumbers exhausted")
	}
	v := r.list[r.i]
	r.i++
	return v
}

func (r *ListRng) Used() int {
	return r.i
}

type SeededRNG struct {
	r    *rand.Rand
	used int
}

func NewSeededRNG(seed int64) *SeededRNG {
	return &SeededRNG{
		r: rand.New(rand.NewSource(seed)),
	}
}

func (r *SeededRNG) Next() float64 {
	r.used++
	return r.r.Float64()
}

func (r *SeededRNG) Used() int {
	return r.used
}

type RNGFactory interface {
	NewReplication(repIdx int) RNG
	Replications() int
}

type listFactory struct {
	nums []float64
}

func (f listFactory) NewReplication(int) RNG {
	return NewListRNG(f.nums)
}

func (f listFactory) Replications() int {
	return 1
}

type seedsFactory struct {
	seeds []int64
}

func (f seedsFactory) NewReplication(i int) RNG {
	return NewSeededRNG(f.seeds[i])
}
func (f seedsFactory) Replications() int {
	return len(f.seeds)
}

func BuildRNGFactory(nums []float64, seeds []int64) RNGFactory {
	if len(seeds) > 0 {
		return seedsFactory{seeds: seeds}
	}
	return listFactory{nums: nums}
}

func Uniform(r RNG, min, max float64) float64 {
	u := r.Next()
	return min + u*(max-min)
}

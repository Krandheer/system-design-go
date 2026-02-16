package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"sync/atomic"
)

// BloomFilter is a probabilistic data structure.
type BloomFilter struct {
	bits      []uint64
	m         uint64 // total bits
	k         uint64 // total hash functions
}


func NewBloomFilter(expectedItems uint64, falsePositiveRate float64) *BloomFilter {
	m := uint64(-float64(expectedItems) *
		math.Log(falsePositiveRate) /
		(math.Ln2 * math.Ln2))

	k := uint64((float64(m) / float64(expectedItems)) * math.Ln2)

	return &BloomFilter{
		bits: make([]uint64, (m+63)/64),
		m:    m,
		k:    k,
	}
}

func hash(data string) (uint64, uint64) {
	h1 := fnv.New64a()
	h1.Write([]byte(data))
	sum1 := h1.Sum64()

	h2 := fnv.New64()
	h2.Write([]byte(data))
	sum2 := h2.Sum64()

	return sum1, sum2
}


func (bf *BloomFilter) setBit(pos uint64) {
	word := pos / 64
	mask := uint64(1) << (pos % 64)
	for {
		old := atomic.LoadUint64(&bf.bits[word])
		if old & mask != 0 {
			return
		}
		if atomic.CompareAndSwapUint64(&bf.bits[word], old, old | mask) {
			return
		}
	}
}

func (bf *BloomFilter) getBit(pos uint64) bool {
	word := pos / 64
	mask := uint64(1) << (pos % 64)
	return (atomic.LoadUint64(&bf.bits[word]) & mask) != 0
}

// Add inserts data into the filter.
func (bf *BloomFilter) Add(data string) {
	h1, h2 := hash(data)

	for i := uint64(0); i < bf.k; i++ {
		pos := (h1 + i*h2) % bf.m
		bf.setBit(pos)
	}
}

func (bf *BloomFilter) Check(data string) bool {
	h1, h2 := hash(data)

	for i := uint64(0); i < bf.k; i++ {
		pos := (h1 + i*h2) % bf.m
		if !bf.getBit(pos) {
			return false // DEFINITELY NOT PRESENT
		}
	}
	return true // MAYBE PRESENT
}


func main() {
	// Create a filter with 1000000 items and a 1% false positive rate
	bf := NewBloomFilter(1000000, 0.01)

	// 1. Add some data
	bf.Add("apple")
	bf.Add("banana")
	bf.Add("orange")

	fmt.Println("\n--- Checking Existence ---")

	// 2. Check for items we added (True Positives)
	fmt.Printf("Is 'apple' in set? %v\n", bf.Check("apple"))     // Expected: true
	fmt.Printf("Is 'banana' in set? %v\n", bf.Check("banana"))   // Expected: true

	// 3. Check for items we did NOT add (True Negatives)
	fmt.Printf("Is 'car' in set? %v\n", bf.Check("car"))         // Expected: false

	// 4. Demonstrate False Positive potential
	// In a real system, you tune the expected items and false positive rate to make this probability < 1%.
	// Here, with 1000000 items and a 1% false positive rate, collisions are unlikely.
	// If this prints 'true', it's a False Positive!
	// (Note: It depends on the hash values, so it might correctly say false too).
	fmt.Printf("Is 'kiwi' in set? %v (Possible False Positive)\n", bf.Check("kiwi"))

}

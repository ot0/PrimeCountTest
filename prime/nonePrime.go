package prime

import (
	"math"
)

func CountNonePrime(stop int) int {
	sstop := int(math.Sqrt(float64(stop)) + 1.0)

	var ss SieveStore
	ss.Init(0, sstop)
	count := CountAtokin(sstop, &ss)

	primes := MakePrimeArray(count, &ss)
	//fmt.Println(primes)

	nonePrime := 0
	for i, num := range primes {
		nonePrime += CountNP(num, i, stop, primes)
	}

	return nonePrime - len(primes) + 1
}

func CountNP(num int, i int, stop int, primes []int) int {
	sum := stop / num
	if i == 0 {
		return sum
	}
	for j := 0; j < i; j++ {
		k := num * primes[j]
		if k > stop {
			return sum
		}
		sum -= CountNP(k, j, stop, primes)
	}
	return sum
}

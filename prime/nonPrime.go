package prime

import (
	"math"
)

func CountNonPrime(stop int) int {
	sstop := int(math.Sqrt(float64(stop)) + 1.0)

	var ss SieveStore
	ss.Init(0, sstop)
	count := CountAtokin(sstop, &ss)

	primes := MakePrimeArray(count, &ss)
	//fmt.Println(primes)

	nonPrime := 0
	for i, num := range primes {
		nonPrime += CountNP(num, i, stop, primes)
	}

	return nonPrime - len(primes) + 1
}

func CountNonPrimeMulti(stop int) int {
	parallel := 8
	ch := make(chan int, parallel)
	for i := 0; i < parallel; i++ {
		ch <- 0
	}

	sstop := int(math.Sqrt(float64(stop)) + 1.0)

	var ss SieveStore
	ss.Init(0, sstop)
	count := CountAtokin(sstop, &ss)

	primes := MakePrimeArray(count, &ss)
	//fmt.Println(primes)

	Calc := func(num int, i int){
		ch <- CountNP(num, i, stop, primes)
	}
	nonPrime := 0
	for i:=len(primes)-1; i>=0; i-- {
		nonPrime += <-ch
		go Calc(primes[i], i)
	}

	for i := 0; i < parallel; i++ {
		nonPrime += <- ch
	}

	return nonPrime - len(primes) + 1
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

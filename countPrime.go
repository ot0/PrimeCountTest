package main

import (
	"fmt"
	"time"

	"./prime"
)

func main() {
	fmt.Println("Hello World!")

	stop := 1000 * 1000 * 1000 * 100
	start := time.Now()
	count := stop - prime.CountNonPrimeMulti(stop)
	fmt.Println(count, time.Now().Sub(start).Seconds())

	//AllTest()
}

func AllTest() {
	var stop int

	fmt.Println("Eratosthenes, normal memory")
	stop = 1
	for i := 1; i <= 9; i++ {
		stop *= 10
		start := time.Now()

		var ss prime.SieveStore
		ss.Init(0, stop)
		count := prime.CountPrime(stop, &ss)

		fmt.Println(i, count, time.Now().Sub(start).Seconds())
	}

	fmt.Println("Eratosthenes, wheel memory")
	stop = 1
	for i := 1; i <= 10; i++ {
		stop *= 10
		start := time.Now()

		var wh prime.Wheel
		base := []int{2, 3, 5, 7, 11, 13}
		wh.Init(base)
		ss := wh.NewStoreF(0, stop)
		count := prime.CountPrime(stop, ss)

		fmt.Println(i, count, time.Now().Sub(start).Seconds())
	}

	fmt.Println("Atokin, normal memory")
	stop = 1
	for i := 1; i <= 9; i++ {
		stop *= 10
		start := time.Now()

		var ss prime.SieveStore
		ss.Init(0, stop)
		count := prime.CountAtokin(stop, &ss)
		fmt.Println(i, count, time.Now().Sub(start).Seconds())
	}

	fmt.Println("Atokin, wheel memory")
	stop = 1
	for i := 1; i <= 10; i++ {
		stop *= 10
		start := time.Now()

		var wh prime.Wheel
		base := []int{2, 3, 5, 7, 11, 13}
		wh.Init(base)
		ss := wh.NewStoreF(0, stop)
		count := prime.CountAtokin(stop, ss)

		fmt.Println(i, count, time.Now().Sub(start).Seconds())
	}
	fmt.Println("Eratosthenes Multiprocess, normal memory")
	stop = 1
	for i := 1; i <= 11; i++ {
		stop *= 10
		start := time.Now()

		count := prime.CountPrimeMulti(stop)

		fmt.Println(i, count, time.Now().Sub(start).Seconds())
	}

	fmt.Println("Atokin Multiprocess, normal memory")
	stop = 1
	for i := 1; i <= 11; i++ {
		stop *= 10
		start := time.Now()

		count := prime.CountAtokinMulti(stop)

		fmt.Println(i, count, time.Now().Sub(start).Seconds())
	}

	fmt.Println("Atokin Multiprocess, Wheel memory")
	stop = 1000
	for i := 4; i <= 11; i++ {
		stop *= 10
		start := time.Now()

		count := prime.CountAtokinWheelMulti(stop)

		fmt.Println(i, count, time.Now().Sub(start).Seconds())
	}

	fmt.Println("Count Non-Prime")
	stop = 1
	for i := 1; i <= 11; i++ {
		stop *= 10
		start := time.Now()

		count := stop - prime.CountNonPrime(stop)

		fmt.Println(i, count, time.Now().Sub(start).Seconds())

	}

	fmt.Println("Count Non-Prime Multiprocess")
	stop = 1
	for i := 1; i <= 12; i++ {
		stop *= 10
		start := time.Now()

		count := stop - prime.CountNonPrimeMulti(stop)

		fmt.Println(i, count, time.Now().Sub(start).Seconds())

	}
}

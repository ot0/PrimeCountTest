package prime

import (
	"fmt"
	"math"
	"sync"
)

func CountPrime(stop int, ss SieveStoreInterface) int {
	SetAllTrue(ss)
	count := 0
	sstop := sqrt(stop)+1
	for i := 2; i < sstop; i++ {
		if ss.Get(i) {
			count++
			for j := i * i; j < stop; j += i {
				ss.Set(j, false)
			}
		}
	}
	for i:=sstop; i<stop; i++ {
		if ss.Get(i) {
			count++
		}
	}	
	return count
}

func CountAtokin(stop int, ss SieveStoreInterface) int {
	Atokin(ss)
	AtokinSqr(ss, ss)
	ss.Set(2, true)
	ss.Set(3, true)
	return CountSieve(ss, stop) + 1 // for 2
}

/*-----------------------------------------*/
/* Sieve Tools*/
/*-----------------------------------------*/
func SetAllTrue(ss SieveStoreInterface) {
	for i := ss.Start(); i < ss.Stop(); i++ {
		ss.Set(i, true)
	}
}

func MakePrimeArray(stop int, ss SieveStoreInterface) []int {
	mem := make([]int, stop)
	count := 0
	for i := 2; count < stop; i++ {
		if ss.Get(i) {
			mem[count] = i
			count++
		}
	}
	return mem
}

func CountSieve(d SieveStoreInterface, max int) int {
	count := 0
	stop := d.Stop()
	if max < stop {
		stop = max + 1
	}
	// d.Start is even
	for i := d.Start() + 1; i < stop; i += 2 {
		if d.Get(i) {
			count++
		}
	}
	return count
}

/*-----------------------------------------*/
/* Multi Eratosthenes Sieve*/
/*-----------------------------------------*/

func Eratosthenes(ss SieveStoreInterface, s0 SieveStoreInterface) {
	SetAllTrue(ss)

	sqrtN := sqrt(ss.Stop())
	if s0.Stop() <= sqrtN {
		sqrtN = s0.Stop() - 1
	}

	// s0.Start() is even. because 2 in base
	for n := s0.Start() + 1; n <= sqrtN; n += 2 {
		if s0.Get(n) {
			a := (ss.Start()-1)/n + 1
			a = a & ^1 + 1 // to odd
			for k := a * n; k < ss.Stop(); k += n * 2 {
				ss.Set(k, false)
			}
		}
	}
}

func CountPrimeMulti(stop int) int {
	parallel := 4
	ch := make(chan int, parallel)

	sstop := sqrt(stop)
	if sstop%2 == 1{
		sstop++
	}
	var d0 SieveStore
	d0.Init(2, sstop)
	count := CountPrime(d0.Stop(), &d0)
	//fmt.Println(sstop, count)

	for i := 0; i < parallel; i++ {
		ch <- 0
	}

	Calc := func(i int) {
		var d SieveStore
		d.Init(i, sstop)
		Eratosthenes(&d, &d0)
		c := CountSieve(&d, stop)
		//fmt.Println(i, c, d.Start())
		ch <- c
	}
	for i := d0.Stop(); i < stop; i += sstop {
		go Calc(i)
		count += <-ch
	}
	for i := 0; i < parallel; i++ {
		count += <-ch
	}

	return count
}

func CountAtokinMulti(stop int) int {
	parallel := 4
	ch := make(chan int, parallel)

	sstop := sqrt(stop)
	if sstop%2 == 1{
		sstop++
	}

	var wg sync.WaitGroup
	wg.Add(1)	

	var d0 SieveStore
	d0.Init(2, sstop)

	count := 0
	go func(){
		Atokin(&d0)
		AtokinSqr(&d0, &d0)
		d0.Set(2, true)
		d0.Set(3, true)
		wg.Done()
		c := CountSieve(&d0, stop) + 1
		ch <- c
	}()

	for i := 0; i < parallel-1; i++ {
		ch <- 0
	}

	Calc := func(i int) {
		var d SieveStore
		d.Init(i, sstop)
		Atokin(&d)
		wg.Wait()
		AtokinSqr(&d, &d0)
		c := CountSieve(&d, stop)
		//fmt.Println(i, c, d.Start())
		ch <- c
	}
	for i := d0.Stop(); i < stop; i += sstop {
		go Calc(i)
		count += <-ch
	}
	for i := 0; i < parallel; i++ {
		count += <-ch
	}

	return count
}

/*-----------------------------------------*/
/* Multi Eratosthenes Sieve*/
/*-----------------------------------------*/

func CountAtokinWheelMulti(stop int) int {
	parallel := 4
	ch := make(chan int, parallel)

	var wh Wheel
	base := []int{2, 3, 5, 7, 11, 13}
	wh.Init(base)

	step := wh.step * wh.count
	//step := 2 * 3 * 5 * 7 * 11 * 13 * 17 * 19 * 32
	/*
		@xeon x3460, parallel=4, base=2~13
		used memory:59,548k  taskmgr.exe
		6 :     2.353 # < step
		7 :     2.303 # < step
		8 :     3.626 # < step
		9 :    11.274
		10:    83.351
		11:   843.551
		12:  8260.395
		13: 83339.415
	*/
	if math.Log10(float64(step))*2 < math.Log10(float64(stop)) {
		fmt.Println(step, math.Log10(float64(step)), math.Log10(float64(stop)))
		fmt.Println("Stop Size Error. need more step")
		return 0
	}

	sum := 0
	var wg sync.WaitGroup
	wg.Add(1)

	d0 := wh.NewStore(0, step)
	go func() {
		Atokin(d0)
		AtokinSqr(d0, d0)
		//d0.Set(2, true)
		//d0.Set(3, true)
		wg.Done()
		c := CountSieve(d0, stop) + len(base)
		ch <- c
		//fmt.Println(0, c)
	}()

	for i := 0; i < parallel-1; i++ {
		ch <- 0
	}

	Calc := func(n int) {
		d := wh.NewStore(n, step)
		Atokin(d)
		wg.Wait()
		AtokinSqr(d, d0)
		c := CountSieve(d, stop)
		ch <- c
		//fmt.Println(d.Start(), c)
	}
	for i := step; i < stop; i += step {

		go Calc(i)

		sum += <-ch
		//fmt.Println("total:", sum)
	}

	for i := 0; i < parallel; i++ {
		sum += <-ch
		//fmt.Println("total:", sum)
	}
	return sum
}

/*-----------------------------------------*/
/* Atokin Sieve*/
/*-----------------------------------------*/

func calc4x(start int, y int) int {
	if start > y*y {
		return int(math.Ceil(math.Sqrt((float64(start) - float64(y*y)) / 4.0)))
	} else {
		return 1
	}
}

func calc3xm(start int, y int, z int) int {
	x := int(math.Ceil(math.Sqrt((float64(start) + float64(y*y)) / 3.0)))
	if x < y {
		return y + 1
	} else {
		if x%2 == z {
			x++
		}
		return x
	}
}

func calc3xp(start int, y int) int {
	x := 1
	if start > y*y {
		x = int(math.Ceil(math.Sqrt((float64(start) - float64(y*y)) / 3.0)))
		if x%2 == 0 {
			x++
		}
	}
	return x
}

func Atokin(d SieveStoreInterface) {
	N := d.Stop()
	sqrtN := sqrt(N)
	for y := 1; y <= sqrtN; y += 2 {
		if y%3 == 0 {
			continue
		}
		for x := calc4x(d.Start(), y); ; x++ {
			n := 4*x*x + y*y
			if n >= N {
				break
			}
			d.Toggle(n)
		}
		for x := calc3xm(d.Start(), y, 1); ; x += 2 {
			n := 3*x*x - y*y
			if n >= N {
				break
			}
			d.Toggle(n)
		}
	}

	for y := 2; y <= sqrtN; y += 2 {
		if y%3 == 0 {
			continue
		}
		for x := calc3xp(d.Start(), y); ; x += 2 {
			n := 3*x*x + y*y
			if n >= N {
				break
			}
			d.Toggle(n)
		}
		for x := calc3xm(d.Start(), y, 0); ; x += 2 {
			n := 3*x*x - y*y
			if n >= N {
				break
			}
			d.Toggle(n)
		}
	}

	for y := 3; y <= sqrtN; y += 6 {
		for x := calc4x(d.Start(), y); ; x++ {
			if x%3 == 0 {
				continue
			}
			n := 4*x*x + y*y
			if n >= N {
				break
			}
			d.Toggle(n)
		}
	}
}

func AtokinSqr(d SieveStoreInterface, dn SieveStoreInterface) {
	sqrtN := sqrt(d.Stop())
	if dn.Stop() <= sqrtN {
		sqrtN = dn.Stop() - 1
	}
	// dn.Start() is even. because 2 in base
	for n := dn.Start() + 1; n <= sqrtN; n += 2 {
		if dn.Get(n) {
			a := (d.Start()-1)/(n*n) + 1
			a = a & ^1 + 1 // to odd
			for k := a * n * n; k < d.Stop(); k += n * n * 2 {
				d.Set(k, false)
			}
		}
	}
}

/*
localfunc
*/

func sqrt(n int) int {
	return int(math.Sqrt(float64(n)))
}

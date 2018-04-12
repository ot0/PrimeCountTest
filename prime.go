package main

import (
	"fmt"
	"math"
	"time"
)

const PrimeSize = 100000000000

func main(){
	fmt.Println("Hello World!")
	//var d SieveStore
	
	var d SieveSieve
	base := []int{2, 3, 5, 7, 11, 13, 17, 19}
	d.Init(PrimeSize, base)
	
	start := time.Now()
	//Eratosthenes(&d)
	
	Atokin1(&d)
	Atokin4(&d)
	
	fmt.Println(time.Now().Sub(start).Seconds())
	PrintSieve(&d, false)
}

func PrintSieve(d SieveStoreInterface, all bool){
	count := 0
	for i:=2; i<d.Length(); i++ {
		if d.Get(i) {
			if all{
				fmt.Print(i)
				fmt.Print(",")
			}
			count++
		}
	}
	if all {
		fmt.Println()
	}
	fmt.Println(count)
}

func Eratosthenes(d SieveStoreInterface){
	for i:=2; i<d.Length(); i++{
		d.Set(i, true)
	}

	N := int(math.Sqrt(float64(d.Length())))+1
	for i:=2; i<N; i++ {
		if d.Get(i) {
			for j:=i*2; j < d.Length(); j+=i {
				d.Set(j, false)
			}
		}
	}
}

func Atokin1(d SieveStoreInterface){
	N := d.Length()
	sqrtN := int(math.Sqrt(float64(N)))
	for z:=1; z <=5; z+=4 {
		for y:=z; y<=sqrtN; y+=6 {
			for x:=1;;x++  {
				n:=4*x*x+y*y
				if n>=N {
					break
				}
				d.Toggle(n)
			}
			for x:=y+1;;x+=2 {
				n:=3*x*x-y*y
				if n>=N {
					break
				}
				d.Toggle(n)
			}
		}
	}
	for z:=2; z <=4; z+=2 {
		for y:=z; y<=sqrtN; y+=6 {
			for x:=1;; x+=2 {
				n:=3*x*x+y*y
				if n>=N {
					break
				}
				d.Toggle(n)
			}
			for x:=y+1;; x+=2 {
				n:=3*x*x-y*y
				if n>=N {
					break
				}
				d.Toggle(n)
			}
		}
	}
	for y:=3; y<=sqrtN; y+=6 {
		for z:=1; z <= 2; z++ {
			for x:=z;; x+=3 {
				n:=4*x*x+y*y
				if n>=N {
					break
				}
				d.Toggle(n)
			}
		}
	}
}

func Atokin4(d SieveStoreInterface){
	N := d.Length()
	sqrtN := int(math.Sqrt(float64(N)))
	for n := 5; n<=sqrtN; n++{
		if d.Get(n){
			for k:=n*n; k<N; k+=n*n{
				d.Set(k, false)
			}
		}
	}
	d.Set(2, true)
	d.Set(3, true)
}


type SieveStoreInterface interface{
	Toggle(i int)
	Set(i int, f bool)
	Get(i int) bool
	Length() int
}

type SieveSieve struct{
	base []int
	size int
	step int
	memory [][]uint64
	modF []SSHole
}

func (p *SieveSieve) Init(size int, basePrime []int){
	p.base = basePrime
	p.size = size
	p.step = 1
	for _, i := range p.base {
		p.step *= i
	} 
	p.modF = make([]SSHole, p.step)
	count := 0
	for j:=0; j<p.step; j++ {
		isPrime := true
		for _, i := range p.base{
			if j % i == 0 {
				p.modF[j] = &FalseHole{}
				isPrime = false
				break
			}
		}
		if isPrime {
			p.modF[j] = &ToggleHole{pos:count/64, bit: 1 << uint64(count%64)}
			//fmt.Println(count/64, 1 << uint64(count%64))
			count++
		}
	}
	// ceil
	memsize := (p.size - 1)/ p.step +1
	cellsize := (count -1) / 64 +1
	p.memory = make([][]uint64, memsize)
	for j:=0; j<memsize; j++ {
		p.memory[j] = make([]uint64, cellsize)
	}
	fmt.Println("step:",p.step, "memory:", memsize * cellsize * 8, "count:", count)
}

func (p *SieveSieve) Get(i int) bool{
	pos := i / p.step
	mod := i % p.step
	//fmt.Println(p.memory[pos-1])
	if pos == 0 {
		for _, j := range p.base {
			if j==i {
				return true
			}
		}
	}
	return p.modF[mod].Get(&p.memory[pos])
}
func (p *SieveSieve) Set(i int, f bool){
	pos := i / p.step
	mod := i % p.step
	p.modF[mod].Set(&p.memory[pos], f)
}
func (p *SieveSieve) Toggle(i int){
	pos := i / p.step
	mod := i % p.step
	p.modF[mod].Toggle(&p.memory[pos])
}
func (p *SieveSieve) Length() int{
	return p.size
}

type SSHole interface {
	Get(*[]uint64) bool
	Set(*[]uint64, bool)
	Toggle(*[]uint64)
}

type FalseHole struct {}
func (p *FalseHole) Get(i *[]uint64) bool{
	return false
}
func (p *FalseHole) Set(i *[]uint64, f bool){}
func (p *FalseHole) Toggle(i *[]uint64){}

type ToggleHole struct{
	pos int;
	bit uint64;
}
func (p *ToggleHole) Get(i *[]uint64) bool{
	return (*i)[p.pos] & p.bit == p.bit
}
func (p *ToggleHole) Set(i *[]uint64, f bool){
	g := p.bit
	if !f {
		g = 0
	}
	//fmt.Println(p.pos)
	//fmt.Println((*i)[p.pos])
	(*i)[p.pos] = (*i)[p.pos] & ^p.bit | g
	//fmt.Println((*i)[p.pos])	
}
func (p *ToggleHole) Toggle(i *[]uint64){
	(*i)[p.pos] = (*i)[p.pos] ^ p.bit
}

type SieveStore struct{
	x [PrimeSize]bool
}

func (p *SieveStore) Toggle(i int){
	p.x[i] = !p.x[i]
}

func (p *SieveStore) Set(i int, f bool){
	p.x[i] = f
}

func (p *SieveStore) Get(i int) bool{
	return p.x[i]
}

func (p *SieveStore) Length() int{
	return len(p.x)
}
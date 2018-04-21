package prime

/*-----------------------------------------*/
/* Interface*/
/*-----------------------------------------*/

type SieveStoreInterface interface {
	Toggle(i int)
	Set(i int, f bool)
	Get(i int) bool
	Start() int
	Stop() int
}

/*-----------------------------------------*/
/* SieveStore*/
/*-----------------------------------------*/

type SieveStore struct {
	start int
	stop  int
	x     []bool
}

func (p *SieveStore) Init(start int, step int) {
	p.start = start
	p.stop = start + step
	p.x = make([]bool, step)
}

func (p *SieveStore) Toggle(i int) {
	p.x[i-p.start] = !p.x[i-p.start]
}

func (p *SieveStore) Set(i int, f bool) {
	p.x[i-p.start] = f
}

func (p *SieveStore) Get(i int) bool {
	return p.x[i-p.start]
}

func (p *SieveStore) Start() int {
	return p.start
}

func (p *SieveStore) Stop() int {
	return p.stop
}

/*-----------------------------------------*/
/* Wheel Factorization Store*/
/*-----------------------------------------*/

type WheelStore struct {
	start  int
	stop   int
	memory [][]uint64
	wh     *Wheel
}

func (p *WheelStore) Get(i int) bool {
	pos := (i - p.start) / p.wh.step
	mod := (i - p.start) % p.wh.step
	return p.wh.modF[mod].Get(&p.memory[pos])
}

func (p *WheelStore) Set(i int, f bool) {
	pos := (i - p.start) / p.wh.step
	mod := (i - p.start) % p.wh.step
	p.wh.modF[mod].Set(&p.memory[pos], f)
}
func (p *WheelStore) Toggle(i int) {
	pos := (i - p.start) / p.wh.step
	mod := (i - p.start) % p.wh.step
	//fmt.Println(i, p.start, p.stop, len(p.memory), pos)
	p.wh.modF[mod].Toggle(&p.memory[pos])
}

func (p *WheelStore) Start() int {
	return p.start
}
func (p *WheelStore) Stop() int {
	return p.stop
}

/*-----------------------------------------*/
/* Wheel Factorization Divive Store at first*/
/*-----------------------------------------*/
type WheelStoreF struct {
	WheelStore
}

func (p *WheelStoreF) Get(i int) bool {
	pos := i / p.wh.step
	mod := i % p.wh.step
	if pos == 0 {
		for _, j := range p.wh.base {
			if j == i {
				return true
			}
		}
	}
	return p.wh.modF[mod].Get(&p.memory[pos])
}

/*-----------------------------------------*/
/* Wheel */
/*-----------------------------------------*/

type Wheel struct {
	base  []int
	step  int
	modF  []Hole
	count int
}

func (p *Wheel) Init(basePrime []int) {
	p.base = basePrime
	p.step = 1
	for _, i := range p.base {
		p.step *= i
	}
	p.modF = make([]Hole, p.step)
	count := 0
	noneHole := &FalseHole{}
	for j := 0; j < p.step; j++ {
		isBase := false
		for _, i := range p.base {
			if j%i == 0 {
				p.modF[j] = noneHole
				isBase = true
				break
			}
		}
		if !isBase {
			p.modF[j] = &ToggleHole{pos: count / 64, bit: 1 << uint64(count%64)}
			//fmt.Println(count/64, 1 << uint64(count%64))
			count++
		}
	}
	p.count = count
	//fmt.Println("step:", p.step, "count:", p.count, "rate", float64(p.step)/float64(p.count))
}

func (p *Wheel) MakeMemory(size int) [][]uint64 {
	// ceil
	memsize := (size-1)/p.step + 1
	cellsize := (p.count-1)/64 + 1
	mem := make([][]uint64, memsize)
	for j := 0; j < memsize; j++ {
		mem[j] = make([]uint64, cellsize)
	}
	return mem
}

func (p *Wheel) NewStoreF(start int, step int) SieveStoreInterface {
	/* if need first prime */
	x := &WheelStoreF{}
	x.wh = p
	x.start = start
	x.stop = start + step
	x.memory = p.MakeMemory(step)
	return x
}

func (p *Wheel) NewStore(start int, step int) SieveStoreInterface {
	x := &WheelStore{}
	x.wh = p
	x.start = start
	x.stop = start + step
	x.memory = p.MakeMemory(step)
	return x
}

/*-----------------------------------------*/
/* Hole*/
/*-----------------------------------------*/

type Hole interface {
	Get(*[]uint64) bool
	Set(*[]uint64, bool)
	Toggle(*[]uint64)
}

type TrueHole struct{}

func (p *TrueHole) Get(i *[]uint64) bool {
	return true
}
func (p *TrueHole) Set(i *[]uint64, f bool) {}
func (p *TrueHole) Toggle(i *[]uint64)      {}

type FalseHole struct{}

func (p *FalseHole) Get(i *[]uint64) bool {
	return false
}
func (p *FalseHole) Set(i *[]uint64, f bool) {}
func (p *FalseHole) Toggle(i *[]uint64)      {}

type ToggleHole struct {
	pos int
	bit uint64
}

func (p *ToggleHole) Get(i *[]uint64) bool {
	return (*i)[p.pos]&p.bit == p.bit
}
func (p *ToggleHole) Set(i *[]uint64, f bool) {
	g := p.bit
	if !f {
		g = 0
	}
	(*i)[p.pos] = (*i)[p.pos] & ^p.bit | g
}
func (p *ToggleHole) Toggle(i *[]uint64) {
	(*i)[p.pos] = (*i)[p.pos] ^ p.bit
}

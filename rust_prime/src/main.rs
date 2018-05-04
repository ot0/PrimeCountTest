use std::vec::Vec;
use std::time::Instant;

fn main() {
    println!("Hello, world!");

    let start = Instant::now();

    const BUF_SIZE: usize = 1000*1000*100 - 1;
    let mut store = SieveStore{
        s: 2,
        x: vec![true; BUF_SIZE]
    };
    
    let mut count = 0;
    for i in store.start()..store.stop(){
        if store.get(i) {
            let mut j = i * i;
            while j < store.stop() {
                store.set(j, false);
                j += i;
            }
            count +=1;
            //print!("{},", i)
        }
    }
    //println!("");

    let end = start.elapsed();
    println!("{}, {:?}", count, end);
}

trait PrimeStore {
    fn get(&self, usize) -> bool;
    fn set(&mut self, usize, bool);
    fn toggle(&mut self, usize);
    fn start(&self) -> usize;
    fn stop(&self) -> usize;
}

struct SieveStore {
    s: usize,
    x: Vec<bool>
}

impl PrimeStore for SieveStore {
    fn get(&self, i:usize) -> bool {
        self.x[i - self.s]
    }
    fn set(&mut self, i:usize, f:bool){
        self.x[i - self.s] = f;
    }
    fn toggle(&mut self, i:usize){
        self.x[i - self.s] = !self.x[i - self.s];
    }
    fn start(&self) -> usize {
        self.s
    }
    fn stop(&self) -> usize {
        self.s + self.x.len()
    }
}
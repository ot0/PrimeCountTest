use std::vec::Vec;
use std::time::Instant;
use std::thread;
use std::sync::{mpsc, Arc};

fn main() {
    println!("Hello, world!");

    let stop: usize = 1000*1000*1000*1;

    let start = Instant::now();
    let count = eratosthenes_single(stop);
    let end = start.elapsed();
    println!("{}, {:?}", count, end);

    let start = Instant::now();
    let count = eratosthenes_multi(stop);
    let end = start.elapsed();
    println!("{}, {:?}", count, end);

}

fn eratosthenes_single(stop: usize) -> usize {
    let mut store = SieveStore::new(2, stop -1);
    eratosthenes(&mut store).iter().count()
}

fn eratosthenes_multi(stop: usize) -> usize {
    let process = 4;
    let step = isqrt(stop);

    let mut ss0 = Arc::new(SieveStore::new(2, step));
    
    
    eratosthenes(Arc::get_mut(&mut ss0).unwrap());
    //eratosthenes(&mut *ss0);
    //ss0.iter().for_each(|x| print!("{},",x));

    let mut count = ss0.iter().count();
    
    let mut start = ss0.stop();

    let (tx, rx) = mpsc::channel();

    for _i in 0..process {
        let tx = tx.clone();
        tx.send(0).unwrap();
    }

    let mut is_continue = true;
    while is_continue{
        count += rx.recv().unwrap();
        let tmp_stop = start + step;       
        let mut ss = if tmp_stop < stop {
            SieveStore::new(start, step)
        }else{
            is_continue = false;
            SieveStore::new(start, stop - start + 1)
        };
        let tx = tx.clone();
        let ss0 = ss0.clone();
        thread::spawn(move ||{
            eratosthenes_other(&mut ss, &(*ss0));
            tx.send(ss.iter().count()).unwrap();
        });
        start = tmp_stop;
    }
    for _i in 0..process {
        count += rx.recv().unwrap();
    }
    //println!("");
    count
}

fn isqrt(i: usize) -> usize {
    (i as f64).sqrt() as usize
}

fn eratosthenes(ps: &mut PrimeStore) ->&mut PrimeStore{
    for i in ps.start()..ps.stop(){
        if ps.get(i) {
            let mut j = i * i;
            while j < ps.stop() {
                ps.set(j, false);
                j += i;
            }
        }
    }   
    ps
}

fn eratosthenes_other(ps: &mut PrimeStore, ps0: &PrimeStore){
    for n in ps0.iter() {
        let s = (ps.start()-1) / n + 1;
        let mut i = n * s;
        while i < ps.stop(){
            ps.set(i, false);
            i += n;
        }
    }
}

trait PrimeStore {
    fn get(&self, usize) -> bool;
    fn set(&mut self, usize, bool);
    fn toggle(&mut self, usize);
    fn start(&self) -> usize;
    fn stop(&self) -> usize;
    fn iter(&self) -> SieveStoreIter;
}

struct SieveStore {
    s: usize,
    x: Vec<bool>,
}

struct SieveStoreIter<'a>{
    ss: &'a PrimeStore,
    i: usize
}

impl SieveStore{
    fn new(start:usize, step:usize) -> SieveStore{
        SieveStore{
            s: start,
            x: vec![true; step],
        } 
    }
}

impl<'a> Iterator for SieveStoreIter<'a>{
    type Item = usize;
    fn next(&mut self) -> Option<usize> {
        loop {
            if self.i >= self.ss.stop() {
                return None;
            }
            if self.ss.get(self.i) {
                break
            }
            self.i += 1;
        }
        let n = self.i;
        self.i += 1;
        Some(n)
    }
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
    fn iter(&self) -> SieveStoreIter {
        SieveStoreIter{
            ss: self,
            i: self.start()
        }
    }
}
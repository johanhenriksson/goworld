package main

import "fmt"

type Random struct {
    state       int32
}

func NewRandom(seed int32) *Random {
    return &Random { state: seed }
}

func (r *Random) Seed(seed int32) {
    r.state = seed
}

func (r *Random) Next() int32 {
    val := ((r.state * 1103515245) + 12345) & 0x7fffffff
    r.state = val
    return val
}

func (r *Random) Range(min, max int32) int32 {
    v := r.Next()
    return v % (max - min) + min
}

func main() {
    rnd := NewRandom(5)
    for i := 0; i < 10; i++ {
        fmt.Println(rnd.Range(10,15))
    }
}

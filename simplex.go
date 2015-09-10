package main

import (
    "fmt"
    opensimplex "github.com/ojrac/opensimplex-go"
)

func main() {
    size := 10
    simplex := opensimplex.NewWithSeed(0)
    for y := 0; y < size; y++ {
        for x := 0; x < size; x++ {
            v := simplex.Eval2(float64(x),float64(y))
            fmt.Printf("%.3f ", v)
        }
        fmt.Printf("\n")
    }
}

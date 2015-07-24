package util

import (
    "os"
    "path/filepath"
)

var ExePath string = "."

func init() {
    var err error
    ExePath, err = filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        panic(err)
    }
}

package engine

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func RunProfilingServer(port int) {
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil); err != nil {
		log.Println("failed to launch profiling http server on port", port)
	} else {
		log.Printf("pprof server available at http://localhost:%d\n", port)
	}
}

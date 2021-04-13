package common

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// SignalCheck catch terminate
func SignalCheck(send chan<- int, fn ...func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	go func() {
		<-sigs
		for _, fnItem := range fn {
			fnItem()
		}
		fmt.Println(`Byeee..`)
		send <- 1
	}()
}

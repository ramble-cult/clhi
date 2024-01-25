package chat

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/context"
)

const timeFormat = "03:04:05 PM"

func ClientLogf(ts time.Time, format string, args ...interface{}) {
	log.Printf("[%s] <<Client>>: "+format, append([]interface{}{ts.Format(timeFormat)}, args...)...)
}

func ServerLogf(ts time.Time, format string, args ...interface{}) {
	log.Printf("[%s] <<Server>>: "+format, append([]interface{}{ts.Format(timeFormat)}, args...)...)
}

func MessageLog(ts time.Time, name, msg string) {
	log.Printf("[%s] %s: %s", ts.Format(timeFormat), name, msg)
}

// func
// 	if !debugMode {
// 		return
// 	}
// 	log.Printf("[%s] <<Debug>>: "+format, append([]interface{}{time.Now().Format(timeFormat)}, args...)...)
// }

func SignalContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		<-sigs

		signal.Stop(sigs)
		close(sigs)
		cancel()
	}()

	return ctx
}

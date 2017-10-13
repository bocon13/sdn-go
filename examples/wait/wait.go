package wait

import (
	"os"
	"os/signal"
	"syscall"
)

// Wait for SIGINT or SIGTERM
func UntilKilled() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alerting"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/watcher"
)

func main() {
	cfgPath := flag.String("config", "", "path to config file (optional)")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("portwatch: failed to load config: %v", err)
	}

	scanner := portscanner.NewScanner(cfg.ProcNetTCP)
	alerter := alerting.NewAlerter(os.Stdout)
	w := watcher.New(cfg, scanner, alerter)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := w.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("portwatch: watcher exited with error: %v", err)
	}
}

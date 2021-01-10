package main

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-ping/ping"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func checkLoss(packetLoss prometheus.Gauge, target string) {
	for {
		time.Sleep(time.Second)

		pinger, err := ping.NewPinger(target)
		if err != nil {
			log.Print("ping.NewPinger error: ", err)
			continue
		}

		if runtime.GOOS == "linux" {
			pinger.SetPrivileged(true)
		}
		pinger.Timeout = time.Second

		err = pinger.Run()
		if err != nil {
			log.Print("pinger.Run error: ", err)
			continue
		}

		stats := pinger.Statistics()
		packetLoss.Set(stats.PacketLoss)
	}
}

func main() {
	registry := prometheus.NewRegistry()
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)

	target := os.Getenv("TARGET_HOST")
	if target == "" {
		log.Fatal("Missing environment variable: TARGET_HOST")
	}

	packetLoss := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "packet_loss",
		Help: "Number of lost packets",
	})
	registry.MustRegister(packetLoss)
	go checkLoss(packetLoss, target)

	log.Print("Starting server")
	http.ListenAndServe(":2112", nil)
}

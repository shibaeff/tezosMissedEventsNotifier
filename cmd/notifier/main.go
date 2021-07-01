package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"tezos/missedEventsNotifier/internal/configs"
	"tezos/missedEventsNotifier/internal/scheduling"
	"tezos/missedEventsNotifier/pkg/api"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()
	config, err := configs.GetConfig("./config/config.yaml")
	if err != nil {
		log.Fatalln("Failed to read config")
	}
	cycle, err := strconv.Atoi(config.Cycle)
	if err != nil {
		log.Fatalln(err)
	}
	tzApi := api.NewApi(config.Host, config.Delegate, cycle)
	scheduler := scheduling.NewScheduler(tzApi)
	scheduler.EndorsementsWg().Add(1)
	scheduler.ScheduleEndorsements()
	scheduler.BakingsWg().Add(1)
	scheduler.ScheduleBakings()
	scheduler.BakingsWg().Wait()
	scheduler.EndorsementsWg().Wait()
}

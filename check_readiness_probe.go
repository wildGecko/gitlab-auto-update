package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
)

func getReadinessProbes() [10]string {
	gitlabUrl := os.Getenv("GITLAB_URL")
	token := os.Getenv("GITLAB_PROBES_TOKEN")
	readinessEndpoint := fmt.Sprintf(gitlabUrl+readinessPartOfUrl, token)
	var readProbe ReadinessProbe

	request, err := http.NewRequest("GET", readinessEndpoint, nil)
	if err != nil {
		log.Error(err)
	}

	result, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error(err)
	}

	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	err = json.Unmarshal(body, &readProbe)
	if err != nil {
		log.Error(err)
	}
	statuses := [10]string{
		readProbe.Status,
		readProbe.MasterCheck[0].Status,
		readProbe.DbCheck[0].Status,
		readProbe.CacheCheck[0].Status,
		readProbe.QueuesCheck[0].Status,
		readProbe.RateLimitingCheck[0].Status,
		readProbe.SessionsCheck[0].Status,
		readProbe.SharedStateCheck[0].Status,
		readProbe.TraceChunksCheck[0].Status,
		readProbe.GitalyCheck[0].Status,
	}
	return statuses
}

func checkReadinessProbes() {
	i := 0
	n := 0
Loop:
	for {
		statuses := getReadinessProbes()
		if i > 1 {
			log.Error("A lot of errors in readiness probes")
			notificationToSlack("no")
			sendAlert()
			break
		} else {
			for _, st := range statuses {
				if st != "ok" && i <= 1 {
					i++
					n++
					log.Error("Failed checking of readiness probe")
					log.Info("Retry check for 5 minutes")
					if i != 2 {
						time.Sleep(5 * time.Minute)
					}
					continue Loop
				}
			}
			break
		}
	}
	if n == 0 {
		log.Info("All readiness probes are ready")
	}

}

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

func getLivenessProbe() string {
	gitlabUrl := os.Getenv("GITLAB_URL")
	token := os.Getenv("GITLAB_PROBES_TOKEN")
	livenessEndpoint := fmt.Sprintf(gitlabUrl+livenessPartOfUrl, token)
	var livProbe LivenessProbe

	request, err := http.NewRequest("GET", livenessEndpoint, nil)
	if err != nil {
		log.Error(err)
	}

	result, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error(err)
	}

	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	err = json.Unmarshal(body, &livProbe)
	if err != nil {
		log.Error(err)
	}

	return livProbe.Status
}

func checkLivenessProbe() {
	i := 0
	n := 0
Loop:
	for {
		status := getLivenessProbe()
		if i > 1 {
			log.Error("A lot of errors in liveness probes")
			notificationToSlack("no")
			sendAlert()
			break
		} else {
			if status != "ok" && i <= 1 {
				i++
				n++
				log.Error("Failed checking of liveness probe")
				log.Info("Retry check for 5 minutes")
				if i != 2 {
					time.Sleep(5 * time.Minute)
				}
				continue Loop
			}
			break
		}
	}
	if n == 0 {
		log.Info("All liveness probes are ready")
	}
}

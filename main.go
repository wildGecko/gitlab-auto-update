package main

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	byt        int64
	gigabyte   float64
	appVersion = ""
)

const (
	checkUrl           = "https://version.gitlab.com/check.svg?gitlab_info="
	livenessPartOfUrl  = `/-/liveness?token=%s`
	readinessPartOfUrl = `/-/readiness?all=1&token=%s`
	madisonUrl         = `https://madison.flant.com/api/events/custom/%s`
	B                  = 1
	KB                 = 1024 * B
	MB                 = 1024 * KB
	GB                 = 1024 * MB
	statusFileName     = "/opt/gitlab-updater/statusUpdate.txt"
	fileWithBackupName = "/opt/gitlab-updater/backupName.txt"
)

func initLog() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	initLog()
	log.Info("Start app...")
	log.Info("Version: ", appVersion)
	initEnv()
	runTask()
}

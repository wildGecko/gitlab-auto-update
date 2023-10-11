package main

import (
	"os"
	"regexp"

	"github.com/bombsimon/logrusr/v2"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func updateCheck() {
	log.Info("Update check...")
	updateListPackages()
	r, _ := regexp.Compile("\\d{1,}.\\d{1,}.\\d{1,}")
	statusUpdate := gitlabCheckUpdate()
	targetVersion := r.FindString(findTargetVersion())
	currentVersion := r.FindString(getCurrentVersion())
	if statusUpdate == "update available" || statusUpdate == "update asap" {
		if targetVersion == currentVersion {
			log.Info("All minor versions are already installed")
			writeToStatusFile("NO")
		} else {
			writeToStatusFile("YES")
			notificationToSlack("needUpdate")
		}
	} else {
		writeToStatusFile("NO")
	}
}

func updateInstall() {
	updateStatus := readStatusFile()
	if updateStatus == "YES" {
		notificationToSlack("startUpdate")
		freeSpace := checkDiskSize()
		backupSize := getBackupSize()
		createBackup(freeSpace, backupSize)
		checkMigrations()
		checkQueue()
		log.Info("START UPDATE")
		updateListPackages()
		updateGitlab()
		checkLivenessProbe()
		checkReadinessProbes()
		notificationToSlack("ok")
		deleteFiles(statusFileName)
		backupName := readFileWithBackupName()
		deleteFiles(backupName)
		checkBackupsFiles()
	} else {
		log.Info("All minor versions are already installed")
		deleteFiles(statusFileName)
	}
}

func runTask() {
	updateCheckTime := os.Getenv("TIME_UPDATE_CHECK")
	updateInstallTime := os.Getenv("TIME_UPDATE_INSTALL")
	logger := logrusr.New(log.StandardLogger(), logrusr.WithName("cron"))
	c := cron.New(
		cron.WithChain(
			cron.SkipIfStillRunning(logger),
			cron.Recover(logger),
		),
		cron.WithLogger(logger),
	)
	if _, err := c.AddFunc(updateCheckTime, func() { updateCheck() }); err != nil {
		log.Fatal(err)
	}
	if _, err := c.AddFunc(updateInstallTime, func() { updateInstall() }); err != nil {
		log.Fatal(err)
	}
	c.Run()
}

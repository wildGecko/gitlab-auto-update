package main

import (
	"bufio"
	logs "github.com/pieterclaerhout/go-log"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strconv"
)

func checkMigrations() {
	log.Info("Checking migrations...")
	gitlabRailsMigrations := "bash"
	arg0 := "-c"
	arg1 := "gitlab-rails runner -e production 'puts Gitlab::BackgroundMigration.remaining'"

	cmd := exec.Command(gitlabRailsMigrations, arg0, arg1)
	r, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	done := make(chan struct{})
	scanner := bufio.NewScanner(r)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			log.Info("Count of migrations: ", line)
			if line == strconv.Itoa(0) {
				log.Info("NO MIGRATIONS")
			} else {
				log.Info("There are some migrations.")
				checkMigrations()
			}
		}
		done <- struct{}{}
	}()

	err := cmd.Start()
	logs.CheckError(err)

	<-done

	err = cmd.Wait()
	logs.CheckError(err)
}

func checkQueue() {
	log.Info("Checking queue...")
	gitlabRailsQueue := "bash"
	arg0 := "-c"
	arg1 := "gitlab-rails runner -e production 'puts Gitlab::Database::BackgroundMigration::BatchedMigration.queued.count'"

	cmd := exec.Command(gitlabRailsQueue, arg0, arg1)

	r, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	done := make(chan struct{})
	scanner := bufio.NewScanner(r)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			log.Info("Count of job queue: ", line)
			if line == strconv.Itoa(0) {
				log.Info("Queue is empty.")
			} else {
				log.Info("Queue is not empty.")
				checkQueue()
			}
		}
		done <- struct{}{}
	}()

	err := cmd.Start()
	logs.CheckError(err)

	<-done

	err = cmd.Wait()
	logs.CheckError(err)
}

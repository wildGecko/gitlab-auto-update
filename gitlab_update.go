package main

import (
	"bufio"
	"fmt"
	logs "github.com/pieterclaerhout/go-log"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func updateListPackages() {
	aptUpdate := "apt"
	arg := "update"

	log.Info("Start update list of packages")
	cmd := exec.Command(aptUpdate, arg)
	r, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	done := make(chan struct{})
	scanner := bufio.NewScanner(r)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			log.Info(line)
		}
		done <- struct{}{}
	}()

	err := cmd.Start()
	logs.CheckError(err)

	<-done

	err = cmd.Wait()
	logs.CheckError(err)
}

func updateGitlab() {
	targetVersion := findTargetVersion()
	bash := "bash"
	arg0 := "-c"
	aptInstall := fmt.Sprintf("apt install -y gitlab-ce=%s", targetVersion)

	log.Info("Start update GitLab")
	cmd := exec.Command(bash, arg0, aptInstall)
	r, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	scanner := bufio.NewScanner(r)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			log.Info(line)
		}
	}()

	if err := cmd.Start(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if code := exiterr.Error(); code != "0" {
				notificationToSlack("no")
			}
		}
		logs.CheckError(err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if code := exiterr.Error(); code != "0" {
				notificationToSlack("no")
			}
		}
		logs.CheckError(err)
	}
}

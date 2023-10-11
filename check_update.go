package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func getCurrentVersion() string {
	var version GilabVersion
	gitlabUrl := os.Getenv("GITLAB_URL")
	versionUrl := gitlabUrl + `/api/v4/version`
	accessToken := os.Getenv("GITLAB_API_TOKEN")

	req, err := http.NewRequest("GET", versionUrl, nil)

	if err != nil {
		log.Error(err)
	}

	req.Header.Add("PRIVATE-TOKEN", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	err = json.Unmarshal(body, &version)
	if err != nil {
		log.Error(err)
	}

	return fmt.Sprintf("{\"version\":\"%s\"}", version.Version)

}

func encodingVersion() string {
	version := getCurrentVersion()
	b64Version := base64.StdEncoding.EncodeToString([]byte(version))
	return b64Version
}

func gitlabCheckUpdate() string {
	version := encodingVersion()
	gitlabUrl := os.Getenv("GITLAB_URL")

	var status GitlabUpdateStatus

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", checkUrl, version), nil)
	if err != nil {
		log.Error(err)
	}

	req.Header.Set("Referer", gitlabUrl)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal(body, &status)

	return status.Text.Text
}

func findTargetVersion() string {
	regCurrentVersion := regexp.MustCompile("[0-9]{2}")
	currentVersion := regCurrentVersion.FindAllString(getCurrentVersion(), -1)[0]

	bash := "bash"
	arg0 := "-c"
	arg1 := fmt.Sprintf("apt-cache madison gitlab-ce | awk '{print $3}' | grep ^%s | head -n 1", currentVersion)

	cmd := exec.Command(bash, arg0, arg1)
	std, _ := cmd.Output()
	ver := strings.Trim(string(std), "\n")

	return ver
}

func writeToStatusFile(status string) {
	file, err := os.Create(statusFileName)
	if err != nil {
		log.Fatal("Failed to create file. ", err)
	}
	defer file.Close()

	_, err = file.WriteString(status)
	if err != nil {
		log.Fatal("Failed to write to file. ", err)
	}
}

func readStatusFile() string {
	updateStatus, err := os.ReadFile(statusFileName)
	if err != nil {
		log.Fatal("Failed to read file. ", err)
	}

	return string(updateStatus)
}

func deleteFiles(name string) {
	deleteFile := os.Remove(name)
	if deleteFile != nil {
		log.Error("Failed to delete file. ", deleteFile)
	}
}

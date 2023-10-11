package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type logger struct{}

func (l logger) Output(i int, out string) error {
	log.WithField("logger", "slack").Print(out)
	return nil
}

func notificationToSlack(status string) {
	token := os.Getenv("SLACK_AUTH_TOKEN")
	channelID := os.Getenv("SLACK_CHANNEL_ID")
	updateInstallTime := os.Getenv("TIME_UPDATE_INSTALL")
	schedule, err := cron.ParseStandard(updateInstallTime)
	if err != nil {
		log.Fatal(err)
	}
	r, _ := regexp.Compile("\\d{1,}.\\d{1,}.\\d{1,}")
	targetVersion := r.FindString(findTargetVersion())
	currentVersion := r.FindString(getCurrentVersion())

	messageUpdate := fmt.Sprintf("New version available. The update will be installed at %s.", schedule.Next(time.Now()).Format(time.RFC822))
	client := slack.New(token, slack.OptionDebug(true), slack.OptionLog(logger{}))

	switch status {
	case "needUpdate":
		attach := slack.Attachment{
			Text:  "GitLab update",
			Color: "#36a64f",
			Fields: []slack.AttachmentField{
				{
					Title: "Status",
					Value: messageUpdate,
				}, {
					Title: "Current version",
					Value: currentVersion,
				}, {
					Title: "Target version",
					Value: targetVersion,
				},
			},
		}
		_, _, err := client.PostMessage(channelID, slack.MsgOptionAttachments(attach))
		if err != nil {
			log.Error(err)
		}
	case "startUpdate":
		attach := slack.Attachment{
			Text:  "GitLab update",
			Color: "#36a64f",
			Fields: []slack.AttachmentField{
				{
					Title: "Status",
					Value: "GitLab update launched",
				},
			},
		}
		_, _, err := client.PostMessage(channelID, slack.MsgOptionAttachments(attach))
		if err != nil {
			log.Error(err)
		}
	case "backupOk":
		attach := slack.Attachment{
			Text:  "Create backup of GitLab",
			Color: "#36a64f",
			Fields: []slack.AttachmentField{
				{
					Title: "Status",
					Value: "Backup successfully completed.",
				},
			},
		}
		_, _, err := client.PostMessage(channelID, slack.MsgOptionAttachments(attach))
		if err != nil {
			log.Error(err)
		}
	case "backupErr":
		attach := slack.Attachment{
			Text:  "Create backup of GitLab",
			Color: "#FF0000",
			Fields: []slack.AttachmentField{
				{
					Title: "Status",
					Value: "Backup failed.",
				},
			},
		}
		_, _, err := client.PostMessage(channelID, slack.MsgOptionAttachments(attach))
		if err != nil {
			log.Error(err)
		}
	case "ok":
		attach := slack.Attachment{
			Text:  "GitLab update",
			Color: "#36a64f",
			Fields: []slack.AttachmentField{
				{
					Title: "Status",
					Value: "Update installation successfully completed.",
				},
			},
		}
		_, _, err := client.PostMessage(channelID, slack.MsgOptionAttachments(attach))
		if err != nil {
			log.Error(err)
		}
	case "no":
		attach := slack.Attachment{
			Text:  "GitLab update",
			Color: "#FF0000",
			Fields: []slack.AttachmentField{
				{
					Title: "Status",
					Value: "Update installation failed.",
				},
			},
		}
		_, _, err := client.PostMessage(channelID, slack.MsgOptionAttachments(attach))
		if err != nil {
			log.Error(err)
		}
	case "cancel":
		attach := slack.Attachment{
			Text:  "GitLab update",
			Color: "#FF0000",
			Fields: []slack.AttachmentField{
				{
					Title: "Status",
					Value: "Update canceled due to backup problems.",
				},
			},
		}
		_, _, err := client.PostMessage(channelID, slack.MsgOptionAttachments(attach))
		if err != nil {
			log.Error(err)
		}
	}
}

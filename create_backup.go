package main

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"syscall"
	"time"

	logs "github.com/pieterclaerhout/go-log"
	log "github.com/sirupsen/logrus"
)

func checkDiskSize() float64 {
	var disk DiskInfo
	log.Info("Check available space...")
	path := os.Getenv("GITLAB_BACKUP_DIR")

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		log.Fatal("Error reading directory: ", err)
	}
	disk.Size = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.Size - disk.Free

	freeSpace := float64(disk.Free) / float64(GB)
	log.Info("Free: ", freeSpace, "GB")
	return freeSpace
}

func getBackupSize() float64 {
	root := os.Getenv("GITLAB_BACKUP_DIR")
	os.Chdir(root)
	var sort_struct []AllFiles

	all_files := make(map[string]string)
	log.Info("Get size of backup...")

	file, err := os.Open(".")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Info("No backup file was found. I will create the first backup")
		} else {
			log.Error(err)
		}
	}

	defer file.Close()

	files, err := file.Readdir(-1)
	if err != nil {
		log.Error(err)
	}

	for _, f := range files {
		info, err := os.Stat(f.Name())
		if err != nil {
			log.Error("No open file: ", err)
		}

		all_files[f.Name()] = info.ModTime().String()
	}

	for name, date := range all_files {
		sort_struct = append(sort_struct, AllFiles{name, date})
	}

	sort.Slice(sort_struct, func(i, j int) bool {
		return sort_struct[i].Date > sort_struct[j].Date
	})

	if len(sort_struct) == 0 {
		log.Info("No backup files")
		return 0
	} else {
		backup, err := os.Open(sort_struct[0].Name)

		defer backup.Close()

		backupSize, err := backup.Stat()
		if err != nil {
			log.Error(err)
		}

		byt = backupSize.Size()
	}

	gigabyte = ((float64)(byt / 1024 / 1024 / 1024))

	log.Info("Last backup file - ", sort_struct[0].Name)
	log.Info("Size of last backup file: ", gigabyte, "GB")
	return gigabyte
}

func createBackup(freeSpace, backupSize float64) {
	factorFreeSpace := os.Getenv("RATE_SIZE")
	rateSize, _ := strconv.ParseFloat(factorFreeSpace, 64)
	bash := "bash"
	arg0 := "-c"
	arg1 := "gitlab-rake gitlab:backup:create SKIP=registry"
	spaceReq := backupSize * rateSize

	log.Info("Space requirements - ", spaceReq, ". Checking...")

	if spaceReq < freeSpace {
		log.Info("Check for running backup and restore tasks")
		for {
			if _, err := os.Stat("/opt/gitlab/embedded/service/gitlab-rails/tmp/backup_restore.pid"); err != nil {
				log.Info(err)
				break
			}
			log.Info("File exists: /opt/gitlab/embedded/service/gitlab-rails/tmp/backup_restore.pid")
			log.Info("There is another backup and restore task in progress. Waiting for completion...")
			time.Sleep(time.Minute * 1)
		}

		log.Info("Check is successful. Creating backup...")
		cmd := exec.Command(bash, arg0, arg1)
		r, _ := cmd.StdoutPipe()
		cmd.Stderr = cmd.Stdout

		done := make(chan struct{})

		scanner := bufio.NewScanner(r)

		go func() {
			i := 0
			for scanner.Scan() {
				line := scanner.Text()
				log.Info(line)
				substr := "Backup\\s\\d{1,}\\w\\d{4}\\w\\d+\\w\\d+\\w\\d+\\.\\d{1,}\\.\\d{1,}\\sis\\sdone"
				substrBackupName := "archive:\\s\\d{1,}_\\d{1,}_\\d{1,}_\\d{1,}_\\d{1,}\\.\\d{1,}\\.\\d{1,}_\\w{1,}\\.tar\\s\\...\\sdone"
				matched, _ := regexp.MatchString(substr, line)
				if matched == true {
					i++
				}
				matchedName, _ := regexp.MatchString(substrBackupName, line)
				if matchedName == true {
					re, _ := regexp.Compile(`(\d{1,}_\d{1,}_\d{1,}_\d{1,}_\d{1,}\.\d{1,}\.\d{1,}_\w{1,}\.tar)`)
					res := re.FindAllString(line, -1)
					log.Info("Name of backup is ", res[0])
					writeFileWithBackupName(res[0])
				}
			}
			if i > 0 {
				notificationToSlack("backupOk")
			} else {
				notificationToSlack("backupErr")
				notificationToSlack("cancel")
				backupName := readFileWithBackupName()
				deleteFiles(backupName)
				checkBackupsFiles()
			}
			done <- struct{}{}
		}()

		err := cmd.Start()
		logs.CheckError(err)

		<-done

		err = cmd.Wait()
		logs.CheckError(err)

	} else {
		notificationToSlack("Cancel")
		log.Fatal("No left space...")
		backupName := readFileWithBackupName()
		deleteFiles(backupName)
		checkBackupsFiles()
	}

}

func writeFileWithBackupName(name string) {
	file, err := os.Create(fileWithBackupName)
	if err != nil {
		log.Error("Failed to create file. ", err)
	}
	defer file.Close()

	_, err = file.WriteString(name)
	if err != nil {
		log.Error("Failed to write to file. ", err)
	}
}

func readFileWithBackupName() string {
	updateStatus, err := os.ReadFile(fileWithBackupName)
	if err != nil {
		log.Error("Failed to read file. ", err)
	}

	return string(updateStatus)
}

func checkBackupsFiles() {
	filesList := [10]string{
		"/var/opt/gitlab/backups/backup_information.yml",
		"/var/opt/gitlab/backups/db",
		"/var/opt/gitlab/backups/uploads.tar.gz",
		"/var/opt/gitlab/backups/builds.tar.gz",
		"/var/opt/gitlab/backups/artifacts.tar.gz",
		"/var/opt/gitlab/backups/pages.tar.gz",
		"/var/opt/gitlab/backups/lfs.tar.gz",
		"/var/opt/gitlab/backups/terraform_state.tar.gz",
		"/var/opt/gitlab/backups/packages.tar.gz",
		"/var/opt/gitlab/backups/ci_secure_files.tar.gz",
	}

	for i := 0; i < len(filesList); i++ {
		_, err := os.Stat(filesList[i])
		if err != nil {
			if os.IsNotExist(err) {
				log.Info("File - ", filesList[i], "does not exists")
			}
		} else {
			log.Info("Delete file - ", filesList[i])
			deleteFiles(filesList[i])
		}

	}

}

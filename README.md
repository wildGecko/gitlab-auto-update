## Утилита для обновления гитлаба

### ВНИМАНИЕ!
Обновление происходит только в рамках текущей мажорной версии, то есть устанавливаются только минорные версии. Мажорная версия автоматически НЕ СТАВИТСЯ.

### Описание переменных конфигурационного файла

- `GITLAB_URL` - адрес сервера GitLab
- `GITLAB_API_TOKEN` - токен для доступа к API. Достаточно выставить *scope* ```read_api```
- `GITLAB_BACKUP_DIR` - указывается каталог, в котором хранятся бэкапы на ВМ с самим GitLab'ом
- `RATE_SIZE` - указывается множитель запаса свободного места на диске от размера последнего бэкапа. Например, если последний бэкап весит 25GB, то при значение `2` свободно должно быть 50GB
- `GITLAB_PROBES_TOKEN` - токен для проверки ```probes```. Берётся по пути `Admin -> Moitoring -> Healt Check`
- `SLACK_AUTH_TOKEN` - указываем auth token
- `SLACK_CHANNEL_ID` - указываем ID канала
- `PROJECT_NAME` - имя проекта, в котором находится GitLab
- `TIME_UPDATE_CHECK` - время проверки наличия обновления (по времени сервера, где установлен GitLab. Обычно это UTC) в формате crontab.
- `TIME_UPDATE_INSTALL` - время установки обновления (по времени сервера, где установлен GitLab. Обычно это UTC) в формате crontab.


### Как пользоваться
1. Добавляем бота (Gitlab Update Bot) в клиентский канал.
2. В GitLab создаём технического пользователя (можно использовать текущего, если такой уже есть).
3. Создаём ему токен со *scope* ```read_api```.
4. Копируем ```gitlab-updater``` на cервер GitLab  в каталог `/opt/gitlab-updater`.
5. Копируем `check-time` на cервер GitLab (любой каталог) и выполняем. Вы увидите в каком виде/формате нужно заполнить переменные `TIME_UPDATE_CHECK` и `TIME_UPDATE_INSTALL`.
6. Создаём файл .env в `/opt/gitlab-updater` и заполняем его по примеру. **ВНИМАНИЕ!!!** Это только пример. Подставляйте свои значения. Расшифровка переменных выше.
```shell
GITLAB_URL=URL
GITLAB_API_TOKEN=XXX88hk89798YXY-sfds
GITLAB_BACKUP_DIR=/var/opt/gitlab/backups
RATE_SIZE=2
GITLAB_PROBES_TOKEN=z_gt678TFvgtg34afa
SLACK_AUTH_TOKEN=xoxb-....PzBG
SLACK_CHANNEL_ID=C...31
PROJECT_NAME=my_project
TIME_UPDATE_CHECK="0 9 * * *"
TIME_UPDATE_INSTALL="0 3 * * *"
```

9. Создаём systemd unit `/etc/systemd/system/gitlab-update.service` с содержимым:

```shell
[Unit]
Description=Gitlab Updater
After=network.target

[Service]
Type=simple
User=root
Group=root

WorkingDirectory=/opt/gitlab-updater
ExecStart=/opt/gitlab-updater/gitlab-updater
SyslogIdentifier=gitlab-updater
Restart=always

[Install]
WantedBy=multi-user.target
```
10. Включаем автозагрузку и запускаем утилиту:
```shell
systemctl enable gitlab-update.service
systemctl start gitlab-update.service
```

11. Посмотреть логи можно так
```shell
journalctl -xe -u gitlab-update
```

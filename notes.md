# exodus

<img src="https://avatars3.githubusercontent.com/u/61350001?s=200&v=4" alt="exodus" width="128"/>

Swiss Army knife to export file(s) to external targets

## Install

```shell script
curl https://raw.githubusercontent.com/kruzio/exodus/master/download.sh | bash
```

# exodus sendfile

```shell script
Send File to one or more destinations

#Send to Slack
echo myfilecontent | bin/exodus sendfile -f -  --target="slack://mychannel?apikey=xoxb-myslackapp-oauth-token&title=My File"

# Send files from a watch directory (one shot)
exodus sendfile  --target=webhook+http://localhost:8080/stuff?content-type=text --watch /tmp/exodus

# Send files from a watch directory (forever)
exodus sendfile  --target=webhook+http://localhost:8080/stuff?content-type=text --watch /tmp/exodus --watch-forever
```

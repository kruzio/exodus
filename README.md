![release](https://img.shields.io/github/v/release/kruzio/exodus?sort=semver)
![Go Version](https://img.shields.io/github/go-mod/go-version/kruzio/exodus)
![Release](https://github.com/kruzio/exodus/workflows/Release/badge.svg)
[![codecov](https://codecov.io/gh/kruzio/exodus/branch/master/graph/badge.svg)](https://codecov.io/gh/kruzio/exodus)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![Tweet](https://img.shields.io/twitter/url?style=social&url=https%3A%2F%2Fgithub.com%2Fkruzio%2Fexodus)

<img src="https://raw.githubusercontent.com/kruzio/artwork/main/logo.png" alt="logo" width="128"/>

# exodus
- Export file(s) to external targets easily
- Send Alert Notification to one or more desitnations

# Supported Targets

```shell script

+--------------+-------------------------------------------------------------------------------------------------------------------------------------------+
|    SCHEME    |                                                                   INFO                                                                    |
+--------------+-------------------------------------------------------------------------------------------------------------------------------------------+
| GS           |   Upload to GCP Cloud Storage              | gs://my-bucket                                                                               |
|              |                                            | For additional information see https://gocloud.dev/howto/blob/#gcs-ctor                      |
|              | -------------------------------------------+--------------------------------------------------------------------------------------------- |
|              |                                                                                                                                           |
| SLACK        |   Post file to a slack channel             | slack://mychannel?apikey=<mykey>[&file-type=json&title=mymsgtitle]                           |
|              |                                            |                                                                                              |
|              |   apikey=<mykey>                           | Slack API token - xoxo-YOURTOKEN                                                             |
|              |                                            | For additional information see https://api.slack.com/apps                                    |
|              |                                            | Note that your app must join the destintation channel                                        |
|              |   file-type=json                           | The content type                                                                             |
|              |   title=mymsgtitle                         | The notification title                                                                       |
|              | -------------------------------------------+--------------------------------------------------------------------------------------------- |
|              |                                                                                                                                           |
| FILE         |   Save file to the local file system using | file:///path/to/dir                                                                          |
|              |                                            | The filename will be the same as name of the inout file name                                 |
|              |                                            | For additional information see https://gocloud.dev/howto/blob/#local                         |
|              | -------------------------------------------+--------------------------------------------------------------------------------------------- |
|              |                                                                                                                                           |
| WEBHOOK+HTTP |   Post file to a webhook                   | For example: webhook+http://myserver?x-headers=X-myheader:myval&token-bearer=1234            |
|              |                                            |                                                                                              |
|              |   >> Authentication Options <<             |                                                                                              |
|              |   token-bearer=<token>                     | Support Authorization Bearer token based authentication                                      |
|              |   username=<username>                      | Basic HTTP Authentication scheme                                                             |
|              |   password=<password>                      | Basic HTTP Authentication scheme                                                             |
|              |                                            |                                                                                              |
|              |   >> Additional Options <<                 |                                                                                              |
|              |   proxy-url=<proxy>                        | The proxy URL the webhook client should connect to                                           |
|              |   content-type=<contentType>               | defaults to json and can be one of: json | text | xml | html | multipart                     |
|              |   x-headers=k1:v1,k2:v2                    | additional custom request headers                                                            |
|              | -------------------------------------------+--------------------------------------------------------------------------------------------- |
|              |                                                                                                                                           |
| WEBHOOK      |   Post file to a webhook                   | For example: webhook://myserver?x-headers=X-myheader:myval&token-bearer=1234                 |
|              |                                            |                                                                                              |
|              |   >> Authentication Options <<             |                                                                                              |
|              |   token-bearer=<token>                     | Support Authorization Bearer token based authentication                                      |
|              |   username=<username>                      | Basic HTTP Authentication scheme                                                             |
|              |   password=<password>                      | Basic HTTP Authentication scheme                                                             |
|              |                                            |                                                                                              |
|              |   >> Additional Options <<                 |                                                                                              |
|              |   proxy-url=<proxy>                        | The proxy URL the webhook client should connect to                                           |
|              |   content-type=<contentType>               | defaults to json and can be one of: json | text | xml | html | multipart                     |
|              |   x-headers=k1:v1,k2:v2                    | additional custom request headers                                                            |
|              |                                            |                                                                                              |
|              |   >> TLS Options <<                        |                                                                                              |
|              |   skip-verify=true                         | If one wished to allow connection to untrusted server                                        |
|              |   ca-file=<path-to-file>                   | CA PEM file                                                                                  |
|              | -------------------------------------------+--------------------------------------------------------------------------------------------- |
|              |                                                                                                                                           |
| SMTP         |   Send file via email (smtp)               | smtp://smtpserver?to=<email>&from=<email>&username=myuser&password=mypass                    |
|              |                                            |                                                                                              |
|              |   to=<target>[,<target>]                   | the destination email address(es) - required                                                 |
|              |   from=<from-email>                        | From email address - required                                                                |
|              |   username=<username>                      | The smtp server authentication information - required                                        |
|              |   password=<password>                      | The smtp server authentication information - required                                        |
|              |                                            |                                                                                              |
|              |   subject=<subject>                        | The Subject line of the email message                                                        |
|              |   skip-verify=true                         | Skip SMTP server TLS verification - (not recommended)                                        |
|              | -------------------------------------------+--------------------------------------------------------------------------------------------- |
|              |                                                                                                                                           |
| S3           |   Upload to AWS S3 bucket                  | s3://bucket-name/subdir?region=us-west-1                                                     |
|              |                                            | For additional information see https://gocloud.dev/howto/blob/#s3                            |
|              | -------------------------------------------+--------------------------------------------------------------------------------------------- |
|              |                                                                                                                                           |
| AZBLOB       |   Upload to Azure Blob storage             | azblob://my-container                                                                        |
|              |                                            | For additional information see https://gocloud.dev/howto/blob/#azure                         |
|              | -------------------------------------------+--------------------------------------------------------------------------------------------- |
|              |                                                                                                                                           |
+--------------+-------------------------------------------------------------------------------------------------------------------------------------------+

```

# Using Exodus as Library

```go
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kruzio/exodus/pkg/sendfile"
)

func main() {
	data, err := ioutil.ReadFile("somefile.json")

	if err != nil {
		os.Exit(255)
	}

	//Let's create our client
	uploadUrl := fmt.Sprintf("webhook://dest.io?skip-verify=true&x-headers=X-myheader:myval&token-bearer=1234")

	uploader, err := sendfile.NewUploader(uploadUrl)
	if err != nil {
		os.Exit(255)
	}

	_ = uploader.SetDestName("somefile")

	err = uploader.Export([]byte(data))
	if err != nil {
		os.Exit(255)
	}
}
```

# Using Exodus as an add-on

```shell script
# Send to Slack
echo myfilecontent | exodus sendfile -f -  --target="slack://mychannel?apikey=xoxb-myslackapp-oauth-token&title=My File"

# Send files from a watch directory (one shot)
exodus sendfile  --target=webhook+http://localhost:8080/stuff?content-type=text --watch /tmp/exodus

# Send files from a watch directory (forever)
exodus sendfile  --target=webhook+http://localhost:8080/stuff?content-type=text --watch /tmp/exodus --watch-forever
```

# Using Exodus in Kubernetes

```yaml
kind: Namespace
metadata:
  name: exodus
---
apiVersion: v1
kind: Secret
metadata:
  name: export-targets
  namespace: exodus
type: Opaque
data:
  # printf "webhook://dest.io?skip-verify=true&token-bearer=1234&content-type=text" | base64
  targets: d2ViaG9vazovL2Rlc3QuaW8/c2tpcC12ZXJpZnk9dHJ1ZSZ0b2tlbi1iZWFyZXI9MTIzNCZjb250ZW50LXR5cGU9dGV4dA==
---
apiVersion: batch/v1
kind: Job
metadata:
  name: example
  namespace: exodus
  labels:
    app.kubernetes.io/name: exodus-example
    app.kubernetes.io/instance: example
    app.kubernetes.io/version: "1.0.0"
spec:
  backoffLimit: 1
  template:
    spec:
      # Pod Security
      automountServiceAccountToken: false
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000590000
        runAsGroup: 1000590000
        fsGroup: 1000590000

      volumes:
        # Our Send Box
        - name: sendbox
          emptyDir: {}
      containers:
        - name: exodus
          image: kruzio/exodus:v0.2.0
          imagePullPolicy: IfNotPresent
          volumeMounts:
            - mountPath: /sendbox
              name: sendbox
          args:
            - "sendfile"
            - "--watch"
            - "/sendbox"
            - "--watch-forever"
            - "false"
            # Debugging
            #- "-v"
            #- "7"

          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities:
              drop:
                - ALL
          env:
            - name: KRUZIO_EXODUS_SENDFILE_TARGETS
              valueFrom:
                secretKeyRef:
                  name: export-targets
                  key: targets
        - name: producer
          image: busybox:latest
          imagePullPolicy: IfNotPresent
          volumeMounts:
            - mountPath: /sendbox
              name: sendbox
          command: ["/bin/sh"]
          args:
              - -c
              - "sleep 3 && echo hello > /sendbox/file-to-send.txt && sleep 3 && ls -la /sendbox/ && exit 0"
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities:
              drop:
                - ALL
      restartPolicy: Never
```

## Contributing

### Bugs

If you think you have found a bug please follow the instructions below.

- Please spend a small amount of time giving due diligence to the issue tracker. Your issue might be a duplicate.
- Open a [new issue](https://github.com/kruzio/exodus/issues/new/choose) if a duplicate doesn't already exist.

### Features

If you have an idea to enhance exodus follow the steps below.

- Open a [new issue](https://github.com/kruzio/exodus/issues/new/choose).
- Remember users might be searching for your issue in the future, so please give it a meaningful title to helps others.
- Clearly define the use case, using concrete examples.
- Feel free to include any technical design for your feature.

### Pull Requests

- Your PR is more likely to be accepted if it focuses on just one change.
- Please include a comment with the results before and after your change. 
- Your PR is more likely to be accepted if it includes tests. 
- You're welcome to submit a draft PR if you would like early feedback on an idea or an approach.


[![Stargazers over time](https://starchart.cc/kruzio/exodus.svg)](https://starchart.cc/kruzio/exodus)

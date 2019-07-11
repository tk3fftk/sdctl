## sd-cli

```
NAME:
   sdctl - Screwdriver.cd API wrapper

USAGE:
   validate yamls, start build locally

VERSION:
   0.1.0

COMMANDS:
     get                    get sdctl settings and Screwdriver.cd information
     set                    set sdctl settings
     context                handle screwdriver contexts
     clear                  clear your setting and set to default
     build, b               start a job. parameters: <pipelieid> <start_from> 
     validate, v            validate your screwdriver.yaml, default to screwdriver.yaml
     validate-template, vt  validate your sd-template.yaml, default to sd-template.yaml
     help, h                Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

COPYRIGHT:
   tk3fftk

```

### Setup
In case of using your screwdriver.cd cluster
- Install sdctl
```
$ go get github.com/tk3fftk/sdctl
```
- Get screwdriver user token from https://<your_screwdrivercd>/user-settings
- Set configurations
```
$ sdctl set token <obtained-token>
$ sdctl set api https://<your_screwdrivercd>
```

### Usage
- start build
```
$ sdctl build <pipelineid> <start_from>
```

- validate screwdriver.yaml
```
$ sdctl validate
or
$ sdctl v
```

- validate sd-template.yaml
```
$ sdctl validate-tempalte
or
$ sdctl vt
```

- get build pages from build id
```
$ sdctl set jwt
$ sdctl get build-pages "156442 156518 323281"
```

- switch another screwdriver.cd
```
$ sdctl context set next
$ sdctl set token <obtained-token>
$ sdctl set api https://<your_screwdrivercd>
```

## sd-cli

```
validate yamls, start build locally

Usage:
  sdctl [flags]
  sdctl [command]

Available Commands:
  build             start a job.
  clear             clear your setting and set to default
  context           handle screwdriver contexts
  get               get sdctl settings and Screwdriver.cd information
  help              Help about any command
  set               set sdctl settings
  validate          validate your screwdriver.yaml, default to screwdriver.yaml
  validate-template validate your sd-template.yaml, default to sd-template.yaml

Flags:
  -h, --help      help for sdctl
      --version   version for sdctl

Use "sdctl [command] --help" for more information about a command.
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

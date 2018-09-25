package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/tk3fftk/sdctl/config"
	"github.com/tk3fftk/sdctl/sdapi"
	"gopkg.in/urfave/cli.v1"
)

var (
	readFile  = ioutil.ReadFile
	writeFile = ioutil.WriteFile
)

// SdctlContexts represents the context of Screwdriver.cds
type SdctlContexts struct {
	CurrentContext string                        `json:"current"`
	Contexts       map[string]config.SdctlConfig `json:"contexts"`
}

func loadYaml(yamlPath string) (yaml string, err error) {
	yamlFile, err := readFile(yamlPath)
	if err != nil {
		return
	}
	yaml = fmt.Sprintf("%q", string(yamlFile[:]))

	return
}

func initDotFile(path string, force bool) (err error) {
	_, err = readFile(path)
	if err != nil || force {
		conf := config.SdctlConfig{
			UserToken: "",
			APIURL:    "",
			SDJWT:     "",
		}
		context := SdctlContexts{
			CurrentContext: "default",
			Contexts:       make(map[string]config.SdctlConfig),
		}
		context.Contexts["default"] = conf

		f, _ := json.Marshal(context)
		err = writeFile(path, f, 0660)
		if err != nil {
			return
		}
	}
	return
}

func failureExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	}
	os.Exit(1)
}

func getConfig(path, context string) (conf config.SdctlConfig, err error) {
	confFile, err := readFile(path)
	if err != nil {
		return
	}

	sdcontexts, err := getSdctlContexts(path)
	if err != nil {
		return
	}

	if err = json.Unmarshal(confFile, &sdcontexts); err != nil {
		return
	}
	if context == "" {
		conf = sdcontexts.Contexts[sdcontexts.CurrentContext]
	} else {
		conf = sdcontexts.Contexts[context]
	}

	return
}

func updateConfig(path, userToken, apiURL, sdJWT, context string) (err error) {
	conf, err := getConfig(path, context)
	if err != nil {
		return
	}

	if userToken != "" {
		conf.UserToken = userToken
	}
	if apiURL != "" {
		conf.APIURL = apiURL
	}
	if sdJWT != "" {
		conf.SDJWT = sdJWT
	}

	err = updateSdctlContexts(path, conf, context)
	return
}

func getSdctlContexts(path string) (sdcontexts SdctlContexts, err error) {
	confFile, err := readFile(path)
	if err != nil {
		return
	}

	if err = json.Unmarshal(confFile, &sdcontexts); err != nil {
		return
	}

	return
}

func updateCurrentContext(path, context string) (err error) {
	sdcontexts, err := getSdctlContexts(path)
	if err != nil {
		return
	}

	current := sdcontexts.CurrentContext
	if current == context {
		return
	}

	for cont := range sdcontexts.Contexts {
		if cont == context {
			sdcontexts.CurrentContext = context
			f, _ := json.Marshal(sdcontexts)
			err = writeFile(path, f, 0660)
			break
		}
	}

	if current == sdcontexts.CurrentContext {
		err = fmt.Errorf("%s does not exist", context)
	}

	return
}

func updateSdctlContexts(path string, newConf config.SdctlConfig, context string) (err error) {
	sdcontexts, err := getSdctlContexts(path)

	var f []byte
	if context == "" {
		sdcontexts.Contexts[sdcontexts.CurrentContext] = newConf
	} else {
		sdcontexts.Contexts[context] = newConf
	}

	f, err = json.Marshal(sdcontexts)
	if err != nil {
		return
	}

	err = writeFile(path, f, 0660)
	return
}

func main() {

	usr, err := user.Current()
	if err != nil {
		failureExit(err)
	}
	var configPath = usr.HomeDir + "/.sdctl"
	err = initDotFile(configPath, false)
	if err != nil {
		failureExit(err)
	}

	app := cli.NewApp()
	app.Name = "sdctl"
	app.Usage = "Screwdriver.cd API wrapper"
	app.UsageText = "validate yamls, start build locally"
	app.Copyright = "tk3fftk"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:  "get",
			Usage: "get sdctl settings and Screwdriver.cd information",
			Subcommands: []cli.Command{
				{
					Name:  "token",
					Usage: "get your user token",
					Action: func(c *cli.Context) error {
						conf, err := getConfig(configPath, "")
						if err != nil {
							failureExit(err)
						}
						fmt.Println(conf.UserToken)
						return nil
					},
				},
				{
					Name:  "api",
					Usage: "get configured api url",
					Action: func(c *cli.Context) error {
						conf, err := getConfig(configPath, "")
						if err != nil {
							failureExit(err)
						}
						fmt.Println(conf.APIURL)
						return nil
					},
				},
				{
					Name:  "jwt",
					Usage: "show your jwt",
					Action: func(c *cli.Context) error {
						conf, err := getConfig(configPath, "")
						if err != nil {
							failureExit(err)
						}
						fmt.Println("Bearer: " + conf.SDJWT)
						return nil
					},
				},
				{
					Name:    "build-pages",
					Usage:   "get build page url",
					Aliases: []string{"bp"},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 {
							return cli.ShowAppHelp(c)
						}
						buildID := c.Args().Get(0)

						api := sdapi.New()
						conf, err := getConfig(configPath, "")
						if err != nil {
							failureExit(err)
						}
						if err := api.GetPipelinePageFromBuildID(conf, buildID); err != nil {
							failureExit(err)
						}
						return nil
					},
				},
			},
		},
		{
			Name:  "set",
			Usage: "set sdctl settings",
			Subcommands: []cli.Command{
				{
					Name:  "token",
					Usage: "set your user token",
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 {
							return cli.ShowAppHelp(c)
						}
						if err = updateConfig(configPath, c.Args().Get(0), "", "", ""); err != nil {
							failureExit(err)
						}
						return nil
					},
				},
				{
					Name:  "api",
					Usage: "set your Screwdriver.cd api url",
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 {
							return cli.ShowAppHelp(c)
						}
						if err = updateConfig(configPath, "", c.Args().Get(0), "", ""); err != nil {
							failureExit(err)
						}
						return nil
					},
				},
				{
					Name:  "jwt",
					Usage: "get and store jwt locally",
					Action: func(c *cli.Context) error {
						api := sdapi.New()
						conf, err := getConfig(configPath, "")
						if err != nil {
							failureExit(err)
						} else if conf.UserToken == "" {
							failureExit(errors.New("you must set user token before getting JWT"))
						}
						token, err := api.GetJwt(conf)
						if err != nil {
							failureExit(err)
						}
						if err = updateConfig(configPath, "", "", token, ""); err != nil {
							failureExit(err)
						}
						println("Bearer " + token)
						return nil
					},
				},
			},
		},
		{
			Name:  "context",
			Usage: "handle screwdriver contexts",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Usage:   "show context list",
					Aliases: []string{"ls"},
					Action: func(c *cli.Context) error {
						contexts, err := getSdctlContexts(configPath)
						if err != nil {
							failureExit(err)
						}
						for key := range contexts.Contexts {
							fmt.Println(key)
						}
						return nil
					},
				},
				{
					Name:  "current",
					Usage: "show current context",
					Action: func(c *cli.Context) error {
						contexts, err := getSdctlContexts(configPath)
						if err != nil {
							failureExit(err)
						}
						fmt.Println(contexts.CurrentContext)
						return nil
					},
				},
				{
					Name:  "set",
					Usage: "set current to context. if it doesn't exist, create new one",
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 {
							return cli.ShowAppHelp(c)
						}
						context := c.Args().Get(0)
						err = updateCurrentContext(configPath, context)
						if err != nil {
							if err = updateConfig(configPath, "", "", "", context); err != nil {
								failureExit(err)
							}
							if err = updateCurrentContext(configPath, context); err != nil {
								failureExit(err)
							}
						}
						return nil
					},
				},
			},
		}, {
			Name:  "clear",
			Usage: "clear your setting and set to default",
			Action: func(c *cli.Context) error {
				if err := initDotFile(configPath, true); err != nil {
					failureExit(err)
				}
				println("Cleared your settings")
				return nil
			},
		}, {
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "start a job. parameters: <pipelieid> <start_from> ",
			Action: func(c *cli.Context) error {
				if len(c.Args()) != 2 {
					return cli.ShowAppHelp(c)
				}
				api := sdapi.New()
				conf, err := getConfig(configPath, "")
				if err != nil {
					failureExit(err)
				}
				if err := api.PostEvent(conf, c.Args().Get(0), c.Args().Get(1), false); err != nil {
					failureExit(err)
				}
				return nil
			},
		}, {
			Name:    "validate",
			Aliases: []string{"v"},
			Usage:   "validate your screwdriver.yaml, default to screwdriver.yaml",
			Action: func(c *cli.Context) error {
				var f string
				if len(c.Args()) == 0 {
					f = "screwdriver.yaml"
				} else {
					f = c.Args().Get(0)
				}
				yaml, err := loadYaml(f)
				if err != nil {
					failureExit(err)
				}
				api := sdapi.New()
				conf, err := getConfig(configPath, "")
				if err != nil {
					failureExit(err)
				}
				if err := api.Validator(conf, yaml, false); err != nil {
					failureExit(err)
				}

				return nil
			},
		}, {
			Name:    "validate-template",
			Aliases: []string{"vt"},
			Usage:   "validate your sd-template.yaml, default to sd-template.yaml",
			Action: func(c *cli.Context) error {
				var f string
				if len(c.Args()) == 0 {
					f = "sd-template.yaml"
				} else {
					f = c.Args().Get(0)
				}
				yaml, err := loadYaml(f)
				if err != nil {
					failureExit(err)
				}
				api := sdapi.New()
				conf, err := getConfig(configPath, "")
				if err != nil {
					failureExit(err)
				}
				if err := api.ValidatorTemplate(conf, yaml, false); err != nil {
					failureExit(err)
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}

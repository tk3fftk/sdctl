package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/tk3fftk/sdctl/pkg/sdapi"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
	"gopkg.in/urfave/cli.v1"
)

var configFileName = ".sdctl"

func failureExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
	}
	os.Exit(1)
}

func readYaml(yamlPath string) (yaml string, err error) {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return
	}
	yaml = fmt.Sprintf("%q", string(yamlFile[:]))

	return
}

func main() {

	usr, err := user.Current()
	if err != nil {
		failureExit(err)
	}

	configPath := usr.HomeDir + "/" + configFileName
	config, err := sdctl_context.LoadConfig(configPath, false)
	if err != nil {
		failureExit(err)
	}
	sdctx := config.SdctlContexts[config.CurrentContext]
	api, err := sdapi.New(sdctx, nil)
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
						config.PrintParam(sdctl_context.UserTokenKey, nil)
						return nil
					},
				},
				{
					Name:  "api",
					Usage: "get configured api url",
					Action: func(c *cli.Context) error {
						config.PrintParam(sdctl_context.APIURLKey, nil)
						return nil
					},
				},
				{
					Name:  "jwt",
					Usage: "show your jwt",
					Action: func(c *cli.Context) error {
						config.PrintParam(sdctl_context.SDJWTKey, nil)
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

						if err := api.GetPipelinePageFromBuildID(buildID); err != nil {
							failureExit(err)
						}
						return nil
					},
				},
				{
					Name:  "banner",
					Usage: "get sd banners",
					Action: func(c *cli.Context) error {
						banners, err := api.GetBanners()
						if err != nil {
							failureExit(err)
						}
						fmt.Println("ID\tIsActive\tMessage")
						for _, b := range banners {
							fmt.Printf("%v\t%v\t%v\n", b.ID, b.IsActive, b.Message)
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

						config.SetParam(sdctl_context.UserTokenKey, c.Args().Get(0), nil)
						config.Update(configPath)
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
						config.SetParam(sdctl_context.APIURLKey, c.Args().Get(0), nil)
						config.Update(configPath)
						return nil
					},
				},
				{
					Name:  "jwt",
					Usage: "get and store jwt locally",
					Action: func(c *cli.Context) error {
						// TODO handle it in sdapi.go
						if sdctx.UserToken == "" {
							failureExit(errors.New("you must set user token before getting JWT"))
						}
						token, err := api.GetJWT()
						if err != nil {
							failureExit(err)
						}

						config.SetParam(sdctl_context.SDJWTKey, token, nil)
						config.Update(configPath)
						println("Bearer " + token)
						return nil
					},
				},
				{
					Name:  "banner",
					Usage: "update sd banner.",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "id, i",
							Usage: "banner `ID` when update",
						},
						cli.StringFlag{
							Name:  "msg, m",
							Usage: "banner `MESSAGE` body (required)",
						},
						cli.StringFlag{
							Name:  "type, t",
							Usage: "banner `TYPE`",
							Value: "info",
						},
						cli.StringFlag{
							Name:  "isActive, a",
							Usage: "banner status `FLAG`",
							Value: "true",
						},
						cli.BoolFlag{
							Name:  "delete, d",
							Usage: "delete banner with id",
						},
					},
					Action: func(c *cli.Context) error {
						id := c.String("id")
						msg := c.String("msg")
						bannerType := c.String("type")
						isActive := c.String("isActive")
						delete := c.Bool("delete")

						if msg == "" && id == "" {
							return cli.ShowAppHelp(c)
						}

						_, err := api.UpdateBanner(id, msg, bannerType, isActive, delete, false)
						if err != nil {
							failureExit(err)
						}
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
						config.PrintParam(sdctl_context.ContextsKey, nil)
						return nil
					},
				},
				{
					Name:  "current",
					Usage: "show current context",
					Action: func(c *cli.Context) error {
						config.PrintParam(sdctl_context.CurrentContextKey, nil)
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
						config.SetParam(sdctl_context.CurrentContextKey, context, nil)
						config.Update(configPath)
						return nil
					},
				},
			},
		}, {
			Name:  "clear",
			Usage: "clear your setting and set to default",
			Action: func(c *cli.Context) error {
				sdctl_context.LoadConfig(configPath, true)
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
				if err := api.PostEvent(c.Args().Get(0), c.Args().Get(1), false); err != nil {
					failureExit(err)
				}
				return nil
			},
		}, {
			Name:    "validate",
			Aliases: []string{"v"},
			Usage:   "validate your screwdriver.yaml, default to screwdriver.yaml",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Value: "screwdriver.yaml",
					Usage: "specify pipeline file",
				},
				cli.BoolFlag{
					Name:  "output",
					Usage: "print velidator result",
				},
			},
			Action: func(c *cli.Context) error {
				yaml, err := readYaml(c.String("file"))
				if err != nil {
					failureExit(err)
				}
				if err := api.Validator(yaml, false, c.Bool("output")); err != nil {
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
				yaml, err := readYaml(f)
				if err != nil {
					failureExit(err)
				}
				if err := api.ValidatorTemplate(yaml, false); err != nil {
					failureExit(err)
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}

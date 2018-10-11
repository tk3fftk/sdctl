package sdapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
)

var transport = &http.Transport{
	TLSClientConfig: &tls.Config{},
}
var client = &http.Client{Transport: transport}

// SDAPI wraps Screwdriver.cd API
type SDAPI interface {
	GetJwt(conf sdctl_context.SdctlContext) (string, error)
	PostEvent(conf sdctl_context.SdctlContext, pipelineID string, startFrom string, retry bool) error
	Validator(conf sdctl_context.SdctlContext, yaml string, retry bool) error
	ValidatorTemplate(conf sdctl_context.SdctlContext, yaml string, retry bool) error
	GetPipelinePageFromBuildID(conf sdctl_context.SdctlContext, buildID string) error
}

type validatorResponse struct {
	Annotations   interface{} `json:"annotations"`
	Errors        []string    `json:"errors"`
	Jobs          interface{} `json:"jobs"`
	Workflow      []string    `json:"workflow"`
	WorkflowGraph interface{} `json:"workflowGraph"`
}

type templateValidatorResponse struct {
	Template interface{}             `json:"template"`
	Errors   []templateValidateError `json:"errors"`
}

type templateValidateError struct {
	Message string      `json:"message"`
	Path    string      `json:"path"`
	Type    string      `json:"type"`
	Context interface{} `json:"context"`
}

type tokenResponse struct {
	JWT string `json:"token"`
}

type buildResponse struct {
	EventID int `json:"eventId"`
}

type eventResponse struct {
	PipelineID int `json:"pipelineId"`
}

type sdapi struct{}

// New creates a SDAPI instance
func New() SDAPI {
	s := sdapi{}
	return SDAPI(s)
}

func (sd sdapi) GetJwt(sdctlContext sdctl_context.SdctlContext) (string, error) {
	jwt := new(tokenResponse)
	u := sdctlContext.APIURL + "/v4/auth/token?api_token=" + sdctlContext.UserToken
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if res.StatusCode != 200 {
		return "", errors.New(strconv.Itoa(res.StatusCode))
	} else if err != nil {
		return "", err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(jwt)

	return jwt.JWT, err
}

func (sd sdapi) PostEvent(sdctlContext sdctl_context.SdctlContext, pipelineID string, startFrom string, retry bool) error {
	u := sdctlContext.APIURL + "/v4/events"
	b := map[string]string{
		"pipelineId": pipelineID,
		"startFrom":  startFrom,
	}
	jsonBody, _ := json.Marshal(b)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+sdctlContext.SDJWT)
	res, err := client.Do(req)
	if res.StatusCode != 201 {
		if retry {
			return errors.New(strconv.Itoa(res.StatusCode))
		}
		sdctlContext.SDJWT, err = sd.GetJwt(sdctlContext)
		if err != nil {
			return err
		}
		return sd.PostEvent(sdctlContext, pipelineID, startFrom, true)
	} else if err != nil {
		return err
	}

	return nil
}

func (sd sdapi) Validator(sdctlContext sdctl_context.SdctlContext, yaml string, retry bool) error {
	vr := new(validatorResponse)
	u := sdctlContext.APIURL + "/v4/validator"
	b := `{"yaml":` + yaml + `}`

	req, err := http.NewRequest("POST", u, bytes.NewBuffer([]byte(b)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+sdctlContext.SDJWT)
	res, err := client.Do(req)
	if res.StatusCode != 200 {
		if retry {
			return errors.New(strconv.Itoa(res.StatusCode))
		}
		sdctlContext.SDJWT, err = sd.GetJwt(sdctlContext)
		if err != nil {
			return err
		}
		return sd.Validator(sdctlContext, yaml, true)
	} else if err != nil {
		return err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(vr)
	if err != nil {
		return err
	}
	if len(vr.Errors) != 0 {
		return errors.New(vr.Errors[0])
	}

	println("🙆")

	return nil
}

func (sd sdapi) ValidatorTemplate(sdctlContext sdctl_context.SdctlContext, yaml string, retry bool) error {
	tvr := new(templateValidatorResponse)
	u := sdctlContext.APIURL + "/v4/validator/template"
	b := `{"yaml":` + yaml + `}`

	req, err := http.NewRequest("POST", u, bytes.NewBuffer([]byte(b)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+sdctlContext.SDJWT)
	res, err := client.Do(req)
	if res.StatusCode != 200 {
		if retry {
			return errors.New(strconv.Itoa(res.StatusCode))
		}
		sdctlContext.SDJWT, err = sd.GetJwt(sdctlContext)
		if err != nil {
			return err
		}
		return sd.ValidatorTemplate(sdctlContext, yaml, true)
	} else if err != nil {
		return err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(tvr)
	if err != nil {
		return err
	}
	if len(tvr.Errors) != 0 {
		i := 0
		for i < len(tvr.Errors) {
			fmt.Printf("%v\n", tvr.Errors[i].Message)
			i++
		}
		return errors.New("invalid template of Screwdriver.cd")
	}

	println("🙆")

	return nil
}

func (sd sdapi) GetPipelinePageFromBuildID(sdctlContext sdctl_context.SdctlContext, buildID string) error {
	buildIDList := strings.Split(strings.Replace(strings.TrimSpace(buildID), "\n", " ", -1), " ")
	buildIDLength := len(buildIDList)

	var wg sync.WaitGroup
	wg.Add(buildIDLength)

	exit := make(chan error, buildIDLength)

	for _, b := range buildIDList {
		go func(b string) {
			defer wg.Done()

			br, err := getBuilds(sdctlContext, b)
			if err != nil {
				exit <- err
				return
			}
			er, err := getEvents(sdctlContext, br.EventID)
			if err != nil {
				exit <- err
				return
			}
			println(strings.Replace(sdctlContext.APIURL, "api-cd", "cd", 1) + "/pipelines/" + strconv.Itoa(er.PipelineID) + "/builds/" + b)
		}(b)
	}

	wg.Wait()

	select {
	case err := <-exit:
		return err
	default:
		return nil
	}
}

func getBuilds(sdctlContext sdctl_context.SdctlContext, buildID string) (*buildResponse, error) {
	br := new(buildResponse)
	getBuildAPI := sdctlContext.APIURL + "/v4/builds/" + buildID + "?token=" + sdctlContext.SDJWT

	req, err := http.NewRequest("GET", getBuildAPI, nil)
	if err != nil {
		return br, err
	}

	res, err := client.Do(req)
	if res.StatusCode != 200 {
		return br, errors.New(strconv.Itoa(res.StatusCode))
	} else if err != nil {
		return br, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(br)

	return br, err
}

func getEvents(sdctlContext sdctl_context.SdctlContext, eventID int) (*eventResponse, error) {
	er := new(eventResponse)
	getBuildAPI := sdctlContext.APIURL + "/v4/events/" + strconv.Itoa(eventID) + "?token=" + sdctlContext.SDJWT

	req, err := http.NewRequest("GET", getBuildAPI, nil)
	if err != nil {
		return er, err
	}

	res, err := client.Do(req)
	if res.StatusCode != 200 {
		return er, errors.New(strconv.Itoa(res.StatusCode))
	} else if err != nil {
		return er, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(er)

	return er, err
}

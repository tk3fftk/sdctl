package sdapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
)

var (
	mockUserToken = "token"
	mockAPIURL    = "localhost"
	mockSDJWT     = "invalid_jwt"
	mockSDContext = sdctl_context.SdctlContext{
		UserToken: mockUserToken,
		APIURL:    mockAPIURL,
		SDJWT:     mockSDJWT,
	}
	mockPipelineID             = "1234"
	mockStartFrom              = "~commit"
	mockYaml                   = "jobs:\r\n  main:\r\n    image: node:10\r\n    steps:\r\n      - echo: echo hoge"
	mockSDJWTResponse          = "testdata/jwt.json"
	mockSDUnauthorizedResponse = "testdata/unauthorized.json"
	mockSDBadRequestResponse   = "testdata/bad_request.json"
)

func TestNew(t *testing.T) {
	invalidAPIURL := "%&"
	invalidSDContext := sdctl_context.SdctlContext{
		APIURL: invalidAPIURL,
	}

	cases := map[string]struct {
		sdctx          sdctl_context.SdctlContext
		expectedResult bool
	}{
		"Creates SDAPI successfully": {
			mockSDContext,
			true,
		},
		"Failure creates SDAPI": {
			invalidSDContext,
			false,
		},
	}

	for k, v := range cases {
		k := k
		v := v

		t.Run(k, func(t *testing.T) {
			_, err := New(v.sdctx, nil)
			switch v.expectedResult {
			case true:
				if err != nil {
					t.Errorf("error should be nil but: '%v'", err)
				}
			case false:
				if err == nil {
					t.Errorf("error should not be nil but nil")
				} else {
					fmt.Printf("%v\n", err)
				}
			}
		})
	}
}

func TestGetJWT(t *testing.T) {
	cases := map[string]struct {
		expectedResult   bool
		expectedResponse string
	}{
		"Get JWT successfully": {
			true,
			mockSDJWTResponse,
		},
		"Failure getting JWT with unauthorized status": {
			false, mockSDUnauthorizedResponse,
		},
	}

	for k, v := range cases {
		k := k
		v := v

		t.Run(k, func(t *testing.T) {

			muxAPI := http.NewServeMux()
			testAPIServer := httptest.NewServer(muxAPI)
			defer testAPIServer.Close()

			path := "/v4/auth/token"
			mockSDContext.APIURL = testAPIServer.URL
			muxAPI.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				if !v.expectedResult {
					w.WriteHeader(http.StatusUnauthorized)
				}
				http.ServeFile(w, r, v.expectedResponse)
			})

			sdapi, err := New(mockSDContext, nil)
			if err != nil {
				t.Fatal("should not cause error")
			}

			jwt, err := sdapi.GetJWT()
			switch v.expectedResult {
			case true:
				if err != nil {
					t.Errorf("error should be nil but: '%v'", err)
				}
				if jwt != "thisissdjwttoken" {
					t.Errorf("'%v' is not expected", jwt)
				}
			case false:
				if err == nil {
					t.Errorf("error should not be nil but nil")
				} else {
					fmt.Printf("%v\n", err)
				}
			}
		})
	}
}

func TestGetBanners(t *testing.T) {
	cases := map[string]struct {
		expectedResult   bool
		expectedResponse string
	}{
		"Get banners successfully": {
			true,
			"testdata/banner.json",
		},
		"Get banners with failure": {
			false,
			mockSDBadRequestResponse,
		},
	}

	for k, v := range cases {
		k := k
		v := v

		t.Run(k, func(t *testing.T) {

			muxAPI := http.NewServeMux()
			testAPIServer := httptest.NewServer(muxAPI)
			defer testAPIServer.Close()

			path := "/v4/banners"
			mockSDContext.APIURL = testAPIServer.URL
			muxAPI.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				if !v.expectedResult {
					w.WriteHeader(http.StatusUnauthorized)
				}
				http.ServeFile(w, r, v.expectedResponse)
			})

			sdapi, err := New(mockSDContext, nil)
			if err != nil {
				t.Fatal("should not cause error")
			}

			banners, err := sdapi.GetBanners()
			switch v.expectedResult {
			case true:
				if err != nil {
					t.Errorf("error should be nil but: '%v'", err)
				}
				b := banners[0]
				expected := "Due to planned upgrade of Kubernetes, Screwdriver will be down"
				if b.ID != 0 || b.Message != expected {
					t.Errorf("response should be equal with dummy date but: '%v' and '%v'", b.ID, b.Message)
				}
			case false:
				if err == nil {
					t.Errorf("error should not be nil but nil")
				} else {
					fmt.Printf("%v\n", err)
				}
			}
		})
	}
}

func TestUpdateBanner(t *testing.T) {
	dummyId := "13"
	dummyMessage := "Due to planned upgrade of Kubernetes, Screwdriver will be down"
	dummyIsActive := "false"
	dummyType := "info"

	cases := map[string]struct {
		id               string
		expectedResult   bool
		expectedResponse string
		delete           bool
	}{
		"Create a banner successfully": {
			"",
			true,
			"testdata/banner_creation.json",
			false,
		},
		"Failed to create a banner": {
			"",
			false,
			"testdata/bad_banner_creation.json",
			false,
		},
		"Update a banner successfully": {
			dummyId,
			true,
			"testdata/banner_patch.json",
			false,
		},
		"Failed to update a banner": {
			dummyId,
			false,
			"testdata/banner_not_found.json",
			false,
		},
		"Delete a banner successfully": {
			dummyId,
			true,
			"",
			true,
		},
		"Failed to delete a banner": {
			dummyId,
			false,
			"testdata/banner_not_found.json",
			true,
		},
	}

	for k, v := range cases {
		k := k
		v := v

		t.Run(k, func(t *testing.T) {
			muxAPI := http.NewServeMux()
			testAPIServer := httptest.NewServer(muxAPI)
			defer testAPIServer.Close()

			path := "/v4/banners/" + v.id
			mockSDContext.APIURL = testAPIServer.URL
			muxAPI.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				if !v.expectedResult {
					if v.id != "" {
						w.WriteHeader(http.StatusNotFound)
					} else {
						w.WriteHeader(http.StatusBadRequest)
					}
				}
				if v.delete {
					w.WriteHeader(http.StatusNoContent)
				}
				http.ServeFile(w, r, v.expectedResponse)
			})

			sdapi, err := New(mockSDContext, nil)
			if err != nil {
				t.Fatal("should not cause error")
			}

			banner, err := sdapi.UpdateBanner(dummyId, dummyMessage, dummyType, dummyIsActive, v.delete, false)
			switch v.expectedResult {
			case true:
				if err != nil {
					t.Errorf("error should be nil but: '%v'", err)
				}

				expctedResponseJsonFile, err := ioutil.ReadFile(v.expectedResponse)
				if err != nil && !v.delete {
					t.Fatal("should not cause error")
				}
				expectedResponseJson := new(bannerResponse)
				err = json.Unmarshal(expctedResponseJsonFile, expectedResponseJson)

				if diff := cmp.Diff(expectedResponseJson, &banner); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			case false:
				if err == nil {
					t.Errorf("error should not be nil but nil")
				} else {
					fmt.Printf("%v\n", err)
				}
			}
		})
	}
}

func TestPostEvent(t *testing.T) {
	cases := map[string]struct {
		expectedResult        bool
		expectedResponse      string
		expectedRetry         bool
		expectedRetryResponse string
	}{
		"POST events successfully": {
			true,
			"testdata/post_event.json",
			false,
			"",
		},
		"Retry successfully after authorization": {
			true,
			mockSDUnauthorizedResponse,
			true,
			"testdata/post_event.json",
		},
		"Failure with bad request": {
			false,
			mockSDBadRequestResponse,
			false,
			"",
		},
		"Failure with bad request after retrying": {
			false,
			mockSDUnauthorizedResponse,
			true,
			mockSDBadRequestResponse,
		},
	}

	for k, v := range cases {
		k := k
		v := v

		t.Run(k, func(t *testing.T) {

			muxAPI := http.NewServeMux()
			testAPIServer := httptest.NewServer(muxAPI)
			defer testAPIServer.Close()

			path := "/v4/events"
			mockSDContext.APIURL = testAPIServer.URL
			muxAPI.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				if v.expectedRetry {
					jwt := r.Header.Get("Authorization")
					if jwt == "Bearer thisissdjwttoken" {
						if v.expectedResult {
							w.WriteHeader(http.StatusCreated)
						} else {
							w.WriteHeader(http.StatusInternalServerError)
						}
						http.ServeFile(w, r, v.expectedRetryResponse)
						return
					}
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					if v.expectedResult {
						w.WriteHeader(http.StatusCreated)
					} else {
						w.WriteHeader(http.StatusInternalServerError)
					}
				}

				http.ServeFile(w, r, v.expectedResponse)
			})
			tokenPath := "/v4/auth/token"
			muxAPI.HandleFunc(tokenPath, func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, mockSDJWTResponse)
			})

			sdapi, err := New(mockSDContext, nil)
			if err != nil {
				t.Fatal("should not cause error")
			}

			err = sdapi.PostEvent(mockPipelineID, mockStartFrom, false)
			switch v.expectedResult {
			case true:
				if err != nil {
					t.Errorf("error should be nil but: '%v'", err)
				}
			case false:
				if err == nil {
					t.Errorf("error should not be nil but nil")
				} else {
					fmt.Printf("%v\n", err)
				}
			}
		})
	}
}

func TestValidator(t *testing.T) {
	cases := map[string]struct {
		expectedHttpResult     bool
		output                 bool
		expectedResponse       string
		expectedRetry          bool
		expectedRetryResponse  string
		expectedValidateResult bool
	}{
		"POST validate successfully": {
			true,
			false,
			"testdata/validate.json",
			false,
			"",
			true,
		},
		"POST validate successfully wich output": {
			true,
			true,
			"testdata/validate.json",
			false,
			"",
			true,
		},
		"Retry successfully after authorization": {
			true,
			false,
			mockSDUnauthorizedResponse,
			true,
			"testdata/validate.json",
			true,
		},
		"Failure with bad request": {
			false,
			false,
			mockSDBadRequestResponse,
			false,
			"",
			false,
		},
		"Failure with invalid yaml": {
			true,
			false,
			"testdata/config_parse_error.json",
			false,
			"",
			false,
		},
		"Failure with bad request after retrying": {
			false,
			false,
			mockSDUnauthorizedResponse,
			true,
			mockSDBadRequestResponse,
			false,
		},
	}

	for k, v := range cases {
		k := k
		v := v

		t.Run(k, func(t *testing.T) {

			muxAPI := http.NewServeMux()
			testAPIServer := httptest.NewServer(muxAPI)
			defer testAPIServer.Close()

			path := "/v4/validator"
			mockSDContext.APIURL = testAPIServer.URL
			muxAPI.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				if v.expectedRetry {
					jwt := r.Header.Get("Authorization")
					if jwt == "Bearer thisissdjwttoken" {
						if !v.expectedHttpResult {
							w.WriteHeader(http.StatusInternalServerError)
						}
						http.ServeFile(w, r, v.expectedRetryResponse)
						return
					}
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					if !v.expectedHttpResult {
						w.WriteHeader(http.StatusInternalServerError)
					}
				}

				http.ServeFile(w, r, v.expectedResponse)
			})
			tokenPath := "/v4/auth/token"
			muxAPI.HandleFunc(tokenPath, func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, mockSDJWTResponse)
			})

			sdapi, err := New(mockSDContext, nil)
			if err != nil {
				t.Fatal("should not cause error")
			}

			err = sdapi.Validator(mockYaml, false, v.output)
			switch v.expectedValidateResult {
			case true:
				if err != nil {
					t.Errorf("error should be nil but: '%v'", err)
				}
			case false:
				if err == nil {
					t.Errorf("error should not be nil but nil")
				} else {
					fmt.Printf("%v\n", err)
				}
			}
		})
	}
}

func TestValidatorTemplate(t *testing.T) {
	cases := map[string]struct {
		expectedHttpResult     bool
		expectedResponse       string
		expectedRetry          bool
		expectedRetryResponse  string
		expectedValidateResult bool
	}{
		"POST validate successfully": {
			true,
			"testdata/validate_template.json",
			false,
			"",
			true,
		},
		"Retry successfully after authorization": {
			true,
			mockSDUnauthorizedResponse,
			true,
			"testdata/validate_template.json",
			true,
		},
		"Failure with bad request": {
			false,
			mockSDBadRequestResponse,
			false,
			"",
			false,
		},
		"Failure with invalid yaml": {
			true,
			"testdata/template_parse_error.json",
			false,
			"",
			false,
		},
		"Failure with bad request after retrying": {
			false,
			mockSDUnauthorizedResponse,
			true,
			mockSDBadRequestResponse,
			false,
		},
	}

	for k, v := range cases {
		k := k
		v := v

		t.Run(k, func(t *testing.T) {

			muxAPI := http.NewServeMux()
			testAPIServer := httptest.NewServer(muxAPI)
			defer testAPIServer.Close()

			path := "/v4/validator/template"
			mockSDContext.APIURL = testAPIServer.URL
			muxAPI.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				if v.expectedRetry {
					jwt := r.Header.Get("Authorization")
					if jwt == "Bearer thisissdjwttoken" {
						if !v.expectedHttpResult {
							w.WriteHeader(http.StatusInternalServerError)
						}
						http.ServeFile(w, r, v.expectedRetryResponse)
						return
					}
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					if !v.expectedHttpResult {
						w.WriteHeader(http.StatusInternalServerError)
					}
				}

				http.ServeFile(w, r, v.expectedResponse)
			})
			tokenPath := "/v4/auth/token"
			muxAPI.HandleFunc(tokenPath, func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, mockSDJWTResponse)
			})

			sdapi, err := New(mockSDContext, nil)
			if err != nil {
				t.Fatal("should not cause error")
			}

			err = sdapi.ValidatorTemplate(mockYaml, false)
			switch v.expectedValidateResult {
			case true:
				if err != nil {
					t.Errorf("error should be nil but: '%v'", err)
				}
			case false:
				if err == nil {
					t.Errorf("error should not be nil but nil")
				} else {
					fmt.Printf("%v\n", err)
				}
			}
		})
	}
}

package sdapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
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
	dummyID := "13"
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
			dummyID,
			true,
			"testdata/banner_patch.json",
			false,
		},
		"Failed to update a banner": {
			dummyID,
			false,
			"testdata/banner_not_found.json",
			false,
		},
		"Delete a banner successfully": {
			dummyID,
			true,
			"",
			true,
		},
		"Failed to delete a banner": {
			dummyID,
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

			banner, err := sdapi.UpdateBanner(dummyID, dummyMessage, dummyType, dummyIsActive, v.delete, false)
			switch v.expectedResult {
			case true:
				if err != nil {
					t.Errorf("error should be nil but: '%v'", err)
				}

				expctedResponseJSONFile, err := ioutil.ReadFile(v.expectedResponse)
				if err != nil && !v.delete {
					t.Fatal("should not cause error")
				}
				expectedResponseJSON := new(BannerResponse)
				err = json.Unmarshal(expctedResponseJSONFile, expectedResponseJSON)

				if diff := cmp.Diff(expectedResponseJSON, &banner); diff != "" {
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
		expectedHTTPResult     bool
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
						if !v.expectedHTTPResult {
							w.WriteHeader(http.StatusInternalServerError)
						}
						http.ServeFile(w, r, v.expectedRetryResponse)
						return
					}
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					if !v.expectedHTTPResult {
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
		expectedHTTPResult     bool
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
						if !v.expectedHTTPResult {
							w.WriteHeader(http.StatusInternalServerError)
						}
						http.ServeFile(w, r, v.expectedRetryResponse)
						return
					}
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					if !v.expectedHTTPResult {
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

func TestGetPipelineSecrets(t *testing.T) {
	pipelineID := 1111

	cases := map[string]struct {
		secretsStatusCode int
		expectSecrets     []Secret
		expectErr         error
	}{
		"Get secrets successfully": {
			http.StatusOK,
			[]Secret{
				{
					ID:         11,
					PipelineID: pipelineID,
					Name:       "name1",
					AllowInPR:  true,
				},
				{
					ID:         12,
					PipelineID: pipelineID,
					Name:       "name2",
					AllowInPR:  false,
				},
			},
			nil,
		},
		"Failed to get secrets because of invalid status code": {
			http.StatusUnauthorized,
			nil,
			fmt.Errorf("GET /v4/pipelines/%d/secrets?token=invalid_jwt status code is not %d: %d", pipelineID, http.StatusOK, http.StatusUnauthorized),
		},
	}

	for k, v := range cases {
		k := k
		v := v
		t.Run(k, func(t *testing.T) {
			muxAPI := http.NewServeMux()
			testAPIServer := httptest.NewServer(muxAPI)
			defer testAPIServer.Close()

			path := fmt.Sprintf("/v4/pipelines/%d/secrets", pipelineID)

			muxAPI.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(v.secretsStatusCode)
				secretsJSON, _ := json.Marshal(v.expectSecrets)
				w.Write(secretsJSON)
			})

			mockSDContext.APIURL = testAPIServer.URL
			sdapi, err := New(mockSDContext, nil)
			if err != nil {
				t.Fatal("should not cause error")
			}
			secrets, err := sdapi.getPipelineSecrets(pipelineID)
			if !reflect.DeepEqual(err, v.expectErr) {
				t.Errorf("err should be %#v, but actual is %#v", v.expectErr, err)
			}
			if !reflect.DeepEqual(secrets, v.expectSecrets) {
				t.Errorf("secrets should be %#v, but actual is %#v", v.expectSecrets, secrets)
			}
		})
	}
}

func TestCreateSecret(t *testing.T) {
	pipelineID := 1111
	key := "secretKey"
	value := "secretValue"
	allowInPR := false

	cases := map[string]struct {
		createdStatusCode int
		expectErr         error
	}{
		"Create a secret successfully": {
			http.StatusCreated,
			nil,
		},
		"Failed to create a secrets because of invalid status code": {
			http.StatusUnauthorized,
			fmt.Errorf("POST /v4/secrets status code is not %d: %d", http.StatusCreated, http.StatusUnauthorized),
		},
	}

	for k, v := range cases {
		k := k
		v := v
		t.Run(k, func(t *testing.T) {
			muxAPI := http.NewServeMux()
			testAPIServer := httptest.NewServer(muxAPI)
			defer testAPIServer.Close()

			path := "/v4/secrets"
			muxAPI.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(v.createdStatusCode)
			})

			mockSDContext.APIURL = testAPIServer.URL
			sdapi, err := New(mockSDContext, nil)
			if err != nil {
				t.Fatal("should not cause error")
			}
			actual := sdapi.createSecret(pipelineID, key, value, allowInPR)
			if !reflect.DeepEqual(actual, v.expectErr) {
				t.Errorf("err should be %#v, but actual is %#v", v.expectErr, actual)
			}
		})
	}
}

func TestUpdateSecret(t *testing.T) {
	secretID := 11
	value := "secretValue"
	allowInPR := false

	cases := map[string]struct {
		updatedStatusCode int
		expectErr         error
	}{
		"Updated a secret successfully": {
			http.StatusOK,
			nil,
		},
		"Failed to update a secret because of invalid status code": {
			http.StatusUnauthorized,
			fmt.Errorf("PUT /v4/secrets/%d status code is not %d: %d", secretID, http.StatusOK, http.StatusUnauthorized),
		},
	}

	for k, v := range cases {
		k := k
		v := v
		t.Run(k, func(t *testing.T) {
			muxAPI := http.NewServeMux()
			testAPIServer := httptest.NewServer(muxAPI)
			defer testAPIServer.Close()

			path := fmt.Sprintf("/v4/secrets/%d", secretID)
			muxAPI.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(v.updatedStatusCode)
			})

			mockSDContext.APIURL = testAPIServer.URL
			sdapi, err := New(mockSDContext, nil)
			if err != nil {
				t.Fatal("should not cause error")
			}
			actual := sdapi.updateSecret(secretID, value, allowInPR)
			if !reflect.DeepEqual(actual, v.expectErr) {
				t.Errorf("err should be %#v, but actual is %#v", v.expectErr, actual)
			}
		})
	}
}

func TestSetSecret(t *testing.T) {
	pipelineID := 1111
	value := "secretValue"
	allowInPR := false
	cases := map[string]struct {
		key     string
		secrets []Secret

		secretsStatusCode int
		createdStatusCode int
		updatedStatusCode int

		createCount int
		updateCount int

		expectErr error
	}{
		"Create a secret successfully": {
			"name3",
			[]Secret{
				{
					ID:         11,
					PipelineID: pipelineID,
					Name:       "name1",
					AllowInPR:  true,
				},
				{
					ID:         12,
					PipelineID: pipelineID,
					Name:       "name2",
					AllowInPR:  false,
				},
			},
			http.StatusOK,
			http.StatusCreated,
			0,
			1,
			0,
			nil,
		},
		"Update a secret successfully": {
			"name1",
			[]Secret{
				{
					ID:         11,
					PipelineID: pipelineID,
					Name:       "name1",
					AllowInPR:  true,
				},
				{
					ID:         12,
					PipelineID: pipelineID,
					Name:       "name2",
					AllowInPR:  false,
				},
			},
			http.StatusOK,
			0,
			http.StatusOK,
			0,
			1,
			nil,
		},
	}

	for k, v := range cases {
		k := k
		v := v
		t.Run(k, func(t *testing.T) {
			muxAPI := http.NewServeMux()
			testAPIServer := httptest.NewServer(muxAPI)
			defer testAPIServer.Close()

			pipelineSecretPATH := fmt.Sprintf("/v4/pipelines/%d/secrets", pipelineID)

			muxAPI.HandleFunc(pipelineSecretPATH, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(v.secretsStatusCode)
				secretsJSON, _ := json.Marshal(v.secrets)
				w.Write(secretsJSON)
			})

			var (
				createCount int
				updateCount int
			)
			if v.createdStatusCode != 0 {
				createSecretPATH := "/v4/secrets"
				muxAPI.HandleFunc(createSecretPATH, func(w http.ResponseWriter, r *http.Request) {
					createCount++
					w.WriteHeader(v.createdStatusCode)
				})
			}

			if v.updatedStatusCode != 0 {
				updateSecretPATH := fmt.Sprintf("/v4/secrets/%d", 11)
				muxAPI.HandleFunc(updateSecretPATH, func(w http.ResponseWriter, r *http.Request) {
					updateCount++
					w.WriteHeader(v.updatedStatusCode)
				})
			}

			mockSDContext.APIURL = testAPIServer.URL
			sdapi, err := New(mockSDContext, nil)
			if err != nil {
				t.Fatal("should not cause error")
			}
			actual := sdapi.SetSecret(pipelineID, v.key, value, allowInPR)
			if !reflect.DeepEqual(err, v.expectErr) {
				t.Errorf("err should be %#v, but actual is %#v", v.expectErr, actual)
			}

			if createCount != v.createCount {
				t.Errorf("create count should be %d, but actual is %d", v.createCount, createCount)
			}
			if updateCount != v.updateCount {
				t.Errorf("update count should be %d, but actual is %d", v.updateCount, updateCount)
			}
		})
	}
}

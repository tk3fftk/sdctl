package config

// SdctlConfig is to use Screwdriver.cd API
type SdctlConfig struct {
	UserToken string `json:"token"`
	APIURL    string `json:"api"`
	SDJWT     string `json:"jwt"`
}

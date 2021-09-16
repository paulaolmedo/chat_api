package auth

const (
	// token generation properties
	fileLocation         = "pkg/auth/auth0.properties"
	urlProperty          = "url"
	clientIDProperty     = "client_id"
	clientSecretProperty = "client_secret"
	audienceProperty     = "audience"
	grantTypeProperty    = "grant_type"
)

type OAuthData struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

type OAuthToken struct {
	Token      string `json:"access_token"`
	Expiration string `json:"expires_in"`
	TokenType  string `json:"token_type"`
}

package auth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/magiconair/properties"
)

// ValidateUser checks for a token and validates it
// before allowing the method to execute
func ValidateUser(_ http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: validate token
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	}
}

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

// loadAuthProperties loads the authentication data to generate a token. It should be validated whether the fields are left blank or not to avoid problems
func loadAuthProperties() (OAuthData, string) {
	p := properties.MustLoadFile("pkg/auth/auth0.properties", properties.UTF8)

	URL := p.MustGetString("url")
	clientID := p.MustGetString("client_id")
	clientSecret := p.MustGetString("client_secret")
	audience := p.MustGetString("audience")
	grantType := p.MustGetString("grant_type")

	// no es la mejor forma de devolver esto...
	return OAuthData{ClientID: clientID, ClientSecret: clientSecret, Audience: audience, GrantType: grantType}, URL
}

func GetBearerToken() (OAuthToken, error) {
	data, url := loadAuthProperties()

	requestByte, _ := json.Marshal(data)
	payload := bytes.NewReader(requestByte)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return OAuthToken{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return OAuthToken{}, err
	} // falta validar el hecho de que DENTRO del body venga un error

	var token OAuthToken
	json.Unmarshal(body, &token)

	return token, nil
}

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
		// INSTEAD OF THIS THE TOKEN VALIDATION WILL BE MADE INSIDE OF THE FILE JWT_MIDDLEWARE.GO,
		//WHICH TAKES A BEARER TOKEN AND VERIFIES IT'S AUDIENCE AND ISSUER
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	}
}

// loadAuthProperties loads the authentication data to generate a token. It should be validated whether the fields are left blank or not to avoid problems
func loadAuthProperties() (OAuthData, string) {
	p := properties.MustLoadFile(fileLocation, properties.UTF8)

	URL := p.MustGetString(urlProperty)
	clientID := p.MustGetString(clientIDProperty)
	clientSecret := p.MustGetString(clientSecretProperty)
	audience := p.MustGetString(audienceProperty)
	grantType := p.MustGetString(grantTypeProperty)

	// no es la mejor forma de devolver esto...
	return OAuthData{ClientID: clientID, ClientSecret: clientSecret, Audience: audience, GrantType: grantType}, URL
}

// GetBearerToken generates a Bearer token given certain credentials
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

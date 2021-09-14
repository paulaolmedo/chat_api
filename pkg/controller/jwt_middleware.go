package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/magiconair/properties"
)

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

//getPemCert returns a valid certificate
func getPemCert(token *jwt.Token, issuer string) (string, error) {
	cert := ""
	resp, err := http.Get(issuer + ".well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

// Configura el middleware para solicitar la autenticación con JWT tokens
func SetMiddlewareJWT() *jwtmiddleware.JWTMiddleware {
	p := properties.MustLoadFile("app/app.properties", properties.UTF8)
	issuer_input := p.MustGetString("issuer")
	audience_input := p.MustGetString("audience")

	if issuer_input == "" || audience_input == "" {
		log.Fatalf("invalid variable names, issuer %v, audience %v", issuer_input, audience_input)
		return nil
	}

	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'audience' claim
			m := token.Claims.(jwt.MapClaims)

			checkAud := m.VerifyAudience(audience_input, false)
			if !checkAud {
				log.Fatalf("invalid audience, got %v, want %v", m["aud"], audience_input)
				return token, errors.New("invalid audience")
			}

			// Verify 'issuer' claim
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(issuer_input, false)
			if !checkIss {
				log.Fatalf("invalid issuer, got %v, want %v", m["iss"], issuer_input)
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token, issuer_input)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
}

package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"main/helpers"
	"net/http"
)

func AuthorizationMixer() ([]*http.Cookie, error) {
	auth, err := json.Marshal(helpers.GlobalConfig.Auth)
	if err != nil {
		log.Println("Can't read the file.")
	}

	req, err := http.Post(helpers.GlobalConfig.MIXAPI.ApiPostAuthTest, "application/json", bytes.NewBuffer(auth))
	if err != nil {
		return nil, err
	} else if req.StatusCode != http.StatusOK {
		bodyBytes, readAuthApiHttpBodyError := io.ReadAll(req.Body)
		if readAuthApiHttpBodyError != nil {
			fmt.Println("Error reading the response body of a rejected authorization request.\n", readAuthApiHttpBodyError)
			return nil, readAuthApiHttpBodyError
		}
		return nil, errors.New("Authorization error.\n" + req.Status + " " + string(bodyBytes))
	}

	return req.Cookies(), nil
}

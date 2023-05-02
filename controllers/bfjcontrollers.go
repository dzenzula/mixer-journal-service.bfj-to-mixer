package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"main/cache"
	"main/config"
	"main/models"
	"net/http"
	"strconv"
)

var client = &http.Client{}

func GetListBf() (nBF []int) {
	var data models.ListBF
	url := config.GlobalConfig.BFJAPI.ApiGetListBF
	req, getListOfBFErr := http.Get(url)
	if getListOfBFErr != nil {
		log.Println(getListOfBFErr.Error())
		return nil
	}

	body, readingErr := io.ReadAll(req.Body)
	if readingErr != nil {
		log.Println(readingErr.Error())
		return nil
	}

	jsonError := json.Unmarshal(body, &data)
	if jsonError != nil {
		log.Println(jsonError.Error())
		return nil
	}

	return data.Name
}

func GetLastBFJJournalsData(nBF []int) (ids map[int][]int) {
	var data models.Journals
	var cookies []*http.Cookie
	ids = map[int][]int{}

	for _, n := range nBF {
		endpoint := fmt.Sprintf(config.GlobalConfig.BFJAPI.ApiGetLastJournals, strconv.Itoa(n))
		err := getBfjApiResponse(endpoint, cookies, &data)
		if err != nil {
			return nil
		}

		for i := 0; i < 1; i++ {
			ids[n] = append(ids[n], data.DataJournals[i].ID)
		}
	}

	cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, ids, nil)
	return ids
}

func GetBFJTappings(journalId int, cookies []*http.Cookie) (tappingIds []models.Tapping) {
	var data []models.Tapping
	endpoint := fmt.Sprintf(config.GlobalConfig.BFJAPI.ApiGetTappings, strconv.Itoa(journalId))
	err := getBfjApiResponse(endpoint, cookies, &data)
	if err != nil {
		return nil
	}
	return data
}

func AuthorizeBFJ() (cookies []*http.Cookie, cookiesErr error) {
	auth, err := json.Marshal(config.GlobalConfig.Auth)
	if err != nil {
		log.Println("Can't read the file.")
	}

	req, err := http.Post(config.GlobalConfig.BFJAPI.ApiPostAuthProd, "application/json", bytes.NewBuffer(auth))
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

func getBfjApiResponse(endpoint string, cookies []*http.Cookie, data interface{}) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Println(err.Error())
	}

	countCookies := len(cookies)
	for i := 0; i < countCookies; i++ {
		req.AddCookie(cookies[i])
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	} else if resp.StatusCode != http.StatusOK {
		cookies, authError := AuthorizeBFJ()
		if authError != nil {
			log.Println("Failed to get new cookies:", authError)
			return authError
		}
		getBfjApiResponse(endpoint, cookies, data)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Error decoding JSON string:", err)
		return err
	}

	return nil
}

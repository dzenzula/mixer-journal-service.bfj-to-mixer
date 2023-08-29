package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/config"
	"main/logger"
	"main/models"
	"net/http"
	"strconv"
	"time"
)

var client = &http.Client{}

func GetListBf() (nBF []int) {
	var data models.ListBF
	url := config.GlobalConfig.BFJAPI.ApiGetListBF
	req, getListOfBFErr := http.Get(url)
	if getListOfBFErr != nil {
		logger.Logger.Println(getListOfBFErr.Error())
		return nil
	}

	body, readingErr := io.ReadAll(req.Body)
	if readingErr != nil {
		logger.Logger.Println(readingErr.Error())
		return nil
	}

	jsonError := json.Unmarshal(body, &data)
	if jsonError != nil {
		logger.Logger.Println(jsonError.Error())
		return nil
	}

	return data.Name
}

func GetLastBFJJournalsData(nBF []int, ids *map[int][]int) {
	var data models.Journals

	if *ids != nil {
		*ids = make(map[int][]int)
	}

	for _, n := range nBF {
		endpoint := fmt.Sprintf(config.GlobalConfig.BFJAPI.ApiGetLastJournals, strconv.Itoa(n))
		err := getBfjApiResponse(endpoint, nil, &data)
		if err != nil {
			return
		}

		for i := 0; i < 2; i++ {
			(*ids)[n] = append((*ids)[n], data.DataJournals[i].ID)
		}
	}
}

func GetBFJTappings(journalId int, cookies *[]*http.Cookie) (tappingIds []models.Tapping) {
	var data []models.Tapping
	endpoint := fmt.Sprintf(config.GlobalConfig.BFJAPI.ApiGetTappings, strconv.Itoa(journalId))
	err := getBfjApiResponse(endpoint, cookies, &data)
	if err != nil {
		return nil
	}
	data = tappingFilterKC(data)
	return data
}

func tappingFilterKC(tappings []models.Tapping) []models.Tapping {
	for i, tapping := range tappings {
		var ladlesKC []models.Ladle
		for _, ladle := range tapping.ListLaldes {
			if ladle.Destination == "КЦ" {
				ladlesKC = append(ladlesKC, ladle)
			}
		}
		tappings[i].ListLaldes = nil
		tappings[i].ListLaldes = ladlesKC
	}
	return tappings
}

func AuthorizeBFJ(cookies *[]*http.Cookie) {
	auth, _ := json.Marshal(config.GlobalConfig.Auth)

	for {
		success := true
		req, err := http.Post(config.GlobalConfig.BFJAPI.ApiPostAuthProd, "application/json", bytes.NewBuffer(auth))
		if err != nil {
			success = false
			logger.Logger.Printf("Failed to send authorization request: %v", err)
		}

		if req != nil {
			defer req.Body.Close()

			if req.StatusCode != http.StatusOK {
				success = false
				bodyBytes, err := io.ReadAll(req.Body)
				if err != nil {
					logger.Logger.Printf("Failed to read authorization response body: %v\n", err)
				}
				logger.Logger.Printf("Rejected authorization request: %s\n", bodyBytes)
			}
		}

		if success {
			*cookies = req.Cookies()
			return
		} else {
			logger.Logger.Println("Next try to authorize will be in a 5 minutes")
			time.Sleep(time.Minute * 5)
		}
	}
}

func getBfjApiResponse(endpoint string, cookies *[]*http.Cookie, data interface{}) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		logger.Logger.Println(err.Error())
	}

	if cookies != nil {
		for _, cookie := range *cookies {
			req.AddCookie(cookie)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Println(err.Error())
	} else if resp.StatusCode != http.StatusOK {
		AuthorizeBFJ(cookies)
		getBfjApiResponse(endpoint, cookies, data)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Println(err.Error())
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Logger.Println("Error decoding JSON string:", err)
		return err
	}

	return nil
}

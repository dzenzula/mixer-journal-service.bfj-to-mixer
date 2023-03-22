package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"main/cache"
	"main/helpers"
	"main/models"
	"net/http"
	"os"
	"strconv"
)

var configAPI *models.ConfigAPI = helpers.LoadAPIConfig("models/json/configAPI.json")
var configPath *models.ConfigPath = helpers.LoadPathConfig("models/json/configPaths.json")
var client = &http.Client{}

func GetListBf() (nBF []int) {
	url := configAPI.ApiGetListBF
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

	var data models.ListBF
	jsonError := json.Unmarshal(body, &data)
	if jsonError != nil {
		log.Println(jsonError.Error())
		return nil
	}

	return data.Name
}

func GetLastJournalsData(nBF []int) (ids map[int][]int) {
	var data models.Journals
	ids = map[int][]int{}

	var yaml *cache.Data = cache.ReadYAMLFile(configPath.CachePath)

	for _, n := range nBF {
		req, err := http.NewRequest("GET", configAPI.ApiGetLastJournals+strconv.Itoa(n), nil)
		if err != nil {
			log.Println(err.Error())
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err.Error())
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err.Error())
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println("Error decoding JSON string:", err)
			return nil
		}

		// out, err := json.MarshalIndent(data, "", "    ")
		// if err != nil {
		// 	log.Println("Error decoding JSON string:", err)
		// 	return nil
		// }

		for _, id := range data.DataJournals {
			if !cache.IdExists(yaml, id.ID) {
				cache.WriteYAMLFile(configPath.CachePath, yaml, id.ID)
				ids[n] = append(ids[n], id.ID)
			}
		}

		// fmt.Println(string(out))
		defer resp.Body.Close()
	}

	return ids
}

func GetJournalData(ids map[int][]int) {
	cookies, authError := authorize()
	for key, values := range ids {
		for _, id := range values {
			if authError != nil {
				log.Println("Authorization error. \n", authError)
				return
			}

			req, err := http.NewRequest("GET", configAPI.ApiGetjournal+strconv.Itoa(id)+"&"+strconv.Itoa(key), nil)
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
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err.Error())
			}

			var data models.Journal
			err = json.Unmarshal(body, &data)
			if err != nil {
				log.Println("Error decoding JSON string:", err)
				return
			}

			// out, err := json.MarshalIndent(data, "", "    ")
			// if err != nil {
			// 	log.Println("Error decoding JSON string:", err)
			// 	return
			// }

			// fmt.Println(string(out))
			resp.Body.Close()
		}
	}
}

func authorize() ([]*http.Cookie, error) {
	file, err := os.Open(configPath.AuthPath)
	if err != nil {
		log.Println(err.Error())
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Println("Can't read the file.")
	}

	req, err := http.Post(configAPI.ApiPostAuth, "application/json", bytes.NewBuffer(fileContent))
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

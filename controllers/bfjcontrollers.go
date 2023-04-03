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
	"strconv"
)

var client = &http.Client{}

func GetListBf() (nBF []int) {
	var data models.ListBF
	url := helpers.GlobalConfig.BFJAPI.ApiGetListBF
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

func GetLastJournalsData(nBF []int) (ids map[int][]int) {
	var data models.Journals
	var cookies []*http.Cookie
	ids = map[int][]int{}

	var yaml *cache.Data = &cache.Data{}

	for _, n := range nBF {
		endpoint := fmt.Sprintf(helpers.GlobalConfig.BFJAPI.ApiGetLastJournals, strconv.Itoa(n))
		err := getAPIResponse(endpoint, cookies, &data)
		if err != nil {
			return nil
		}

		// out, err := json.MarshalIndent(data, "", "    ")
		// if err != nil {
		// 	log.Println("Error decoding JSON string:", err)
		// 	return nil
		// }

		for i := 0; i < 2; i++ {
			ids[n] = append(ids[n], data.DataJournals[i].ID)
		}
		// fmt.Println(string(out))
	}

	cache.WriteYAMLFile(helpers.GlobalConfig.Path.CachePath, yaml, ids)

	return ids
}

func GetJournalDatas(ids map[int][]int, cookies []*http.Cookie) {
	var data models.Journal
	for key, values := range ids {
		for _, id := range values {
			endpoint := fmt.Sprintf(helpers.GlobalConfig.BFJAPI.ApiGetjournal, strconv.Itoa(id), strconv.Itoa(key))
			err := getAPIResponse(endpoint, cookies, &data)
			if err != nil {
				log.Println("Couldn't get datas journal")
			}

			// out, err := json.MarshalIndent(data, "", "    ")
			// if err != nil {
			// 	log.Println("Error decoding JSON string:", err)
			// 	return
			// }

			// fmt.Println(string(out))
		}
	}
}

func GetJournalData(nBF int, journalId int, cookies []*http.Cookie) (idJournal models.Journal) {
	var data models.Journal
	endpoint := fmt.Sprintf(helpers.GlobalConfig.BFJAPI.ApiGetjournal, strconv.Itoa(journalId), strconv.Itoa(nBF))
	err := getAPIResponse(endpoint, cookies, &data)
	if err != nil {
		return
	}

	// out, err := json.MarshalIndent(data, "", "    ")
	// if err != nil {
	// 	log.Println("Error decoding JSON string:", err)
	// 	return
	// }

	// fmt.Println(string(out))

	return data
}

func GetChemCoxes(journalId int, cookies []*http.Cookie) (chemCoxes []models.ChemCoxe) {
	var data []models.ChemCoxe
	endpoint := fmt.Sprintf(helpers.GlobalConfig.BFJAPI.ApiGetChemCoxes, strconv.Itoa(journalId))
	err := getAPIResponse(endpoint, cookies, &data)
	if err != nil {
		return nil
	}

	// out, err := json.MarshalIndent(data, "", "    ")
	// if err != nil {
	// 	log.Println("Error decoding JSON string:", err)
	// 	return
	// }

	// fmt.Println(string(out))
	return data
}

func GetChemicalSlags(journalId int, cookies []*http.Cookie) (chemicalSlags []models.ChemicalSlag) {
	var data []models.ChemicalSlag
	endpoint := fmt.Sprintf(helpers.GlobalConfig.BFJAPI.ApiGetChemicalsSlags, strconv.Itoa(journalId))
	err := getAPIResponse(endpoint, cookies, &data)
	if err != nil {
		return nil
	}

	out, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Println("Error decoding JSON string:", err)
		return
	}

	fmt.Println(string(out))
	return data
}

func GetChemMaterials(journalId int, cookies []*http.Cookie) (chemicalMaterials []models.ChemMaterial) {

	var data []models.ChemMaterial
	endpoint := fmt.Sprintf(helpers.GlobalConfig.BFJAPI.ApiGetChemMaterials, strconv.Itoa(journalId))
	err := getAPIResponse(endpoint, cookies, &data)
	if err != nil {
		return nil
	}

	// out, err := json.MarshalIndent(data, "", "    ")
	// if err != nil {
	// 	log.Println("Error decoding JSON string:", err)
	// 	return
	// }

	// fmt.Println(string(out))
	return data
}

func GetTappings(journalId int, cookies []*http.Cookie) (tappingIds []int) {
	var data []models.Tapping
	endpoint := fmt.Sprintf(helpers.GlobalConfig.BFJAPI.ApiGetTappings, strconv.Itoa(journalId))
	err := getAPIResponse(endpoint, cookies, &data)
	if err != nil {
		return nil
	}

	var Ids []int
	for _, tapping := range data {
		for _, ladle := range tapping.ListLaldes {
			tappingIds = append(Ids, ladle.IDTapping)
		}
	}

	// out, err := json.MarshalIndent(data, "", "    ")
	// if err != nil {
	// 	log.Println("Error decoding JSON string:", err)
	// 	return
	// }

	// fmt.Println(string(out))
	return tappingIds
}

func AuthorizeProd() (cookies []*http.Cookie, cookiesErr error) {
	auth, err := json.Marshal(helpers.GlobalConfig.Auth)
	if err != nil {
		log.Println("Can't read the file.")
	}

	req, err := http.Post(helpers.GlobalConfig.BFJAPI.ApiPostAuthProd, "application/json", bytes.NewBuffer(auth))
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

func getAPIResponse(endpoint string, cookies []*http.Cookie, data interface{}) error {
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
		cookies, authError := AuthorizeProd()
		if authError != nil {
			log.Println("Failed to get new cookies:", authError)
			return authError
		}
		getAPIResponse(endpoint, cookies, data)
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

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

var client = &http.Client{}

func GetListBf() (nBF []int) {
	url := helpers.CfgAPI.ApiGetListBF
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
	var isDelete bool = false
	ids = map[int][]int{}

	var yaml *cache.Data = cache.ReadYAMLFile(helpers.CfgPath.CachePath)

	for _, n := range nBF {
		req, err := http.NewRequest("GET", fmt.Sprintf(helpers.CfgAPI.ApiGetLastJournals, strconv.Itoa(n)), nil)
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

		for i := 0; i < 4; i++ {
			if !cache.IdExists(yaml, data.DataJournals[i].ID) {
				ids[n] = append(ids[n], data.DataJournals[i].ID)
				isDelete = true
			}
		}

		if len(ids) == 0 {
			for i := 0; i < 4; i++ {
				ids[n] = append(ids[n], data.DataJournals[i].ID)
			}
		}
		// fmt.Println(string(out))
		defer resp.Body.Close()
	}

	if isDelete {
		yaml = cache.DeleteIds(helpers.CfgPath.CachePath, yaml)
		isDelete = false
	}

	cache.WriteYAMLFile(helpers.CfgPath.CachePath, yaml, ids)

	return ids
}

func GetJournalDatas(ids map[int][]int, countCookies int) {
	cookies, authError := AuthorizeProd()
	if authError != nil {
		log.Println("Authorization error. \n", authError)
		return
	}
	for key, values := range ids {
		for _, id := range values {
			req, err := http.NewRequest("GET", fmt.Sprintf(helpers.CfgAPI.ApiGetjournal, strconv.Itoa(id), strconv.Itoa(key)), nil)
			if err != nil {
				log.Println(err.Error())
			}

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
			defer resp.Body.Close()
		}
	}
}

func GetJournalData(nBF int, id int, cookies []*http.Cookie) (idJournal int) {
	req, err := http.NewRequest("GET", fmt.Sprintf(helpers.CfgAPI.ApiGetjournal, strconv.Itoa(id), strconv.Itoa(nBF)), nil)
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
			return
		}
		return GetJournalData(nBF, id, cookies)
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
	defer resp.Body.Close()

	return data.DataJournals.ID
}

func GetChemCoxes(id int, cookies []*http.Cookie) (chemCoxes []models.ChemCoxe) {
	req, err := http.NewRequest("GET", fmt.Sprintf(helpers.CfgAPI.ApiGetChemCoxes, strconv.Itoa(id)), nil)
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
			return
		}
		return GetChemCoxes(id, cookies)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()

	var data []models.ChemCoxe
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
	return data
}

func GetChemicalSlags(id int, cookies []*http.Cookie) (chemicalSlags []models.ChemicalSlag) {
	req, err := http.NewRequest("GET", fmt.Sprintf(helpers.CfgAPI.ApiGetChemicalsSlags, strconv.Itoa(id)), nil)
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
			return
		}
		return GetChemicalSlags(id, cookies)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()

	var data []models.ChemicalSlag
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
	return data
}

func GetChemMaterials(id int, cookies []*http.Cookie) (chemicalMaterials []models.ChemMaterial) {
	req, err := http.NewRequest("GET", fmt.Sprintf(helpers.CfgAPI.ApiGetChemMaterials, strconv.Itoa(id)), nil)
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
			return
		}
		return GetChemMaterials(id, cookies)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()

	var data []models.ChemMaterial
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
	return data
}

func GetTappings(journalId int, cookies []*http.Cookie) (tappingIds []int) {
	req, err := http.NewRequest("GET", fmt.Sprintf(helpers.CfgAPI.ApiGetTappings, strconv.Itoa(journalId)), nil)
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
			return
		}
		return GetTappings(journalId, cookies)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()

	var data []models.Tapping
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Error decoding JSON string:", err)
		return
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

// func authorizeTest() ([]*http.Cookie, error) {
// 	file, err := os.Open(helpers.CfgPath.AuthPath)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// 	defer file.Close()

// 	fileContent, err := io.ReadAll(file)
// 	if err != nil {
// 		log.Println("Can't read the file.")
// 	}

// 	req, err := http.Post(helpers.CfgAPI.ApiPostAuthTest, "application/json", bytes.NewBuffer(fileContent))
// 	if err != nil {
// 		return nil, err
// 	} else if req.StatusCode != http.StatusOK {
// 		bodyBytes, readAuthApiHttpBodyError := io.ReadAll(req.Body)
// 		if readAuthApiHttpBodyError != nil {
// 			fmt.Println("Error reading the response body of a rejected authorization request.\n", readAuthApiHttpBodyError)
// 			return nil, readAuthApiHttpBodyError
// 		}
// 		return nil, errors.New("Authorization error.\n" + req.Status + " " + string(bodyBytes))
// 	}

// 	return req.Cookies(), nil
// }

func AuthorizeProd() ([]*http.Cookie, error) {
	file, err := os.Open(helpers.CfgPath.AuthPath)
	if err != nil {
		log.Println(err.Error())
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Println("Can't read the file.")
	}

	req, err := http.Post(helpers.CfgAPI.ApiPostAuthProd, "application/json", bytes.NewBuffer(fileContent))
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

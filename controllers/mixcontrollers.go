package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"main/config"
	"main/models"
	"net/http"
	"strconv"
)

func GetLastMIXJournalsData(nMIX []int) (ids map[int][]int) {
	var data models.MixJournals
	var cookies []*http.Cookie
	ids = map[int][]int{}

	for _, n := range nMIX {
		endpoint := fmt.Sprintf(config.GlobalConfig.MIXAPI.ApiGetLastJournals, strconv.Itoa(n))
		err := getMixApiResponse(endpoint, cookies, &data)
		if err != nil {
			return nil
		}

		if len(data.DataJournals) > 0 {
			for i := 0; i < 1; i++ {
				ids[n] = append(ids[n], data.DataJournals[i].ID)
			}
		}
	}

	return ids
}

func PostMixChemical(tapping models.Tapping, cookies []*http.Cookie) {
	listLadles := tapping.ListLaldes
	postMixChemical(listLadles, cookies)
}

func PostMixChemicalList(listLadles []models.Ladle, cookies []*http.Cookie) {
	postMixChemical(listLadles, cookies)
}

func postMixChemical(listLadles []models.Ladle, cookies []*http.Cookie) {
	endpoint := fmt.Sprintf(config.GlobalConfig.MIXAPI.ApiPostChemical)
	for nMix := 1; nMix < 5; nMix++ {
		for _, ladle := range listLadles {
			chem := models.ChemicalDTO{
				NMix:       nMix,
				Ladle:      ladle.Ladle,
				Proba:      int(ladle.Chemical.Proba),
				NumTaphole: ladle.Chemical.NumTaphole,
				DT:         checkChemDate(ladle),
				Si:         int(ladle.Chemical.Si),
				Mn:         int(ladle.Chemical.Mn),
				S:          int(ladle.Chemical.S),
				P:          int(ladle.Chemical.P),
				Belong:     "LadleMovement",
			}
			postMixApiRequest(endpoint, cookies, chem)
		}
	}
}

func checkChemDate(ldl models.Ladle) *string {
	if ldl.Chemical.Dt != "" {
		return &ldl.Chemical.Dt
	}
	return nil
}

func PostMixLadleMovement(nBf int, tapping models.Tapping) models.LadleMovement {
	var ldlMvm models.LadleMovement

	ldlMvm.Date = tapping.DtCloseTaphole
	ldlMvm.NumDp = nBf
	ldlMvm.NumTapping = tapping.NumTapping
	ldlMvm.DtCloseTaphole = tapping.DtCloseTaphole
	ldlMvm.TemperExhaustIron = int(tapping.Temper)

	return ldlMvm
}

func PostMixListLadles(listLadles []models.Ladle, ldlMvm models.LadleMovement, mixIds map[int][]int, cookies []*http.Cookie) {
	for i := 0; i < len(listLadles); i++ {
		for _, keys := range mixIds {
			ldlMvm.LadleTapping = listLadles[i].Ladle
			ldlMvm.MassCastIron = int(listLadles[i].Weight)

			for _, key := range keys {
				endpoint := fmt.Sprintf(config.GlobalConfig.MIXAPI.ApiPostLadleMovement, strconv.Itoa(key))
				postErr := postMixApiRequest(endpoint, cookies, ldlMvm)
				if postErr != nil {
					fmt.Println(endpoint + " Success!")
				}

			}
		}
	}
}

func AuthorizeMix() ([]*http.Cookie, error) {
	auth, err := json.Marshal(config.GlobalConfig.Auth)
	if err != nil {
		log.Println("Can't marshal the data")
	}

	req, err := http.Post(config.GlobalConfig.MIXAPI.ApiPostAuthTest, "application/json", bytes.NewBuffer(auth))
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

func postMixApiRequest(endpoint string, cookies []*http.Cookie, requestData interface{}) error {
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		log.Println("Error encoding request JSON:", err)
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return err
	}
	req.Header.Set("content-type", "application/json")

	countCookies := len(cookies)
	for i := 0; i < countCookies; i++ {
		req.AddCookie(cookies[i])
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error executing HTTP request:", err)
		return err
	} else if resp.StatusCode != http.StatusOK {
		cookies, authError := AuthorizeMix()
		if authError != nil {
			log.Println("Failed to get new cookies:", authError)
			return authError
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err.Error())
		}
		defer resp.Body.Close()

		fmt.Println(string(body))

		if string(body) == "NotAuthorized" {
			postMixApiRequest(endpoint, cookies, requestData)
		} else {
			return nil
		}

	}

	return nil
}

func getMixApiResponse(endpoint string, cookies []*http.Cookie, data interface{}) error {
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
		cookies, authError := AuthorizeMix()
		if authError != nil {
			log.Println("Failed to get new cookies:", authError)
			return authError
		}
		getMixApiResponse(endpoint, cookies, data)
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

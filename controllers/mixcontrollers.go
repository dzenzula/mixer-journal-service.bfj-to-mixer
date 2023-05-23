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

func GetLastMIXJournalsData(nMIX []int, cookies *[]*http.Cookie, ids *map[int][]int) {
	var data models.MixJournals

	if *ids != nil {
		*ids = make(map[int][]int)
	}

	for _, n := range nMIX {
		endpoint := fmt.Sprintf(config.GlobalConfig.MIXAPI.ApiGetLastJournals, strconv.Itoa(n))
		err := getMixApiResponse(endpoint, cookies, &data)
		if err != nil {
			return
		}

		if len(data.DataJournals) > 0 {
			for i := 0; i < 1; i++ {
				(*ids)[n] = append((*ids)[n], data.DataJournals[i].ID)
			}
		}
	}
}

func PostMixChemical(tapping models.Tapping, cookies *[]*http.Cookie) {
	listLadles := tapping.ListLaldes
	postMixChemical(listLadles, cookies)
}

func PostMixChemicalList(listLadles []models.Ladle, cookies *[]*http.Cookie) {
	postMixChemical(listLadles, cookies)
}

func postMixChemical(listLadles []models.Ladle, cookies *[]*http.Cookie) {
	endpoint := fmt.Sprintf(config.GlobalConfig.MIXAPI.ApiPostChemical)
	for nMix := 1; nMix < 5; nMix++ {
		for _, ladle := range listLadles {
			chem := models.ChemicalDTO{
				NMix:       nMix,
				Ladle:      ladle.Ladle,
				NumSample:  int(ladle.Chemical.Proba),
				NumTaphole: ladle.Chemical.NumTaphole,
				DT:         checkChemDate(ladle),
				Si:         float64(ladle.Chemical.Si),
				Mn:         float64(ladle.Chemical.Mn),
				S:          float64(ladle.Chemical.S),
				P:          float64(ladle.Chemical.P),
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

func PostMixListLadles(listLadles []models.Ladle, ldlMvm models.LadleMovement, mixIds *map[int][]int, cookies *[]*http.Cookie) {
	for i := 0; i < len(listLadles); i++ {
		for _, keys := range *mixIds {
			ldlMvm.LadleTapping = listLadles[i].Ladle
			ldlMvm.MassCastIron = int(listLadles[i].Weight)

			for _, key := range keys {
				endpoint := fmt.Sprintf(config.GlobalConfig.MIXAPI.ApiPostLadleMovement, strconv.Itoa(key))
				postErr := postMixApiRequest(endpoint, cookies, ldlMvm)
				if postErr != nil {
					logger.Logger.Println(postErr.Error())
					
				}

			}
		}
	}
}

func AuthorizeMix(cookies *[]*http.Cookie) {
	authJSON, _ := json.Marshal(config.GlobalConfig.Auth)

	for {
		success := true
		req, err := http.Post(config.GlobalConfig.MIXAPI.ApiPostAuthTest, "application/json", bytes.NewBuffer(authJSON))
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
			logger.Logger.Println("Authorization MIX success!")
			*cookies = req.Cookies()
			return
		} else {
			logger.Logger.Println("Next try to authorize will be in a 5 minutes")
			time.Sleep(time.Minute * 5)
		}
	}
}

func postMixApiRequest(endpoint string, cookies *[]*http.Cookie, requestData interface{}) error {
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		logger.Logger.Println("Error encoding request JSON:", err)
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Logger.Println("Error creating HTTP request:", err)
		return err
	}
	req.Header.Set("content-type", "application/json")

	for _, cookie := range *cookies {
		req.AddCookie(cookie)
	}

	resp, errResp := client.Do(req)

	if errResp != nil {
		logger.Logger.Println("Error executing HTTP request:", err)
		return err
	} else if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Logger.Println(err.Error())
		}
		defer resp.Body.Close()

		logger.Logger.Println("Respond PostMIX error: ", string(body))

		if string(body) == "NotAuthorized" {
			AuthorizeMix(cookies)
			postMixApiRequest(endpoint, cookies, requestData)
		} else {
			return nil
		}
	}

	return nil
}

func getMixApiResponse(endpoint string, cookies *[]*http.Cookie, data interface{}) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		logger.Logger.Println(err.Error())
	}

	for _, cookie := range *cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)

	body, errResp := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Println(err.Error())
	}
	defer resp.Body.Close()

	if errResp != nil {
		logger.Logger.Println(err.Error())
		return errResp
	} else if resp.StatusCode != http.StatusOK {
		logger.Logger.Println("Respond GetMIX error: ", string(body))

		if string(body) == "NotAuthorized" {
			AuthorizeMix(cookies)
			getMixApiResponse(endpoint, cookies, data)
		} else {
			return nil
		}
		return nil
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Logger.Println("Error decoding JSON string:", err)
		return err
	}

	return nil
}

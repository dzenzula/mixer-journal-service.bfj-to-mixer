package main

import (
	"fmt"
	"log"
	"main/cache"
	"main/config"
	"main/controllers"
	"main/models"
	"net/http"
	"time"
)

func main() {
	bfjCookies, bfjErr := controllers.AuthorizeBFJ()
	if bfjErr != nil {
		log.Println("Authorization error. \n", bfjErr)
	}

	mixCookies, mixErr := controllers.AuthorizeMix()
	if mixErr != nil {
		log.Println("Authorization error. \n", mixErr)
	}

	var nBF []int = controllers.GetListBf()
	var nMix []int = []int{1, 2, 3, 4}

	bfjIds := controllers.GetLastBFJJournalsData(nBF)
	mixIds := controllers.GetLastMIXJournalsData(nMix)

	duration := time.Until(time.Now().Truncate(time.Minute).Add(time.Minute))
	time.Sleep(duration)
	ticker := time.NewTicker(1 * time.Minute).C
	fmt.Println(time.Now().String())

	for {
		service(nBF, &bfjIds, bfjCookies, mixIds, mixCookies, ticker)
	}
}

func service(nBF []int, bfjIds *map[int][]int, bfjCookies []*http.Cookie, mixIds map[int][]int, mixCookies []*http.Cookie, ticker <-chan time.Time) {
	now := <-ticker
	if (now.Hour() == 8 && now.Minute() == 0) || (now.Hour() == 20 && now.Minute() == 0) {
		var clear []map[int]int
		newIds := controllers.GetLastBFJJournalsData(nBF)
		if len(newIds) != 0 {
			bfjIds = &newIds
		}
		cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, newIds, clear)
		fmt.Println(now.String(), "Worked fine shift")
	} else if now.Minute() == 0 {
		tIds := cache.ReadYAMLFile(config.GlobalConfig.Path.CachePath)
		for nBf, values := range *bfjIds {
			for _, idJournal := range values {

				tappings := controllers.GetBFJTappings(idJournal, bfjCookies)

				for _, tapping := range tappings {
					sendLadleMovements(nBf, tIds, tapping, mixIds, mixCookies)
				}

			}
		}
		cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, nil, tIds.Tappings)
		fmt.Println(now.String(), "Worked fine hour check")
	}
}

func sendLadleMovements(nBf int, tIds *cache.Data, tapping models.Tapping, mixIds map[int][]int, mixCookies []*http.Cookie) {
	if !cache.TappingIdExists(tIds, tapping.ID) {
		ldlMvm := controllers.PostMixLadleMovement(nBf, tapping)
		controllers.PostMixListLadles(tapping.ListLaldes, ldlMvm, mixIds, mixCookies)
		controllers.PostMixChemical(tapping, mixCookies)
		tIds.Tappings = append(tIds.Tappings, map[int]int{tapping.ID: len(tapping.ListLaldes)})
	} else {
		tVal := cache.FindTappingIdValue(tIds, tapping.ID)
		if len(tapping.ListLaldes) != tVal {
			numMissingLadles := len(tapping.ListLaldes) - tVal
			missingLadles := tapping.ListLaldes[len(tapping.ListLaldes)-numMissingLadles:]
			ldlMvm := controllers.PostMixLadleMovement(nBf, tapping)
			controllers.PostMixListLadles(missingLadles, ldlMvm, mixIds, mixCookies)
			controllers.PostMixChemicalList(missingLadles, mixCookies)
			cache.UpdateTappingValue(tIds, tapping.ID, len(tapping.ListLaldes))
			cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, nil, tIds.Tappings)
		}
	}
}

package main

import (
	"fmt"
	"log"
	"main/cache"
	"main/config"
	"main/controllers"
	"main/models"
	"net/http"
	"reflect"
	"time"
)

var currList map[int][]models.Ladle = make(map[int][]models.Ladle)

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

	data := cache.ReadYAMLFile(config.GlobalConfig.Path.CachePath)
	cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, bfjIds, data.Tappings)

	duration := time.Until(time.Now().Truncate(time.Minute).Add(time.Minute))
	time.Sleep(duration)

	fmt.Println(time.Now().String())

	for {
		service(nBF, &bfjIds, bfjCookies, nMix, &mixIds, mixCookies)

		ticker := time.NewTicker(time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)))
		<-ticker.C
	}
}

func service(nBF []int, bfjIds *map[int][]int, bfjCookies []*http.Cookie, nMix []int, mixIds *map[int][]int, mixCookies []*http.Cookie) {
	now := time.Now().Truncate(time.Minute)
	if (now.Hour() == 8 && now.Minute() == 0) || (now.Hour() == 20 && now.Minute() == 0) {
		currList = make(map[int][]models.Ladle)
		fmt.Println(time.Now().String(), "Shift started")

		newBfjIds := controllers.GetLastBFJJournalsData(nBF)
		newMixIds := controllers.GetLastMIXJournalsData(nMix)

		if len(newBfjIds) != 0 {
			bfjIds = &newBfjIds
		}
		if len(newMixIds) != 0 {
			mixIds = &newMixIds
		}

		hourCheck(bfjIds, bfjCookies, mixIds, mixCookies)

		cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, *bfjIds, nil)
		fmt.Println(time.Now().String(), "First data shift transfer")
	} else if now.Second() == 0 {
		fmt.Println(time.Now().String(), "Hour check started")

		hourCheck(bfjIds, bfjCookies, mixIds, mixCookies)

		fmt.Println(time.Now().String(), "Hour check finished")
	}
}

func hourCheck(bfjIds *map[int][]int, bfjCookies []*http.Cookie, mixIds *map[int][]int, mixCookies []*http.Cookie) {
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
}

func sendLadleMovements(nBf int, tIds *cache.Data, tapping models.Tapping, mixIds *map[int][]int, mixCookies []*http.Cookie) {
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

	if currList[tapping.ID] != nil {
		if !reflect.DeepEqual(currList[tapping.ID], tapping.ListLaldes) {
			fmt.Println(time.Now().Truncate(time.Minute).String(), "Chemical changed!")
			oldLadles := currList[tapping.ID]
			currLadles := tapping.ListLaldes
			changedLadles := []models.Ladle{}
			for _, newLadle := range currLadles {
				for _, currLadle := range oldLadles {
					if newLadle.Ladle == currLadle.Ladle {
						if !reflect.DeepEqual(newLadle.Chemical, currLadle.Chemical) {
							changedLadles = append(changedLadles, newLadle)
						}
					}
				}
			}
			controllers.PostMixChemicalList(changedLadles, mixCookies)
			currList[tapping.ID] = tapping.ListLaldes
		}
	} else {
		currList[tapping.ID] = append(currList[tapping.ID], tapping.ListLaldes...)
	}
}

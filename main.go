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

	morningShiftTick, eveningShiftTick, oneHourTick := initializeTimers()

	bfjIds := controllers.GetLastBFJJournalsData(nBF)
	mixIds := controllers.GetLastMIXJournalsData(nMix)

	for {
		select {
		case <-morningShiftTick:
			var clear []map[int]int
			newIds := controllers.GetLastBFJJournalsData(nBF)
			if len(newIds) != 0 {
				bfjIds = newIds
			}
			cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, nil, clear)
			fmt.Println("Morning shift works")
		case <-eveningShiftTick:
			var clear []map[int]int
			newIds := controllers.GetLastBFJJournalsData(nBF)
			if len(newIds) != 0 {
				bfjIds = newIds
			}
			cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, nil, clear)
			fmt.Println("Evening shift works")
		case tm := <-oneHourTick:
			tIds := cache.ReadYAMLFile(config.GlobalConfig.Path.CachePath)
			for nBf, values := range bfjIds {
				for _, idJournal := range values {
					go func(nBf int, idJournal int, bfjCookies []*http.Cookie) {
						tappings := controllers.GetBFJTappings(idJournal, bfjCookies)

						for _, tapping := range tappings {
							sendLadleMovements(nBf, tIds, tapping, mixIds, mixCookies)
						}
					}(nBf, idJournal, bfjCookies)
				}
			}
			cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, nil, tIds.Tappings)
			fmt.Println(tm.String() + " 1hour works")

			/* case tm := <-minuteTick:
			tIds := cache.ReadYAMLFile(config.GlobalConfig.Path.CachePath)
			for nBf, values := range bfjIds {
				for _, idJournal := range values {
					//go func(nBf int, idJournal int, bfjCookies []*http.Cookie) {
					tappings := controllers.GetBFJTappings(idJournal, bfjCookies)

					for _, tapping := range tappings {
						sendLadleMovements(nBf, tIds, tapping, mixIds, mixCookies)
					}
					//}(nBf, idJournal, bfjCookies)
				}
			}
			cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, nil, tIds.Tappings)
			fmt.Println(tm.String() + " 1min works")*/
		}
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

func initializeTimers() (morningShiftTick, eveningShiftTick, oneHourTick <-chan time.Time) {
	duration := time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)) // Calculate the duration until the next minute starts
	checktime := time.Time{}.Add(duration)
	fmt.Println(checktime)
	time.Sleep(duration) // Wait until the next minute starts

	oneHourDuration, _ := time.ParseDuration(config.GlobalConfig.Time.OneHourInterval)
	//oneMinDuration, _ := time.ParseDuration(config.GlobalConfig.Time.OneMinuteInterval)
	mTime, _ := time.Parse("15:04", config.GlobalConfig.Time.MorningShift)
	eTime, _ := time.Parse("15:04", config.GlobalConfig.Time.EveningShift)

	now := time.Now()
	var nextHour time.Time
	if now.Minute() == 0 {
		nextHour = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	} else {
		nextHour = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
	}

	nextMorning := time.Date(nextHour.Year(), nextHour.Month(), nextHour.Day(), mTime.Hour(), mTime.Minute(), 0, 0, time.Now().Location())
	nextEvening := time.Date(nextHour.Year(), nextHour.Month(), nextHour.Day(), eTime.Hour(), eTime.Minute(), 0, 0, time.Now().Location())

	if time.Now().After(nextEvening) {
		nextEvening = nextEvening.AddDate(0, 0, 1)
	}
	if time.Now().After(nextMorning) {
		nextMorning = nextMorning.AddDate(0, 0, 1)
	}

	morningShiftTick = time.NewTimer(time.Until(nextMorning)).C
	eveningShiftTick = time.NewTimer(time.Until(nextEvening)).C

	oneHourTick = time.NewTicker(oneHourDuration).C
	//minuteTick = time.NewTicker(oneMinDuration).C

	return morningShiftTick, eveningShiftTick, oneHourTick
}

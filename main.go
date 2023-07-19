package main

import (
	"fmt"
	"main/cache"
	"main/config"
	"main/controllers"
	"main/logger"
	"main/models"
	"net/http"
	"reflect"
	"time"
)

var (
	currList   map[int][]models.Ladle
	bfjCookies []*http.Cookie
	mixCookies []*http.Cookie
	bfjIds     map[int][]int
	mixIds     map[int][]int
	data       *cache.Data
	nBF        []int
	nBlock     []int
)

func main() {
	logger.InitLogger()
	initialize()

	fmt.Println(time.Now().String())

	for {
		service()
		waitForNextMinute()
	}
}

// initialize выполняет начальную настройку.
func initialize() {
	currList = make(map[int][]models.Ladle)
	bfjIds = map[int][]int{}
	mixIds = map[int][]int{}
	data = cache.ReadYAMLFile(config.GlobalConfig.Path.CachePath)

	controllers.AuthorizeBFJ(&bfjCookies)
	controllers.AuthorizeMix(&mixCookies)

	nBF = controllers.GetListBf()
	nBlock = []int{1, 2}

	controllers.GetLastBFJJournalsData(nBF, &bfjIds)
	controllers.GetLastBlockJournalsData(nBlock, &mixCookies, &mixIds)

	cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, bfjIds, data.Tappings)
}

// service является основной функцией, которая выполняет логику обслуживания.
func service() {
	now := time.Now().Truncate(time.Minute)
	if (now.Hour() == 8 && now.Minute() == 0) || (now.Hour() == 20 && now.Minute() == 0) {
		currList = make(map[int][]models.Ladle)
		fmt.Println(now.Format("2006-01-02 15:04:05"), "Shift started")

		controllers.GetLastBFJJournalsData(nBF, &bfjIds)
		controllers.GetLastBlockJournalsData(nBlock, &mixCookies, &mixIds)

		cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, bfjIds, nil)

		minuteCheck()
	} else if now.Second() == 0 {
		fmt.Println(now.Format("2006-01-02 15:04:05"), "Minute check started")

		minuteCheck()
	}
}

// minuteCheck выполняет проверку каждую минуту.
func minuteCheck() {
	data = cache.ReadYAMLFile(config.GlobalConfig.Path.CachePath)
	for nBf, values := range bfjIds {
		for _, idJournal := range values {
			tappings := controllers.GetBFJTappings(idJournal, &bfjCookies)
			for _, tapping := range tappings {
				sendLadleMovements(nBf, data, tapping, &mixIds, &mixCookies)
			}
		}
	}
	cache.WriteYAMLFile(config.GlobalConfig.Path.CachePath, nil, data.Tappings)
}

// sendLadleMovements отправляет движения ковшей.
func sendLadleMovements(nBf int, tIds *cache.Data, tapping models.Tapping, mixIds *map[int][]int, mixCookies *[]*http.Cookie) {
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

	handleChemicalChanges(tapping, tapping.ListLaldes, mixCookies)
}

// handleChemicalChanges обрабатывает изменения в составе химикатов.
func handleChemicalChanges(tapping models.Tapping, currLadles []models.Ladle, mixCookies *[]*http.Cookie) {
	// Обработка изменений в составе химикатов
	if currList[tapping.ID] == nil {
		currList[tapping.ID] = tapping.ListLaldes
	} else if !reflect.DeepEqual(currList[tapping.ID], tapping.ListLaldes) {
		oldLadles := currList[tapping.ID]
		currLadles := tapping.ListLaldes
		changedLadles := []models.Ladle{}

		for _, newLadle := range currLadles {
			for _, currLadle := range oldLadles {
				currLadle.Chemical.DtUpdate = newLadle.Chemical.DtUpdate
				if newLadle.Ladle == currLadle.Ladle && !reflect.DeepEqual(newLadle.Chemical, currLadle.Chemical) {
					changedLadles = append(changedLadles, newLadle)
					logger.Logger.Println(time.Now().Truncate(time.Minute).String(), newLadle.Ladle, " changed!")
				}
			}
		}

		if len(changedLadles) != 0 {
			controllers.PostMixChemicalList(changedLadles, mixCookies)
		}

		currList[tapping.ID] = tapping.ListLaldes
	}
}

// waitForNextMinute ждет до следующей минуты.
func waitForNextMinute() {
	duration := time.Until(time.Now().Truncate(time.Minute).Add(time.Minute))
	time.Sleep(duration)
}

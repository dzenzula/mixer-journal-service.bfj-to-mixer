package main

import (
	"fmt"
	"log"
	"main/controllers"
	"net/http"
	"time"
)

func main() {
	cookiesBFJ, authBFJError := controllers.AuthorizeProd()
	if authBFJError != nil {
		log.Println("Authorization error. \n", authBFJError)
	}

	// cookiesMIX, authMIXError := controllers.AuthorizationMixer()
	// if authMIXError != nil {
	// 	log.Println("Authorization error. \n", authMIXError)
	// }
	// fmt.Println(cookiesMIX)

	var nBF []int = controllers.GetListBf()

	duration := time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)) // Calculate the duration until the next minute starts
	checktime := time.Time{}.Add(duration)
	fmt.Println(checktime)
	time.Sleep(duration) // Wait until the next minute starts

	fiveMinuteTick := time.NewTicker(time.Minute * 5).C
	minuteTick := time.NewTicker(time.Minute).C

	ids := controllers.GetLastJournalsData(nBF)

	for {
		select {
		case tm := <-fiveMinuteTick:
			newIds := controllers.GetLastJournalsData(nBF)
			if len(newIds) != 0 {
				ids = newIds
			}

			for key, values := range ids {
				for _, idJournal := range values {
					go func(key int, idJournal int, cookiesBFJ []*http.Cookie) {
						controllers.GetJournalData(key, idJournal, cookiesBFJ)
						controllers.GetChemCoxes(idJournal, cookiesBFJ)
						controllers.GetChemMaterials(idJournal, cookiesBFJ)
						controllers.GetChemicalSlags(idJournal, cookiesBFJ)
						controllers.GetTappings(idJournal, cookiesBFJ)
					}(key, idJournal, cookiesBFJ)
				}
			}
			//controllers.GetJournalDatas(ids)
			fmt.Println(tm)

		case tm := <-minuteTick:
			newIds := controllers.GetLastJournalsData(nBF)
			if len(newIds) != 0 {
				ids = newIds
			}

			for key, values := range ids {
				for _, idJournal := range values {
					//go func(key int, idJournal int, cookiesBFJ []*http.Cookie) {
					controllers.GetJournalData(key, idJournal, cookiesBFJ)
					controllers.GetChemCoxes(idJournal, cookiesBFJ)
					controllers.GetChemMaterials(idJournal, cookiesBFJ)
					controllers.GetChemicalSlags(idJournal, cookiesBFJ)
					controllers.GetTappings(idJournal, cookiesBFJ)
					//}(key, idJournal, cookiesBFJ)
				}
			}

			fmt.Println(tm)
		}

	}
}

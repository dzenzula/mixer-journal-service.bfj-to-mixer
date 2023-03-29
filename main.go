package main

import (
	"fmt"
	"log"
	"main/controllers"
	"net/http"
	"time"
)

func main() {
	cookies, authError := controllers.AuthorizeProd()
	if authError != nil {
		log.Println("Authorization error. \n", authError)
	}

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
				for _, id := range values {
					go func(key int, id int, cookies []*http.Cookie) {
						idJournal := controllers.GetJournalData(key, id, cookies)
						controllers.GetChemCoxes(idJournal, cookies)
						controllers.GetChemMaterials(idJournal, cookies)
						controllers.GetChemicalSlags(idJournal, cookies)
					}(key, id, cookies)
				}
			}
			//controllers.GetJournalDatas(ids)
			fmt.Println(tm)

		case tm := <-minuteTick:
			for key, values := range ids {
				for _, id := range values {
					go func(key int, id int, cookies []*http.Cookie) {
						idJournal := controllers.GetJournalData(key, id, cookies)
						controllers.GetChemCoxes(idJournal, cookies)
						controllers.GetChemMaterials(idJournal, cookies)
						controllers.GetChemicalSlags(idJournal, cookies)
						controllers.GetTappings(idJournal, cookies)
					}(key, id, cookies)
				}
			}

			fmt.Println(tm)
		}

	}
}

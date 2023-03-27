package main

import (
	"fmt"
	"log"
	"main/controllers"
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

	minuteTick := time.NewTicker(time.Minute).C

	for {
		select {
		case tm := <-minuteTick:
			ids := controllers.GetLastJournalsData(nBF)

			for key, values := range ids {
				for _, id := range values {
					idJournal := controllers.GetJournalData(key, id, cookies)
					controllers.GetChemCoxes(idJournal, cookies)
					controllers.GetChemMaterials(idJournal, cookies)
					controllers.GetChemicalSlags(idJournal, cookies)
				}
			}

			//controllers.GetJournalDatas(ids)
			fmt.Println(tm)
		}
	}
}

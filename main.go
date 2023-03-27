package main

import (
	"fmt"
	"main/controllers"
	"time"
)

func main() {
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
					idJournal := controllers.GetJournalData(key, id)
					controllers.GetChemCoxes(idJournal)
					controllers.GetChemMaterials(idJournal)
					controllers.GetChemicalSlags(idJournal)
				}
			}

			//controllers.GetJournalDatas(ids)
			fmt.Println(tm)
		}
	}
}

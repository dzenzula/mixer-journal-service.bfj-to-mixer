package main

import (
	"fmt"

	//"main/cache"
	"main/controllers"
	"time"
)

func main() {
	var nBF []int = controllers.GetListBf()
	// cache.WriteYAMLFile(configPath.CachePath, yaml, 4)

	duration := time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)) // Calculate the duration until the next minute starts
	checktime := time.Time{}.Add(duration)
	fmt.Println(checktime)
	time.Sleep(duration) // Wait until the next minute starts

	minuteTick := time.NewTicker(time.Minute).C

	for {
		select {
		case tm := <-minuteTick:
			ids := controllers.GetLastJournalsData(nBF)
			controllers.GetJournalData(ids)
			fmt.Println(tm)
		}
	}
}

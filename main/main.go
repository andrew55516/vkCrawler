package main

import (
	"log"
	"vkCrawler/utils"
)

func main() {
	//scraped := db.GetLastPost().ID
	//for scraped < 9000 {
	//	var wg sync.WaitGroup
	//	cromedpWorker.Crawl(&wg)
	//	wg.Wait()
	//	scraped = db.GetLastPost().ID
	//	if scraped < 9000 {
	//		time.Sleep(10 * time.Minute)
	//	}
	//}
	//
	//postScraper.ScrapeAllPosts()
	//cromedpWorker.ScrapeUnknownPosts()
	//
	//err := utils.FillAllNodesEdges()
	//if err != nil {
	//	log.Fatal(err)
	//}

	err := utils.LikesOnlyUsersToXLSX()
	if err != nil {
		log.Fatal(err)
	}
}

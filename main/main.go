package main

import (
	"log"
	"sync"
	"time"
	db "vkCrawler/db/sqlc"
	"vkCrawler/internal/cromedpWorker"
	"vkCrawler/internal/postScraper"
)

func main() {
	scraped := db.GetLastPost().ID
	for scraped < 9000 {
		var wg sync.WaitGroup
		cromedpWorker.Crawl(&wg)
		wg.Wait()
		scraped = db.GetLastPost().ID
		if scraped < 9000 {
			time.Sleep(10 * time.Minute)
		}
	}

	postScraper.ScrapeAllPosts()
	cromedpWorker.ScrapeUnknownPosts()

	err := db.FillAllNodes()
	if err != nil {
		log.Fatal(err)
	}

	err = db.FillAllEdges()
	if err != nil {
		log.Fatal(err)
	}

	err = db.FillLikesNodes()
	if err != nil {
		log.Fatal(err)
	}

	err = db.FillLikesEdges()
	if err != nil {
		log.Fatal(err)
	}

	err = db.FillCommentsNodes()
	if err != nil {
		log.Fatal(err)
	}

	err = db.FillCommentsEdges()
	if err != nil {
		log.Fatal(err)
	}
}

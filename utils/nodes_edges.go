package utils

import (
	db "vkCrawler/db/sqlc"
)

func FillAllNodesEdges() (err error) {
	err = db.FillAllNodes()
	if err != nil {
		return
	}

	err = db.FillAllEdges()
	if err != nil {
		return
	}

	err = db.FillLikesNodes()
	if err != nil {
		return
	}

	err = db.FillLikesEdges()
	if err != nil {
		return
	}

	err = db.FillCommentsNodes()
	if err != nil {
		return
	}

	err = db.FillCommentsEdges()
	if err != nil {
		return
	}

	err = db.FillWeightedAllEdges()
	if err != nil {
		return
	}

	err = db.FillWeightedLikesEdges()
	if err != nil {
		return
	}

	err = db.FillWeightedCommentsEdges()
	if err != nil {
		return
	}

	return
}

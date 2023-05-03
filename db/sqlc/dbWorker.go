package db

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:pwd123@localhost:5432/crawler_db?sslmode=disable"
)

var Q *Queries

func init() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}

	Q = New(conn)
}

func WriteDownPost(link string, owner string, likes, comments, reposts int64, created_at time.Time) error {
	if isExists := checkIfPostExistsByLink(link); isExists {
		return nil
	}
	params := WriteDownPostParams{
		Link:      link,
		Owner:     owner,
		Likes:     likes,
		Comments:  comments,
		Reposts:   reposts,
		CreatedAt: created_at,
	}

	_, err := Q.WriteDownPost(context.Background(), params)

	return err
}

func GetLastPost() Post {
	p, _ := Q.GetLastPost(context.Background())
	return p
}

func GetAllPosts() []Post {
	posts, err := Q.GetAllPosts(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return posts
}

func UpdatePost(owner string, likes int, comments int, reposts int, id int) {
	args := UpdatePostParams{
		Owner:    owner,
		Likes:    int64(likes),
		Comments: int64(comments),
		Reposts:  int64(reposts),
		ID:       int64(id),
	}

	_, err := Q.UpdatePost(context.Background(), args)
	if err != nil {
		log.Fatal(err)
	}
}

func WriteDownLike(postID int, owner string) {
	args := WriteDownLikeParams{
		PostID: int64(postID),
		Owner:  owner,
	}

	if checkIfLikeExists(args.PostID, args.Owner) {
		return
	}

	_, err := Q.WriteDownLike(context.Background(), args)
	if err != nil {
		log.Fatal(err)
	}
}

func WriteDownComment(postID int, owner string, threadOwner string, created time.Time) {
	args := WriteDownCommentParams{
		PostID:      int64(postID),
		Owner:       owner,
		ThreadOwner: threadOwner,
		CreatedAt:   created,
	}

	//if checkIfCommentExists(args.PostID, args.Owner) {
	//	return
	//}

	_, err := Q.WriteDownComment(context.Background(), args)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUnknownPosts() []Post {
	posts, err := Q.GetUnknownPosts(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return posts
}

func checkIfPostExistsByLink(link string) bool {
	post, _ := Q.GetPostByLink(context.Background(), link)
	emptyPost := Post{}
	if post != emptyPost {
		return true
	}
	return false
}

func checkIfLikeExists(postID int64, owner string) bool {
	args := GetLikeParams{
		PostID: postID,
		Owner:  owner,
	}
	like, _ := Q.GetLike(context.Background(), args)
	emptyLike := Like{}
	if like != emptyLike {
		return true
	}
	return false
}

func checkIfCommentExists(postID int64, owner string) bool {
	args := GetCommentParams{
		PostID: postID,
		Owner:  owner,
	}
	like, _ := Q.GetComment(context.Background(), args)
	emptyComment := Comment{}
	if like != emptyComment {
		return true
	}
	return false
}

func FillAllNodes() error {
	return Q.FillAllNodes(context.Background())

}

func FillAllEdges() error {
	return Q.FillAllEdges(context.Background())
}

func FillLikesNodes() error {
	return Q.FillLikesNodes(context.Background())
}

func FillLikesEdges() error {
	return Q.FillLikesEdges(context.Background())
}

func FillCommentsNodes() error {
	return Q.FillCommentsNodes(context.Background())
}

func FillCommentsEdges() error {
	return Q.FillCommentsEdges(context.Background())
}

func FillWeightedAllEdges() error {
	return Q.FillWeightedAllEdges(context.Background())
}

func FillWeightedLikesEdges() error {
	return Q.FillWeightedLikesEdges(context.Background())
}

func FillWeightedCommentsEdges() error {
	return Q.FillWeightedCommentsEdges(context.Background())
}

func GetLikesOnlyUsers(toTime time.Time) ([]string, error) {
	return Q.GetLikesOnlyUsers(context.Background(), toTime)
}

func GetAmountOfUsers(toTime time.Time) (int, error) {
	amount, err := Q.GetAmountOfUsers(context.Background(), toTime)
	return int(amount), err
}

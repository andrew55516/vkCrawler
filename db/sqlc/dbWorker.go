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
	err := Q.FillAllNodes(context.Background())
	return err
}

func FillAllEdges() error {
	err := Q.FillAllEdges(context.Background())
	return err
}

func FillLikesNodes() error {
	err := Q.FillLikesNodes(context.Background())
	return err
}

func FillLikesEdges() error {
	err := Q.FillLikesEdges(context.Background())
	return err
}

func FillCommentsNodes() error {
	err := Q.FillCommentsNodes(context.Background())
	return err
}

func FillCommentsEdges() error {
	err := Q.FillCommentsEdges(context.Background())
	return err
}

func FillWeightedAllEdges() error {
	err := Q.FillWeightedAllEdges(context.Background())
	return err
}

func FillWeightedLikesEdges() error {
	err := Q.FillWeightedLikesEdges(context.Background())
	return err
}

func FillWeightedCommentsEdges() error {
	err := Q.FillWeightedCommentsEdges(context.Background())
	return err
}

func GetLikesOnlyUsers(to_time time.Time) ([]string, error) {
	return Q.GetLikesOnlyUsers(context.Background(), to_time)
}

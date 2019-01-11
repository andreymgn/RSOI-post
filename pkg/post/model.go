package post

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	errNotFound   = errors.New("post not found")
	errNotCreated = errors.New("post not created")
)

// Post describes a post
type Post struct {
	UID        uuid.UUID
	UserUID    uuid.UUID
	Title      string
	URL        string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type datastore interface {
	getAll(int32, int32) ([]*Post, error)
	getOne(uuid.UUID) (*Post, error)
	create(string, string, uuid.UUID) (*Post, error)
	update(uuid.UUID, string, string) error
	delete(uuid.UUID) error
	checkExists(uuid.UUID) (bool, error)
	getOwner(uuid.UUID) (string, error)
}

type db struct {
	*sql.DB
}

func newDB(connString string) (*db, error) {
	postgres, err := sql.Open("postgres", connString)
	return &db{postgres}, err
}

func (db *db) getAll(pageSize, pageNumber int32) ([]*Post, error) {
	query := "SELECT * FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2"
	lastRecord := pageNumber * pageSize
	rows, err := db.Query(query, pageSize, lastRecord)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	result := make([]*Post, 0)
	for rows.Next() {
		post := new(Post)
		var uid, userUID string
		err := rows.Scan(&uid, &userUID, &post.Title, &post.URL, &post.CreatedAt, &post.ModifiedAt)
		if err != nil {
			return nil, err
		}

		post.UID, err = uuid.Parse(uid)
		if err != nil {
			return nil, err
		}

		post.UserUID, err = uuid.Parse(userUID)

		result = append(result, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *db) getOne(uid uuid.UUID) (*Post, error) {
	query := "SELECT * FROM posts WHERE uid=$1"
	row := db.QueryRow(query, uid.String())
	result := new(Post)
	var stringUID, stringUserUID string
	switch err := row.Scan(&stringUID, &stringUserUID, &result.Title, &result.URL, &result.CreatedAt, &result.ModifiedAt); err {
	case nil:
		result.UID = uid
		userUID, err := uuid.Parse(stringUserUID)
		if err != nil {
			return nil, err
		}

		result.UserUID = userUID
		return result, nil
	case sql.ErrNoRows:
		return nil, errNotFound
	default:
		return nil, err
	}
}

func (db *db) create(title, url string, userUID uuid.UUID) (*Post, error) {
	post := new(Post)

	query := "INSERT INTO posts (uid, user_uid, title, url, created_at, modified_at) VALUES ($1, $2, $3, $4, $5, $6)"
	uid := uuid.New()

	now := time.Now()

	post.UID = uid
	post.UserUID = userUID
	post.Title = title
	post.URL = url
	post.CreatedAt = now
	post.ModifiedAt = now

	result, err := db.Exec(query, post.UID.String(), userUID.String(), post.Title, post.URL, post.CreatedAt, post.ModifiedAt)
	nRows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if nRows == 0 {
		return nil, errNotCreated
	}

	return post, nil
}

func (db *db) update(uid uuid.UUID, title, url string) error {
	query := "UPDATE posts SET title=COALESCE(NULLIF($1,''), title), url=COALESCE(NULLIF($2,''), url), modified_at=$3 WHERE uid=$4"
	result, err := db.Exec(query, title, url, time.Now(), uid.String())
	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if nRows == 0 {
		return errNotFound
	}

	return nil
}

func (db *db) delete(uid uuid.UUID) error {
	query := "DELETE FROM posts WHERE uid=$1"
	result, err := db.Exec(query, uid.String())
	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if nRows == 0 {
		return errNotFound
	}

	return nil
}

func (db *db) checkExists(uid uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM posts WHERE uid=$1)"
	row := db.QueryRow(query, uid.String())
	var result bool
	switch err := row.Scan(&result); err {
	case nil:
		return result, nil
	case sql.ErrNoRows:
		return false, errNotFound
	default:
		return false, err
	}
}

func (db *db) getOwner(uid uuid.UUID) (string, error) {
	query := "SELECT user_uid FROM posts WHERE uid=$1"
	row := db.QueryRow(query, uid.String())
	var result string
	switch err := row.Scan(&result); err {
	case nil:
		return result, nil
	case sql.ErrNoRows:
		return "", errNotFound
	default:
		return "", err
	}
}

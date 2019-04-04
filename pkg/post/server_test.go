package post

import (
	"errors"
	"testing"
	"time"

	pb "github.com/andreymgn/RSOI-post/pkg/post/proto"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

var (
	errDummy     = errors.New("dummy")
	dummyUID     = uuid.New()
	nilUIDString = uuid.Nil.String()
)

type mockdb struct{}

func (mdb *mockdb) getAll(pageSize, pageNumber int32) ([]*Post, error) {
	result := make([]*Post, 0)
	uid1 := uuid.New()
	uid2 := uuid.New()
	uid3 := uuid.New()

	result = append(result, &Post{uid1, uid2, "First post", "google.com", time.Now(), time.Now()})
	result = append(result, &Post{uid2, uid3, "Second post", "", time.Now(), time.Now().Add(time.Second * 10)})
	result = append(result, &Post{uid3, uid1, "Third post", "yandex.ru", time.Now(), time.Now()})
	return result, nil
}

func (mdb *mockdb) getOne(uid uuid.UUID) (*Post, error) {
	if uid == uuid.Nil {
		uid := uuid.New()

		return &Post{uid, uid, "First post", "google.com", time.Now(), time.Now()}, nil
	}

	return nil, errDummy
}

func (mdb *mockdb) create(title, url string, userUID uuid.UUID) (*Post, error) {
	if title == "success" {
		uid := uuid.New()

		return &Post{uid, userUID, "First post", "google.com", time.Now(), time.Now()}, nil
	}

	return nil, errDummy
}

func (mdb *mockdb) update(uid uuid.UUID, title, url string) error {
	if uid == uuid.Nil {
		return nil
	}

	return errDummy
}

func (mdb *mockdb) delete(uid uuid.UUID) error {
	if uid == uuid.Nil {
		return nil
	}

	return errDummy
}

func (mdb *mockdb) checkExists(uid uuid.UUID) (bool, error) {
	if uid == uuid.Nil {
		return true, nil
	}

	return false, errDummy
}

func (mdb *mockdb) checkToken(token string) (bool, error) {
	return true, nil
}

func (mdb *mockdb) getOwner(uid uuid.UUID) (string, error) {
	return nilUIDString, nil
}

func TestListPosts(t *testing.T) {
	s := &Server{&mockdb{}}
	var pageSize int32 = 3
	req := &pb.ListPostsRequest{PageSize: pageSize, PageNumber: 1}
	res, err := s.ListPosts(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(res.Posts) != int(pageSize) {
		t.Errorf("unexpected number of posts: got %v want %v", len(res.Posts), pageSize)
	}
}

func TestGetPost(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.GetPostRequest{Uid: nilUIDString}
	_, err := s.GetPost(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestGetPostFail(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.GetPostRequest{Uid: ""}
	_, err := s.GetPost(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestCreatePost(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.CreatePostRequest{Title: "success", UserUid: nilUIDString}
	_, err := s.CreatePost(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCreatePostFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.CreatePostRequest{Title: ""}
	_, err := s.CreatePost(context.Background(), req)
	if err != statusNoPostTitle {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.CreatePostRequest{Title: "fail"}
	_, err = s.CreatePost(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestUpdatePost(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.UpdatePostRequest{Uid: nilUIDString}
	_, err := s.UpdatePost(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestUpdatePostFail(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.UpdatePostRequest{Uid: ""}
	_, err := s.UpdatePost(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestDeletePost(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.DeletePostRequest{Uid: nilUIDString}
	_, err := s.DeletePost(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestDeletePostFail(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.DeletePostRequest{Uid: ""}
	_, err := s.DeletePost(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestCheckExists(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.CheckExistsRequest{Uid: nilUIDString}
	_, err := s.CheckExists(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCheckExistsFail(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.CheckExistsRequest{Uid: ""}
	_, err := s.CheckExists(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

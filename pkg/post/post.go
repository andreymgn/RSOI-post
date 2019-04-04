package post

import (
	pb "github.com/andreymgn/RSOI-post/pkg/post/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	statusNoPostTitle  = status.Error(codes.InvalidArgument, "post title is required")
	statusNotFound     = status.Error(codes.NotFound, "post not found")
	statusInvalidUUID  = status.Error(codes.InvalidArgument, "invalid UUID")
	statusInvalidToken = status.Errorf(codes.Unauthenticated, "invalid token")
)

func internalError(err error) error {
	return status.Error(codes.Internal, err.Error())
}

// SinglePost converts Post to SinglePost
func (p *Post) SinglePost() (*pb.SinglePost, error) {
	createdAtProto, err := ptypes.TimestampProto(p.CreatedAt)
	if err != nil {
		return nil, internalError(err)
	}

	modifiedAtProto, err := ptypes.TimestampProto(p.ModifiedAt)
	if err != nil {
		return nil, internalError(err)
	}

	res := new(pb.SinglePost)
	res.Uid = p.UID.String()
	res.UserUid = p.UserUID.String()
	res.Title = p.Title
	res.Url = p.URL
	res.CreatedAt = createdAtProto
	res.ModifiedAt = modifiedAtProto

	return res, nil
}

// ListPosts returns newest posts
func (s *Server) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	posts, err := s.db.getAll(pageSize, req.PageNumber)
	if err != nil {
		return nil, internalError(err)
	}
	res := new(pb.ListPostsResponse)
	for _, post := range posts {
		postResponse, err := post.SinglePost()
		if err != nil {
			return nil, err
		}

		res.Posts = append(res.Posts, postResponse)
	}

	res.PageSize = pageSize
	res.PageNumber = req.PageNumber

	return res, nil
}

// GetPost returns single post by ID
func (s *Server) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.SinglePost, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	post, err := s.db.getOne(uid)
	switch err {
	case nil:
		return post.SinglePost()
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// CreatePost creates a new post
func (s *Server) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.SinglePost, error) {
	if req.Title == "" {
		return nil, statusNoPostTitle
	}

	userUID, err := uuid.Parse(req.UserUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	post, err := s.db.create(req.Title, req.Url, userUID)
	if err != nil {
		return nil, internalError(err)
	}

	return post.SinglePost()
}

// UpdatePost updates post by ID
func (s *Server) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.db.update(uid, req.Title, req.Url)
	switch err {
	case nil:
		return new(pb.UpdatePostResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// DeletePost deletes post by ID
func (s *Server) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.db.delete(uid)
	switch err {
	case nil:
		return new(pb.DeletePostResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// CheckExists checks if post with ID exists in DB
func (s *Server) CheckExists(ctx context.Context, req *pb.CheckExistsRequest) (*pb.CheckExistsResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	result, err := s.db.checkExists(uid)
	switch err {
	case nil:
		res := new(pb.CheckExistsResponse)
		res.Exists = result
		return res, nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// GetOwner returns post owner
func (s *Server) GetOwner(ctx context.Context, req *pb.GetOwnerRequest) (*pb.GetOwnerResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	result, err := s.db.getOwner(uid)
	switch err {
	case nil:
		res := new(pb.GetOwnerResponse)
		res.OwnerUid = result
		return res, nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

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
	statusNoPostTitle    = status.Error(codes.InvalidArgument, "post title is required")
	statusNotFound       = status.Error(codes.NotFound, "post not found")
	statusInvalidUUID    = status.Error(codes.InvalidArgument, "invalid UUID")
	statusInvalidToken   = status.Errorf(codes.Unauthenticated, "invalid token")
	statusNoCategoryName = status.Error(codes.InvalidArgument, "category name is required")
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
	res.CategoryUid = p.CategoryUID.String()
	res.Title = p.Title
	res.Url = p.URL
	res.CreatedAt = createdAtProto
	res.ModifiedAt = modifiedAtProto

	return res, nil
}

// SingleCategory converts Category to SingleCategory
func (c *Category) SingleCategory() *pb.SingleCategory {
	res := new(pb.SingleCategory)
	res.Uid = c.UID.String()
	res.UserUid = c.UserUID.String()
	res.Name = c.Name

	return res
}

// ListPosts returns newest posts
func (s *Server) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	uid, err := uuid.Parse(req.CategoryUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	posts, err := s.db.getAllPosts(uid, pageSize, req.PageNumber)
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

	post, err := s.db.getOnePost(uid)
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

	categoryUID, err := uuid.Parse(req.CategoryUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	post, err := s.db.createPost(req.Title, req.Url, userUID, categoryUID)
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

	err = s.db.updatePost(uid, req.Title, req.Url)
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

	err = s.db.deletePost(uid)
	switch err {
	case nil:
		return new(pb.DeletePostResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// CheckPostExists checks if post with ID exists in DB
func (s *Server) CheckPostExists(ctx context.Context, req *pb.CheckPostExistsRequest) (*pb.CheckPostExistsResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	result, err := s.db.checkPostExists(uid)
	switch err {
	case nil:
		res := new(pb.CheckPostExistsResponse)
		res.Exists = result
		return res, nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// GetPostOwner returns post owner
func (s *Server) GetPostOwner(ctx context.Context, req *pb.GetPostOwnerRequest) (*pb.GetPostOwnerResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	result, err := s.db.getPostOwner(uid)
	switch err {
	case nil:
		res := new(pb.GetPostOwnerResponse)
		res.OwnerUid = result
		return res, nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// ListCategories returns categories
func (s *Server) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	categories, err := s.db.getAllCategories(pageSize, req.PageNumber)
	if err != nil {
		return nil, internalError(err)
	}

	res := new(pb.ListCategoriesResponse)
	for _, category := range categories {
		categoryResponse := category.SingleCategory()

		res.Categories = append(res.Categories, categoryResponse)
	}

	res.PageSize = pageSize
	res.PageNumber = req.PageNumber

	return res, nil
}

// CreateCategory creates a new post category
func (s *Server) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.SingleCategory, error) {
	if req.Name == "" {
		return nil, statusNoCategoryName
	}

	userUID, err := uuid.Parse(req.UserUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	category, err := s.db.createCategory(req.Name, userUID)
	if err != nil {
		return nil, internalError(err)
	}

	return category.SingleCategory(), nil
}

package post

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/post"
	"forum/internal/service/comment"
	"forum/internal/service/image"
	"forum/internal/service/tag"
	"forum/internal/service/user"
	"strconv"
)

type IPostService interface {
	SavePost(entity.PostCreateForm) (int, error)
	GetPost(int) (entity.PostView, error)
	GetAllPosts() (*[]entity.PostView, error)
	GetAllPostsByTagId(int) (*[]entity.PostView, error)
	GetAllPostsByUserId(int) (*[]entity.PostView, error)
	GetAllPostsByUserReaction(int) (*[]entity.PostView, error)
	GetAllCommentedPostsWithComments(userID int) (*[]entity.PostView, *[][]entity.CommentView, error)
	ExistsPost(postID int) (bool, error)
	CheckPostAttrs(*entity.PostCreateForm, bool) (bool, error)
	DeletePost(postID int, userID int) error
	DeletePostPrivileged(postID int, userID int, userRole string) error
	GetAuthorID(postID int) (int, error)
}

type postService struct {
	imgService     image.IImageService
	tagService     tag.ITagService
	commentService comment.ICommentService
	userService    user.IUserService
	postRepo       post.IPostRepository
}

// Constructor for post service
func NewPostsService(r post.IPostRepository, is image.IImageService, ts tag.ITagService, cs comment.ICommentService, us user.IUserService) *postService {
	return &postService{
		imgService:     is,
		tagService:     ts,
		commentService: cs,
		userService:    us,
		postRepo:       r,
	}
}

var _ IPostService = (*postService)(nil)

func (ps *postService) SavePost(p entity.PostCreateForm) (int, error) {
	var tagIDs []int
	for _, tagIDStr := range p.Tags {
		tagID, _ := strconv.Atoi(tagIDStr) // Don't handle error because we know Ids are valid (checked before)
		tagIDs = append(tagIDs, tagID)
	}

	id, err := ps.postRepo.Insert(p, tagIDs)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (ps *postService) GetPost(postId int) (entity.PostView, error) {
	post, err := ps.postRepo.Get(postId)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			return entity.PostView{}, entity.ErrInvalidPostID
		}
		return entity.PostView{}, err
	}

	imgName, err := ps.imgService.Get(postId)
	if err != nil {
		return entity.PostView{}, err
	}

	tags := ConvertToStrArr(post.PostTags)
	pView := entity.PostView{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		CreatedAt:   post.CreatedAt,
		Username:    post.Username,
		Likes:       post.Likes,
		Dislikes:    post.Dislikes,
		CommentsLen: post.CommentsLen,
		PostTags:    tags,
		ImageName:   imgName,
	}

	return pView, nil
}

func (ps *postService) GetAllPosts() (*[]entity.PostView, error) {
	posts, err := ps.postRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return ConvertEntitiesToViews(posts)
}

func (ps *postService) GetAllPostsByTagId(tagID int) (*[]entity.PostView, error) {
	posts, err := ps.postRepo.GetAllByTagId(tagID)
	if err != nil {
		return nil, err
	}

	return ConvertEntitiesToViews(posts)
}

func (ps *postService) GetAllPostsByUserId(userID int) (*[]entity.PostView, error) {
	posts, err := ps.postRepo.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	return ConvertEntitiesToViews(posts)
}

func (ps *postService) GetAllPostsByUserReaction(userID int) (*[]entity.PostView, error) {
	posts, err := ps.postRepo.GetAllByUserReaction(userID)
	if err != nil {
		return nil, err
	}

	return ConvertEntitiesToViews(posts)
}

func (ps *postService) GetAllCommentedPostsWithComments(userID int) (*[]entity.PostView, *[][]entity.CommentView, error) {
	postsEntities, err := ps.postRepo.GetAllCommentedPosts(userID)
	if err != nil {
		return nil, nil, err
	}
	posts, _ := ConvertEntitiesToViews(postsEntities)

	var allComments [][]entity.CommentView

	for _, p := range *posts {
		comments, err := ps.commentService.GetAllUserCommentsForPost(userID, p.ID)
		if err != nil {
			return nil, nil, err
		}

		allComments = append(allComments, *comments)
	}

	return posts, &allComments, nil
}

func (ps *postService) ExistsPost(postID int) (bool, error) {
	return ps.postRepo.Exists(postID)
}

func (ps *postService) CheckPostAttrs(p *entity.PostCreateForm, withImage bool) (bool, error) {
	if !IsRightPost(p, withImage) {
		return false, nil
	}

	areTagsExist, err := ps.tagService.AreTagsExist(p.Tags)
	if !areTagsExist || err != nil {
		return false, err
	}

	return true, nil
}

func (ps *postService) GetAuthorID(postID int) (int, error) {
	return ps.postRepo.GetAuthorID(postID)
}

func (ps *postService) DeletePost(postID int, userID int) error {
	exists, err := ps.postRepo.Exists(postID)
	if err != nil {
		return err
	}
	if !exists {
		return entity.ErrPostNotFound
	}

	return ps.postRepo.Delete(postID, userID)
}

func (ps *postService) DeletePostPrivileged(postID int, userID int, userRole string) error {
	exists, err := ps.postRepo.Exists(postID)
	if err != nil {
		return err
	}
	if !exists {
		return entity.ErrPostNotFound
	}

	authorID, err := ps.postRepo.GetAuthorID(postID)
	if err != nil {
		return err
	}

	err = ps.postRepo.DeleteByPrivileged(postID)
	if err != nil {
		return err
	}

	notificaiton := entity.Notification{
		Type:     entity.DELETE_POST,
		UserFrom: userID,
		UserTo:   authorID,
	}
	if userRole == entity.MODERATOR {
		notificaiton.Content = ". Reason: obscene"
	}

	err = ps.userService.SendNotification(notificaiton)
	if err != nil {
		return err
	}

	return nil
}

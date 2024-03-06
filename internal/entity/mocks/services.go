package mocks

import (
	"forum/internal/entity/mocks/comment"
	"forum/internal/entity/mocks/image"
	"forum/internal/entity/mocks/post"
	"forum/internal/entity/mocks/reaction"
	"forum/internal/entity/mocks/tag"
	"forum/internal/entity/mocks/user"
	"forum/internal/repository"
	"forum/internal/service"
)

func NewServicesMock(r *repository.Repositories) *service.Services {
	return &service.Services{
		Post:     post.NewPostServiceMock(r.Post, image.NewImageServiceMock(r.Image), tag.NewTagServiceMock(r.Tag)),
		User:     user.NewUserServiceMock(r.User),
		Tag:      tag.NewTagServiceMock(r.Tag),
		Comment:  comment.NewTagServiceMock(r.Comment),
		Reaction: reaction.NewReactionServiceMock(r.Reaction),
	}
}

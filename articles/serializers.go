package articles

import (
	"golang-gin-realworld-example-app/users"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type TagSerializer struct {
	C *gin.Context
	TagModel
}

type TagsSerializer struct {
	C    *gin.Context
	Tags []TagModel
}

func (s *TagSerializer) Response() string {
	return s.TagModel.Tag
}

func (s *TagsSerializer) Response() []string {
	response := []string{}
	for _, tag := range s.Tags {
		serializer := TagSerializer{s.C, tag}
		response = append(response, serializer.Response())
	}
	return response
}

type ArticleUserSerializer struct {
	C *gin.Context
	ArticleUserModel
}

func (s *ArticleUserSerializer) Response() users.ProfileResponse {
	response := users.ProfileSerializer{s.C, s.ArticleUserModel.UserModel}
	return response.Response()
}

type ArticleSerializer struct {
	C *gin.Context
	ArticleModel
}

type ArticleResponse struct {
	ID             uint                  `json:"-"`
	Title          string                `json:"title"`
	Slug           string                `json:"slug"`
	Description    string                `json:"description"`
	Body           string                `json:"body"`
	CreatedAt      string                `json:"createdAt"`
	UpdatedAt      string                `json:"updatedAt"`
	Author         users.ProfileResponse `json:"author"`
	Tags           []string              `json:"tagList"`
	Favorite       bool                  `json:"favorited"`
	FavoritesCount uint                  `json:"favoritesCount"`
}

type ArticlesSerializer struct {
	C        *gin.Context
	Articles []ArticleModel
}

func (s *ArticleSerializer) Response() ArticleResponse {
	myUserModel := s.C.MustGet("my_user_model").(users.UserModel)
	authorSerializer := ArticleUserSerializer{s.C, s.Author}
	response := ArticleResponse{
		ID:          s.ID,
		Slug:        slug.Make(s.Title),
		Title:       s.Title,
		Description: s.Description,
		Body:        s.Body,
		CreatedAt:   s.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		//UpdatedAt:      s.UpdatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:      s.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		Author:         authorSerializer.Response(),
		Favorite:       s.isFavoriteBy(GetArticleUserModel(myUserModel)),
		FavoritesCount: s.favoritesCount(),
	}
	response.Tags = make([]string, 0)
	for _, tag := range s.Tags {
		serializer := TagSerializer{s.C, tag}
		response.Tags = append(response.Tags, serializer.Response())
	}
	return response
}

func (s *ArticlesSerializer) Response() []ArticleResponse {
	response := []ArticleResponse{}
	for _, article := range s.Articles {
		serializer := ArticleSerializer{s.C, article}
		response = append(response, serializer.Response())
	}
	return response
}

type CommentSerializer struct {
	C *gin.Context
	CommentModel
}

type CommentsSerializer struct {
	C        *gin.Context
	Comments []CommentModel
}

type CommentResponse struct {
	ID        uint                  `json:"id"`
	Body      string                `json:"body"`
	CreatedAt string                `json:"createdAt"`
	UpdatedAt string                `json:"updatedAt"`
	Author    users.ProfileResponse `json:"author"`
}

func (s *CommentSerializer) Response() CommentResponse {
	authorSerializer := ArticleUserSerializer{s.C, s.Author}
	response := CommentResponse{
		ID:        s.ID,
		Body:      s.Body,
		CreatedAt: s.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt: s.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		Author:    authorSerializer.Response(),
	}
	return response
}

func (s *CommentsSerializer) Response() []CommentResponse {
	response := []CommentResponse{}
	for _, comment := range s.Comments {
		serializer := CommentSerializer{s.C, comment}
		response = append(response, serializer.Response())
	}
	return response
}

type VoteSerializer struct {
	vote CommentModelVote
}

type VoteResponse struct {
	UserID    uint `json:"user_id"`
	CommentID uint `json:"comment_id"`
	UpVote    bool `json:"up_vote"`
	DownVote  bool `json:"down_vote"`
}

func (self VoteSerializer) Response() VoteResponse {
	voteModel := self.vote
	voteResponse := VoteResponse{
		UserID:    voteModel.UserID,
		CommentID: voteModel.CommentID,
		UpVote:    voteModel.UpVote,
		DownVote:  voteModel.DownVote,
	}
	return voteResponse
}

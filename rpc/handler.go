package rpctransport

import (
	"context"

	"github.com/jgillard/practising-go-tdd/internal"
)

type Server struct {
	categoryStore internal.CategoryStore
	questionStore internal.QuestionStore
}

func NewServer(categoryStore internal.CategoryStore, questionStore internal.QuestionStore) *Server {
	p := new(Server)

	p.categoryStore = categoryStore
	p.questionStore = questionStore

	return p
}

func (s *Server) GetStatus(ctx context.Context, req *EmptyRequest) (*StatusReply, error) {
	status := internal.GetStatus()
	reply := &StatusReply{Status: status}
	return reply, nil
}

func (s *Server) ListCategories(ctx context.Context, req *EmptyRequest) (*ListCategoryReply, error) {
	categoryList := s.categoryStore.ListCategories()

	categoryReplies := []*GetCategoryReply{}

	for _, category := range categoryList.Categories {
		categoryReplies = append(categoryReplies,
			&GetCategoryReply{
				ID: category.ID, Name: category.Name, ParentID: category.ParentID,
			},
		)
	}

	reply := &ListCategoryReply{Categories: categoryReplies}

	return reply, nil
}

func (s *Server) GetCategory(ctx context.Context, req *GetCategoryRequest) (*GetCategoryReply, error) {
	category := s.categoryStore.GetCategory(req.CategoryID)
	reply := &GetCategoryReply{
		ID:       category.ID,
		Name:     category.Name,
		ParentID: category.ParentID,
	}
	return reply, nil
}

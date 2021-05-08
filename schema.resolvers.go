package main

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"lumic/graph/generated"
	"lumic/graph/model"
	"time"

	"github.com/latolukasz/orm"
	log "github.com/sirupsen/logrus"
)

func (r *mutationResolver) CreateImage(ctx context.Context, input model.NewImage) (*model.Image, error) {
	img := image{}
	_ = img.SetCategory(input.Category)
	err := img.Create(input.UserName, input.Caption, input.Link)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	graphImage := model.Image{
		ID:       fmt.Sprintf("%v", img.ID),
		Caption:  img.Caption,
		UserName: img.UserName,
		Created:  img.Created.Format(time.RFC3339),
	}
	return &graphImage, nil
}

func (r *queryResolver) Images(ctx context.Context, search *string) ([]*model.Image, error) {
	query := &orm.RedisSearchQuery{}

	if search != nil {
		query.Query(*search)
	}

	var entities []*image
	engine.RedisSearch(&entities, query, orm.NewPager(1, 50), "Category")
	var res []*model.Image
	for _, img := range entities {
		res = append(res, &model.Image{
			ID:       fmt.Sprintf("%v", img.ID),
			Caption:  img.Caption,
			UserName: img.UserName,
			Created:  img.Created.Format(time.RFC3339),
			Category: &model.ImageCategory{Name: img.Category.Name, ID: fmt.Sprintf("%v", img.Category.ID)},
		})
	}

	return res, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

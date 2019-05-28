package graph

import (
	"context"
)

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Message(ctx context.Context) (*Message, error) {
	return &Message{Body: "Spoon!"}, nil
}

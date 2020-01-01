package graph

import (
	"context"
	"math/rand"

	"github.com/sjansen/strongbox-elm/backend/api/auth"
)

var quotes = []string{
	"Despite all my rage I am still just a rat in a cage ",
	"I get knocked down, but I get up again.",
	"I'm the Scatman!",
	"If you want to destroy my sweater, hold this thread as I walk away.",
	"Kilroy was here.",
	"Make a little birdhouse in your soul.",
	"Movin' to the country, gonna eat a lot of peaches.",
	"Soy un perdedor.",
	"Spoon!",
	"The world is a vampire.",
}

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Message(ctx context.Context) (*Message, error) {
	quote := quotes[rand.Intn(len(quotes))]
	return &Message{Body: quote}, nil
}

func (r *queryResolver) Whoami(ctx context.Context) (*User, error) {
	var resp *User

	if user := auth.ForContext(ctx); user != nil {
		resp = &User{
			GivenName:  user.GivenName,
			FamilyName: user.FamilyName,
			Email:      user.Email,
		}
	}

	return resp, nil
}

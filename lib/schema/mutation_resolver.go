package schema

import (
	"context"

	"github.com/jasongwartz/carbon-offset-backend/lib/schema/generated"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Purchase(ctx context.Context) (*generated.MakePurchase, error) {
	return &generated.MakePurchase{}, nil
}

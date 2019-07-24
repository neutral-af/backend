package resolvers

import generated "github.com/jasongwartz/carbon-offset-backend/lib/graphql-generated"

type Resolver struct{}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Estimate() generated.EstimateResolver {
	return &estimateResolver{r}
}

func (r *Resolver) GetEstimate() generated.GetEstimateResolver {
	return &getEstimateResolver{r}
}

func (r *Resolver) MakePurchase() generated.MakePurchaseResolver {
	return &makePurchaseResolver{r}
}

func (r *Resolver) PaymentActions() generated.PaymentActionsResolver {
	return &paymentActionsResolver{r}
}

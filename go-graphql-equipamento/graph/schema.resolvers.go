package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"go-graphql-equipamento/graph/generated"
	"go-graphql-equipamento/graph/model"
)

func (r *mutationResolver) UpdateComputador(ctx context.Context, id string, input model.UpdateComputador) (*model.Computador, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateGpu(ctx context.Context, id string, input model.NovoGpu) (*model.Gpu, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateCPU(ctx context.Context, id string, input model.NovoCPU) (*model.CPU, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateRAM(ctx context.Context, id string, input model.NovoRAM) (*model.RAM, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CriarComputador(ctx context.Context, input model.NovoComputador) (*model.Computador, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CriarItem(ctx context.Context, input model.NovoItem) (*model.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CriarCPU(ctx context.Context, input model.NovoCPU) (*model.CPU, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CriarGpu(ctx context.Context, input model.NovoGpu) (*model.Gpu, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CriarRAM(ctx context.Context, input model.NovoRAM) (*model.RAM, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetComputadores(ctx context.Context) ([]*model.Computador, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetComputador(ctx context.Context, id string) (*model.Computador, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetCPU(ctx context.Context, id string) (*model.CPU, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetCPUS(ctx context.Context) ([]*model.CPU, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetGpu(ctx context.Context, id string) (*model.Gpu, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetGpus(ctx context.Context) ([]*model.Gpu, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetRAM(ctx context.Context, id string) (*model.RAM, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetRAMS(ctx context.Context) ([]*model.RAM, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetItems(ctx context.Context) ([]*model.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetItem(ctx context.Context, id string) (*model.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

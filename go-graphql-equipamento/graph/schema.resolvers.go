package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"go-graphql-equipamento/graph/generated"
	"go-graphql-equipamento/graph/model"
	"go-graphql-equipamento/redisconf"
	"log"
	"os"
)

var addressR = os.Getenv("REDISADDRESS")
var portR = os.Getenv("REDISPORT")
var passwordR = os.Getenv("REDISPASSWORD")

var redisClienteDB = redisconf.NovoClienteRedis(addressR, portR, passwordR)

// TODO UpdateComputador
func (r *mutationResolver) UpdateComputador(ctx context.Context, id string, input model.UpdateComputador) (*model.Computador, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO UpdateGPU
func (r *mutationResolver) UpdateGpu(ctx context.Context, id string, input model.NovoGpu) (*model.Gpu, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO UpdateCPU
func (r *mutationResolver) UpdateCPU(ctx context.Context, id string, input model.NovoCPU) (*model.CPU, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO UpdateRAM
func (r *mutationResolver) UpdateRAM(ctx context.Context, id string, input model.NovoRAM) (*model.RAM, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO CriarComputador
func (r *mutationResolver) CriarComputador(ctx context.Context, input model.NovoComputador) (*model.Computador, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO CriarItem
func (r *mutationResolver) CriarItem(ctx context.Context, input model.NovoItem) (*model.Item, error) {
	var novoItem model.Item

	novoItem.Marca = input.Marca
	novoItem.Modelo = &input.Modelo
	novoItem.Nome = *input.Nome
	novoItem.PaginaWeb = input.PaginaWeb

	item, err := json.Marshal(novoItem)
	if err != nil {
		return nil, err
	}

	err = redisClienteDB.Set(context.Background(), "id1", item, 0).Err()
	if err != nil {
		return nil, err
	}

	val, getErr := redisClienteDB.Get(ctx, "id1").Result()
	if getErr != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println(val)

	return &novoItem, nil
}

// TODO CriarCPU
func (r *mutationResolver) CriarCPU(ctx context.Context, input model.NovoCPU) (*model.CPU, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO CriarGpu
func (r *mutationResolver) CriarGpu(ctx context.Context, input model.NovoGpu) (*model.Gpu, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO CriarRAM
func (r *mutationResolver) CriarRAM(ctx context.Context, input model.NovoRAM) (*model.RAM, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO Implementar GetComputadores
func (r *queryResolver) GetComputadores(ctx context.Context) ([]*model.Computador, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO Implementar GetComputador
func (r *queryResolver) GetComputador(ctx context.Context, id string) (*model.Computador, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO Implementar GetCPU
func (r *queryResolver) GetCPU(ctx context.Context, id string) (*model.CPU, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO Implementar GetCPUS
func (r *queryResolver) GetCPUS(ctx context.Context) ([]*model.CPU, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO Implementar GetCPU
func (r *queryResolver) GetGpu(ctx context.Context, id string) (*model.Gpu, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO Implementar GetGpus
func (r *queryResolver) GetGpus(ctx context.Context) ([]*model.Gpu, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO Implementar GetRAM
func (r *queryResolver) GetRAM(ctx context.Context, id string) (*model.RAM, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO Implementar GetRAMS
func (r *queryResolver) GetRAMS(ctx context.Context) ([]*model.RAM, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO Implementar GetItems
func (r *queryResolver) GetItems(ctx context.Context) ([]*model.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

// TODO Implementar GetItem
func (r *queryResolver) GetItem(ctx context.Context, id string) (*model.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

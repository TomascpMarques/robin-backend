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

var resolverLogger = log.New(os.Stdout, "GraphQL-Resolver (*) ", log.LstdFlags)
var redisClienteDB = redisconf.NovoClienteRedis(redisconf.AddressRed, redisconf.PortRed, redisconf.PasswordRed)

func (r *mutationResolver) UpdateComputador(ctx context.Context, id string, input model.UpdateComputador) (*model.ComputadorAtualizado, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateGpu(ctx context.Context, id string, input model.NovoGpu) (*model.ComponenteAtualizado, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateCPU(ctx context.Context, id string, input model.NovoCPU) (*model.ComponenteAtualizado, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateRAM(ctx context.Context, id string, input model.NovoRAM) (*model.ComponenteAtualizado, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateItem(ctx context.Context, id string, input model.NovoItem) (*model.ItemAtualizado, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateSoftware(ctx context.Context, id string, input model.NovoSoftware) (*model.SoftwareAtualizado, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CriarSoftware(ctx context.Context, input model.NovoSoftware) (*model.SoftwareCriado, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CriarComputador(ctx context.Context, input model.NovoComputador) (*model.ComputadorCriado, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CriarItem(ctx context.Context, input model.NovoItem) (*model.ItemCriado, error) {
	var novoItem model.ItemCriado

	// Atribuição do input a estutura de dados correta
	novoItem.Nome = input.Nome

	// Criação de chave/index relacional aos existentes
	chave := redisconf.IDHandlerBD(&redisClienteDB)
	idChave := "Item" + chave
	novoItem.ID = idChave

	// Tradução do input da mutation para json através da indentação do mesmo
	item, err := json.MarshalIndent(input, "", "\t")
	if err != nil {
		return nil, err
	}

	// Inserçaõ do novo registo na base de dados
	err = redisClienteDB.Set(context.Background(), idChave, item, 0).Err()
	if err != nil {
		return nil, err
	}

	// Escreve no ecrã o registo insserido para verificação da insserção
	// e visualização do novo registo
	val, getErr := redisClienteDB.Get(ctx, idChave).Result()
	if getErr != nil {
		resolverLogger.Fatalf("Erro: %v", err)
		return nil, err
	}
	resolverLogger.Printf("[$] Valor insserido: %v", val)

	return &novoItem, nil
}

func (r *mutationResolver) CriarCPU(ctx context.Context, input model.NovoCPU) (*model.ComponenteCriado, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CriarGpu(ctx context.Context, input model.NovoGpu) (*model.ComponenteCriado, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CriarRAM(ctx context.Context, input model.NovoRAM) (*model.ComponenteCriado, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ApagarComponente(ctx context.Context, id string) (*model.ComponenteApagado, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ApagarItem(ctx context.Context, id string) (*model.RegistoApagado, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ApagarSoftware(ctx context.Context, id string) (*model.RegistoApagado, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ApagarComputador(ctx context.Context, id string) (*model.RegistoApagado, error) {
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

func (r *queryResolver) GetSoftwares(ctx context.Context) ([]*model.Software, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetSoftware(ctx context.Context, id string) (*model.Software, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

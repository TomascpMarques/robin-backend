package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"go-graphql-equipamento/graph/generated"
	"go-graphql-equipamento/graph/model"
	"go-graphql-equipamento/loggers"
	"go-graphql-equipamento/redishandle"
)

var resolverLogger = loggers.ResolverLogger

// RedisClienteDB - cliente que establece a conexão ao serviço redis
var RedisClienteDB = redishandle.NovoClienteRedis(redishandle.AddressRed, redishandle.PortRed, redishandle.PasswordRed)

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
	var registo redishandle.RegistoRedisDB
	var NovoSoftware model.SoftwareCriado

	// CriaEstruturaRegisto - Metodo do package redishandle em relação à estrutura RegistoRedisDB
	registo.CriaEstruturaRegisto(&RedisClienteDB, input, "Software")

	// Inserção do novo registo na base de dados
	// ? Todo o error cheking é feito na função
	redishandle.SetRegistoBD(&RedisClienteDB, registo)

	// Atribuição do input a estutura de dados correta para o retorno
	NovoSoftware.ID = registo.Key
	NovoSoftware.Nome = &input.Nome

	return &NovoSoftware, nil
}

func (r *mutationResolver) CriarComputador(ctx context.Context, input model.NovoComputador) (*model.ComputadorCriado, error) {
	var registo redishandle.RegistoRedisDB
	var NovoComputador model.ComputadorCriado

	// CriaEstruturaRegisto - Metodo do package redishandle em relação à estrutura RegistoRedisDB
	// passamos input.conteudo como o corpo do registo para podermos mais tarde mapear só
	// a info relacionada ao computador em sí (Ram, Gpu, Cpu...), para a struct defenida na schema do GraphQl
	registo.CriaEstruturaRegisto(&RedisClienteDB, input, "Computador")

	// Inserção do novo registo na base de dados
	// ? Todo o error cheking é feito na função
	redishandle.SetRegistoBD(&RedisClienteDB, registo)

	// Atribuição do input a estutura de dados correta para o retorno
	NovoComputador.ID = registo.Key
	NovoComputador.Nome = &input.Nome

	return &NovoComputador, nil
}

func (r *mutationResolver) CriarItem(ctx context.Context, input model.NovoItem) (*model.ItemCriado, error) {
	var registo redishandle.RegistoRedisDB
	var novoItem model.ItemCriado

	// CriaEstruturaRegisto - Metodo do package redishandle em relação à estrutura RegistoRedisDB
	registo.CriaEstruturaRegisto(&RedisClienteDB, input, "Item")

	// Inserção do novo registo na base de dados
	// ? Todo o error cheking é feito na função
	redishandle.SetRegistoBD(&RedisClienteDB, registo)

	// Atribuição do input a estutura de dados correta para o retorno
	novoItem.ID = registo.Key
	novoItem.Nome = input.Nome

	return &novoItem, nil
}

func (r *mutationResolver) CriarCPU(ctx context.Context, input model.NovoCPU) (*model.ComponenteCriado, error) {
	var registo redishandle.RegistoRedisDB
	var novoCPU model.ComponenteCriado

	// CriaEstruturaRegisto - Metodo do package redishandle em relação à estrutura RegistoRedisDB
	registo.CriaEstruturaRegisto(&RedisClienteDB, input, "CPU")

	// Inserção do novo registo na base de dados
	// ? Todo o error cheking é feito na função
	redishandle.SetRegistoBD(&RedisClienteDB, registo)

	// Atribuição do input a estutura de dados correta para o retorno
	novoCPU.ID = registo.Key
	novoCPU.Modelo = &input.Modelo
	novoCPU.TipoComponente = "CPU"

	return &novoCPU, nil
}

func (r *mutationResolver) CriarGpu(ctx context.Context, input model.NovoGpu) (*model.ComponenteCriado, error) {
	var registo redishandle.RegistoRedisDB
	var novoGPU model.ComponenteCriado

	// CriaEstruturaRegisto - Metodo do package redishandle em relação à estrutura RegistoRedisDB
	registo.CriaEstruturaRegisto(&RedisClienteDB, input, "GPU")

	// Inserção do novo registo na base de dados
	// ? Todo o error cheking é feito na função
	redishandle.SetRegistoBD(&RedisClienteDB, registo)

	// Atribuição do input a estutura de dados correta para o retorno
	novoGPU.ID = registo.Key
	novoGPU.Modelo = input.Modelo
	novoGPU.TipoComponente = "GPU"

	return &novoGPU, nil
}

func (r *mutationResolver) CriarRAM(ctx context.Context, input model.NovoRAM) (*model.ComponenteCriado, error) {
	var registo redishandle.RegistoRedisDB
	var novaRAM model.ComponenteCriado

	// CriaEstruturaRegisto - Metodo do package redishandle em relação à estrutura RegistoRedisDB
	registo.CriaEstruturaRegisto(&RedisClienteDB, input, "RAM")

	// Inserção do novo registo na base de dados
	// ? Todo o error cheking é feito na função
	redishandle.SetRegistoBD(&RedisClienteDB, registo)

	// Atribuição do input a estutura de dados correta para o retorno
	novaRAM.ID = registo.Key
	novaRAM.Modelo = input.Modelo
	novaRAM.TipoComponente = "RAM"

	return &novaRAM, nil
}

func (r *mutationResolver) CriarMboard(ctx context.Context, input model.NovaMBoard) (*model.ComponenteCriado, error) {
	var registo redishandle.RegistoRedisDB
	var NovaMBoard model.ComponenteCriado

	// CriaEstruturaRegisto - Metodo do package redishandle em relação à estrutura RegistoRedisDB
	registo.CriaEstruturaRegisto(&RedisClienteDB, input, "MBoard")

	// Inserção do novo registo na base de dados
	// ? Todo o error cheking é feito na função
	redishandle.SetRegistoBD(&RedisClienteDB, registo)

	// Atribuição do input a estutura de dados correta para o retorno
	NovaMBoard.ID = registo.Key
	NovaMBoard.Modelo = input.Modelo
	NovaMBoard.TipoComponente = "Mother Board"

	return &NovaMBoard, nil
}

func (r *mutationResolver) CriarStorage(ctx context.Context, input model.NovoStorage) (*model.ComponenteCriado, error) {
	var registo redishandle.RegistoRedisDB
	var NovoStorage model.ComponenteCriado

	// CriaEstruturaRegisto - Metodo do package redishandle em relação à estrutura RegistoRedisDB
	registo.CriaEstruturaRegisto(&RedisClienteDB, input, "Storage")

	// Inserção do novo registo na base de dados
	// ? Todo o error cheking é feito na função
	redishandle.SetRegistoBD(&RedisClienteDB, registo)

	// Atribuição do input a estutura de dados correta para o retorno
	NovoStorage.ID = registo.Key
	NovoStorage.Modelo = input.Modelo
	NovoStorage.TipoComponente = "Armazenamento"

	return &NovoStorage, nil
}

func (r *mutationResolver) CriarMicrofone(ctx context.Context, input model.NovoMicrofone) (*model.ComponenteCriado, error) {
	var registo redishandle.RegistoRedisDB
	var NovoMicrofone model.ComponenteCriado

	// CriaEstruturaRegisto - Metodo do package redishandle em relação à estrutura RegistoRedisDB
	registo.CriaEstruturaRegisto(&RedisClienteDB, input, "Microfone")

	// Inserção do novo registo na base de dados
	// ? Todo o error cheking é feito na função
	redishandle.SetRegistoBD(&RedisClienteDB, registo)

	// Atribuição do input a estutura de dados correta para o retorno
	NovoMicrofone.ID = registo.Key
	NovoMicrofone.Modelo = &input.Modelo
	NovoMicrofone.TipoComponente = "Microfone"

	return &NovoMicrofone, nil
}

func (r *mutationResolver) CriarCamera(ctx context.Context, input model.NovaCamera) (*model.ComponenteCriado, error) {
	var registo redishandle.RegistoRedisDB
	var NovaCamera model.ComponenteCriado

	// CriaEstruturaRegisto - Metodo do package redishandle em relação à estrutura RegistoRedisDB
	registo.CriaEstruturaRegisto(&RedisClienteDB, input, "Camera")

	// Inserção do novo registo na base de dados
	// ? Todo o error cheking é feito na função
	redishandle.SetRegistoBD(&RedisClienteDB, registo)

	// Atribuição do input a estutura de dados correta para o retorno
	NovaCamera.ID = registo.Key
	NovaCamera.Modelo = &input.Modelo
	NovaCamera.TipoComponente = "Camera"

	return &NovaCamera, nil
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
	var item model.Computador
	valorDB, err := redishandle.GetRegistoBD(&RedisClienteDB, id)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(valorDB), &item)
	if err != nil {
		resolverLogger.Printf("[!] Erro ao mapear JSON para a Estrutura fornecida: %v\n\t[!] OU os campos da estrutura e do contéudo não são compatíveis", &item)
		return nil, err
	}

	return &item, nil
}

func (r *queryResolver) GetCPUS(ctx context.Context) ([]*model.CPU, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetCPU(ctx context.Context, id string) (*model.CPU, error) {
	var item model.CPU
	valorDB, err := redishandle.GetRegistoBD(&RedisClienteDB, id)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(valorDB), &item)
	if err != nil {
		resolverLogger.Panicf("[!] Erro ao mapear JSON para a Estrutura fornecida: %v\n\t[!] OU os campos da estrutura e do contéudo não são compatíveis", item)
		return nil, err
	}

	return &item, nil
}

func (r *queryResolver) GetGpus(ctx context.Context) ([]*model.Gpu, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetGpu(ctx context.Context, id string) (*model.Gpu, error) {
	var item model.Gpu
	valorDB, err := redishandle.GetRegistoBD(&RedisClienteDB, id)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(valorDB), &item)
	if err != nil {
		resolverLogger.Panicf("[!] Erro ao mapear JSON para a Estrutura fornecida: %v\n\t[!] OU os campos da estrutura e do contéudo não são compatíveis", item)
		return nil, err
	}

	return &item, nil
}

func (r *queryResolver) GetRAMS(ctx context.Context) ([]*model.RAM, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetRAM(ctx context.Context, id string) (*model.RAM, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetItems(ctx context.Context) ([]*model.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetItem(ctx context.Context, id string) (*model.Item, error) {
	var item model.Item
	valorDB, err := redishandle.GetRegistoBD(&RedisClienteDB, id)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(valorDB), &item)
	if err != nil {
		resolverLogger.Panicf("[!] Erro ao mapear JSON para a Estrutura fornecida: %v [!] OU os campos da estrutura e do contéudo não são compatíveis", item)
		return nil, err
	}

	return &item, nil
}

func (r *queryResolver) GetSoftwares(ctx context.Context) ([]*model.Software, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetSoftware(ctx context.Context, id string) (*model.Software, error) {
	var item model.Software
	valorDB, err := redishandle.GetRegistoBD(&RedisClienteDB, id)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(valorDB), &item)
	if err != nil {
		resolverLogger.Panicf("[!] Erro ao mapear JSON para a Estrutura fornecida: %v\n\t[!] OU os campos da estrutura e do contéudo não são compatíveis", item)
		return nil, err
	}

	return &item, nil
}

func (r *queryResolver) GetMicrofones(ctx context.Context) ([]*model.Microfone, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetMicrofone(ctx context.Context, id string) (*model.Microfone, error) {
	var item model.Microfone
	valorDB, err := redishandle.GetRegistoBD(&RedisClienteDB, id)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(valorDB), &item)
	if err != nil {
		resolverLogger.Panicf("[!] Erro ao mapear JSON para a Estrutura fornecida: %v\n\t[!] OU os campos da estrutura e do contéudo não são compatíveis", item)
		return nil, err
	}

	return &item, nil
}

func (r *queryResolver) GetCameras(ctx context.Context) ([]*model.Camera, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetCamera(ctx context.Context, id string) (*model.Camera, error) {
	var item model.Camera
	valorDB, err := redishandle.GetRegistoBD(&RedisClienteDB, id)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(valorDB), &item)
	if err != nil {
		resolverLogger.Panicf("[!] Erro ao mapear JSON para a Estrutura fornecida: %v\n\t[!] OU os campos da estrutura e do contéudo não são compatíveis", item)
		return nil, err
	}

	return &item, nil
}

func (r *queryResolver) GetStorages(ctx context.Context) ([]*model.Storage, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetStorage(ctx context.Context, id string) (*model.Storage, error) {
	var item model.Storage
	valorDB, err := redishandle.GetRegistoBD(&RedisClienteDB, id)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(valorDB), &item)
	if err != nil {
		resolverLogger.Panicf("[!] Erro ao mapear JSON para a Estrutura fornecida: %v\n\t[!] OU os campos da estrutura e do contéudo não são compatíveis", item)
		return nil, err
	}

	return &item, nil
}

func (r *queryResolver) GetMBoards(ctx context.Context) ([]*model.MBoard, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetMBoard(ctx context.Context, id string) (*model.MBoard, error) {
	var item model.MBoard
	valorDB, err := redishandle.GetRegistoBD(&RedisClienteDB, id)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(valorDB), &item)
	if err != nil {
		resolverLogger.Panicf("[!] Erro ao mapear JSON para a Estrutura fornecida: %v\n\t[!] OU os campos da estrutura e do contéudo não são compatíveis", item)
		return nil, err
	}

	return &item, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

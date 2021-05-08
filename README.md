# Backend Projeto Robin
> PAP Tomás Marques

![Robin Logo](https://github.com/TomascpMarques/robin-backend/blob/master/robin-logo.png)

## Intro - Arquitétura
O programa em sí, é constituído por outros microserviços que tentam ao máximo fazer o seu trabalho sem depender dos outros, são lançados como clusters ex: o serviço robinequipamento precis de uma base-de-dados, mas não precisa de ter conexão com o sistema de autenticação, para validar pedidos e ações. Logo os serviços são lançados através de um docker-compose file, que cira a própria rede virtual interna, e os serviços conectam aos outros que forem necessários para o funcionamento.

## Serviço de gestão de equipamento
É um serviço que permite inserir, atualizar e apagar registos de equipamentos como funções básicas.
Mas adiciona funções inspiradas em GraphQl, que permitem buscar os registos da base de dados de uma maneira simples e minimalista na reposta ao mesmo. Permite ao conssumidor da API que especifique parametros, que indicam que tipo de rgisto e seus atributuos, devem ser devolvidos na resposta.

Este serviço implementa as funções de auticação do serviço de login, sem necessitar de conexão ao mesmo.

## Serviço de autenticação
Este Serviço só têm como depedência um outro, a base de dados redis, para guardar users. O serviço disponibiliza a criação, autenticação e verificação de tokens de utilizadores.

## Serviço de informação de utilizador
Este seviço fornece informação mais detalhada sobre o utilizador, tal como o nome próprio, email, especialidades, contribuições, etc.
Conecta com um serviço mongodb, como base de dados, e pode atualizar, eliminar e inserir conteúdo na bd.

## Serviço de documentação/tutoriais para manutenções
Este serviço permite criar repositórios de ficheiros para que assim se possa guardar informação em formato de texto.
Vai ser constituido por dois serviços, um que gere os repos e outro que gere os ficheiros.

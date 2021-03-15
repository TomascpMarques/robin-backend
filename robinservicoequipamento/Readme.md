# Serviço Backend Robin:
## Robin Serviço Equipamento    
Este serviço fica responsável pelo equipamento que o Robin gere.

Disponibiliza funções para poder adicionar, remover, atualizar, extrair informação específica, dos registos existentes na base de dados ligada ao serviço. Estas funções estão disponivéis como _endpoints_ acessivéis pela API (Dynamic Queries GO) logo no index do servidor web na porta 8000 (https://xxxx:8000/).

### Os endpoints são:
* AdicionarRegisto - Adiciona um registo numa base de dados e coleção especifícada

* ApagarRegistoDeItem  - Apaga um registo pelo seu ObjectID, na bd e coleção fornecida

* AtualizararRegistoDeItem 

* BuscarRegistoPorObjID - Busca um registo na base de dados pelo ID especificado

* BuscarRegistosQueryCustom - Toma um nome de uma bd e uma coleção como alvos do query. O query em sí é um map, que vai fornecer os valores ao filtro do tipo bson.M. Toma uma token para autorização

* BuscarInfoItemQuery - Busca um registo e devolve os campos especificados no query

* BuscarInfoItems - BBusca vários registos, com campos que satisfazem o queryBD e devolve os campos especificados no query


Todas estas funções devolvem um `map[string]interface{}`, com o propósito do resultado de cada endpoint ser convertido para json, e enviado logo em seguida, cada função trata da autorização e permissões para o pedido, quando são chamadas, logo todas tem um parametro `token string`, que aceita a JWT, para autenticação do pedido.
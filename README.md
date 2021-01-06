# Backend-Graphql-Equipamento
PAP Tomás Marques

## Intro
Aqui vão ficar as back-ends de cada sistema, agora só tem a do equipamento, mas terá outras como:
  1. **Login.**
  2. **Manutenção:** <br>
    2.1. Atualização <br>
    2.2. Substituições <br>
    2.3. Inspeções <br>
    2.4. _todos os ões_ <br>
  3. **Gerenciamento da infrastrutura fisíca:** <br>
    3.1. Salas de aulas <br>
    3.2. Armazéns <br>
    3.3. Aulas de informática <br>
  4. **Equipamento**

## Arquitetura
Cada parte da backe-end do serviço vai ser colocada num docker container, com a sua respetiva GraphQl-API(em Go) e base-de-dados redis. <br>
Como é que os containers vão falar entre sí é que não sei ainda, talvez grpc? ou uma outra api web?

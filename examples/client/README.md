# Documentação sobre clientes

Requisitos: go versão 1.18

O Hyperledger Fabric é uma plataforma para soluções de contabilidade distribuída sustentadas por uma arquitetura modular que oferece altos graus de confidencialidade, resiliência, flexibilidade e escalabilidade. Ele foi projetado para oferecer suporte a implementações conectáveis ​​de diferentes componentes e acomodar a complexidade e os meandros existentes em todo o ecossistema econômico.

Este projeto blockchain possui dois clientes:

- Outro cliente recebe dados da API do FieldClimate os insere no blockchain

- Um cliente realiza um web scrape, recebendo os dados das estações em tempo real do website [Alerta Rio](http://alertario.rio.rj.gov.br/tabela-de-dados/), e os inserindo no blockchain (chaincode em construção)

Para mais detalhes, a documentação oficial do Hyperledger Fabric será um otimo guia com tutoriais de como a rede funciona -  <link>https://hyperledger-fabric.readthedocs.io/en/latest/</link>

Para utilizar os clientes é necessário primeiro levantar a rede e instalar o seu respectivo chaincode/smart contract. O chaincode utilizado no exemplo é o do FieldClimate.

# Configuração inicial

- Para começar, abra o terminal e acesse a pasta de um dos clientes com o comando ````cd <fieldclimate/alertario>``.

- Abra o arquivo "connection-org.yaml" e substitua o seu conteúdo pelo conteúdo do arquivo encontrado em "../resources/network.yaml" no tutorial inicial.

- Feito isso, dentro do arquivo "connection-org.yaml", navegue até a seção "Organizations" e, na subseção INMETROMSP, INSIRA o seguinte campo
```
    certificateAuthorities:
      - inmetro-ca.default
```

- Agora, na seção "client" campo "organization", deixe o seguinte valor

```
    client:
        organization: Org1MSP
```

Com isso você está pronto para executar os clientes

# Tutorial de como usar o cliente FieldClimate

OBS: É necessário inserir suas próprias chaves pública e privada em fieldclimate/modules/api.go, comentando a linha 31 e alterando as variáveis PUBLIC_KEY, PRIVATE_KEY pelas suas respectivas chaves privada e publica.

Acesse a pasta fieldclimate com o comando:

```
    cd fieldclimate
```

E execute o cliente com o comando:

```
    go run main.go
```

Ele irá buscar os dados de uma estação na API e os inserir dando um Invoke no chaincode. Após isso, fará um Query buscando esses mesmos dados no ledger.

# Tutorial de como usar o cliente FieldClimate

Acesse a pasta alertario com o comando:

```
    cd alertario
```

E execute o cliente com o comando:

```
    go run main.go
```

Ele irá buscar os dados de uma estação no site Alerta Rio e os inserir dando um Invoke no chaincode. Após isso, fará um Query buscando esses mesmos dados no ledger.
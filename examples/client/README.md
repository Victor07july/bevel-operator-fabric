#Documentação sobre clientes

O Hyperledger Fabric é uma plataforma para soluções de contabilidade distribuída sustentadas por uma arquitetura modular que oferece altos graus de confidencialidade, resiliência, flexibilidade e escalabilidade. Ele foi projetado para oferecer suporte a implementações conectáveis ​​de diferentes componentes e acomodar a complexidade e os meandros existentes em todo o ecossistema econômico.

Este projeto blockchain possui dois clientes:

- Outro cliente recebe dados da API do FieldClimate os insere no blockchain

- Um cliente realiza um web scrape, recebendo os dados das estações em tempo real do website [Alerta Rio](http://alertario.rio.rj.gov.br/tabela-de-dados/), e os inserindo no blockchain (chaincode em construção)

Para mais detalhes, a documentação oficial do Hyperledger Fabric será um otimo guia com tutoriais de como a rede funciona -  <link>https://hyperledger-fabric.readthedocs.io/en/latest/</link>

Para utilizar os clientes é necessário primeiro instalar o seu respectivo chaincode/smart contract. O chaincode utilizado no exemplo é o do FieldClimate.

# Tutorial de como usar o cliente FieldClimate

OBS: É necessário inserir suas próprias chaves pública e privada em fieldclimate/modules/api.go, comentando a linha 31 e alterando as variáveis PUBLIC_KEY, PRIVATE_KEY pelas suas respectivas chaves privada e publica.

Para começar, abra o terminal e acesse a pasta do projeto.
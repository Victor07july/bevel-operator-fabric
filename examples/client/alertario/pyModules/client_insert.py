import requests
import sys
from hfc.fabric import Client as client_fabric
from hfc.fabric_network.gateway import Gateway
from hfc.fabric_network.network import Network
from hfc.fabric_network.contract import Contract

import asyncio
from datetime import datetime
import time
#from funcPrevisao import getPrevisao
from funcEstacao import getEstacao

domain = ["connection-org1"]
channel_name = "mychannel"
cc_name = "fabpki"
cc_version = "1.0"
callpeer = ['peer0.org1.example.com', 'peer0.org2.example.com']

URL = "http://alertario.rio.rj.gov.br/upload/TempoReal.html"

# arrays para receber os dados
arrayPrecipitacao = []
arrayDados = []

# inicializando variáveis
horaLeituraP = "Indisponível"
totalUltimaHora = "Indisponível"
situacao = "Indisponível"

horaLeituraD = "Indisponível"
direcaoVentoGraus = "Indisponível"
velocidadeVento = "Indisponível"
temperatura = "Indisponível"
pressao = "Indisponível"
umidade = "Indisponível"

# ------ VERIFICAÇÃO E EXECUÇÃO --------

if __name__ == "__main__":

    #test if the city name was informed as argument
    if len(sys.argv) != 2: # o primeiro  argumento sempre vai ser o chamado do python
        print("Usage:",sys.argv[0], "<\"ID da Estação \"> ")
        exit(1)

    # recebe o id da estação como variavel
    idEstacao = str(sys.argv[1])

    # envia o id inserido para a função de verificar estações
    arrayPrecipitacao, arrayDados, ultimaAtualizacaoE = getEstacao(URL, idEstacao)
    
    # ------ PEGANDO O HORÁRIO DE EXECUÇÃO DO CLIENTE --------
    # Horário de execução do cliente em unix
    timestampCliente = time.time()
    print(f'Horário de execução do cliente em UNIX: {timestampCliente}')

    # ------ PEGANDO O HORÁRIO DE ULTIMA ATUALIZAÇÃO DOS DADOS DAS ESTAÇÕES -------
    # Extrair a hora e a data da string
    hora_data = ultimaAtualizacaoE.split(": ")[-1]  # "10:23 - 18/07/2023"
    
    # Separar a hora e a data da ultima atualização
    hora, data = hora_data.split(" - ")
    
    # Converter para o formato UNIX
    formato = "%H:%M - %d/%m/%Y"
    data_hora = datetime.strptime(hora_data, formato)
    timestampEstacao = int(data_hora.timestamp())
    print(f'Última atualização dos dados das estações em UNIX: {timestampEstacao}')

    # pode ser que uma das estações tenha apenas um dos dados disponiveis (ou nenhum)
    # verifica se pelo menos um dos dados estão disponíveis
    situacao = ""

    if arrayPrecipitacao or arrayDados:
        print("Sucesso! Uma ou ambas as arrays possui os dados necessários")

        if arrayPrecipitacao:
            horaLeitura = arrayPrecipitacao[2]
            totalUltimaHora = arrayPrecipitacao[4]
            situacao = "Somente precipitacao disponivel"

        if arrayDados:
            direcaoVentoGraus = arrayDados[3]
            velocidadeVento = arrayDados[4]
            temperatura = arrayDados[5]
            pressao = arrayDados[6]
            umidade = arrayDados[7]
            if situacao == "Somente precipitacao disponivel":
                situacao = "Precipitacao e dados disponiveis"
            else:
                situacao = "Somente dados disponiveis"

        print(situacao)
        print("Iniciando o chaincode...")
        loop = asyncio.get_event_loop()
        #creates a loop object to manage async transactions
        
        new_gateway = Gateway() # Creates a new gateway instance

        
        c_hlf = client_fabric('/home/stephanie/Inmetrochain-Vehicle/blockchain/gateway/connection-org1.json')
        user = c_hlf.get_user('org1.example.com', 'User1')
        admin = c_hlf.get_user('org1.example.com', 'Admin')
       # print(admin)
        peers = []
        peer = c_hlf.get_peer('peer0.org1.example.com')
        peers.append(peer)
        options = {'wallet': ''}
    
        c_hlf.new_channel(channel_name)
        
        response = loop.run_until_complete(c_hlf.chaincode_invoke(
               requestor=admin,
               channel_name=channel_name,
               peers=['peer0.org1.example.com', 'peer0.org2.example.com'],
               args=[idEstacao,horaLeitura,totalUltimaHora,situacao, direcaoVentoGraus, velocidadeVento, temperatura, pressao, umidade, str(timestampEstacao), str(timestampCliente)],
               fcn= 'insertStationData',
               cc_name=cc_name,
               wait_for_event=True,  # optional, for private data
               # for being sure chaincode invocation has been commited in the ledger, default is on tx event
               #cc_pattern="^invoked*"  # if you want to wait for chaincode event and you have a `stub.SetEvent("invoked", value)` in your chaincode
               ))
        print(response)
    
    else:
        print("Falha! Ambas as arrays estão vazias! Você escolheu um ID válido?")


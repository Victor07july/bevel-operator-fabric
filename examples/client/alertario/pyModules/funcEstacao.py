from bs4 import BeautifulSoup
import requests
import json
import sys

# https://www.reddit.com/r/learnpython/comments/f4v9ex/bs4_looping_through_table/
# https://python-forum.io/thread-27991.html
# https://stackoverflow.com/questions/1936466/how-to-scrape-only-visible-webpage-text-with-beautifulsoup

if __name__ == "__main__":

    #test if the city name was informed as argument
    if len(sys.argv) != 2: # o primeiro  argumento sempre vai ser o chamado do python
        print("Usage:",sys.argv[0], "<\"ID da Estação \"> ")
        exit(1)

    URL = "http://alertario.rio.rj.gov.br/upload/TempoReal.html"

    # faz o request da url
    response = requests.get(URL)
    print(response.status_code)

    # pega o conteúdo e faz o parse
    soup = BeautifulSoup(response.content, 'lxml')

    # seletores CSS
    tabelaPrecipitacao = soup.select_one('body > table:nth-child(2) > tbody:nth-child(2)') 
    tabelaDados = soup.select_one('body > table:nth-child(3) > tbody:nth-child(2)')

    # horário de atualização
    ultimaAtualizacao = soup.select_one('body > p:nth-child(1) > font:nth-child(1)')
    #print(ultimaAtualizacao.text)

    #idEstacao = input("Insira o ID da estação: ")
    idEstacao = sys.argv[1]

    td_list = []
    linhaPrecipitacao = []
    linhaDados = []

    # flag para saber se a estação existe
    flag = False
    oneExists = False # para saber se a estação tem apenas uma das tabelas disponíveis

    for tr in tabelaPrecipitacao.find_all('tr'):
        #td_list.append(tr.find_all('td')[0].text.strip())

        if tr.find_all('td')[0].text.strip() == str(idEstacao):
            flag = True
            #linhaPrecipitacao.append(tr.get_text(", ", strip=True))                
            
            # pega todos os elementos <td> dentro de <tr>, itera sobre eles e insere em uma array
            tds = tr.find_all('td')     
            for td in tds:
                linhaPrecipitacao.append(td.get_text())
            print(f"Dados de precipitação para a estação {idEstacao}, {tr.find_all('td')[1].text.strip()}: ")
            print(linhaPrecipitacao)
            print(f"TESTE: {linhaPrecipitacao[2]}")
            oneExists = True


    if not flag:
        print("Esta estação não está disponível na tabela de dados de precipitação.")
        print("Os seguintes ID's estão disponíveis na tabela de precipitação: ")
        for tr in tabelaPrecipitacao.find_all('tr'):
            print(tr.find_all('td')[0].text.strip(), end=' | ')
        print(' ')

    print('----------' * 3)

    # reseta a flag e td_list
    flag = False
    td_list = []

    for tr in tabelaDados.find_all('tr'):

        if tr.find_all('td')[0].text.strip() == str(idEstacao):
            flag = True
            
            tds = tr.find_all('td')        
            for td in tds:
                linhaDados.append(td.get_text())
            print(f"Dados meteorológicos para a estação {idEstacao}, {tr.find_all('td')[1].text.strip()}")
            print(linhaDados)
            oneExists = True


    if not flag:
        print("Esta estação não está disponível na tabela de dados meteorológicos.")
        print("Os seguintes ID's estão disponíveis na tabela de dados meteorológicos: ")
        for tr in tabelaDados.find_all('tr'):
            print(tr.find_all('td')[0].text.strip(), end=' | ')
        print(' ')

    # Extrair a hora e a data da string ultima atualizacao
    hora_data = ultimaAtualizacao.text.split(": ")[-1]  # "10:23 - 18/07/2023"

    # Separar a hora e a data da ultima atualização
    hora, data = hora_data.split(" - ")

    
    ## Verificação de qual tabela está disponível e atribuição de valores
    # Caso as 2 estejam disponíveis
    if linhaPrecipitacao and linhaDados:
        horaLeitura =              linhaPrecipitacao[2]
        precipitacaoUltimaHora =   linhaPrecipitacao[4]
        direcaoVentoGraus =        linhaDados[3]
        velocidadeVento =          linhaDados[4]
        temperatura =              linhaDados[5]
        pressao =                  linhaDados[6]
        umidade =                  linhaDados[7]
    
    # Caso apenas a tabela de precipitação esteja disponível (caso a estação exista, esta sempre estará disponível)
    elif linhaPrecipitacao:
        print("Lembrete: Tabela de dados não disponível para esta estação")
        horaLeitura =              linhaPrecipitacao[2]
        precipitacaoUltimaHora =   linhaPrecipitacao[4]
        direcaoVentoGraus =        "Indisponivel"
        velocidadeVento =          "Indisponivel"
        temperatura =              "Indisponivel"
        pressao =                  "Indisponivel"
        umidade =                  "Indisponivel"

    # Caso nenhuma das tabelas esteja disponível (isso é, caso a estação não exista), encerra o programa
    if oneExists == False:
        estruturaDados = {
            "OneExists":          oneExists,
        }

        with open("../json/alertario.json", "w") as outfile:
            json.dump(estruturaDados, outfile, indent=4)
        raise Exception("Erro: A estação não existe ou não está disponível")

    # criação da estrutura json com os dados obtidos
    estruturaDados = {
        "HoraLeitura":        horaLeitura,
        "PrecipitacaoUltHr":  precipitacaoUltimaHora,
        "DirVentoGraus":      direcaoVentoGraus,
        "VelocidadeVento":    velocidadeVento,
        "Temperatura":        temperatura,
        "Pressao":            pressao,
        "Umidade":            umidade,
        "UltimaAtualizacao":  hora + " " + data,
        "OneExists":          oneExists,

    }

    with open("../json/alertario.json", "w") as outfile:
        json.dump(estruturaDados, outfile, indent=4)
    print("JSON criado com sucesso!")

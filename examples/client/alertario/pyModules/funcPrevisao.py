from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.firefox.options import Options

# ESTA PREVISÃO É PARA AS PRÓXIMAS 24 HORAS E EU PEGUEI SOMENTE PARA A PRÓXIMA MADRUGADA E PRÓXIMA MANHÃ
# HÁ TAMBÉM UMA PREVISÃO GERAL PARA OS PRÓXIMOS 4 DIAS
# http://alertario.rio.rj.gov.br/4-dias/

# para poder escolher printar ou não
def getPrevisao(print_output = True):
    # personaliza configurações do driver
    options = Options()
    options.add_argument("--headless") # abrir o navegador em segundo plano

    # usa a engine do firefox (geckodriver) como driver para buscar os elementos
    driver = webdriver.Firefox(options=options)

    # definindo um tempo de espera implícito de 5 segundos, para esperar que a página carregue
    #driver.implicitly_wait(3)

    # acessa a página do Alerta Rio
    driver.get("http://alertario.rio.rj.gov.br/24-horas/")

    # os dados da tabela estão dentro de um iframe (JavaScript), portanto, é necssário mudar o driver para o modo iframe
    iframe = driver.find_element(By.CLASS_NAME, "iframe-class")
    driver.switch_to.frame(iframe)

    # busca o elemento HTML através de seu XPATH
    titulo = driver.find_element(By.XPATH, "/html/body/table/thead[1]/tr[1]/th/b")

    # variáveis de previsao para proxima madrugada
    ceuProxMadrugada          = driver.find_element(By.XPATH, "/html/body/table/tbody[1]/tr[2]/td[4]")
    precipitacaoProxMadrugada = driver.find_element(By.XPATH, "/html/body/table/tbody[1]/tr[3]/td[4]")
    ventoProxMadrugada        = driver.find_element(By.XPATH, "/html/body/table/tbody[1]/tr[4]/td[4]")
    tendenciaProxMadrugada    = driver.find_element(By.XPATH, "/html/body/table/tbody[1]/tr[4]/td[4]")

    # variáveis de previsão para próxima manhã
    ceuProxManha              = driver.find_element(By.XPATH, "/html/body/table/tbody[1]/tr[2]/td[5]") 
    precipitacaoProxManha     = driver.find_element(By.XPATH, "/html/body/table/tbody[1]/tr[3]/td[5]")
    ventoProxManha            = driver.find_element(By.XPATH, "/html/body/table/tbody[1]/tr[4]/td[5]")
    tendenciaProxManha        = driver.find_element(By.XPATH, "/html/body/table/tbody[1]/tr[5]/td[5]")

    # previsao de temperatura para o dia atual
    ultimaAtualizacao         = driver.find_element(By.XPATH, "/html/body/p")
    barraJacarepagua          = driver.find_element(By.XPATH, "/html/body/table/tbody[2]/tr/td[1]")
    centroGrandeTijuca        = driver.find_element(By.XPATH, "/html/body/table/tbody[2]/tr/td[2]")
    zonaNorte                 = driver.find_element(By.XPATH, "/html/body/table/tbody[2]/tr/td[3]")
    zonaOeste                 = driver.find_element(By.XPATH, "/html/body/table/tbody[2]/tr/td[4]")
    zonaSul                   = driver.find_element(By.XPATH, "/html/body/table/tbody[2]/tr/td[5]")

    # previsao sinótica
    previsaoSinotica          = driver.find_element(By.XPATH, "/html/body/table/tbody[4]/tr/td")

    # imprime as previsoes
    if print_output:
        print('----------------------------------' * 3)
        # print("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~" + titulo.text + "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
        print('----------------------------------' * 3)
        print('-> PRÓXIMA MADRUGADA')
        print('----------------------------------' * 3)
        print("Céu na próxima madrugada: " + ceuProxMadrugada.text)
        print("Precipitação na próxima madrugada: " + precipitacaoProxMadrugada.text)
        print("Vento na próxima madrugada: " + ventoProxMadrugada.text)
        print("Tendencia da Temperatura na próxima madrugada: " + tendenciaProxMadrugada.text)

        print('----------------------------------' * 3)
        print('-> PRÓXIMA MANHÃ')
        print('----------------------------------' * 3)
        print("Céu na próxima manha: " + ceuProxManha.text)
        print("Precipitação na próxima manha: " + precipitacaoProxManha.text)
        print("Vento na próxima manha: " + ventoProxManha.text)
        print("Tendencia da Temperatura na próxima manha: " + tendenciaProxManha.text)

        print('----------------------------------' * 3)
        print('-> PREVISÃO DE TEMPERATURA DO DIA ATUAL PARA AS REGIÕES DO RJ')
        print('----------------------------------' * 3)
        print('-> ' + ultimaAtualizacao.text)
        print('Barra/Jacarepaguá: \n' + barraJacarepagua.text + "\n")
        print('Centro/Grande Tijuca: \n' + centroGrandeTijuca.text + "\n")
        print('Zona Norte: \n' + zonaNorte.text + "\n")
        print('Zona Oeste: \n' + zonaOeste.text + "\n")
        print('Zona Sul: \n' + zonaSul.text + "\n")

        print('----------------------------------' * 3)
        print('-> PREVISAO SINÓTICA (de modo resumido)')
        print(previsaoSinotica.text)

        print('----------------------------------' * 3)
        print('-> AS PREVISÕES ACIMA SE REFERM ÀS PRÓXIMAS 24 HORAS')
        print('-> Fonte: Alerta Rio: http://alertario.rio.rj.gov.br/24-horas/')
        print('----------------------------------' * 3)

    previsoes =  [
        # titulo.text, 
        ceuProxMadrugada.text, 
        precipitacaoProxMadrugada.text, 
        ventoProxMadrugada.text, 
        tendenciaProxMadrugada.text, 
        ceuProxManha.text, 
        precipitacaoProxManha.text, 
        ventoProxManha.text, 
        tendenciaProxManha.text, 
        ultimaAtualizacao.text, 
        barraJacarepagua.text, 
        centroGrandeTijuca.text, 
        zonaNorte.text, 
        zonaOeste.text, 
        zonaSul.text, 
        previsaoSinotica.text
    ]

    # sair do iframe
    driver.switch_to.default_content()

    # encerrar o driver
    driver.quit()

    return previsoes
    
#getPrevisao(print_output=True)
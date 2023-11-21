package modules

import (
	"encoding/json"
	"fmt"
	"os"
	log "github.com/sirupsen/logrus"

)

type AlertaRio struct {
	HoraLeitura         string      `json:"horaleitura"`
	PrecipitacaoUltHora string 		`json:"precipitacaoulthr"`
	DirVentoGraus       string 		`json:"dirventograus"`
	VelVento            string      `json:"velocidadevento"`
	Temperatura         string      `json:"temperatura"`
	Pressao             string      `json:"pressao"`
	Umidade             string      `json:"umidade"`
	UltimaAtualizacao   interface{} `json:"ultimaatualizacao"`
	OneExists 		 	bool        `json:"oneexists"`
}

func JSONRead() (string, string, string, string, string, string, string, string) {
	fmt.Println("Lendo JSON...")
	file, err := os.Open("json/alertario.json")
	if err != nil {
		panic(err)
	}
	fmt.Println("Leitura realizada com sucesso!")
	defer file.Close()

	// criando um decoder JSON para o arquivo lido
	decoder := json.NewDecoder(file)

	// decodificando json na estrutura
	alertario := new(AlertaRio)
	err = decoder.Decode(&alertario)
	if err != nil {
		fmt.Println("Erro ao decodificar JSON: ", err)
		log.Fatal(err)
	}


	// verificando se há dados para esta estação
	oneExists := alertario.OneExists
	if !oneExists {
		fmt.Println("Não há dados para esta estação! Verifique se ela está disponível. Encerrando programa...")
		log.Info("Esta estação não existe ou está indisponível. Encerrando...")
		os.Exit(0)
	}


	// separando os dados da estrutura em variáveis
	horaLeitura := alertario.HoraLeitura
	precipitacaoUltHora := alertario.PrecipitacaoUltHora
	dirVentoGraus := alertario.DirVentoGraus
	velVento := alertario.VelVento
	temperatura := alertario.Temperatura
	pressao := alertario.Pressao
	umidade := alertario.Umidade
	ultimaAtualizacao := fmt.Sprintf("%v", alertario.UltimaAtualizacao)

	log.Info("Dados encontrados. Retornando...")
	return horaLeitura, precipitacaoUltHora, dirVentoGraus, velVento, temperatura, pressao, umidade, ultimaAtualizacao
}

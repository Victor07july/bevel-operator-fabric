package main

import (
	"alertario/modules"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"os"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	// inicializando o log
	file, err := os.OpenFile("logs/log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Info("Iniciando cliente...")

	fmt.Print("Insira o ID da estação: ")
	var stationID string
	fmt.Scanln(&stationID)
	log.Info("Estação buscada: ", stationID)

	//configFilePath := os.Args[1]
	configFilePath := "connection-org.yaml"
	channelName := "demo"
	mspID := "INMETROMSP"
	chaincodeName := "alertario"

	log.Info("Conectando-se ao Alerta Rio...")
	modules.CallPy(stationID)
	log.Info("Conexão realizada. Estação buscada: ", stationID + ". Dados da busca salvos em json/alertario.json")

	log.Info("Lendo JSON...")
	/*pega os dados retornados (lidos no JSON) */
	horaLeitura, precipitacaoUltHora, dirVentoGraus, velVento, temperatura, pressao, umidade, ultimaAtualizacao := modules.JSONRead()

	/*transformando horario em unix*/
	// Obtém a data e hora atual e converte em unix
	dataAtual := time.Now()
	timestampClient := dataAtual.Unix()

	// obtém a data e hora dos dados da api e converte em unix
	layout := "15:04 02/01/2006" // Define o layout da string de data
	// Transforma ultimaAtualizacao em horário Unix
	ultimaAtualizacaoTime, err := time.Parse(layout, ultimaAtualizacao)
	if err != nil {
		log.Error("Failed to parse ultimaAtualizacao: ", err)
	}
	timestampAtualizacao := ultimaAtualizacaoTime.Unix()

	fmt.Println("Hora Leitura: ", horaLeitura)
	fmt.Println("Precipitacao Ultima Hora: ", precipitacaoUltHora)
	fmt.Println("Direcao Vento Graus: ", dirVentoGraus)
	fmt.Println("Velocidade Vento: ", velVento)
	fmt.Println("Temperatura: ", temperatura)
	fmt.Println("Pressao: ", pressao)
	fmt.Println("Umidade: ", umidade)
	fmt.Println("Ultima Atualizacao: ", ultimaAtualizacao)
	fmt.Println("Timestamp Atualizacao: ", timestampAtualizacao)
	fmt.Println("Timestamp Client: ", timestampClient)


	enrollID := randomString(10)
	registerEnrollUser(configFilePath, enrollID, mspID)

	/* O invoke pode ser feito com o gateway (gw) (recomendado) ou sem */
	log.Info("Realizando invoke...")
	modules.InvokeCCgw(configFilePath, channelName, enrollID, mspID, chaincodeName, "InsertStationData", []string{
		stationID,
		horaLeitura,
		precipitacaoUltHora,
		dirVentoGraus,
		velVento,
		temperatura,
		pressao,
		umidade,
		strconv.FormatInt(timestampAtualizacao, 10),
		strconv.FormatInt(timestampClient, 10),
		})
	log.Info("Realizando query")
	modules.QueryCCgw(configFilePath, channelName, enrollID, mspID, chaincodeName, "QueryStation", []string{stationID})
}

func registerEnrollUser(configFilePath, enrollID, mspID string) {
	log.Info("Registering User : ", enrollID)
	sdk, err := fabsdk.New(config.FromFile(configFilePath))
	ctx := sdk.Context()
	caClient, err := mspclient.New(
		ctx,
		mspclient.WithCAInstance("inmetro-ca.default"),
		mspclient.WithOrg(mspID),
	)

	if err != nil {
		log.Error("Failed to create msp client: %s\n", err)
	}

	if caClient != nil {
		log.Info("ca client created")
	}
	enrollmentSecret, err := caClient.Register(&mspclient.RegistrationRequest{
		Name:           enrollID,
		Type:           "client",
		MaxEnrollments: -1,
		Affiliation:    "",
		Attributes:     nil,
		Secret:         enrollID,
	})
	if err != nil {
		log.Error(err)
	}
	err = caClient.Enroll(
		enrollID,
		mspclient.WithSecret(enrollmentSecret),
		mspclient.WithProfile("tls"),
	)
	if err != nil {
		log.Error(errors.WithMessage(err, "failed to register identity"))
	}

	wallet, err := gateway.NewFileSystemWallet(fmt.Sprintf("wallet/%s", mspID))

	signingIdentity, err := caClient.GetSigningIdentity(enrollID)
	key, err := signingIdentity.PrivateKey().Bytes()
	identity := gateway.NewX509Identity(mspID, string(signingIdentity.EnrollmentCertificate()), string(key))

	err = wallet.Put(enrollID, identity)
	if err != nil {
		log.Error(err)
	}

}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

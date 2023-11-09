package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"fieldclimate/modules"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	//configFilePath := os.Args[1]
	configFilePath := "connection-org.yaml"
	channelName := "demo"
	mspID := "INMETROMSP"
	chaincodeName := "fieldclimate"

	enrollID := randomString(10)
	registerEnrollUser(configFilePath, enrollID, mspID)

	/* conecta-se a API e insere dados em um json */
	modules.APIConnect("00206C61")
	
	// lê os dados do json
	stationID := "00206C61"
	deviceName, deviceValues, deviceUnit, deviceDate := modules.JSONRead(stationID, "HC Air temperature")

	/* CONVERSÃO DE DATAS EM UNIX */
	// Obtém a data e hora atual e converte em unix
	dataAtual := time.Now()
	clientExecutionUnix := dataAtual.Unix()

	// obtém a data e hora dos dados da api e converte em unix
	layout := "2006-01-02 15:04:05"

	parsedDeviceDate, err := time.Parse(layout, deviceDate)
    if err != nil {
        fmt.Println("Erro ao analisar a data:", err)
        return
    }
	deviceDateUnix := parsedDeviceDate.Unix()


	// passando valores para string
	deviceValuesString := ""
	for key, values := range deviceValues {
		deviceValuesString += key + ": " + fmt.Sprint(values) + " "
	}

	deviceDateUnixStr := strconv.FormatInt(deviceDateUnix, 10)
	clientExecutionUnixStr := strconv.FormatInt(deviceDateUnix, 10)

	fmt.Println("Nome do dispositivo: ", deviceName)
	fmt.Println("Dados enviados por ele: ", deviceValuesString)
	fmt.Println("Unidade de medição: ", deviceUnit)
	fmt.Println("Horário em que os dados foram atualizados na API: ", deviceDate)
	fmt.Println("Horário de inserção dos dados na API em Unix: ", deviceDateUnix)
	fmt.Println("Horário de execução do cliente: ", dataAtual)
	fmt.Println("Horário de execução do cliente em unix: ", clientExecutionUnix)

	/* O invoke pode ser feito com o gateway (gw) (recomendado) ou sem */
	modules.InvokeCCgw(configFilePath, channelName, enrollID, mspID, chaincodeName, "InsertDeviceData", []string{
		stationID, 
		deviceName, 
		deviceUnit.(string), 
		deviceValuesString, 
		deviceDateUnixStr, 
		clientExecutionUnixStr,
		})
	modules.QueryCCgw(configFilePath, channelName, enrollID, mspID, chaincodeName, "QueryDevice", []string{stationID})
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
		Attributes: nil,
		Secret:     enrollID,
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

package main

import (
	"fmt"
	"math/rand"
	"time"

	"fieldclimate/modules"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	// configurações do cliente
	//configFilePath := os.Args[1]
	configFilePath := "connection-org.yaml"
	channelName := "demo"
	mspID := "INMETROMSP"
	chaincodeName := "fieldclimate"

	enrollID := randomString(10)
	registerEnrollUser(configFilePath, enrollID, mspID)

	// puxando dados da api
	// conecta-se a API e insere dados em um json
	//modules.APIConnect("00206C61")
	
	// lê os dados do json
	deviceName, deviceValues, deviceUnit, deviceDate := modules.JSONRead("00206C61", "HC Air temperature")

	// CONVERSÃO DE DATAS EM UNIX
	// Obtém a data e hora atual e converte em unix
	dataAtual := time.Now()
	unixTimestampAtual := dataAtual.Unix()

	// obtém a data e hora dos dados da api e converte em unix
	layout := "2006-01-02 15:04:05"

	parsedDeviceDate, err := time.Parse(layout, deviceDate)
    if err != nil {
        fmt.Println("Erro ao analisar a data:", err)
        return
    }

	deviceDateUnix := parsedDeviceDate.Unix()

	fmt.Println("Nome do dispositivo: ", deviceName)
	fmt.Println("Dados enviados por ele: ", deviceValues)
	fmt.Println("Unidade de medição: ", deviceUnit)
	fmt.Println("Horário de inserção dos dados na API em Unix: ", deviceDateUnix)
	fmt.Println("Horário de execução do cliente em unix: ", unixTimestampAtual)

	/* O invoke pode ser feito com o gateway (recomendado) ou sem */
	//invokeCC(configFilePath, channelName, enrollID, mspID, chaincodeName, "ReadStationData")
	invokeCCgw(configFilePath, channelName, enrollID, mspID, chaincodeName, "ReadStationData") //stationid, devicename, ...
	//queryCC(configFilePath, channelName, enrollID, mspID, chaincodeName, "QueryAllCars")
	//queryCCgw(configFilePath, channelName, enrollID, mspID, chaincodeName, "QueryAllCars")
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
	//mspclient.WithCAInstance("hq-guild-ca.fabric"),
	//mspclient.WithOrg(mspID),

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
		// CAName:         "Org1MSP",
		Attributes: nil,
		Secret:     enrollID,
	})
	if err != nil {
		//fmt.Println("VERIFICAÇÃO")
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
func invokeCCgw(configFilePath, channelName, userName, mspID, chaincodeName, fcn string) {

	configBackend := config.FromFile(configFilePath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		log.Error(err)
	}

	wallet, err := gateway.NewFileSystemWallet(fmt.Sprintf("wallet/%s", mspID))

	gw, err := gateway.Connect(
		gateway.WithSDK(sdk),
		gateway.WithUser(userName),
		gateway.WithIdentity(wallet, userName),
	)
	if err != nil {
		log.Error("Failed to create new Gateway: %s", err)
	}
	defer gw.Close()
	nw, err := gw.GetNetwork(channelName)
	if err != nil {
		log.Error("Failed to get network: %s", err)
	}

	contract := nw.GetContract(chaincodeName)

	// aqui ele chama a função com os parametros!
	//resp, err := contract.SubmitTransaction(fcn, userName, "a", "b", "1", "ewdscwds")
	resp, err := contract.SubmitTransaction(fcn, "parametro1", "parametro2")
	
	if err != nil {
		log.Error("Failed submit transaction: %s", err)
	}
	log.Info(resp)

}
func invokeCC(configFilePath, channelName, userName, mspID, chaincodeName, fcn string) {

	userName = "admin"

	configBackend := config.FromFile(configFilePath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		log.Error(err)
	}

	chContext := sdk.ChannelContext(
		channelName,
		fabsdk.WithUser(userName),
		fabsdk.WithOrg(mspID),
	)

	ch, err := channel.New(chContext)
	if err != nil {
		log.Error(err)
	}

	var args [][]byte

	inputArgs := []string{userName, "23", "234", "2324", "234"}
	for _, arg := range inputArgs {
		args = append(args, []byte(arg))
	}
	response, err := ch.Execute(
		channel.Request{
			ChaincodeID:     chaincodeName,
			Fcn:             fcn,
			Args:            args,
			TransientMap:    nil,
			InvocationChain: nil,
			IsInit:          false,
		},
	)

	if err != nil {
		log.Error(err)
	}

	log.Infof("txid=%s", response.TransactionID)
}

func queryCC(configFilePath, channelName, userName, mspID, chaincodeName, fcn string) {
	userName = "admin"

	configBackend := config.FromFile(configFilePath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		log.Error(err)
	}
	log.Println(sdk)
	chContext := sdk.ChannelContext(
		channelName,
		fabsdk.WithUser(userName),
		fabsdk.WithOrg(mspID),
	)

	ch, err := channel.New(chContext)
	if err != nil {
		log.Error(err)
	}

	response, err := ch.Query(
		channel.Request{
			ChaincodeID:     chaincodeName,
			Fcn:             fcn,
			Args:            nil,
			TransientMap:    nil,
			InvocationChain: nil,
			IsInit:          false,
		},
	)

	if err != nil {
		log.Error(err)
	}
	log.Infof("response=%s", response.Payload)
}

func queryCCgw(configFilePath, channelName, userName, mspID, chaincodeName, fcn string) {

	configBackend := config.FromFile(configFilePath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		log.Error(err)
	}
	gw, err := gateway.Connect(
		gateway.WithSDK(sdk),
		gateway.WithUser(userName),
	)

	if err != nil {
		log.Error("Failed to create new Gateway: %s", err)
	}
	defer gw.Close()
	nw, err := gw.GetNetwork(channelName)
	if err != nil {
		log.Error("Failed to get network: %s", err)
	}

	contract := nw.GetContract(chaincodeName)

	resp, err := contract.EvaluateTransaction(fcn)

	if err != nil {
		log.Error("Failed submit transaction: %s", err)
	}
	log.Info(string(resp))

}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

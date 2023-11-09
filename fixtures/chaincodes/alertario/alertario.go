package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *DadosEstacaoAR
}

type DadosEstacaoAR struct {
	// ID da estação é a chave e não entra no struct
	HoraLeitura       	  string `json:"horaleitura"`
	PrecipitacaoUltHora   string `json:"totalultimahora"`
	DirecaoVentoGraus 	  string `json:"direcaoventograus"`
	VelocidadeVento   	  string `json:"velocidadevento"`
	Temperatura       	  string `json:"temperatura"`
	Pressao           	  string `json:"pressao"`
	Umidade           	  string `json:"umidade"`
	TimestampEstacao  	  string `json:"timestampestacao"`
	TimestampCliente  	  string `json:"timestampcliente"`
}

const compositeKey = ""

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	stationData := []DadosEstacaoAR{
		DadosEstacaoAR{HoraLeitura: "12:00", PrecipitacaoUltHora: "0", DirecaoVentoGraus: "0", VelocidadeVento: "0", Temperatura: "0", Pressao: "0", Umidade: "0", TimestampEstacao: "0", TimestampCliente: "0"},
		DadosEstacaoAR{HoraLeitura: "13:00", PrecipitacaoUltHora: "1", DirecaoVentoGraus: "1", VelocidadeVento: "1", Temperatura: "1", Pressao: "1", Umidade: "1", TimestampEstacao: "1", TimestampCliente: "1"},
		DadosEstacaoAR{HoraLeitura: "14:00", PrecipitacaoUltHora: "2", DirecaoVentoGraus: "2", VelocidadeVento: "2", Temperatura: "2", Pressao: "2", Umidade: "2", TimestampEstacao: "2", TimestampCliente: "2"},
		DadosEstacaoAR{HoraLeitura: "15:00", PrecipitacaoUltHora: "3", DirecaoVentoGraus: "3", VelocidadeVento: "3", Temperatura: "3", Pressao: "3", Umidade: "3", TimestampEstacao: "3", TimestampCliente: "3"},
		DadosEstacaoAR{HoraLeitura: "16:00", PrecipitacaoUltHora: "4", DirecaoVentoGraus: "4", VelocidadeVento: "4", Temperatura: "4", Pressao: "4", Umidade: "4", TimestampEstacao: "4", TimestampCliente: "4"},
	}

	for i, data := range stationData {
		deviceAsBytes, _ := json.Marshal(data)
		err := ctx.GetStub().PutState("DEVICE"+strconv.Itoa(i), deviceAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

func (s *SmartContract) InsertStationData(ctx contractapi.TransactionContextInterface, stationID string, horaLeitura string, precipitacaoUltHora string, direcaoVentoGraus string, temperatura string, pressao string, umidade string, timestampEstacao string, timestampCliente string) error {

	stationData := DadosEstacaoAR{
		HoraLeitura: 	 		horaLeitura,
		PrecipitacaoUltHora: 	precipitacaoUltHora,
		DirecaoVentoGraus: 		direcaoVentoGraus,
		VelocidadeVento: 		"0",
		Temperatura: 			temperatura,
		Pressao: 				pressao,
		Umidade: 				umidade,
		TimestampEstacao: 		timestampEstacao,
		TimestampCliente:		timestampCliente,
	}

	stationDataAsBytes, _ := json.Marshal(stationData)


	return ctx.GetStub().PutState(stationID, stationDataAsBytes)
	
}

func (s *SmartContract) QueryStation(ctx contractapi.TransactionContextInterface, stationID string) (*DadosEstacaoAR, error) {
	stationDataAsBytes, err := ctx.GetStub().GetState(stationID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if stationDataAsBytes == nil {
		return nil, fmt.Errorf("%s estação com o nome inserido não existe", stationID)
	}

	stationData := new(DadosEstacaoAR)
	_ = json.Unmarshal(stationDataAsBytes, stationData)

	return stationData, nil
}


// QueryAllDevices returns all devices found in world state
func (s *SmartContract) QueryAllStations(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		stationData := new(DadosEstacaoAR)
		_ = json.Unmarshal(queryResponse.Value, stationData)

		queryResult := QueryResult{Key: queryResponse.Key, Record: stationData}
		results = append(results, queryResult)
	}

	return results, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}

/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ver git tag github

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}
/*
// id chave 
type Station {
	credits
	owner
}
*/
// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Device
}

// igual o device atual sem o values
// o que não muda vai pro register
/*
	deviceregister

	devicename


*/

// o que muda vai pro data
/*
devicedata

	Values              string `json:"values"`
	LastUpdateUnix      string `json:"lastupdateunix"`
	ClientExecutionUnix string `json:"clientexecutionunix"`
*/

// dispositivo de uma estação. Identificador será station id
type Device struct {
	DeviceName          string `json:"devicename"`
	Unit                string `json:"unit"`
	Values              string `json:"values"`
	LastUpdateUnix      string `json:"lastupdateunix"`
	ClientExecutionUnix string `json:"clientexecutionunix"`
	// credits
}

const compositeKey = ""

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	devices := []Device{
		Device{DeviceName: "HC Air temperature", Unit: "%C", Values: "12, 13, 24", LastUpdateUnix: "123123", ClientExecutionUnix: "234234"},
		Device{DeviceName: "HC Air temperature", Unit: "%C", Values: "=3, 20, 23", LastUpdateUnix: "123124", ClientExecutionUnix: "234235"},
		Device{DeviceName: "HC Air temperature", Unit: "%C", Values: "14, 25, 22", LastUpdateUnix: "123125", ClientExecutionUnix: "234236"},
		Device{DeviceName: "HC Air temperature", Unit: "%C", Values: "15, 30, 21", LastUpdateUnix: "123126", ClientExecutionUnix: "234237"},
		Device{DeviceName: "HC Air temperature", Unit: "%C", Values: "16, 35, 20", LastUpdateUnix: "123127", ClientExecutionUnix: "234238"},
	}

	for i, device := range devices {
		deviceAsBytes, _ := json.Marshal(device)
		err := ctx.GetStub().PutState("DEVICE"+strconv.Itoa(i), deviceAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// criar função de registrar esta~çao

func (s *SmartContract) InsertDeviceData(ctx contractapi.TransactionContextInterface, stationID string, deviceName string, unit string, values string, lastupdateunix string, clientexecutionunix string) error {
	device := Device{
		DeviceName:          deviceName,
		Unit:                unit,
		Values:              values,
		LastUpdateUnix:      lastupdateunix,
		ClientExecutionUnix: clientexecutionunix,
	}

	deviceAsBytes, _ := json.Marshal(device)

	// chave composta
	key, err := ctx.GetStub().CreateCompositeKey(compositeKey, []string{stationID, deviceName})

	if err != nil {
		return fmt.Errorf("Failed to create composite key: %s", err.Error())
	}

	return ctx.GetStub().PutState(key, deviceAsBytes)
}

func (s *SmartContract) QueryDeviceByCompositeKey(ctx contractapi.TransactionContextInterface, stationID string, deviceName string) (*Device, error) {
	// cria a chave composta usando o id da estação e nome do dispositivo
	key, err := ctx.GetStub().CreateCompositeKey(compositeKey, []string{stationID, deviceName})

	if err != nil {
		return nil, fmt.Errorf("Failed to create composite key: %s", err)
	}

	// buscar o valor associado à chave composta
	/* tentar usar getStateByPartialCompositeKey futuramente */
	valueAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %s", err)
	}

	if valueAsBytes == nil {
		return nil, fmt.Errorf("Nenhum valor encontrado para a chave", key)
	}

	// criando estrutura
	device := new(Device)

	// realiza o unmarshal do valor encontrado
	_ = json.Unmarshal(valueAsBytes, device)

	return device, nil

}

// procura dados do dispositivo em um tempo unix específico
func (s *SmartContract) QueryByHistory(ctx contractapi.TransactionContextInterface, stationID string, deviceName string, unixTime string) (*Device, error) {
	key, err := ctx.GetStub().CreateCompositeKey(compositeKey, []string{stationID, deviceName})

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	historyIer, err := ctx.GetStub().GetHistoryForKey(key)

	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	for historyIer.HasNext() {
		queryResponse, err := historyIer.Next()
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		deviceAsBytes := queryResponse.Value
		device := new(Device)
		json.Unmarshal(deviceAsBytes, device)
		lastUpdateUnix := device.LastUpdateUnix

		if unixTime == lastUpdateUnix {
			fmt.Printf("Found device data with timestamp: ", unixTime)
			return device, nil
		}

	}
	historyIer.Close()
	return nil, fmt.Errorf("Device data with timestamp: " + unixTime + " not found")

}

// QueryAllDevices returns all devices found in world state
func (s *SmartContract) QueryAllDevices(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
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

		device := new(Device)
		_ = json.Unmarshal(queryResponse.Value, device)

		queryResult := QueryResult{Key: queryResponse.Key, Record: device}
		results = append(results, queryResult)
	}

	return results, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fieldclimate chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fieldclimate chaincode: %s", err.Error())
	}
}

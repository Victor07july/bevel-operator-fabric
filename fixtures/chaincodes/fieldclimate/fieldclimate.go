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

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Device
}

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

	// verifica se já existe com a mesma chave composta (id e updatetimeunix)
	// exists, err := s.AssetExists(ctx, key)
	// if err != nil {
	// 	return err
	// }

	// if exists {
	// 	return fmt.Errorf("o objeto %s já existe com o mesmo id", key)
	// }

	return ctx.GetStub().PutState(key, deviceAsBytes)
	
}

func (s *SmartContract) QueryDevice(ctx contractapi.TransactionContextInterface, stationID string) (*Device, error) {
	deviceAsBytes, err := ctx.GetStub().GetState(stationID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if deviceAsBytes == nil {
		return nil, fmt.Errorf("%s dispositivo com o nome inserido não existe", stationID)
	}

	device := new(Device)
	_ = json.Unmarshal(deviceAsBytes, device)

	return device, nil
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
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}

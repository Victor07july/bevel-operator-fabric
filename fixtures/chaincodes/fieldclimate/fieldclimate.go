/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// identificador será station id
type StationData struct {
	DeviceName          string `json:"name"`
	Unit                string `json:"unit"`
	Values              string `json:"values"`
	LastUpdateUnix      string `json:"lastupdateunix"`
	ClientExecutionUnix string `json:"clientexecutionunix"`
	// credits
}

const compositeKey = ""

// chave composta id e timestamp unix
func (s *SmartContract) InsertStationData(ctx contractapi.TransactionContextInterface,
	stationid string,
	devicename string,
	unit string,
	values string,
	lastupdateunix string,
	clientexecutionunix string,
) error {

	stationdata := StationData{
		DeviceName:          devicename,
		Unit:                unit,
		Values:              values,
		LastUpdateUnix:      lastupdateunix,
		ClientExecutionUnix: clientexecutionunix,
	}

	stationdataJSON, err := json.Marshal(stationdata)
	if err != nil {
		return err
	}

	// chave composta
	key, err := ctx.GetStub().CreateCompositeKey(compositeKey, []string{stationid, devicename})

	// verifica se já existe com a mesma chave composta (id e updatetimeunix)
	// exists, err := s.AssetExists(ctx, key)
	// if err != nil {
	// 	return err
	// }

	// if exists {
	// 	return fmt.Errorf("o objeto %s já existe com o mesmo id", key)
	// }

	return ctx.GetStub().PutState(key, stationdataJSON)
}

func (s *SmartContract) ReadStationData(ctx contractapi.TransactionContextInterface, id string, devicename string) (*StationData, error) {
	// cria a chave composta usando o id
	key, err := ctx.GetStub().CreateCompositeKey(compositeKey, []string{id, devicename})
	if err != nil {
		return nil, err
	}

	// buscar o valor associado à chave composta
	// *tentar usar getStateByPartialCompositeKey*
	value, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, err
	}

	// verifique se o valor retornado é nulo
	if value == nil {
		return nil, fmt.Errorf("Nenhum valor encontrado para a chave %s", key)
	}

	// cria estrutura para receber manipular os dados recebidos
	var stationdata StationData

	// realiza o unmarshal dos dados recebidos na estrutura
	err = json.Unmarshal(value, &stationdata)
	if err != nil {
		return nil, err
	}

	return &stationdata, nil
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

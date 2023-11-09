package modules

import (
	"encoding/json"
	"fmt"
	"os"
)

// Define uma estrutura que corresponde à estrutura do JSON
type SensorValue struct {
	Name         string                   `json:"name"`
	NameOriginal string                   `json:"name_original"`
	Type         string                   `json:"type"`
	Decimals     int                      `json:"decimals"`
	Unit         interface{}              `json:"unit"`
	Ch           int                      `json:"ch"`
	Code         int                      `json:"code"`
	Group        int                      `json:"group"`
	Serial       interface{}              `json:"serial"`
	MAC          string                   `json:"mac"`
	Registered   string                   `json:"registered"`
	Vals         struct{}                 `json:"vals"`
	Aggr         []string                 `json:"aggr"`
	Values       map[string][]interface{} `json:"values"`
}

type Data struct {
	Dates []string      `json:"dates"`
	Data  []SensorValue `json:"data"`
}

func JSONRead(stationid string, stationdevice string) (string, map[string][]interface{}, interface{},string) {

	// Abra o arquivo JSON
	file, err := os.Open("stations/" + stationid + ".json")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo JSON:", err)
		return "", nil, nil, ""
	}
	defer file.Close()

	// Crie um decoder JSON para o arquivo
	decoder := json.NewDecoder(file)

	// Decodificar o JSON na estrutura de dados
	var data Data
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Erro ao decodificar o JSON:", err)
		return "", nil, nil, ""
	}

	// buscando data de atualização do json
	date := data.Dates[0]

	// busca do dispositivo desejado no campo data
	flag := false
	dispositivo := stationdevice
	for _, item := range data.Data {
		if item.Name == dispositivo {
			fmt.Println("Dispositivo " + stationdevice + " encontrado.")
			// fmt.Println("Name:", item.Name)
			// fmt.Println("Unit: ", item.Unit)
			// fmt.Println("Values: ", item.Values["avg"])
			// Se você quiser acessar outros campos, faça-o aqui
			flag = true
			return item.Name, item.Values, item.Unit, date
		}
	}

	if flag == false {
		fmt.Println("Dispositivo não encontrada")
		return "", nil, nil, ""
	}

	return "", nil, nil, ""

}

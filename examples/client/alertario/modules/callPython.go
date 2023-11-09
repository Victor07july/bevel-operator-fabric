package modules

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

// https://stackoverflow.com/questions/41415337/running-external-python-in-golang-catching-continuous-exec-command-stdout

// funcao de executar o python via terminal
func CallPy(stationid string) {
	fmt.Println("Python Call")

	// executa o comando via terminal dentro da pasta pymodules
	cmd := exec.Command("python3", "-u", "funcEstacao.py", stationid)
	cmd.Dir = ("./pyModules")

	// captura saida
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	// captura um possivel erro
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	// inicia o comando
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	// captura as saídas e espera a finalização
	go copyOutput(stdout)
	go copyOutput(stderr)
	cmd.Wait()
}

// função para capturar saídas
func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

// alternativa mais simples
func SimpleCallPy() {
	cmd := exec.Command("python", "-u", "funcEstacao.py")
	cmd.Dir = ("./pyModules")
	out, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}

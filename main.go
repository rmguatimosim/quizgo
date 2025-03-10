package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type GameState struct {
	Name      string
	Points    int
	Questions []Question
}

func (g *GameState) Init(ch chan int) {
	fmt.Println("Seja bem vindo(a) ao quiz!")
	fmt.Println("Neste quiz, você responderá 10 perguntas sobre o tema escolhido.")
	fmt.Println("Cada resposta correta valerá 10 pontos.")
	fmt.Println("Você precisará somar pelo menos 70 pontos para passar no teste.")
	fmt.Println("Escreva o seu nome: ")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		panic("Erro ao ler a string")
	}
	g.Name = name
	fmt.Printf("Vamos ao jogo, %s", g.Name)
	defineTema(ch)
	fmt.Printf("Lembre-se, você tem 10 segundos para responder cada pergunta.\n")
	fmt.Print("Pressione \"Enter\" para começar.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

}

func (g *GameState) ProcessCSV(ch chan int) {
	selecao := <-ch
	var tema string
	switch selecao {
	case 1:
		tema = "ingles.csv"
	case 2:
		tema = "historia.csv"
	case 3:
		tema = "conhecimentos.csv"
	}
	f, err := os.Open(tema)
	if err != nil {
		panic("Erro ao ler arquivo.")
	}
	defer f.Close()
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		panic("Erro ao ler o arquivo.")
	}

	for index, record := range records {
		if index > 0 {
			correta, _ := toInt(record[5])
			question := Question{
				Text:    record[0],
				Options: record[1:5],
				Answer:  correta,
			}
			g.Questions = append(g.Questions, question)
		}
	}
}

func (g *GameState) Run() {
	// Exibe a pergunta pro usuário
	for index, question := range g.Questions {
		fmt.Printf("\033[33m %d . %s \033[0m\n", index+1, question.Text)
		//  iterar sobre as opções que temos no game state
		// e exibir no terminal
		for j, option := range question.Options {
			fmt.Printf("[%d] %s\n", j+1, option)
		}
		fmt.Println("Digite o número da alternativa correta:")
		//Coleta a entrada do usuário
		//Valida o caractere inserido
		//Se for errado, usuário insere novamente.
		var answer int
		var err error
		alternativas := []int{1, 2, 3, 4}

		for {
			//timer que controla o tempo de resposta do usuário.
			//Após 10 segundos ele fecha o programa.
			timer := time.NewTimer(time.Second * 10)
			go func() {
				<-timer.C
				fmt.Println("Tempo esgotado.")
				os.Exit(1)
			}()

			reader := bufio.NewReader(os.Stdin)
			read, _ := reader.ReadString('\n')
			answer, err = toInt(read[:len(read)-1])
			if err != nil {
				fmt.Println(err.Error())
				timer.Stop()
				continue
			}
			if !slices.Contains(alternativas, answer) {
				fmt.Println("Alternativa inválida. Tente novamente.")
				timer.Stop()
				continue
			}
			timer.Stop()
			break

		}
		//Validar a resposta e exibir mensagem
		//calcular pontuação
		if answer == question.Answer {
			fmt.Println("Parabéns, você acertou!")
			g.Points += 10
		} else {
			fmt.Println("Resposta errada.")
			fmt.Println("-----------------")
		}
	}

}

func main() {
	ch := make(chan int)
	game := &GameState{}
	go game.ProcessCSV(ch)
	game.Init(ch)
	game.Run()
	checaPontuacao(game.Points)
}

func toInt(str string) (int, error) {
	str = strings.TrimSpace(str)
	i, erro := strconv.Atoi(str)

	if erro != nil {
		return 0, errors.New("opção inválida. Tente novamente")
	}
	return i, nil
}

func checaPontuacao(p int) {
	// função que verifica a pontuação do usuário e imprime no console o seu resultado
	fmt.Printf("Fim de jogo, você fez %d pontos\n", p)
	if p < 70 {
		fmt.Println("Pontuação insuficiente.")
	} else if p <= 100 {
		fmt.Println("Parabéns! Nota máxima!")
	} else {
		fmt.Println("Parabéns, você passou!")
	}
}

func defineTema(ch chan int) {
	fmt.Println("Escolha o tema das questões.")
	fmt.Println("[1] - Inglês")
	fmt.Println("[2] - História do Brasil")
	fmt.Println("[3] - Conhecimentos gerais")
	var selecao int
	var err error
	opcoes := []int{1, 2, 3}
	for {
		reader := bufio.NewReader(os.Stdin)
		read, _ := reader.ReadString('\n')
		selecao, err = toInt(read[:len(read)-1])
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if !slices.Contains(opcoes, selecao) {
			fmt.Println("Opção inválida.")
			continue
		}
		break
	}
	ch <- selecao
	close(ch)
}

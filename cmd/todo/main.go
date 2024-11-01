package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const toDoFile string = "todo.json"

type Tarefa struct {
	Titulo     string `json:"titulo"`
	Descricao  string `json:"descricao"`
	Prioridade int    `json:"prioridade"`
	Categoria  string `json:"categoria"`
	Finalizada bool   `json:"finalizada"`
}

func writeToDoFile(tarefas []Tarefa) {
	toDoFileWriter, err := os.OpenFile(toDoFile, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo \"%s\" para leitura: \"%s\"", toDoFile, err)
	}
	jsonEncoder := json.NewEncoder(toDoFileWriter)
	err = jsonEncoder.Encode(tarefas)
	if err != nil {
		log.Fatalf("Erro ao gravar no arquivo \"%s\": \"%s\"", toDoFile, err)
	}
	defer toDoFileWriter.Close()
}

func main() {

	mainCmd := flag.NewFlagSet("", flag.ExitOnError)
	add := mainCmd.Bool("add", false, "Adicionar tarefa.\n\tEx: --add")
	list := mainCmd.Bool("list", false, "Listar as tarefas.\n\tEx: --list")
	edit := mainCmd.Bool("edit", false, "Editar uma tarefa.\n\tEx: --edit")
	del := mainCmd.Bool("del", false, "Deletar uma tarefa.\n\tEx: --del")

	if len(os.Args) == 1 {
		fmt.Println("Usage:")
		mainCmd.PrintDefaults()
		return
	}

	mainCmd.Parse(os.Args[1:2])

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addTitulo := addCmd.String("titulo", "", "Título da tarefa.")
	addDescricao := addCmd.String("descricao", "default", "Descricao da tarefa.")
	addPrioridade := addCmd.Int("prioridade", 0, "Prioridade da tarefa.")
	addCategoria := addCmd.String("categoria", "default", "Categoria da tarefa.")
	addFinalizada := addCmd.Bool("finalizada", false, "Status da tarefa.")

	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	editIndice := editCmd.Int("indice", 0, "Indice da tarefa.")
	editTitulo := editCmd.String("titulo", "", "Título da tarefa.")
	editDescricao := editCmd.String("descricao", "", "Descricao da tarefa.")
	editPrioridade := editCmd.Int("prioridade", 0, "Prioridade da tarefa.")
	editCategoria := editCmd.String("categoria", "", "Categoria da tarefa.")
	editFinalizada := editCmd.Bool("finalizada", false, "Status da tarefa.")

	delCmd := flag.NewFlagSet("del", flag.ExitOnError)
	delIndice := delCmd.Int("indice", 0, "Indice da tarefa.")

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Erro ao obter o caminho do dirtório atual:", err)
	}

	_, err = os.Stat(toDoFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Criando o arquivo \"%s\" em: \"%s\"\n", toDoFile, currentDir)
			_, err := os.Create(toDoFile)
			if err != nil {
				log.Fatalf("Erro ao criar arquivo \"%s\": \"%s\"\n", toDoFile, err)
			}
			fmt.Printf("Arquivo \"%s/%s\" criado com sucesso! \n", currentDir, toDoFile)
		} else {
			log.Fatalf("Erro ao abrir o arquivo \"%s\": \"%s\"\n", toDoFile, err)
		}
	}

	toDoFileReader, err := os.Open(toDoFile)
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo \"%s\": \"%s\"\n", toDoFile, err)
	}
	defer toDoFileReader.Close()

	var tarefas []Tarefa

	jsonDecoder := json.NewDecoder(toDoFileReader)
	err = jsonDecoder.Decode(&tarefas)
	if err != nil && err != io.EOF {
		log.Fatalf("Erro ao carregar o arquivo \"%s\": \"%s\"\n", toDoFile, err)
	}

	if *add {
		addCmd.Parse(os.Args[2:])
		var tarefa Tarefa
		if *addTitulo != "" {
			tarefa.Titulo = *addTitulo
		} else {
			log.Fatal("O título da tarefa é obrigatótio. Ex: --add --titulo=\"alguma tarefa\"")
		}
		tarefa.Descricao = *addDescricao
		tarefa.Prioridade = *addPrioridade
		tarefa.Categoria = *addCategoria
		if *addFinalizada {
			tarefa.Finalizada = true
		}
		tarefas = append(tarefas, tarefa)
		writeToDoFile(tarefas)
	}

	if *list {
		if len(tarefas) == 0 {
			log.Fatalf("Nenhuma tarefa encontrada no arquivo \"%s/%s\".\n", currentDir, toDoFile)
		}
		fmt.Printf("Lista de tarefas:\n")
		for i, tarefa := range tarefas {
			fmt.Printf("%d - %s, %s, %d, %s, %t\n", i+1, tarefa.Titulo, tarefa.Descricao, tarefa.Prioridade, tarefa.Categoria, tarefa.Finalizada)
		}
	}

	if *edit {
		editCmd.Parse(os.Args[2:])
		if *editIndice != 0 {
			if *editIndice > len(tarefas) {
				log.Fatal("O indice informado é maior do que o total de tarefas.")
			}
			indiceEncontrado := false
			for i := range tarefas {
				if i+1 == *editIndice {
					indiceEncontrado = true
					if *editTitulo != "" {
						tarefas[i].Titulo = *editTitulo
					}
					if *editDescricao != "" {
						tarefas[i].Descricao = *editDescricao
					}
					if *editPrioridade != 0 {
						tarefas[i].Prioridade = *editPrioridade
					}
					if *editCategoria != "" {
						tarefas[i].Categoria = *editCategoria
					}
					if *editFinalizada {
						tarefas[i].Finalizada = *editFinalizada
					}
					writeToDoFile(tarefas)
					break
				}
			}
			if !indiceEncontrado {
				log.Fatal("Nenhuma tarefa possui o indice fornecido.")
			}
		} else {
			log.Fatal("O indice da tarefa é obrigatótio. Ex: --edit --indice=\"123\"")
		}
	}

	if *del {
		delCmd.Parse(os.Args[2:])
		if *delIndice != 0 {
			if *delIndice > len(tarefas) {
				log.Fatal("O indice informado é maior do que o total de tarefas.")
			}
			indiceEncontrado := false
			for i := range tarefas {
				if i+1 == *delIndice {
					indiceEncontrado = true
					tarefas = append(tarefas[:i], tarefas[i+1:]...)
					writeToDoFile(tarefas)
					break
				}
			}
			if !indiceEncontrado {
				log.Fatal("Nenhuma tarefa possui o indice fornecido.")
			}
		} else {
			log.Fatal("O indice da tarefa é obrigatótio. Ex: --del --indice=\"123\"")
		}
	}
}

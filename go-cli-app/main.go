package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/prisma/prisma-examples/go-cli-app/prisma-client"
)

func main() {
	db := prisma.DB{
		// Debug: true,
	}

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		todos := allTodos(db)
		printTodos(todos)
		os.Exit(0)
	}
	command := argsWithoutProg[0]

	if len(argsWithoutProg) == 1 {
		if command == "count" {
			fmt.Println("Total Todos: ", len(allTodos(db)))
		} else {
			fmt.Println("Invalid command: ", command, " printing todos")
			todos := allTodos(db)
			printTodos(todos)
		}
		os.Exit(0)
	}

	if len(argsWithoutProg) == 2 {
		v1 := argsWithoutProg[1]
		if command == "create" {
			todo := createTodo(db, v1)
			printTodo(todo, nil)
		} else if command == "delete" {
			todo := deleteTodo(db, v1)
			printTodo(todo, nil)
		} else if command == "get" {
			todo := getTodo(db, v1)
			printTodo(todo, nil)
		} else if command == "get-user" {
			user := getTodoUser(db, v1)
			printUser(user, nil)
		} else if command == "list" {
			v1Int, err := strconv.ParseInt(v1, 10, 32)
			if err != nil {
				panic(err)
			}
			todos := someTodos(db, int32(v1Int))
			printTodos(todos)
		} else if command == "search" {
			todos := searchTodos(db, v1)
			printTodos(todos)
		} else {
			fmt.Println("Invalid command: ", command, " printing todos")
			allTodos(db)
		}
		os.Exit(0)
	}
	fmt.Println("Invalid command: ", command, " printing todos")
	todos := allTodos(db)
	printTodos(todos)
}

func printUser(user prisma.User, key *int) {
	if key != nil {
		fmt.Println("User #", *key, ": ", user.Name, " - ", user.ID)
	} else {
		fmt.Println("User ", user.Name, " - ", user.ID)
	}
}

func printTodo(todo prisma.Todo, key *int) {
	if key != nil {
		fmt.Println("Todo #", *key, ": ", todo.Text, " - ", todo.ID)
	} else {
		fmt.Println("Todo ", todo.Text, " - ", todo.ID)
	}
}

func printTodos(todos []prisma.Todo) {
	for key, todo := range todos {
		printTodo(todo, &key)
	}
}

func allTodos(db prisma.DB) []prisma.Todo {
	fmt.Println("All Todos:")
	todos := db.Todos(&prisma.TodosParams{}).Exec()
	return todos
}

func someTodos(db prisma.DB, first int32) []prisma.Todo {
	fmt.Println("Some Todos:", first)
	todos := db.Todos(&prisma.TodosParams{
		First: &first,
	}).Exec()
	return todos
}

func searchTodos(db prisma.DB, q string) []prisma.Todo {
	fmt.Println("Search Todos:", q)
	orderBy := prisma.CreatedAtDescTodoOrderByInput
	todos := db.Todos(&prisma.TodosParams{
		Where: &prisma.TodoWhereInput{
			TextContains: &q,
		},
		OrderBy: &orderBy,
	}).Exec()
	return todos
}

func createTodo(db prisma.DB, text string) prisma.Todo {
	fmt.Println("Create Todo")
	done := false
	userID := "cjlhytekx005708701pice3uj"
	todo := db.CreateTodo(&prisma.TodoCreateInput{
		Done: &done,
		Text: &text,
		User: &prisma.UserCreateOneInput{
			Connect: &prisma.UserWhereUniqueInput{
				ID: &userID,
			},
		},
	}).Exec()
	return todo
}

func deleteTodo(db prisma.DB, id string) prisma.Todo {
	fmt.Println("Create Todo")
	todo := db.DeleteTodo(&prisma.TodoWhereUniqueInput{
		ID: &id,
	}).Exec()
	return todo
}

func getTodo(db prisma.DB, id string) prisma.Todo {
	fmt.Println("Get Todo")

	exists := db.Exists.Todo(&prisma.TodoWhereUniqueInput{
		ID: &id,
	})
	if exists {
		fmt.Println("Todo exists")
	} else {
		fmt.Println("Todo dos not exist")
	}

	todo := db.Todo(&prisma.TodoWhereUniqueInput{
		ID: &id,
	}).Exec()
	return todo
}

func getTodoUser(db prisma.DB, id string) prisma.User {
	fmt.Println("Get Todo User")

	exists := db.Exists.Todo(&prisma.TodoWhereUniqueInput{
		ID: &id,
	})
	if exists {
		fmt.Println("Todo exists")
	} else {
		fmt.Println("Todo dos not exist")
	}

	todo := db.Todo(&prisma.TodoWhereUniqueInput{
		ID: &id,
	}).User(&prisma.UserWhereInput{}).Exec()
	return todo
}

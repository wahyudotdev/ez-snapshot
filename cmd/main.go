package main

import (
	"context"
	"ez-snapshot/internal/deps"
	"ez-snapshot/internal/usecase"
	"fmt"

	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
)

type Command struct {
	Name string
	Run  func(ctx context.Context) error
}

func main() {
	ctx := context.Background()
	// define available commands
	commands := []Command{
		{
			Name: "backup",
			Run: func(ctx context.Context) error {
				fmt.Println("Running database backup...")
				uc := usecase.NewBackupDatabaseUseCase(
					deps.NewBackupRepo(ctx),
					deps.NewStorageRepo(ctx),
				)
				err := uc.Execute(ctx)
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name: "restore",
			Run: func(ctx context.Context) error {
				fmt.Println("Running database restore...")
				return nil
			},
		},
		{
			Name: "exit",
			Run: func(ctx context.Context) error {
				fmt.Println("Bye 👋")
				return fmt.Errorf("exit")
			},
		},
		{
			Name: "list",
			Run: func(ctx context.Context) error {
				fmt.Println("Listing backups...")
				uc := usecase.NewListDatabaseUseCase(deps.NewStorageRepo(ctx))
				list, err := uc.Execute(ctx)
				if err != nil {
					return err
				}

				if len(list) == 0 {
					fmt.Println("No backup(s) found")
					return nil
				}

				for i, d := range list {
					fmt.Printf("[%d]: %s\n", i, d.Name)
				}

				return nil
			},
		},
	}

	commandMap := make(map[string]Command)
	for _, c := range commands {
		commandMap[c.Name] = c
	}

	completer := func(d prompt.Document) []prompt.Suggest {
		s := []prompt.Suggest{}
		for _, c := range commands {
			s = append(s, prompt.Suggest{Text: c.Name})
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}

	fmt.Println("Welcome to MySQL Backup CLI (type 'exit' to quit)")
	for {
		input := prompt.Input("> ", completer)

		if cmd, ok := commandMap[input]; ok {
			err := cmd.Run(ctx)
			if err != nil {
				if err.Error() == "exit" {
					break
				}
				log.Error(err)
			}
		} else {
			fmt.Println("Unknown command:", input)
		}
	}
}

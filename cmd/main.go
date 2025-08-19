package main

import (
	"context"
	"ez-snapshot/internal/deps"
	"ez-snapshot/internal/usecase"
	"fmt"
	"os"
	"strconv"

	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
)

type Command struct {
	Name        string
	Description string
	Run         func(ctx context.Context) error
}

func main() {
	ctx := context.Background()
	depUc := usecase.NewDependencyChecker(deps.NewStorageRepo(ctx))
	if err := depUc.Check(); err != nil {
		log.Fatal(err)
	}

	// define available commands
	commands := []Command{
		{
			Name:        "backup",
			Description: "Create a new database backup",
			Run: func(ctx context.Context) error {
				fmt.Println("Running database backup...")
				uc := usecase.NewBackupDatabaseUseCase(
					deps.NewBackupRepo(ctx),
					deps.NewStorageRepo(ctx),
				)
				return uc.Execute(ctx)
			},
		},
		{
			Name:        "restore",
			Description: "Restore database from a selected backup",
			Run: func(ctx context.Context) error {
				fmt.Println("Listing backups...")
				listDbUc := usecase.NewListDatabaseUseCase(deps.NewStorageRepo(ctx))
				list, err := listDbUc.Execute(ctx)
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

				completer := func(d prompt.Document) []prompt.Suggest {
					var s []prompt.Suggest
					for i := range list {
						s = append(s, prompt.Suggest{Text: strconv.Itoa(i)})
					}
					return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
				}

				input := prompt.Input("Select backup number >", completer)

				index, err := strconv.Atoi(input)
				if err != nil {
					return err
				}

				if index < 0 || index >= len(list) {
					return fmt.Errorf("invalid backup number")
				}

				backupKey := list[index].Path

				uc := usecase.NewRestoreDatabaseUseCase(deps.NewBackupRepo(ctx), deps.NewStorageRepo(ctx))
				return uc.Execute(ctx, backupKey)
			},
		},
		{
			Name:        "list",
			Description: "List available backups",
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
		{
			Name:        "help",
			Description: "Show help message",
			Run: func(ctx context.Context) error {
				printHelp()
				return nil
			},
		},
		{
			Name:        "exit",
			Description: "Exit the CLI",
			Run: func(ctx context.Context) error {
				fmt.Println("Bye 👋")
				return fmt.Errorf("exit")
			},
		},
	}

	commandMap := make(map[string]Command)
	for _, c := range commands {
		commandMap[c.Name] = c
	}

	// 🔹 1. Check if arguments provided
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if len(arg) > 2 && arg[:2] == "--" {
			arg = arg[2:] // remove `--`
		}

		if cmd, ok := commandMap[arg]; ok {
			if err := cmd.Run(ctx); err != nil {
				if err.Error() == "exit" {
					os.Exit(0)
				}
				log.Error(err)
				os.Exit(1)
			}
			return
		} else {
			fmt.Println("Unknown command:", os.Args[1])
			printHelp()
			os.Exit(1)
		}
	}

	// 🔹 2. If no arguments → start interactive mode
	completer := func(d prompt.Document) []prompt.Suggest {
		var s []prompt.Suggest
		for _, c := range commands {
			s = append(s, prompt.Suggest{Text: c.Name, Description: c.Description})
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}

	fmt.Println("Welcome EZ-Snapshot CLI (type 'exit' to quit)")
	printHelp()
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
			printHelp()
		}
	}
}

func printHelp() {
	fmt.Println("\nUsage:")
	fmt.Println("  ez-snapshot --<command>\n")
	fmt.Println("Available commands:")
	fmt.Println("  --backup     Create a new database backup")
	fmt.Println("  --restore    Restore database from a selected backup")
	fmt.Println("  --list       List available backups")
	fmt.Println("  --help       Show this help message")
	fmt.Println("  --exit       Exit the CLI (interactive mode only)")
	fmt.Println()
}

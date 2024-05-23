package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "listdir",
	Short: "Listdir is a CLI tool to list files and directories",
	Long:  `Listdir is a CLI tool to list files and directories, with options to ignore specific directories.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		// Preguntar por el directorio a listar
		fmt.Print("Enter the directory to list (default is current directory): ")
		dir, _ := reader.ReadString('\n')
		dir = strings.TrimSpace(dir)
		if dir == "" {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error getting current directory: %v\n", err)
				return
			}
			dir = cwd
		}

		// Preguntar por los directorios a ignorar
		fmt.Print("Enter directories to ignore, separated by commas (default is .git,node_modules,vendor,.idea,.vsc): ")
		ignore, _ := reader.ReadString('\n')
		ignore = strings.TrimSpace(ignore)
		if ignore == "" {
			ignore = ".git,node_modules,vendor,.idea,.vsc"
		}

		// Convertir la lista de directorios a ignorar en un mapa
		ignoredDirs := make(map[string]bool)
		for _, d := range strings.Split(ignore, ",") {
			ignoredDirs[d] = true
		}

		fmt.Println(dir + "/")
		err := listFiles(dir, "", ignoredDirs)
		if err != nil {
			fmt.Printf("Error listing files: %v\n", err)
		}
	},
}

// listFiles recorre el directorio y lista todos los archivos y subdirectorios.
func listFiles(dir string, prefix string, ignoredDirs map[string]bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if ignoredDirs[entry.Name()] {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			return err
		}

		fmt.Printf("%s|-- %s\n", prefix, entry.Name())

		if info.IsDir() {
			newPrefix := prefix + "|   "
			_ = listFiles(path, newPrefix, ignoredDirs)
		}
	}
	return nil
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Aquí puedes agregar flags y configuración global si es necesario
}

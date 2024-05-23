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

		// Crear archivo para escribir la salida
		outputFile, err := os.Create("list_dir_output.txt")
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			return
		}
		defer func(outputFile *os.File) {
			_ = outputFile.Close()
		}(outputFile)
		writer := bufio.NewWriter(outputFile)

		fmt.Println(dir + "/")
		_, _ = writer.WriteString(dir + "/\n")

		err = listFiles(dir, "", ignoredDirs, writer)
		if err != nil {
			fmt.Printf("Error listing files: %v\n", err)
		}

		_ = writer.Flush()
	},
}

// listFiles recorre el directorio y lista todos los archivos y subdirectorios.
func listFiles(dir string, prefix string, ignoredDirs map[string]bool, writer *bufio.Writer) error {
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

		output := fmt.Sprintf("%s|-- %s\n", prefix, entry.Name())
		fmt.Print(output)
		_, _ = writer.WriteString(output)

		if info.IsDir() {
			newPrefix := prefix + "|   "
			_ = listFiles(path, newPrefix, ignoredDirs, writer)
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

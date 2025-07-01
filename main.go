package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	blue  = "\033[34m"
	green = "\033[32m"
	reset = "\033[0m"
)

const (
	symbolContinue = "├── "
	symbolFinal    = "└── "
	symbolBar      = "│   "
	symbolSpace    = "    "
)

// Para contar los archivos y directorios
type counts struct {
	files int
	dirs  int
}

// Comprueba el numero que devuelve info.Mode() en octal
// con 0111 que seria permisos de ejecución del dueño, grupo y otros
// Si algún bit coincide es que tiene algún permiso de ejecución
// Owner -> 0100 (64)
// Group -> 0010 (8)
// Other -> 0001 (1)
// 0100 + 0010 + 0111 = 0111
// Se puede cambiar el numero para solo comprobar permisos concretos
// La o es para indicar que es octal (desde go 1.13)
func isExecutable(mode os.FileMode) bool {
	return mode&0o111 != 0
}

func readFiles(pathName string, prefix string, level int, count *counts) error {
	files, err := os.ReadDir(pathName)
	if err != nil {
		return err
	}

	for i, file := range files {
		name := file.Name()
		fullPath := filepath.Join(pathName, name)

		isLast := i == len(files)-1
		// Conector por defecto
		connector := symbolContinue
		if isLast {
			// Conector si es el ultimo
			connector = symbolFinal
		}

		if file.IsDir() {
			count.dirs++
			fmt.Printf("%s%s%s%s%s\n", prefix, connector, blue, name, reset)

			// Preparemos el prefijo para la siguiente llamada recursiva
			// Conservamos la indentación porque sumamos los símbolos al prefijo anterior
			// no lo sobrescribimos
			newPrefix := prefix
			if isLast {
				// Si es el ultimo solo ponemos espacios
				newPrefix += symbolSpace
			} else {
				// Si no agregamos la barra al principio
				newPrefix += symbolBar
			}

			// Volvemos a llamar a la función para que recorra el directorio
			if err := readFiles(fullPath, newPrefix, level+1, count); err != nil {
				return err
			}
			continue
		}
		count.files++

		// Obtener la información (metadatos) de un archivo para comprobar si es ejecutable
		info, err := os.Stat(fullPath)
		if err != nil {
			fmt.Printf("%s%s%s (stat error: %v)\n", prefix, connector, name, err)
			continue
		}

		// Mode devuelve los permisos del archivo en bits
		if isExecutable(info.Mode()) {
			fmt.Printf("%s%s%s%s%s\n", prefix, connector, green, name, reset)
		} else {
			fmt.Printf("%s%s%s\n", prefix, connector, name)
		}
	}
	return nil
}

func main() {
	// Redefinimos el comportamiento de la bandera -h (help)
	// Para que muestre un mensaje de uso personalizado
	flag.Usage = func() {
		fmt.Printf("Usage: gotree [directory]\n\nLists files like Unix tree.\n")
		flag.PrintDefaults()
	}
	// Obtiene las banderas y los argumentos de la linea de comandos
	flag.Parse()

	path := "."
	// NArg es para ver los argumentos que no sean bandera (-h, -v, etc)
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	fmt.Printf("%s%s%s\n", blue, path, reset)

	var c counts
	if err := readFiles(path, "", 0, &c); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n%d directories, %d files\n", c.dirs+1, c.files)
}

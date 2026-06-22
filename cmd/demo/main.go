// Demo de consola del Radix Trie: carga unas palabras y prueba el
// autocompletado. Es solo para verificar que la estructura funciona;
// la app con base de datos (Entregable 3) y la simulación web
// (Entregable 4) se construirán encima del paquete trie.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/JulioCarrizales/Proyecto_AED_7/trie"
)

func main() {
	t := trie.New()
	for _, w := range []string{
		"casa", "casaca", "caso", "castor", "catedral",
		"perro", "pera", "perla", "persona",
		"sol", "solar", "soledad",
	} {
		t.Insert(w)
	}

	fmt.Printf("Radix Trie cargado con %d palabras.\n", t.Len())
	fmt.Println("Escribe un prefijo para autocompletar (vacío o Ctrl+C para salir):")

	sc := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for sc.Scan() {
		prefijo := strings.TrimSpace(sc.Text())
		if prefijo == "" {
			break
		}
		sugerencias := t.Autocomplete(prefijo)
		if len(sugerencias) == 0 {
			fmt.Println("  (sin coincidencias)")
		} else {
			fmt.Printf("  %s\n", strings.Join(sugerencias, ", "))
		}
		fmt.Print("> ")
	}
}

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
	// Se usa como conjunto de palabras: el valor asociado es struct{} (vacío).
	t := trie.New[struct{}]()
	for _, w := range []string{
		"casa", "casaca", "caso", "castor", "catedral",
		"perro", "pera", "perla", "persona",
		"sol", "solar", "soledad",
	} {
		t.Insert(w, struct{}{})
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
			palabras := make([]string, 0, len(sugerencias))
			for _, m := range sugerencias {
				palabras = append(palabras, m.Key)
			}
			fmt.Printf("  %s\n", strings.Join(palabras, ", "))
		}
		fmt.Print("> ")
	}
}

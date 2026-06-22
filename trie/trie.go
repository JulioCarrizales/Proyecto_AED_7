// Package trie implementa un Radix Trie (árbol PATRICIA) desde cero.
//
// Referencia: Morrison (1968), "PATRICIA — Practical Algorithm To Retrieve
// Information Coded In Alphanumeric", JACM.
//
// Un Radix Trie es un trie comprimido: las cadenas de nodos con un solo hijo
// se fusionan en una sola arista etiquetada con varios caracteres. Esto reduce
// el consumo de memoria y acelera las búsquedas frente a un trie clásico.
//
// Caso de uso del proyecto: autocompletado de palabras por prefijo.
package trie

import (
	"sort"
	"strings"
)

// node es un nodo del Radix Trie.
type node struct {
	label    string           // segmento de arista que llega a este nodo desde su padre
	isWord   bool             // true si una palabra termina exactamente en este nodo
	children map[byte]*node   // hijos indexados por el primer byte de su label
}

func newNode(label string) *node {
	return &node{label: label, children: make(map[byte]*node)}
}

// Trie es la estructura principal. La raíz tiene label vacío.
type Trie struct {
	root *node
	size int // cantidad de palabras almacenadas
}

// New crea un Radix Trie vacío.
func New() *Trie {
	return &Trie{root: newNode("")}
}

// Len devuelve cuántas palabras hay almacenadas.
func (t *Trie) Len() int { return t.size }

// Insert agrega una palabra al trie. Si ya existe, no hace nada.
// Complejidad: O(L) donde L es la longitud de la palabra.
func (t *Trie) Insert(word string) {
	cur := t.root
	rem := word
	for {
		if rem == "" {
			if !cur.isWord {
				cur.isWord = true
				t.size++
			}
			return
		}
		first := rem[0]
		child, ok := cur.children[first]
		if !ok {
			// No hay arista que empiece con este carácter: creamos una hoja.
			leaf := newNode(rem)
			leaf.isWord = true
			cur.children[first] = leaf
			t.size++
			return
		}

		cpl := commonPrefixLen(child.label, rem)
		if cpl == len(child.label) {
			// La arista coincide por completo: descendemos y seguimos.
			cur = child
			rem = rem[cpl:]
			continue
		}

		// Coincidencia parcial: hay que partir la arista del hijo.
		// Ej.: child.label = "team", rem = "tea" -> partir en "tea" + "m".
		mid := newNode(child.label[:cpl])
		child.label = child.label[cpl:]
		mid.children[child.label[0]] = child
		cur.children[first] = mid

		tail := rem[cpl:]
		if tail == "" {
			mid.isWord = true
		} else {
			leaf := newNode(tail)
			leaf.isWord = true
			mid.children[tail[0]] = leaf
		}
		t.size++
		return
	}
}

// Search indica si la palabra exacta está almacenada.
// Complejidad: O(L).
func (t *Trie) Search(word string) bool {
	cur := t.root
	rem := word
	for rem != "" {
		child, ok := cur.children[rem[0]]
		if !ok || !strings.HasPrefix(rem, child.label) {
			return false
		}
		rem = rem[len(child.label):]
		cur = child
	}
	return cur.isWord
}

// Autocomplete devuelve todas las palabras que empiezan con prefix,
// ordenadas alfabéticamente. Es la operación clave de la demo del proyecto.
// Complejidad: O(P + K) donde P es la longitud del prefijo y K el número de
// caracteres recorridos para reunir las coincidencias.
func (t *Trie) Autocomplete(prefix string) []string {
	cur := t.root
	path := "" // cadena desde la raíz hasta cur (incluye su label)
	rem := prefix
	for rem != "" {
		child, ok := cur.children[rem[0]]
		if !ok {
			return nil
		}
		if strings.HasPrefix(child.label, rem) {
			// El prefijo termina dentro de esta arista: todo el subárbol completa.
			out := collect(child, path+child.label, nil)
			sort.Strings(out)
			return out
		}
		if strings.HasPrefix(rem, child.label) {
			path += child.label
			rem = rem[len(child.label):]
			cur = child
			continue
		}
		return nil
	}
	out := collect(cur, path, nil)
	sort.Strings(out)
	return out
}

// Keys devuelve todas las palabras del trie ordenadas alfabéticamente.
func (t *Trie) Keys() []string {
	out := collect(t.root, "", nil)
	sort.Strings(out)
	return out
}

// collect recorre el subárbol acumulando las palabras. acc es la cadena
// completa desde la raíz hasta n (incluyendo el label de n).
func collect(n *node, acc string, out []string) []string {
	if n.isWord {
		out = append(out, acc)
	}
	for _, c := range n.children {
		out = collect(c, acc+c.label, out)
	}
	return out
}

// commonPrefixLen devuelve la longitud del prefijo común entre a y b.
func commonPrefixLen(a, b string) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	i := 0
	for i < n && a[i] == b[i] {
		i++
	}
	return i
}

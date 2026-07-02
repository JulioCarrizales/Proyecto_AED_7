// Package trie implementa un Radix Trie (árbol PATRICIA) genérico desde cero.
//
// Referencia: Morrison (1968), "PATRICIA — Practical Algorithm To Retrieve
// Information Coded In Alphanumeric", JACM.
//
// Un Radix Trie es un trie comprimido: las cadenas de nodos con un solo hijo
// se fusionan en una sola arista etiquetada con varios caracteres. Esto reduce
// el consumo de memoria y acelera las búsquedas frente a un trie clásico.
//
// El trie es genérico: cada clave (string) lleva asociado un valor de tipo V.
// Así puede usarse como un mapa ordenado por prefijos (p. ej. ciudad ->
// departamento) o como un simple conjunto usando Trie[struct{}].
package trie

import (
	"sort"
	"strings"
)

// node es un nodo del Radix Trie. V es el tipo del valor asociado a la clave.
type node[V any] struct {
	label    string            // segmento de arista que llega a este nodo desde su padre
	isWord   bool              // true si una clave termina exactamente en este nodo
	value    V                 // valor asociado (solo válido si isWord)
	children map[byte]*node[V] // hijos indexados por el primer byte de su label
}

func newNode[V any](label string) *node[V] {
	return &node[V]{label: label, children: make(map[byte]*node[V])}
}

// Trie es la estructura principal. La raíz tiene label vacío.
type Trie[V any] struct {
	root *node[V]
	size int // cantidad de claves almacenadas
}

// Match es una coincidencia devuelta por Autocomplete: una clave y su valor.
type Match[V any] struct {
	Key   string
	Value V
}

// New crea un Radix Trie vacío cuyas claves llevan un valor de tipo V.
func New[V any]() *Trie[V] {
	return &Trie[V]{root: newNode[V]("")}
}

// Len devuelve cuántas claves hay almacenadas.
func (t *Trie[V]) Len() int { return t.size }

// Insert agrega (o actualiza) una clave con su valor asociado.
// Complejidad: O(L) donde L es la longitud de la clave.
func (t *Trie[V]) Insert(key string, value V) {
	cur := t.root
	rem := key
	for {
		if rem == "" {
			if !cur.isWord {
				cur.isWord = true
				t.size++
			}
			cur.value = value
			return
		}
		first := rem[0]
		child, ok := cur.children[first]
		if !ok {
			// No hay arista que empiece con este carácter: creamos una hoja.
			leaf := newNode[V](rem)
			leaf.isWord = true
			leaf.value = value
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
		mid := newNode[V](child.label[:cpl])
		child.label = child.label[cpl:]
		mid.children[child.label[0]] = child
		cur.children[first] = mid

		tail := rem[cpl:]
		if tail == "" {
			mid.isWord = true
			mid.value = value
		} else {
			leaf := newNode[V](tail)
			leaf.isWord = true
			leaf.value = value
			mid.children[tail[0]] = leaf
		}
		t.size++
		return
	}
}

// Search busca una clave exacta. Devuelve su valor y true si existe, o el valor
// cero de V y false si no.
// Complejidad: O(L).
func (t *Trie[V]) Search(key string) (V, bool) {
	cur := t.root
	rem := key
	for rem != "" {
		child, ok := cur.children[rem[0]]
		if !ok || !strings.HasPrefix(rem, child.label) {
			var cero V
			return cero, false
		}
		rem = rem[len(child.label):]
		cur = child
	}
	if cur.isWord {
		return cur.value, true
	}
	var cero V
	return cero, false
}

// Delete elimina una clave del trie. Devuelve true si existía (y se eliminó) o
// false si no estaba. Tras eliminar, el árbol se re-comprime.
// Complejidad: O(L).
func (t *Trie[V]) Delete(key string) bool {
	if deleteRec(t.root, key) {
		t.size--
		return true
	}
	return false
}

func deleteRec[V any](n *node[V], rem string) bool {
	if rem == "" {
		if !n.isWord {
			return false
		}
		n.isWord = false
		var cero V
		n.value = cero
		return true
	}

	child, ok := n.children[rem[0]]
	if !ok || !strings.HasPrefix(rem, child.label) {
		return false
	}
	if !deleteRec(child, rem[len(child.label):]) {
		return false
	}

	switch {
	case !child.isWord && len(child.children) == 0:
		// Hoja que ya no es clave: se elimina.
		delete(n.children, rem[0])
	case !child.isWord && len(child.children) == 1:
		// Nodo intermedio no-clave con un solo hijo: se fusiona con él
		// para mantener el árbol comprimido.
		var unico *node[V]
		for _, c := range child.children {
			unico = c
		}
		child.label += unico.label
		child.isWord = unico.isWord
		child.value = unico.value
		child.children = unico.children
	}
	return true
}

// Autocomplete devuelve todas las coincidencias cuya clave empieza con prefix,
// ordenadas alfabéticamente por clave. Es la operación clave de la demo.
// Complejidad: O(P + K) donde P es la longitud del prefijo y K el número de
// caracteres recorridos para reunir las coincidencias.
func (t *Trie[V]) Autocomplete(prefix string) []Match[V] {
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
			sortMatches(out)
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
	sortMatches(out)
	return out
}

// Keys devuelve todas las claves del trie ordenadas alfabéticamente.
func (t *Trie[V]) Keys() []string {
	matches := collect(t.root, "", nil)
	claves := make([]string, 0, len(matches))
	for _, m := range matches {
		claves = append(claves, m.Key)
	}
	sort.Strings(claves)
	return claves
}

// collect recorre el subárbol acumulando las coincidencias. acc es la cadena
// completa desde la raíz hasta n (incluyendo el label de n).
func collect[V any](n *node[V], acc string, out []Match[V]) []Match[V] {
	if n.isWord {
		out = append(out, Match[V]{Key: acc, Value: n.value})
	}
	for _, c := range n.children {
		out = collect(c, acc+c.label, out)
	}
	return out
}

func sortMatches[V any](m []Match[V]) {
	sort.Slice(m, func(i, j int) bool { return m[i].Key < m[j].Key })
}

// NodeView es la representación exportable de un nodo, pensada para serializar
// el árbol a JSON (por ejemplo, para la visualización del Entregable 4).
type NodeView struct {
	Label    string      `json:"label"`
	IsWord   bool        `json:"isWord"`
	Children []*NodeView `json:"children"`
}

// View devuelve una copia serializable del árbol completo, con los hijos
// ordenados por su etiqueta para obtener una representación estable.
func (t *Trie[V]) View() *NodeView {
	return viewNode(t.root)
}

func viewNode[V any](n *node[V]) *NodeView {
	v := &NodeView{Label: n.label, IsWord: n.isWord, Children: []*NodeView{}}

	etiquetas := make([]string, 0, len(n.children))
	porEtiqueta := make(map[string]*node[V], len(n.children))
	for _, c := range n.children {
		etiquetas = append(etiquetas, c.label)
		porEtiqueta[c.label] = c
	}
	sort.Strings(etiquetas)
	for _, e := range etiquetas {
		v.Children = append(v.Children, viewNode(porEtiqueta[e]))
	}
	return v
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

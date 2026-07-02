package trie

import (
	"sort"
	"strings"
)

// node es un nodo del Radix Trie.
type node struct {
	label    string         // segmento de arista que llega a este nodo desde su padre
	isWord   bool           // true si una palabra termina exactamente en este nodo
	children map[byte]*node // hijos indexados por el primer byte de su label
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

// Delete elimina una palabra del trie. Devuelve true si la palabra existía
// (y por lo tanto se eliminó) o false si no estaba.
// Tras eliminar, el árbol se re-comprime: los nodos que quedan sin uso se
// borran y los que quedan con un solo hijo se fusionan.
// Complejidad: O(L).
func (t *Trie) Delete(word string) bool {
	if deleteRec(t.root, word) {
		t.size--
		return true
	}
	return false
}

// deleteRec elimina rem a partir del nodo n y, al volver, compacta el hijo
// afectado (lo borra si quedó inútil o lo fusiona si quedó con un solo hijo).
func deleteRec(n *node, rem string) bool {
	if rem == "" {
		if !n.isWord {
			return false
		}
		n.isWord = false
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
		// Hoja que ya no es palabra: se elimina.
		delete(n.children, rem[0])
	case !child.isWord && len(child.children) == 1:
		// Nodo intermedio no-palabra con un solo hijo: se fusiona con él
		// para mantener el árbol comprimido.
		var unico *node
		for _, c := range child.children {
			unico = c
		}
		child.label += unico.label
		child.isWord = unico.isWord
		child.children = unico.children
	}
	return true
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

// NodeView es la representación exportable de un nodo, pensada para serializar
// el árbol a JSON (por ejemplo, para la visualización del Entregable 4).
type NodeView struct {
	Label    string      `json:"label"`
	IsWord   bool        `json:"isWord"`
	Children []*NodeView `json:"children"`
}

// View devuelve una copia serializable del árbol completo, con los hijos
// ordenados por su etiqueta para obtener una representación estable.
func (t *Trie) View() *NodeView {
	return viewNode(t.root)
}

func viewNode(n *node) *NodeView {
	v := &NodeView{Label: n.label, IsWord: n.isWord, Children: []*NodeView{}}

	etiquetas := make([]string, 0, len(n.children))
	porEtiqueta := make(map[string]*node, len(n.children))
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

package trie

import (
	"reflect"
	"testing"
)

// claves extrae solo las claves de una lista de coincidencias, para comparar.
func claves[V any](ms []Match[V]) []string {
	out := make([]string, 0, len(ms))
	for _, m := range ms {
		out = append(out, m.Key)
	}
	return out
}

func TestInsertSearch(t *testing.T) {
	tr := New[int]()
	palabras := []string{"casa", "casaca", "caso", "perro", "pera"}
	for i, w := range palabras {
		tr.Insert(w, i) // el valor asociado es el índice
	}

	for i, w := range palabras {
		v, ok := tr.Search(w)
		if !ok {
			t.Errorf("Search(%q) = _, false; se esperaba true", w)
		}
		if v != i {
			t.Errorf("Search(%q) valor = %d; se esperaba %d", w, v, i)
		}
	}

	for _, w := range []string{"cas", "ca", "per", "gato", ""} {
		if _, ok := tr.Search(w); ok {
			t.Errorf("Search(%q) = _, true; se esperaba false", w)
		}
	}

	if tr.Len() != len(palabras) {
		t.Errorf("Len() = %d; se esperaba %d", tr.Len(), len(palabras))
	}
}

func TestValorGenerico(t *testing.T) {
	// El trie como mapa ciudad -> departamento.
	tr := New[string]()
	tr.Insert("Cusco", "Cusco")
	tr.Insert("Miraflores", "Lima")

	if dep, ok := tr.Search("Miraflores"); !ok || dep != "Lima" {
		t.Errorf("Search(\"Miraflores\") = %q, %v; se esperaba \"Lima\", true", dep, ok)
	}
	// Valor cero cuando no existe.
	if dep, ok := tr.Search("Tokio"); ok || dep != "" {
		t.Errorf("Search(\"Tokio\") = %q, %v; se esperaba \"\", false", dep, ok)
	}
}

func TestInsertDuplicado(t *testing.T) {
	tr := New[int]()
	tr.Insert("hola", 1)
	tr.Insert("hola", 2) // mismo clave: actualiza el valor, no crece
	if tr.Len() != 1 {
		t.Errorf("Len() = %d tras insertar duplicado; se esperaba 1", tr.Len())
	}
	if v, _ := tr.Search("hola"); v != 2 {
		t.Errorf("valor de \"hola\" = %d; se esperaba 2 (actualizado)", v)
	}
}

func TestAutocomplete(t *testing.T) {
	tr := New[int]()
	for i, w := range []string{"casa", "casaca", "caso", "castor", "perro", "pera"} {
		tr.Insert(w, i)
	}

	got := claves(tr.Autocomplete("cas"))
	want := []string{"casa", "casaca", "caso", "castor"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Autocomplete(\"cas\") = %v; se esperaba %v", got, want)
	}

	got = claves(tr.Autocomplete("per"))
	want = []string{"pera", "perro"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Autocomplete(\"per\") = %v; se esperaba %v", got, want)
	}

	if res := tr.Autocomplete("xyz"); res != nil {
		t.Errorf("Autocomplete(\"xyz\") = %v; se esperaba nil", res)
	}
}

func TestDelete(t *testing.T) {
	tr := New[int]()
	for i, w := range []string{"casa", "casaca", "caso", "castor"} {
		tr.Insert(w, i)
	}

	if !tr.Delete("caso") {
		t.Fatal("Delete(\"caso\") = false; se esperaba true")
	}
	if _, ok := tr.Search("caso"); ok {
		t.Error("\"caso\" sigue existiendo tras eliminarla")
	}
	for _, w := range []string{"casa", "casaca", "castor"} {
		if _, ok := tr.Search(w); !ok {
			t.Errorf("%q desapareció tras borrar \"caso\"", w)
		}
	}
	if tr.Len() != 3 {
		t.Errorf("Len() = %d tras borrar una; se esperaba 3", tr.Len())
	}

	if tr.Delete("gato") {
		t.Error("Delete(\"gato\") = true; se esperaba false")
	}
	if tr.Len() != 3 {
		t.Errorf("Len() = %d tras borrar inexistente; se esperaba 3", tr.Len())
	}
}

func TestDeletePrefijoQueEsPalabra(t *testing.T) {
	// "cas" es clave y además prefijo de "caso".
	tr := New[int]()
	tr.Insert("cas", 1)
	tr.Insert("caso", 2)

	if !tr.Delete("cas") {
		t.Fatal("Delete(\"cas\") = false; se esperaba true")
	}
	if _, ok := tr.Search("cas"); ok {
		t.Error("\"cas\" sigue existiendo tras eliminarla")
	}
	if _, ok := tr.Search("caso"); !ok {
		t.Error("\"caso\" desapareció al borrar \"cas\"")
	}
}

func TestKeys(t *testing.T) {
	tr := New[int]()
	for i, w := range []string{"sol", "sal", "luna"} {
		tr.Insert(w, i)
	}
	got := tr.Keys()
	want := []string{"luna", "sal", "sol"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Keys() = %v; se esperaba %v", got, want)
	}
}

package trie

import (
	"reflect"
	"testing"
)

func TestInsertSearch(t *testing.T) {
	tr := New()
	palabras := []string{"casa", "casaca", "caso", "perro", "pera"}
	for _, w := range palabras {
		tr.Insert(w)
	}

	for _, w := range palabras {
		if !tr.Search(w) {
			t.Errorf("Search(%q) = false; se esperaba true", w)
		}
	}

	for _, w := range []string{"cas", "ca", "per", "gato", ""} {
		if tr.Search(w) {
			t.Errorf("Search(%q) = true; se esperaba false", w)
		}
	}

	if tr.Len() != len(palabras) {
		t.Errorf("Len() = %d; se esperaba %d", tr.Len(), len(palabras))
	}
}

func TestInsertDuplicado(t *testing.T) {
	tr := New()
	tr.Insert("hola")
	tr.Insert("hola")
	if tr.Len() != 1 {
		t.Errorf("Len() = %d tras insertar duplicado; se esperaba 1", tr.Len())
	}
}

func TestAutocomplete(t *testing.T) {
	tr := New()
	for _, w := range []string{"casa", "casaca", "caso", "castor", "perro", "pera"} {
		tr.Insert(w)
	}

	got := tr.Autocomplete("cas")
	want := []string{"casa", "casaca", "caso", "castor"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Autocomplete(\"cas\") = %v; se esperaba %v", got, want)
	}

	got = tr.Autocomplete("per")
	want = []string{"pera", "perro"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Autocomplete(\"per\") = %v; se esperaba %v", got, want)
	}

	if res := tr.Autocomplete("xyz"); res != nil {
		t.Errorf("Autocomplete(\"xyz\") = %v; se esperaba nil", res)
	}
}

func TestDelete(t *testing.T) {
	tr := New()
	for _, w := range []string{"casa", "casaca", "caso", "castor"} {
		tr.Insert(w)
	}

	// Borrar una palabra existente.
	if !tr.Delete("caso") {
		t.Fatal("Delete(\"caso\") = false; se esperaba true")
	}
	if tr.Search("caso") {
		t.Error("\"caso\" sigue existiendo tras eliminarla")
	}
	// Las demás deben seguir intactas.
	for _, w := range []string{"casa", "casaca", "castor"} {
		if !tr.Search(w) {
			t.Errorf("%q desapareció tras borrar \"caso\"", w)
		}
	}
	if tr.Len() != 3 {
		t.Errorf("Len() = %d tras borrar una; se esperaba 3", tr.Len())
	}

	// Borrar algo que no existe no debe cambiar nada.
	if tr.Delete("gato") {
		t.Error("Delete(\"gato\") = true; se esperaba false")
	}
	if tr.Len() != 3 {
		t.Errorf("Len() = %d tras borrar inexistente; se esperaba 3", tr.Len())
	}
}

func TestDeletePrefijoQueEsPalabra(t *testing.T) {
	// "cas" es palabra y además prefijo de "caso".
	tr := New()
	tr.Insert("cas")
	tr.Insert("caso")

	if !tr.Delete("cas") {
		t.Fatal("Delete(\"cas\") = false; se esperaba true")
	}
	if tr.Search("cas") {
		t.Error("\"cas\" sigue existiendo tras eliminarla")
	}
	if !tr.Search("caso") {
		t.Error("\"caso\" desapareció al borrar \"cas\"")
	}
}

func TestKeys(t *testing.T) {
	tr := New()
	for _, w := range []string{"sol", "sal", "luna"} {
		tr.Insert(w)
	}
	got := tr.Keys()
	want := []string{"luna", "sal", "sol"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Keys() = %v; se esperaba %v", got, want)
	}
}

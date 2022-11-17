package main

import (
	"encoding/json"
	"os"
	"path"
	"testing"
)

func TestLogZip(t *testing.T) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent(" ", " ")
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	p := path.Join(wd, "..", "test", "ac_o00407-0139200010477.logjez")
	rs := openAndParseLogDeUrnaZip(p)
	if len(rs) != 1 {
		t.Fatalf("expected only one log, recieved %d\n", len(rs))
	}
	r := rs[0]
	if r.VotosPRQtd != 206 {
		t.Fatalf("expected VotosPRQtd %d, got %d", 206, r.VotosPRQtd)
	}
	if r.QtdTeclaIndevida != 28 {
		t.Fatalf("expected QtdTeclaIndevida %d, got %d", 28, r.QtdTeclaIndevida)
	}
	if r.QtdAlertas != 22 {
		t.Fatalf("expected QtdAlertas %d, got %d", 22, r.QtdAlertas)
	}
	if r.Versao != "8.26.0.0" {
		t.Fatalf("expected Versao %s, got %s", "8.26.0.0", r.Versao)
	}
	if r.Modelo != "UE2020" {
		t.Fatalf("expected Modelo %s, got %s", "UE2020", r.Modelo)
	}
	if r.Secao != 477 {
		t.Fatalf("expected Secao %d, got %d", 477, r.Secao)
	}
	if r.Zona != 1 {
		t.Fatalf("expected Zona %d, got %d", 1, r.Zona)
	}
	if r.Municipio != 1392 {
		t.Fatalf("expected Municipio %d, got %d", 1392, r.Municipio)
	}
}

func TestLogZipUrnaLigadaDia28(t *testing.T) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent(" ", " ")
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	p := path.Join(wd, "..", "test", "am_o00407-0253400060103.logjez")
	rs := openAndParseLogDeUrnaZip(p)
	if len(rs) != 1 {
		t.Fatalf("expected only one log, recieved %d\n", len(rs))
	}
	r := rs[0]
	if r.Versao != "8.26.0.0" {
		t.Fatalf("expected Versao %s, got %s", "8.26.0.0", r.Versao)
	}
	if r.Modelo != "UE2010" {
		t.Fatalf("expected Modelo %s, got %s", "UE2020", r.Modelo)
	}
	if r.Secao != 103 {
		t.Fatalf("expected Secao %d, got %d", 477, r.Secao)
	}
	if r.Zona != 6 {
		t.Fatalf("expected Zona %d, got %d", 1, r.Zona)
	}
	if r.Municipio != 2534 {
		t.Fatalf("expected Municipio %d, got %d", 1392, r.Municipio)
	}
}

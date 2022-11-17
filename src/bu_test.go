package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"
)

func TestProcessaBUApenasPresidente(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	bupath := path.Join(wd, "..", "test", "ac_o00407-0100700090001.bu")
	bu := processaBU(bupath)
	fmt.Printf("bu = %+v\n", bu)
	if bu.DataGeracao.Unix() != 1667142202 {
		t.Fatal("dataGeracao failed")
	}
	if bu.DataHoraEmissao.Unix() != 1667142167 {
		t.Fatal("dataHoraEmissao failed")
	}
	if bu.Municipio != 1007 {
		t.Fatal("municipio failed")
	}
	if bu.Zona != 9 {
		t.Fatal("zona failed")
	}
	if bu.Local != 1104 {
		t.Fatal("local failed")
	}
	if bu.Secao != 1 {
		t.Fatal("secao failed")
	}
	if bu.IDEleicao != 545 {
		t.Fatal("idEleicao failed")
	}
	if bu.QtdEleitoresAptos != 335 {
		t.Fatal("qtdEleitoresAptos failed")
	}
	if bu.QtdComparecimento != 261 {
		t.Fatal("qtdComparecimento failed")
	}
	if len(bu.Votos) != 4 {
		t.Fatal("len votos failed")
	}
	if bu.Votos[0].Qtd != 75 || bu.Votos[0].Voto != 13 {
		t.Fatal("votos lula failed")
	}
	if bu.Votos[1].Qtd != 175 || bu.Votos[1].Voto != 22 {
		t.Fatal("votos bolso failed")
	}
}

func TestProcessBUVotacaoComGovernador(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	bupath := path.Join(wd, "..", "test", "rs_o00407-8473500030006.bu")
	bu := processaBU(bupath)
	fmt.Printf("bu = %+v\n", bu)
	if bu.DataGeracao.Unix() != 1667149505 {
		t.Fatal("dataGeracao failed")
	}
	if bu.DataHoraEmissao.Unix() != 1667149438 {
		t.Fatal("dataHoraEmissao failed")
	}
	if bu.Municipio != 84735 {
		t.Fatal("municipio failed")
	}
	if bu.Zona != 3 {
		t.Fatal("zona failed")
	}
	if bu.Local != 1171 {
		t.Fatal("local failed")
	}
	if bu.Secao != 6 {
		t.Fatal("secao failed")
	}
	if bu.IDEleicao != 545 {
		t.Fatal("idEleicao failed")
	}
	if bu.QtdEleitoresAptos != 371 {
		t.Fatal("qtdEleitoresAptos failed")
	}
	if bu.QtdComparecimento != 333 {
		t.Fatal("qtdComparecimento failed")
	}
	if len(bu.Votos) != 4 {
		t.Fatal("len votos failed")
	}
	if bu.Votos[0].Qtd != 136 || bu.Votos[0].Voto != 13 {
		t.Fatal("votos lula failed")
	}
	if bu.Votos[1].Qtd != 182 || bu.Votos[1].Voto != 22 {
		t.Fatal("votos bolso failed")
	}
}

func TestBUJson(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	bupath := path.Join(wd, "..", "test", "rs_o00407-8473500030006.bu")
	bu := readBU(bupath, "/home/dpsigor/eleicoes/formatoArquivos/python/bu_dump.py", "/home/dpsigor/eleicoes/formatoArquivos/spec/bu.asn1")
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent(" ", " ")
	enc.Encode(bu)
}

func TestBUComOitoCamposEntidadeBU(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	bupath := path.Join(wd, "..", "test", "ba_o00407-3787701050083.bu")
	bu := processaBU(bupath)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent(" ", " ")
	enc.Encode(bu)
}

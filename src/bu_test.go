package main

import (
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
	if bu.dataGeracao.Unix() != 1667142202 {
		t.Fatal("dataGeracao failed")
	}
	if bu.dataHoraEmissao.Unix() != 1667142167 {
		t.Fatal("dataHoraEmissao failed")
	}
	if bu.municipio != 1007 {
		t.Fatal("municipio failed")
	}
	if bu.zona != 9 {
		t.Fatal("zona failed")
	}
	if bu.local != 1104 {
		t.Fatal("local failed")
	}
	if bu.secao != 1 {
		t.Fatal("secao failed")
	}
	if bu.idEleicao != 545 {
		t.Fatal("idEleicao failed")
	}
	if bu.qtdEleitoresAptos != 335 {
		t.Fatal("qtdEleitoresAptos failed")
	}
	if bu.qtdComparecimento != 261 {
		t.Fatal("qtdComparecimento failed")
	}
	if len(bu.votos) != 4 {
		t.Fatal("len votos failed")
	}
	if bu.votos[0].qtd != 75 || bu.votos[0].voto != 13 {
		t.Fatal("votos lula failed")
	}
	if bu.votos[1].qtd != 175 || bu.votos[1].voto != 22 {
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
	if bu.dataGeracao.Unix() != 1667149505 {
		t.Fatal("dataGeracao failed")
	}
	if bu.dataHoraEmissao.Unix() != 1667149438 {
		t.Fatal("dataHoraEmissao failed")
	}
	if bu.municipio != 84735 {
		t.Fatal("municipio failed")
	}
	if bu.zona != 3 {
		t.Fatal("zona failed")
	}
	if bu.local != 1171 {
		t.Fatal("local failed")
	}
	if bu.secao != 6 {
		t.Fatal("secao failed")
	}
	if bu.idEleicao != 545 {
		t.Fatal("idEleicao failed")
	}
	if bu.qtdEleitoresAptos != 371 {
		t.Fatal("qtdEleitoresAptos failed")
	}
	if bu.qtdComparecimento != 333 {
		t.Fatal("qtdComparecimento failed")
	}
	if len(bu.votos) != 4 {
		t.Fatal("len votos failed")
	}
	if bu.votos[0].qtd != 136 || bu.votos[0].voto != 13 {
		t.Fatal("votos lula failed")
	}
	if bu.votos[1].qtd != 182 || bu.votos[1].voto != 22 {
		t.Fatal("votos bolso failed")
	}
}

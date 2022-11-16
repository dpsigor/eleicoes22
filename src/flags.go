package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/dpsigor/hltrnty"
)

func mustExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		flag.Usage()
		log.Fatalf("'%s' não existe\n", path)
	}
}

func parseFlags() (string, string, string, string, string) {
	dirPtr := flag.String("dir", "", "Caminho do diretório que contém os arquivos .bu e .logjez")
	faURL := "https://www.tre-mt.jus.br/eleicoes/eleicoes-2022/documentacao-tecnica-do-software-da-urna-eletronica"
	faDsc := fmt.Sprintf("Caminho do diretório que contém o conteúdo do zip 'Formato dos arquivos de BU, RDV e assinatura digital'. Obter em %s", faURL)
	formatoArquivosPtr := flag.String("fa", "", faDsc)
	ufPtr := flag.String("uf", "", fmt.Sprintf("UF a extrair dados. Valores: %s", ufs))
	destinoPtr := flag.String("dest", "stdout", fmt.Sprintf("Destino. Valores: %s.", destinos))
	flag.Parse()
	if *dirPtr == "" {
		flag.Usage()
		log.Fatal("flag dir é obrigatória")
	}
	mustExist(*dirPtr)
	if *formatoArquivosPtr == "" {
		flag.Usage()
		log.Fatal("flag 'fa' é obrigatória")
	}
	if *ufPtr == "" {
		flag.Usage()
		log.Fatal("flag uf é obrigatória")
	}
	if !hltrnty.Some(ufs, func(uf string) bool {
		return uf == *ufPtr
	}) {
		flag.Usage()
		log.Fatalf("uf inválido: '%s'\n", *ufPtr)
	}
	if !hltrnty.Some(destinos, func(dest string) bool {
		return dest == *destinoPtr
	}) {
		flag.Usage()
		log.Fatalf("destino inválido: '%s'\n", *destinoPtr)
	}
	buDump := path.Join(*formatoArquivosPtr, "python", "bu_dump.py")
	mustExist(buDump)
	buSpec := path.Join(*formatoArquivosPtr, "spec", "bu.asn1")
	mustExist(buSpec)
	return *dirPtr, *ufPtr, *destinoPtr, buDump, buSpec
}

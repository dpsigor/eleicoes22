package main

import (
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/bodgit/sevenzip"
)

func parseLogDeUrna(f io.ReadCloser) Resultado {
	alertas := make([][]byte, 0)
	versoes := make([][]byte, 0)
	teclaIndevida := 0
	modelo := []byte{}
	secao := 0
	zona := 0
	municipio := 0
	bs, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := bytes.Split(bs, LF)
	votosPRQtd := 0
	for _, line := range lines {
		if !bytes.HasPrefix(line, TS) && !bytes.HasPrefix(line, TSMinus1) {
			continue
		}
		if bytes.Contains(line, SE) {
			secao = logDoisPontosDigitos(line)
			continue
		}
		if bytes.Contains(line, ZN) {
			zona = logDoisPontosDigitos(line)
			continue
		}
		if bytes.Contains(line, MU) {
			municipio = logDoisPontosDigitos(line)
			continue
		}
		if bytes.Contains(line, VRS) {
			versoes = append(versoes, line)
			continue
		}
		if bytes.Contains(line, ModUrna) {
			matches := modeloRgx.FindSubmatch(line)
			if len(matches) == 2 {
				modelo = matches[1]
			}
			continue
		}
		if bytes.Contains(line, TCLIND) {
			teclaIndevida++
			continue
		}
		if bytes.Contains(line, AL) {
			alertas = append(alertas, line)
			continue
		}
		if bytes.Contains(line, PR) {
			votosPRQtd++
			continue
		}
	}
	versao := []byte{}
	for _, versaoLine := range versoes {
		ms := versaoRgx.FindSubmatch(versaoLine)
		if len(ms) == 2 {
			if string(versao) != "" && string(versao) != string(ms[1]) {
				log.Fatalf("versões diferentes: %s e %s", versao, ms[1])
			}
			versao = ms[1]
		}
	}
	return Resultado{
		VotosPRQtd:       votosPRQtd,
		QtdTeclaIndevida: teclaIndevida,
		QtdAlertas:       len(alertas),
		Versao:           string(versao),
		Modelo:           string(modelo),
		Secao:            secao,
		Zona:             zona,
		Municipio:        municipio,
	}
}

func openAndParseLogDeUrnaZip(zippath string) []Resultado {
	reader, err := sevenzip.OpenReader(zippath)
	if err != nil {
		log.Fatalf("erro ao abrir 7z reader %s: %s\n", zippath, err)
	}
	defer reader.Close()
	rs := make([]Resultado, 0)
	for k := range reader.File {
		if !strings.HasSuffix(reader.File[k].Name, ".dat") {
			// log.Printf("dentro de '%s' tem um não .dat: '%s'\n", zippath, reader.File[k].Name)
			continue
		}
		f, err := reader.File[k].Open()
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		rs = append(rs, parseLogDeUrna(f))
	}
	return rs
}

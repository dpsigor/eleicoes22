package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

// atoi panics if can't turn str into int
func atoi(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

var doisPontosRgx = regexp.MustCompile(`: (\d{4,5})\s`)

// logDoisPontosDigitos converte linha de log com texto "Descrição: xxxx" em int xxx
// Exemplo:
// 30/10/2022 05:30:00     INFO    67305985        GAP     Se��o Eleitoral: 0010   C3D6B3F42568FCFD
// Retorna int "0010"
func logDoisPontosDigitos(line []byte) int {
	ms := doisPontosRgx.FindSubmatch(line)
	if ms == nil {
		log.Fatalf("falhou em obter valor após ':' em %s\n", line)
	}
	return atoi(string(ms[1]))
}

// fPathSecao recebe ac_o00407-0149000020032.bu e devolve valores
// de municipio, zona e secao. Esses valores serão validados, posteriormente,
// com dados internos do boletim de urna e do log de urna
func fPathSecao(s string) (string, string, string, int, int, int) {
	n := strings.LastIndex(s, "/")
	s = s[n:]
	if len(s) < 25 {
		log.Fatalf("recebeu arquivo com menos de 25 caracteres no nome\n")
	}
	muStr := s[11:16]
	mu, err := strconv.Atoi(muStr)
	if err != nil {
		log.Fatal(err)
	}
	znStr := s[16:20]
	zn, err := strconv.Atoi(znStr)
	if err != nil {
		log.Fatal(err)
	}
	scStr := s[20:24]
	sc, err := strconv.Atoi(scStr)
	if err != nil {
		log.Fatal(err)
	}
	return muStr, znStr, scStr, mu, zn, sc
}

func assertEqual(x, y int, descr string) {
	if x != y {
		log.Fatalf("%s: %d != %d", descr, x, y)
	}
}

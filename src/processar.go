package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/dpsigor/hltrnty"
)

type urna struct {
	buname            string
	logname           string
	municipio         string
	zona              string
	secao             string
	bolso             int64
	lula              int64
	brancos           int64
	nulos             int64
	qtdComparecimento int64
	qtdVotosPR        int64
	qtdTeclaIndevida  int64
	qtdAlertas        int64
	versao            string
	modelo            string
}

type Resultado struct {
	VotosPRQtd       int    `json:"votos_pr_qtd"`
	QtdTeclaIndevida int    `json:"qtd_tecla_indevida"`
	QtdAlertas       int    `json:"qtd_alertas"`
	Versao           string `json:"versao"`
	Modelo           string `json:"modelo"`
	Secao            int    `json:"secao"`
	Zona             int    `json:"zona"`
	Municipio        int    `json:"municipio"`
}

var modeloRgx = regexp.MustCompile(`Modelo de Urna:\s(.+?)\s`)
var versaoRgx = regexp.MustCompile(`Vers.o da aplica..o:\s+(.+?)\s+`)

func preencherUrnaComBU(u *urna, bu BU) error {
	u.qtdComparecimento = bu.QtdComparecimento
	for _, v := range bu.Votos {
		if v.TipoVoto == tipoVotoBranco {
			u.nulos = v.Qtd
		} else if v.TipoVoto == tipoVotoNulo {
			u.brancos = v.Qtd
		} else {
			if v.Voto == bolso {
				u.bolso = v.Qtd
			} else if v.Voto == lula {
				u.lula = v.Qtd
			}
		}
	}
	return nil
}

func processUrna(dir string, entry fs.DirEntry) urna {
	name := entry.Name()
	fpath := fmt.Sprintf("%s/%s", dir, name)
	muStr, znStr, scStr, mu, zn, sc := fPathSecao(fpath)
	u := urna{
		buname:    name,
		municipio: muStr,
		zona:      znStr,
		secao:     scStr,
	}
	// Boletim de Urna
	// bu := readBU(fpath, buDump, buSpec)
	bu := processaBU(fpath)
	// validar interior do BU com nome do arquivo
	buMU := int(bu.Municipio)
	buZN := int(bu.Zona)
	buSC := int(bu.Secao)
	assertEqual(mu, buMU, fmt.Sprintf("bu municipio para %s", name))
	assertEqual(zn, buZN, fmt.Sprintf("bu zona para %s", name))
	assertEqual(sc, buSC, fmt.Sprintf("bu sc para %s", name))
	// Preencher urna com votos
	err := preencherUrnaComBU(&u, bu)
	if err != nil {
		log.Fatalf("erro ao parse BU %s: %s", name, err)
	}
	// Log de Urna
	logFpath := strings.Replace(fpath, ".bu", ".logjez", 1)
	u.logname = strings.Replace(name, ".bu", ".logjez", 1)
	rs := openAndParseLogDeUrnaZip(logFpath)
	if len(rs) != 1 {
		log.Fatalf("devolveu %d resultados para o log zip %s", len(rs), fpath)
	}
	logData := rs[0]
	// validar interior do LOG com nome do arquivo
	assertEqual(mu, logData.Municipio, fmt.Sprintf("log municipio para %s", u.logname))
	assertEqual(zn, logData.Zona, fmt.Sprintf("log zona para %s", u.logname))
	assertEqual(sc, logData.Secao, fmt.Sprintf("log sc para %s", u.logname))
	// Preencher urna com dados do LOG
	u.qtdVotosPR = int64(logData.VotosPRQtd)
	u.qtdTeclaIndevida = int64(logData.QtdTeclaIndevida)
	u.qtdAlertas = int64(logData.QtdAlertas)
	u.versao = logData.Versao
	u.modelo = logData.Modelo
	return u
}

func processEntries(dir string, entries []fs.DirEntry, ch chan<- urna) {
	workers := 4 * 2 * runtime.GOMAXPROCS(0)
	entriesCh := make(chan fs.DirEntry)
	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		go func(entries <-chan fs.DirEntry, results chan<- urna) {
			wg.Add(1)
			for entry := range entries {
				results <- processUrna(dir, entry)
			}
			wg.Done()
		}(entriesCh, ch)
	}
	for k := range entries {
		name := entries[k].Name()
		if strings.Contains(name, "am_o00407-0255000680797") {
			// AM mu 02550 zona 0068 secao 0797 não tem dados
			continue
		}
		if strings.Contains(name, "ba_o00407-3413401710232") {
			// BA mu 34134 zona 0171 secao 0232 não tem dados no dia da eleição segundo turno
			continue
		}
		if strings.HasSuffix(name, ".bu") {
			entriesCh <- entries[k]
		}
	}
	close(entriesCh)
	wg.Wait()
	close(ch)
}

func processUF(dir string, ch chan<- urna, uf string, logger *Logger) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	entries = hltrnty.Filter(entries, func(i fs.DirEntry) bool {
		n := i.Name()
		return strings.HasPrefix(n, uf) && strings.HasSuffix(n, ".bu")
	})
	logger.totalUrnas = len(entries)
	processEntries(dir, entries, ch)
}

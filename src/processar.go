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

func preencherUrnaComBU(u *urna, bu bujson) error {
	achouEleicao := false
	for _, eleicao := range bu.EntidadeBoletimUrna.ResultadosVotacaoPorEleicao {
		if eleicao.IDEleicao != 545 {
			continue
		}
		achouEleicao = true
		for _, resultado := range eleicao.ResultadosVotacao {
			u.qtdComparecimento = resultado.QtdComparecimento
			for _, cargo := range resultado.TotaisVotosCargo {
				for _, voto := range cargo.VotosVotaveis {
					if voto.TipoVoto == "nulo" {
						u.nulos = voto.QuantidadeVotos
					} else if voto.TipoVoto == "branco" {
						u.brancos = voto.QuantidadeVotos
					} else {
						if voto.IdentificacaoVotavel.Codigo == 13 {
							u.lula = voto.QuantidadeVotos
						} else if voto.IdentificacaoVotavel.Codigo == 22 {
							u.bolso = voto.QuantidadeVotos
						} else {
							return fmt.Errorf("bu com voto que não é nulo, nem branco, nem 13 nem 22")
						}
					}
				}
			}
		}
	}
	if !achouEleicao {
		return fmt.Errorf("não achou eleição")
	}
	return nil
}

func processUrna(dir, buDump, buSpec string, entry fs.DirEntry) urna {
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
	bu := readBU(fpath, buDump, buSpec)
	// validar interior do BU com nome do arquivo
	buMU := int(bu.EntidadeBoletimUrna.IdentificacaoSecao.MunicipioZona.Municipio)
	buZN := int(bu.EntidadeBoletimUrna.IdentificacaoSecao.MunicipioZona.Zona)
	buSC := int(bu.EntidadeBoletimUrna.IdentificacaoSecao.Secao)
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

func processEntries(dir, buDump, buSpec string, entries []fs.DirEntry, ch chan<- urna) {
	workers := 2 * runtime.GOMAXPROCS(0)
	entriesCh := make(chan fs.DirEntry)
	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		go func(entries <-chan fs.DirEntry, results chan<- urna) {
			wg.Add(1)
			for entry := range entries {
				results <- processUrna(dir, buDump, buSpec, entry)
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
		if strings.HasSuffix(name, ".bu") {
			entriesCh <- entries[k]
		}
	}
	close(entriesCh)
	wg.Wait()
	close(ch)
}

func processUF(dir, buDump, buSpec string, ch chan<- urna, uf string, logger *Logger) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	entries = hltrnty.Filter(entries, func(i fs.DirEntry) bool {
		n := i.Name()
		return strings.HasPrefix(n, uf) && strings.HasSuffix(n, ".bu")
	})
	logger.totalUrnas = len(entries)
	processEntries(dir, buDump, buSpec, entries, ch)
}

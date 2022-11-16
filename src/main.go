package main

import (
	"fmt"
	"log"
	"os"
)

var ufs = []string{"ac", "al", "am", "ap", "ba", "ce", "df", "es", "go", "ma", "mg", "ms", "mt", "pa", "pb", "pe", "pi", "pr", "rj", "rn", "ro", "rr", "rs", "sc", "se", "sp", "to"}
var destinos = []string{"stdout", "file"}

var bu2json = "./bu2json.py"

func main() {
	dir, uf, destino, buDump, buSpec := parseFlags()

	logger := Logger{}
	ch := make(chan urna)
	go processUF(dir, buDump, buSpec, ch, uf, &logger)

	logProgress := destino == "file"

	w := os.Stdout
	if destino == "file" {
		os.Mkdir("./results", 0777)
		f, err := os.Create("./results/eleicoes_" + uf)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		w = f
	}

	w.WriteString("bu_filename;log_filename;bolso;lula;brancos;nulos;qtd_comparecimento;qtd_votos_pr;qtd_tecla_indevida;qtd_alertas;versao;modelo\n")
	for v := range ch {
		str := fmt.Sprintf("%s;%s;%s;%s;%s;%d;%d;%d;%d;%d;%d;%d;%d;%s;%s\n", v.buname, v.logname, v.municipio, v.zona, v.secao, v.bolso, v.lula, v.brancos, v.nulos, v.qtdComparecimento, v.qtdVotosPR, v.qtdTeclaIndevida, v.qtdAlertas, v.versao, v.modelo)
		w.WriteString(str)
		if logProgress {
			logger.progress(uf)
		}
	}
}

package main

import "fmt"

type Logger struct {
	totalUrnas int
	done       int
}

func (l *Logger) progress(uf string) {
	if l.totalUrnas == 0 {
		return
	}
	l.done++
	fmt.Printf("%s %d de %d (%.02f%%)\n", uf, l.done, l.totalUrnas, 100*float64(l.done)/float64(l.totalUrnas))
}

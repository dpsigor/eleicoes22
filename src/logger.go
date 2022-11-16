package main

import "fmt"

type Logger struct {
	totalUrnas int
	done       int
}

func (l *Logger) progress() {
	if l.totalUrnas == 0 {
		return
	}
	l.done++
	fmt.Printf("%d de %d (%.02f%%)\n", l.done, l.totalUrnas, 100*float64(l.done)/float64(l.totalUrnas))
}

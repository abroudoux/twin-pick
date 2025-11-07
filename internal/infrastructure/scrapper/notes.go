package scrapper

import "slices"

var letterBoxdNotes = []float64{0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5}

func isValidNote(note float64) bool {
	return slices.Contains(letterBoxdNotes, note)
}

package api

import "fmt"

type Annotations []AnnotationData

type AnnotationData struct {
	Range   []int  `json:"range"`
	Message string `json:"message"`
}

// Return annotation.Range as readable format MM:SS
func (annotation AnnotationData) GetReadableRange() []string {
	readableRange := make([]string, len(annotation.Range))

	for i, annotationTime := range annotation.Range {
		readableRange[i] = fmt.Sprintf("%02d:%02d", annotationTime/60, annotationTime%60)
	}

	return readableRange
}

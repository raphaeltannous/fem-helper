package api

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAnnotationData_GetReadableRange(t *testing.T) {
	annotationTests := []struct {
		Range []int
		want  []string
	}{
		{
			Range: []int{323, 453},
			want:  []string{"05:23", "07:33"},
		},
		{
			Range: []int{0, 59},
			want:  []string{"00:00", "00:59"},
		},
		{
			Range: []int{60, 120},
			want:  []string{"01:00", "02:00"},
		},
		{
			Range: []int{3599, 3601},
			want:  []string{"59:59", "60:01"},
		},
		{
			Range: []int{125, 130},
			want:  []string{"02:05", "02:10"},
		},
		{
			Range: []int{1439, 1440},
			want:  []string{"23:59", "24:00"},
		},
	}

	for _, c := range annotationTests {
		testName := fmt.Sprintf("%v", c.Range)
		t.Run(testName, func(t *testing.T) {
			answer := AnnotationData{
				Range: c.Range,
			}.GetReadableRange()

			if !reflect.DeepEqual(answer, c.want) {
				t.Errorf("got %v, want %v", answer, c.want)
			}
		})
	}
}

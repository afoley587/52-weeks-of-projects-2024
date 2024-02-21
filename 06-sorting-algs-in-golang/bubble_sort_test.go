package main

import (
	"reflect"
	"testing"
)

func BenchmarkBubbleSort(t *testing.B) {
	tests := []struct {
		input, sorted []int
	}{
		{[]int{}, []int{}},
		{[]int{1}, []int{1}},
		{[]int{1, 2}, []int{1, 2}},
		{[]int{2, 1}, []int{1, 2}},
		{[]int{5, 4, 2, 1, 3}, []int{1, 2, 3, 4, 5}},
		{
			[]int{
				250, 222, 471, 813, 812, 865, 425,
				277, 13, 498, 458, 379, 737, 486,
				58, 773, 328, 284, 420, 510, 970,
				177, 101, 74, 263, 533, 925, 782,
				242, 576, 771, 613, 962, 918, 791,
				920, 662, 775, 876, 826, 995, 797,
				391, 278, 82, 93, 945, 556, 758, 664,
			},
			[]int{
				13, 58, 74, 82, 93, 101, 177, 222,
				242, 250, 263, 277, 278, 284, 328,
				379, 391, 420, 425, 458, 471, 486,
				498, 510, 533, 556, 576, 613, 662,
				664, 737, 758, 771, 773, 775, 782,
				791, 797, 812, 813, 826, 865, 876,
				918, 920, 925, 945, 962, 970, 995,
			}},
	}
	for i, test := range tests {
		Bubble(test.input)
		if !reflect.DeepEqual(test.input, test.sorted) {
			t.Fatalf("Failed test case #%d. Want %v got %v", i, test.sorted, test.input)
		}
	}
}

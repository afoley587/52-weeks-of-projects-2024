package main

func Heapit(arr []int, n, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n && arr[largest] < arr[left] {
		largest = left
	}

	if right < n && arr[largest] < arr[right] {
		largest = right
	}

	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		Heapit(arr, n, largest)
	}

}

func Heap(arr []int) {
	n := len(arr)

	for i := n/2 - 1; i >= 0; i-- {
		Heapit(arr, n, i)
	}

	for i := n - 1; i > 0; i-- {
		arr[i], arr[0] = arr[0], arr[i]
		Heapit(arr, i, 0)
	}

}

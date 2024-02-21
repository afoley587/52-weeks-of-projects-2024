package main

func Partition(arr []int, low, high int) int {

	pivot := arr[high]

	i := low - 1

	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

func Quick(arr []int, low, high int) {
	if low < high {
		p := Partition(arr, low, high)
		Quick(arr, low, p-1)
		Quick(arr, p+1, high)
	}
}

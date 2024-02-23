# DSA In GoLang - Sorting
## Implementing Some Of The Common Sorting Algos In GoLang

![Thumbnail](./imgs/thumbnail.png)

I'm going to go on a limb and say that you're like me - it's been a while
since you've took a good hard look at data structures and algorithms (DSA).
I haven't formally been testing on these things since I was in school and,
even though I do use these principles in my own work quite often, it's nice
to sometimes have a refresher course. Especially if you're learning a new 
language. When I first started learning Go, I spent a few hours googling 
"Good GoLang Projects For Beginners" and a whole bunch of stuff came up.
But, sometimes it's hard to find a good project. What are the inputs? What
are the outputs? Where are the bugs? Etc... That's why I just re-picked up
some of the standard DSA type things. You know what the inputs are, you know
what the outputs should be, and all of these algorithms are so extensively
researched and documented that finding the bug and seeing any underlying
issues becomes blatanly obvious. So, let's do two things today:

1. A brief refresher of some of the more well-known sorting algorithms 
    (Bubble Sort, Insertion Sort, Heap Sort, and Quick Sort)
2. Implementing these in GoLang and then running test cases on them

## A Naive Approach - Bubble Sort
### Overview

Let's first talk about the algorithm that we all learned at the
start of our CompSci journey - Bubble Sort!

Bubble Sort is a simple sorting algorithm that 
repeatedly steps through the list to be sorted, compares each pair of 
adjacent items, and swaps them if they are in the wrong order. This process 
is repeated until no swaps are needed, indicating that the list is sorted. 
The algorithm gets its name from the way smaller elements 
"bubble" to the top of the list with each iteration.

Here's a basic outline of how the bubble sort algorithm works:

1. Start at the beginning of the list.
2. Compare the first two elements. If the first is greater than the second, swap them.
3. Move to the next pair of elements and repeat the comparison and swap if necessary.
4. Continue this process, moving one element to the right each time, until you reach the end of the list.
5. Repeat steps 1-4 until no more swaps are needed, indicating that the list is sorted.

Bubble Sort is not efficient for large lists because its average and worst-case 
time complexity is O(n^2), where n is the number of elements in the list. 
However, it's easy to understand and implement, making it useful for 
educational purposes or for sorting small lists.

### Code

```golang
func Bubble(arr []int) {
	for i := 0; i < len(arr)-1; i++ {
		for j := 0; j < len(arr)-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}
```

## An Improved Approach - Insertion Sort
### Overview

Let's moved on to a somewhat improved method - Insertion Sort.

Insertion Sort is anoter simple sorting algorithm that 
builds the final sorted list one item at a time. It 
works by iteratively taking an unsorted element and inserting 
it into its correct position within a sorted sublist. Initially, 
the first element of the list is considered to be a sorted sublist of 
one element. Then, for each subsequent element, the algorithm 
compares it with the elements in the sorted sublist, moving larger 
elements one position to the right to make space for the new element.

Here's a basic description of how the insertion sort algorithm works:

1. Start with the second element of the list (assuming the first element is already sorted).
2. Compare this element with the elements in the sorted sublist, moving from right to left.
3. If the current element is smaller than the element being compared, shift the compared element one position to the right.
4. Repeat step 3 until you find the correct position for the current element in the sorted sublist.
5. Insert the current element into its correct position in the sorted sublist.
6. Repeat steps 2-5 for each remaining unsorted element until the entire list is sorted.

Insertion Sort is an efficient algorithm for sorting small lists or lists that are 
almost sorted. Its average and worst-case time complexity is O(n^2), 
where n is the number of elements in the list. However, its best-case time complexity 
is O(n) when the list is already sorted. Additionally, insertion sort has a 
space complexity of O(1), making it suitable for situations with limited memory resources.

### Code

```golang
func Insertion(arr []int) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1

		for j >= 0 && key < arr[j] {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}
```

## A Data Structure Approach - Heap Sort
### Overview

So we've done bubble and insertion sort. Is there more we can do by possibly
modifying our thinking and introducing some different data structures into
the mix? Let's try something new... Heap Sort!

Heap Sort is a comparison-based sorting algorithm that leverages the properties 
of a binary heap data structure to efficiently sort elements. It consists of 
two main phases: heap construction and heapification.

**Heap Construction:**

* The first phase involves constructing a max-heap from the input array. 
    A max-heap is a binary tree where the value of each parent node is 
    greater than or equal to the values of its children.
* Starting from the last non-leaf node (index n/2 - 1, where n is the 
    number of elements), the algorithm sifts down each node to ensure 
    that the entire tree satisfies the max-heap property.

**Heapification (Sorting):**

* Once the max-heap is constructed, the largest element (at the root) is 
    swapped with the last element in the heap.
* The heap size is reduced by one, and the algorithm restores the max-heap 
    property by performing a sift-down operation on the root node.
* This process is repeated iteratively until the heap size reduces to 1. 
    After each iteration, the largest remaining element is placed at the end of the array.
* The resulting array is sorted in ascending order.

Heap Sort has a time complexity of O(n log n) for all cases, where 
n is the number of elements in the array. This makes it an efficient 
sorting algorithm, especially for large datasets. Additionally, heap sort 
is an in-place algorithm, meaning it does not require additional memory 
proportional to the size of the input array, except for a constant amount 
used for temporary storage during the sorting process. However, heap sort is 
not stable, meaning it may change the relative order of elements with equal keys.

### Code

```golang
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
```

## A Recursive Approach - Quick Sort
### Overview

Now, our approach with Heap Sort was pretty good. But, what are most
other libraries using as their internal sorting methods? This brings
us to our final sorting method of the day - Quick Sort.

Quick Sort is a highly efficient sorting algorithm that follows the 
divide-and-conquer strategy. It works by selecting a 'pivot' element 
from the array and partitioning the other elements into two sub-arrays 
according to whether they are less than or greater than the pivot. 
The sub-arrays are then recursively sorted independently.

Here's how quick sort typically works:

**Partitioning:**

* Choose a pivot element from the array. Commonly, the pivot is selected 
    as the last element in the array, but other strategies such as selecting
    the first, middle, or a random element can also be used.
* Reorder the array so that all elements with values less than the pivot 
    come before the pivot, and all elements with values greater than 
    the pivot come after it. After this partitioning, the pivot is 
    in its final sorted position.
* This step is often implemented using the partitioning subroutine, 
    which rearranges the elements such that all elements less than 
    the pivot are moved to its left, and all elements greater than 
    the pivot are moved to its right.

**Recursively Sort Sub-arrays:**

* Recursively apply quick sort to the sub-array of elements 
    with smaller values (i.e., those elements to the left of the pivot) 
    and the sub-array of elements with greater values 
    (i.e., those elements to the right of the pivot).
* The base case for the recursion is when the sub-array has 
    zero or one element, in which case it is already sorted.

**Combine Results:**

* As the recursive calls return, the sub-arrays are combined 
    into a single sorted array.

Quick Sort is generally faster in practice compared to other 
sorting algorithms for large datasets, primarily due to its in-place 
partitioning and relatively simple implementation. It typically has an 
average-case time complexity of O(n log n), making it very efficient. 
However, its worst-case time complexity is O(n^2) when the pivot selection 
consistently results in unbalanced partitions, such as when the input array 
is already sorted or nearly sorted. To mitigate this, various techniques like 
randomizing pivot selection or using median-of-three pivot selection can be 
employed. Overall, quick sort is widely used in practice and serves as a 
foundational sorting algorithm in many programming languages and libraries.

### Code

```golang
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
```

## Running

Let's run these with `go test -bench` (more on benchmarking [here](https://www.practical-go-lessons.com/chap-34-benchmarks)).

If you're curious, each of the files have a test file. We will do a quick
benchmark using `go test -bench`:

```shell
$ go test -bench .
goos: darwin
goarch: arm64
pkg: github.com/afoley587/52-weeks-of-projects/06-sorting-algs-in-golang
BenchmarkBubbleSort-10          1000000000               0.0000041 ns/op
BenchmarkHeapSort-10            1000000000               0.0000030 ns/op
BenchmarkInsertionSort-10       1000000000               0.0000030 ns/op
BenchmarkQuickSort-10           1000000000               0.0000025 ns/op
PASS
ok      github.com/afoley587/52-weeks-of-projects/06-sorting-algs-in-golang     0.178s
$ 
```

The third column represents `Nanoseconds per operation`. In short,
it gives you an idea of how fast on average our functions run. We can see that 
Quick Sort tends to be the fastest over the 1000000000 operations. Now, that
might not always be the case. ALways remember that different algorithms can
be more relevant and useful on different sets of data!

I hope you enjoyed following along! As always, all code can always
be found [here](https://github.com/afoley587/52-weeks-of-projects-2024/tree/main/06-sorting-algs-in-golang) on my GitHub repo!
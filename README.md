# memdebug

go mem debugger

## Install

`go get github.com/steveoc64/memdebug`

## Usage - Track Memory Allocations (slow)

In your code :

```
t1 := time.Now() // get the current time

... do stuff

memdebug.Print(t1, "entrant updated", data1, data2, ...)
```

produces output similar to this :

![example](example.png)

`801.598µs (  139224): 100190264: 144013560:   1900544 <-   [entrant updated]`

Data columns are :

- Time elapsed since t1

- (effect on memory)  Yellow = consumed, Green = freed

- Total memory application on the heap

- Total memory pool allocated from the OS (heap is a subset of this).  Should be the same number most of the time, but will print in red when it grows, grabbing more memory from the OS)

- Total stack size (should be pretty constant)

NOTE:

Using memdebug.Print() will slow the machine down quite a bit, because it does a full GC on each print, and a call to free OS memory, so as to get a super accurate picture of what is actually allocated and freed.

## Usage - CPU Profiling (not slow)

If you are more interested in CPU profiling, add this to your main.go:

```
	memdebug.Profile()
	defer memdebug.WriteProfile()
```

This will turn off the GC on each print, and the subsequent calls to free OS memory, so you want see an accurate picture of memory allocations and frees (just actual allocations pre-GC)

But it will give an accurate picture of CPU usage.

When the app exits, you will see some extra files in the root dir - namely

`cpu.pdf`

Which is a pproff`d PDF file showing a CPU trace where all the time is spent.

Nice !



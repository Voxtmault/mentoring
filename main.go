package main

import (
	"fmt"
	"time"
)

func dataType() {
	// Primitive data types
	// int, float, bool

	// Hover over the variable to see the detail
	var sampleInt int = -1   // Integer can be negative
	var sampleUint uint = 10 // Unsigned integer can't be negative
	fmt.Println("Integer", sampleInt)
	fmt.Println("Unsigned Integer", sampleUint)

	// The are a lot of variations for the above data types, eg
	// int8, int16, int32, int64
	// uint8, uint16, uint32, uint64
	// The only different being the size of the data type or the maximum ammount of data it can hold
	// for example
	// int8 can hold -128 to 127
	// uint8 can hold 0 to 255
	// A simple rule of thumb is, 2 to the power of n, where n is the number of bits (int8 where the n is 8)

	var sampleFloat32 float32 = 3.14
	var sampleFloat64 float64 = 3.14
	// The difference between the two is the precision

	fmt.Println("Float", sampleFloat32)
	fmt.Println("Float", sampleFloat64)

	var sampleString string = "Hello World"
	fmt.Println("String", sampleString)

	// You can access each individual character of the string
	fmt.Println("First character", sampleString[0])

	var sampleBool bool = true
	fmt.Println("Boolean", sampleBool)
	// There are are a lot of other data types, including complex one like struct, array, slice, map, etc
	// You can learn more about them yourself. Google is your best friend here
}

func basicLogic() {
	// In Go, well generally in programming, we have the basic logic
	// If statements
	if 1 == 1 {
		fmt.Println("1 is equal to 1")
	}

	if 1 != 0 {
		fmt.Println("1 is not equal to 0")
	} else if 1 == 0 {
		fmt.Println("1 is equal to 0")
	} else {
		fmt.Println("1 is not equal to 1")
	}

	// For loops
	// Generally speaking there are 3 types of loops; for, while, do while
	// But in Go, there is only one type of loop, the for loop
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	// We can do some variation like
	var iter int = 0
	for iter < 5 {
		fmt.Println(iter)
		iter++
	}

	for {
		fmt.Println("I'm not stuck here with you, you're stuck here with me")
		time.Sleep(3 * time.Second)
		break
	}

	// The most usefull one so far, for range
	sampleString := "never gonna give you up"
	for iteration, value := range sampleString {
		fmt.Println("The first returned value is the current loop iteration", iteration)
		fmt.Println("The second returned value is a COPY of the data", string(value))

		// You can think of rune as an ASCII
		value = 'a'
	}
	fmt.Println("Final Value", sampleString)

	// Switches and Selects
	// Switch, We use switch basically for more readable if-else statements...
	// And it's is way faster than if-else
	day := "Tuesday"

	switch day {
	case "Monday":
		fmt.Println("Start of the work week")
	case "Tuesday":
		fmt.Println("Second day of the work week")
	case "Wednesday":
		fmt.Println("Midweek")
	case "Thursday":
		fmt.Println("Almost there")
	case "Friday":
		fmt.Println("Last work day")
	default:
		fmt.Println("Weekend!")
	}

	// Selects, We use selects if we to deal with multiple channel operations, more on that later
	// It's like a manager or a mom who's waiting for all of her kids to finish their chores
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- "Message from channel 1"
	}()

	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- "Message from channel 2"
	}()

	select {
	case msg1 := <-ch1:
		fmt.Println(msg1)
	case msg2 := <-ch2:
		fmt.Println(msg2)
	}
}

func moarDataTypes() {
	// Arrays
	var sampleArr [5]uint
	fmt.Println("Array", sampleArr)

	// Slices
	var sampleSlice []uint
	fmt.Println("Slice", sampleSlice)

	// The only differences is that in arrays, we HAVE to declare the size of the array
	// For slices, we dont need to...

	// Maps
	var sampleMap map[int]string // For map, we need to do things a little bit different
	// If you try to access or modify an uninitialized map, well golang isn't going to like that and panics...
	// We dont want that so we need to initialize it

	sampleMap = make(map[int]string)
	// or
	sampleMap = map[int]string{}

	fmt.Println("Map", sampleMap)

	// Pointers, the dread of all programming languages
	var sampleVar int = 1
	pointerToSV := &sampleVar           // Well this is basically how we do pointers in Go...
	fmt.Println("Pointer", pointerToSV) // Print the memory address
	fmt.Println("Value", *pointerToSV)  // Print the value

	// We use pointers to make our program runs as efficient as possible...
	// In Go, when you pass a value to a function parameter, the compiler doesn't actually use the same value / variable
	// It makes a copy of that value and then pass the copied value to the function
	// That wouldn't make that much of a difference for a small scale program. But think of a program that deals with huge ammount
	// of data every single second, we can use pointers to make the program more efficient by only creating object / reservices memory
	// address space as needed. We'll get to that later

	// Interfaces
	var sampleInterface interface{}
	fmt.Println("Interface", sampleInterface)

	// Any
	var sampleAny any
	fmt.Println("Any", sampleAny)

	// Interfaces and Any are basically the same. No kidding, you can see in the source code and see that any is just an interface{} :D
	// We use this if we don't know what kind of data we're dealing with or you just don't care about the data types. Usefull for
	// http responses

}

func main() {
	dataType()
	moarDataTypes()
	basicLogic()

	// See func.go for next lesson
}

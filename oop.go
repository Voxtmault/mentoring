package main

import "fmt"

// Hello darkness my old friend
// I have come to learn inheritance and polymorphism again :D

// This is a "class" in Go, you can have as many of these in your projects.
// But hey, who would want that many class right :D
// (said the guy who yapped all night fixing his shitty codes because he could have used a class)
type MyStruct struct {
	UserID      int
	Username    string
	Address     string
	PhoneNumber string
}

type MyDumbStruct struct {
	UserID   int
	Username string
	Password string // Man, what an idiot. Who would put get a password from the DB :D
	// JK, but please dont do this
}

// Methods, Like OOP in general. We can do something like public and private here in Go
// If you notice, all of the func name from the previous file all started using lowercase
// Using lowercase means that function is a private function and is ony accessible through the same file
// or other files in the same directory
func hiThere() {
	// I can access the dataType here since this file is in the same directory as main.go
	dataType()
}

// See example and example2 for function methods

// Now you're thinking, if you can do that with functions... can we do it with classes ?
// Well, turns out... you ABSOLUTELY CAN !!! :D

func (m MyStruct) Default() Player {
	return &MyStruct{
		UserID:      1,
		Username:    "rick",
		Address:     "rolled",
		PhoneNumber: "lmao",
	}
}

func (m MyStruct) superSecretFunction() {
	fmt.Println("well, i'll be damned")
}

// Lastly, we're going to talk about interfaces
// Not that one, but interfaces that's a contract

type Player interface {
	Default() Player
	superSecretFunction()
}

// Interfaces in this context are a kind of working contract, when 1 class or struc want to implement
// an interface, it needs to satisfy the required contract, in this case all of the functions
// defined in the interface...
var _ Player = MyStruct{}     // This is accepted since MyStruct implements / satisfy the Player contract
var _ Player = MyDumbStruct{} // This is not accepted because the dumb idiot forgot to create a signature
// function to satisfy the player contract

// When do we use interface ?
// Well, we generally use that for collaboration project where multiple people works on different part
// of the project and each module / part interact with each other. To ease programmers, we define a
// list of interface or contracts that must be fulfilled in order to run the project. That way
// we dont have to wait for our teammates to finish their work before finding out that he / she
// changed the function parameter or return type :D (deffinitely not speaking from experience :DDDDD)

// Interface are also usefull if you want to create a version of the same stuff, like for example you have
// MyStructV2{} that implements and even more robust stuff inside, it wont affect your teammate since
// as long as it is defined correctly and satisfied all that's going on doesn't matter

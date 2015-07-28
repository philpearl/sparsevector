package sparsevector

import (
	"fmt"
)

func ExampleStringIndex() {

	// Vector for user's ratings of Star Wars
	starwars := NewGenSparseVector(StringIndex{
		"Brian",
		"Liz",
		"Jane",
		"Jenny",
		"Mike",
		"David",
	}, []Value{
		1.0,
		4.5,
		3.5,
		2.0,
		4.0,
		5.0,
	})

	// Vector for user's ratings of Battlestar Galactica
	battlestargalactica := NewGenSparseVector(StringIndex{
		"Brian",
		"Liz",
		"Jane",
		"Mike",
		"David",
		"Penny",
	}, []Value{
		3.0,
		3.5,
		4.5,
		3.0,
		5.0,
	})

	// Naive similarity measure between these two films
	cos := starwars.Cos(battlestargalactica)
	fmt.Sprintf("Similarity is %f", cos)
}

package examples

import (
	"fmt"
	"log"

	"github.com/detailyang/traz-go/traz"
)

type Foo struct {
	A string        `default:"abcd"`
	B int           `default"100"`
	C traz.Duration `default:"5s"`
	D traz.Size     `default:"5KB"`
	E string        `required:"true"`
}

func Exampletraz() {
	mz := traz.NewMultiTraz()
	dtz := traz.NewDefaultTagtraz("default")
	rtz := traz.NewRequiredTagtraz("required")

	str := `
a = "defg"
e = "cc"
`

	tmz := traz.NewTOMLtraz(str)
	mz.Append(dtz, tmz, rtz)

	foo := &Foo{}
	if err := mz.Apply(foo); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", foo)
	// Output: &{A:defg B:0 C:5s D:5KB E:cc}
}

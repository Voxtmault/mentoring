package example

import (
	"github.com/voxtmault/mentoring/example2"
)

func example() {
	dataType() // I can't access dataType functions since this file is not in the same folder

	example2.WellHelloThere() // I can access this function because it is ussing an uppercase letter
}

/*
Feed events to proxyd which it will proxy to dragond.
*/
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	N := 1000
	for i := 0; i < N; i++ {
		fmt.Fprintf(os.Stdout, `{"uid":1, "march_id":1, "at":%d, "op":"arrive"}`,
			time.Now().Add(time.Duration(rand.Intn(1000))*time.Second).Unix())
		fmt.Fprintln(os.Stdout)
		time.Sleep(1 * time.Second)
	}

}

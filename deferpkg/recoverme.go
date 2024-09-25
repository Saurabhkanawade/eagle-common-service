package deferpkg

import "fmt"

func RecoverMe() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from the panic")
		}
	}()
}

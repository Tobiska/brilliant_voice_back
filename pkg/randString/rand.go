package randString

import "math/rand"

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZRSTUHIJKLMNOPQ"

func RandStringBytes(n int) string { //todo rand code room
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

package main

import "github.com/thak1411/rnlog"

func main() {
	err := rnlog.Init("./log", "rnlog.log")
	if err != nil {
		panic(err)
	}
	defer rnlog.Close()

	rnlog.Log("This is Log")
	rnlog.Debug("This is Debug")
	rnlog.Info("This is Info")
	rnlog.Warn("This is Warn")
	rnlog.Error("This is Error")
	rnlog.Fatal("This is Fatal")
	rnlog.Log("RnLog is not kill program")
}

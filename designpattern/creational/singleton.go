package main

import "sync"

type Singleton struct{}

var once sync.Once
var instance *Singleton

// Singleton Usage: Loggers, database connections, or shared configurations.

func GetInstance() *Singleton {
	once.Do(func() {
		instance = &Singleton{}
	})
	return instance
}

func main() {
	GetInstance()
}

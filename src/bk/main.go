package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	fmt.Printf("%v", "BeastKeeper\n")

	configFileName := flag.String("c", "config.bk", "config file to load")

	fmt.Printf("Config File: %v\n", *configFileName)

	configFile, err := os.Open(*configFileName)
	_ = configFile
	if err != nil {
		fmt.Printf("%v\n", "Config File not found")
	} else {
		fmt.Printf("%v\n", "Config File found")
	}

}

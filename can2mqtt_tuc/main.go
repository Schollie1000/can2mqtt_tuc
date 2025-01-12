// minimal main package, only commandline argument parsing and a printHelp() function
package main

import (
	"fmt" // printfoo
	"os"  // args

	C2M "github.com/Schollie1000/can2mqtt_tuc"
)

// Parses commandline arguments
func main() {
	conf := true
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-v":
			C2M.SetDbg(true)
		case "-c":
			i++
			C2M.SetCi(os.Args[i])
		case "-m":
			i++
			C2M.SetCs(os.Args[i])
		case "-f":
			i++
			C2M.SetC2mf(os.Args[i])
		case "-d":
			i++
			C2M.SetConfDirMode(os.Args[i])
		default:
			i = len(os.Args)
			conf = false
			printHelp()
		}
	}
	if conf {
		C2M.Start()
	}
}

// help function (obvious...)
func printHelp() {
	fmt.Printf("Test Drillbotics ")
	fmt.Printf("welcome to the CAN2MQTT Drillbotics edit bridge!\n\n")
	fmt.Printf("Usage: can2mqtt [-f <file>] [-c <CAN-Interface>] [-m <MQTT-Connect>] [-v] [-h]\n")
	fmt.Printf("<file>: a can2mqtt.csv file\n")
	fmt.Printf("<CAN-Interface>: a CAN-Interface e.g. can0\n")
	fmt.Printf("<MQTT-Connect>: connectstring for MQTT. e.g.: tcp://[user:pass@]localhost:1883\n")
}

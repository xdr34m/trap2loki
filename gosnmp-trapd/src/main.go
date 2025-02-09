package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	g "github.com/gosnmp/gosnmp"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("   %s\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	// Create a channel to pass SNMP traps
	trapChan := make(chan *g.SnmpPacket, 100) // Buffered channel to prevent blocking

	// Create a channel to listen for OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	tl := g.NewTrapListener()
	tl.OnNewTrap = func(packet *g.SnmpPacket, addr *net.UDPAddr) {
		//fmt.Printf("Received trap from %s\n", addr.IP)
		trapChan <- packet // Send trap to the channel
	}
	tl.Params = g.Default
	//debug logging
	//tl.Params.Logger = g.NewLogger(log.New(os.Stdout, "", 0))

	go func() {
		err := tl.Listen("0.0.0.0:9162")
		if err != nil {
			log.Panicf("error in listen: %s", err)
		}
	}()

	//parseTraps()

	// Start logging goroutine
	go logTraps(trapChan)

	// Wait for termination signal
	<-signalChan
	fmt.Println("Shutdown signal received, cleaning up...")

	// Close trap listener
	tl.Close()

	// Close the channel to stop the logging goroutine
	close(trapChan)

	fmt.Println("Shutdown complete.")
}

// logTraps reads from trapChan and logs the received traps // should also be used to pass the traps to snmptt or gosnmp-smi
func logTraps(trapChan chan *g.SnmpPacket) {
	for packet := range trapChan {
		for _, v := range packet.Variables {
			switch v.Type {
			case g.OctetString:
				b := v.Value.([]byte)
				log.Printf("OctetStringTrap: OID: %s, string: %x\n", v.Name, b)
			default:
				log.Printf("Trap: %+v\n", v)
			}
		}
	}
}

// func parseTraps() {
// 	module, err := parser.ParseFile("../../mibs/test.mib")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	_ = module
// 	//repr.Println(module)
// 	fmt.Println(module.Body.Nodes)
// }

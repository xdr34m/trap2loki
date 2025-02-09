package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	g "github.com/gosnmp/gosnmp"
)

func sendv1trap(stV snmptrapValues) error {
	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	g.Default.Target = stV.target
	g.Default.Port = stV.port
	g.Default.Version = g.Version1
	//g.Default.Community =

	err := g.Default.Connect()
	if err != nil {
		log.Printf("Connect() err: %v", err)
		return err
	}
	defer g.Default.Conn.Close()

	pdu := g.SnmpPDU{
		Name:  stV.pduName,
		Type:  stV.pduType,
		Value: stV.pduValue,
	}

	trap := g.SnmpTrap{
		Variables:    []g.SnmpPDU{pdu},
		Enterprise:   stV.enterprise,
		AgentAddress: stV.agentAdress,
		GenericTrap:  stV.genericTrap,
		SpecificTrap: stV.specificTrap,
		Timestamp:    stV.timestamp,
	}

	_, err = g.Default.SendTrap(trap)
	if err != nil {
		log.Printf("SendTrap() err: %v", err)
		return err
	}
	return nil
}

func flags() (snmptrapValues, error) {
	target := flag.String("target", "127.0.0.1", "ip/dns for the trap to be send to")
	pduName := flag.String("pduname", "1.3.6.1.2.1.1.6", "pdu oid")
	enterprise := flag.String("enterprise", ".1.3.6.1.6.3.1.1.5.1", "enterpriseOID prefix .")
	agentAdress := flag.String("agentadress", "127.0.0.1", "sender adress YOU!")
	port := flag.Uint64("port", 9116, "port for the trap to be sendto")
	pdutype := flag.String("pdutype", "counter64", "one of: octetstring, counter64, counter32, integer")
	pduValue := flag.String("pduvalue", "64", "value i.e. 64, 'Octetstring', must be parseable by pdutype")
	genericTrap := flag.Int("generictrap", 0, "genericTrap? 0,1,2,3,4,5? (default 0)")
	specificTrap := flag.Int("specifictrap", 0, "specificTrap? 0,1,2,3,4,5? (default 0)")
	timestamp := flag.Uint("timestamp", 300, "timestamp when alarm happend")
	flag.Parse()
	uport := uint16(*port)
	var pduType g.Asn1BER
	var pduConValue interface{}
	//var err *error
	switch *pdutype {
	case "octetstring":
		pduType = g.OctetString
		pduConValue = *pduValue
	case "counter64":
		pduType = g.Counter64
		pduPreValue, err := strconv.Atoi(*pduValue)
		if err != nil {
			return snmptrapValues{}, fmt.Errorf("conversion error of pduvalue, err: %v", err)
		}
		pduConValue = uint64(pduPreValue)
	case "counter32":
		pduType = g.Counter32
		pduPreValue, err := strconv.Atoi(*pduValue)
		if err != nil {
			return snmptrapValues{}, fmt.Errorf("conversion error of pduvalue, err: %v", err)
		}
		pduConValue = uint32(pduPreValue)
	case "integer":
		pduType = g.Integer
		pduPreValue, err := strconv.Atoi(*pduValue)
		if err != nil {
			return snmptrapValues{}, fmt.Errorf("conversion error of pduvalue, err: %v", err)
		}
		pduConValue = pduPreValue
	default:
		return snmptrapValues{}, fmt.Errorf("could not assign pdutype, allowed only: octetstring, counter64, counter32, integer, bitstring")
	}

	flag.Parse()
	f := snmptrapValues{target: *target, pduName: *pduName, enterprise: *enterprise, agentAdress: *agentAdress, port: uport, pduType: pduType, pduValue: pduConValue, genericTrap: *genericTrap, specificTrap: *specificTrap, timestamp: *timestamp}
	return f, nil
}

type snmptrapValues struct {
	target, pduName, enterprise, agentAdress string
	port                                     uint16
	pduType                                  g.Asn1BER
	pduValue                                 interface{}
	genericTrap, specificTrap                int
	timestamp                                uint
}

func main() {
	f, err := flags()
	if err != nil {
		log.Fatalf("could not parse flags, err: %v", err)
	}
	//err := sendv1trap("127.0.0.1", 9116, "1.3.6.1.2.1.1.6", g.OctetString, "a Trap, i hate it", ".1.3.6.1.6.3.1.1.5.1", "127.0.0.1", 0, 0, 1000000)
	err = sendv1trap(f)
	if err != nil {
		log.Fatalf("Error in sendv1trap, err: %v", err)
	}
}

// Copyright 2020 bacnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package main

import (
	"log"
	"net"
	"time"

	"github.com/Nortech-ai/bacnet"
	"github.com/Nortech-ai/bacnet/common"
	"github.com/Nortech-ai/bacnet/services"
	"github.com/spf13/cobra"
)

func init() {
	whoIsCmd.Flags().IntVar(&wiPeriod, "period", 1, "Period, in seconds, between WhoIs requests.")
	whoIsCmd.Flags().IntVar(&nWhoIs, "messages", 1, "Number of messages to send, being 0 unlimited.")
}

var (
	wiPeriod int
	nWhoIs   int

	whoIsCmd = &cobra.Command{
		Use:   "whois",
		Short: "Send WhoIs requests.",
		Long: "There's not much more really. This command sends a configurable number of\n" +
			"WhoIs requests with a configurable period. That's pretty much it.",
		Args: argValidation,
		Run:  whoIsExample,
	}
)

func whoIsExample(cmd *cobra.Command, args []string) {
	remoteUDPAddr, err := net.ResolveUDPAddr("udp", rAddr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %s", err)
	}

	ifaceAddrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("couldn't get interface information: %v\n", err)
	}

	listenConn, err := net.ListenPacket("udp", bAddr)
	if err != nil {
		log.Fatalf("failed to begin listening for packets: %v\n", err)
	}
	defer listenConn.Close()

	mWhoIs, err := bacnet.NewWhois()
	if err != nil {
		log.Fatalf("error generating initial WhoIs: %v\n", err)
	}

	replyRaw := make([]byte, 1024)
	sentRequests := 0
	for {
		listenConn.SetDeadline(time.Now().Add(5 * time.Second))
		if _, err := listenConn.WriteTo(mWhoIs, remoteUDPAddr); err != nil {
			log.Fatalf("Failed to write Unconfimed request WhoIs packet: %s\n", err)
		}

		log.Printf("sent: %x", mWhoIs)

		var nBytes int
		var remoteAddr net.Addr
		for {
			nBytes, remoteAddr, err = listenConn.ReadFrom(replyRaw)
			if err != nil {
				log.Fatalf("error reading incoming packet: %v\n", err)
			}
			if !common.IsLocalAddr(ifaceAddrs, remoteAddr) {
				break
			}
			log.Printf("got our own broadcast, back to listening...\n")
			break
		}

		log.Printf("read %d bytes from %s: %x\n", nBytes, remoteAddr, replyRaw[:nBytes])

		serviceMsg, err := bacnet.Parse(replyRaw[:nBytes])
		if err != nil {
			log.Fatalf("error parsing the received message: %v\n", err)
		}
		// switch between recieved messages
		t := serviceMsg.GetType()
		switch t {

		}
		iAmMessage, ok := serviceMsg.(*services.UnconfirmedIAm)
		if !ok {
			log.Fatalf("we didn't receive an IAm reply...\n")
		}

		log.Printf("unmarshalled BVLC: %#v\n", iAmMessage.BVLC)
		log.Printf("unmarshalled NPDU: %#v\n", iAmMessage.NPDU)

		decodedIAm, err := iAmMessage.Decode()
		if err != nil {
			log.Fatalf("couldn't decode the IAm reply: %v\n", err)
		}

		printIAm(&decodedIAm)

		sentRequests++

		if sentRequests == nWhoIs {
			break
		}

		time.Sleep(time.Duration(wiPeriod) * time.Second)
	}
}

// Copyright 2020 bacnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package main

import (
	"log"
	"net"
	"time"

	"github.com/Nortech-ai/bacnet"
	"github.com/Nortech-ai/bacnet/services"
	"github.com/spf13/cobra"
)

func init() {
	WritePropertyClientCmd.Flags().Uint16Var(&wpObjectType, "object-type", 1, "Object type to read.")
	WritePropertyClientCmd.Flags().Uint32Var(&wpInstanceId, "instance-id", 0, "Instance ID to read.")  // Analog-input
	WritePropertyClientCmd.Flags().Uint16Var(&wpPropertyId, "property-id", 85, "Property ID to read.") // Current-value
	WritePropertyClientCmd.Flags().Float32Var(&wpReal, "real", 0.0, "Real value to write.")
	WritePropertyClientCmd.Flags().StringVar(&wpString, "string", "", "String value to write.")
	WritePropertyClientCmd.Flags().StringVar(&wpValue, "value", "string", "Value to write.")
	WritePropertyClientCmd.Flags().IntVar(&wpPeriod, "period", 1, "Period, in seconds, between requests.")
	WritePropertyClientCmd.Flags().IntVar(&wpN, "messages", 1, "Number of requests to send, being 0 unlimited.")
}

var (
	wpObjectType uint16
	wpInstanceId uint32
	wpPropertyId uint16
	wpReal       float32
	wpString     string
	wpValue      string
	wpPeriod     int
	wpN          int

	WritePropertyClientCmd = &cobra.Command{
		Use:   "wpc",
		Short: "Send WriteProperty requests.",
		Long:  "There's not much more really. This command sends a configurable WriteProperty request.",
		Args:  argValidation,
		Run:   WritePropertyClientExample,
	}
)

func WritePropertyClientExample(cmd *cobra.Command, args []string) {
	remoteUDPAddr, err := net.ResolveUDPAddr("udp", rAddr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %s", err)
	}

	listenConn, err := net.ListenPacket("udp", bAddr)
	if err != nil {
		log.Fatalf("failed to begin listening for packets: %v\n", err)
	}
	defer listenConn.Close()

	var data interface{}
	switch wpValue {
	case "real":
		data = wpReal
	case "string":
		data = wpString
	}

	mWriteProperty, err := bacnet.NewWriteProperty(wpObjectType, wpInstanceId, wpPropertyId, data)
	if err != nil {
		log.Fatalf("error generating initial WriteProperty: %v\n", err)
	}

	replyRaw := make([]byte, 1024)
	sentRequests := 0
	for {
		if _, err := listenConn.WriteTo(mWriteProperty, remoteUDPAddr); err != nil {
			log.Fatalf("failed to write the request: %v\n", err)
		}

		log.Printf("sent: %x", mWriteProperty)

		nBytes, remoteAddr, err := listenConn.ReadFrom(replyRaw)
		if err != nil {
			log.Fatalf("error reading back the reply: %v\n", err)
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
		sACKEnc, ok := serviceMsg.(*services.SimpleACK)
		if !ok {
			log.Fatalf("we didn't receive a SACK reply...\n")
		}

		log.Printf("unmarshalled BVLC: %#v\n", sACKEnc.BVLC)
		log.Printf("unmarshalled NPDU: %#v\n", sACKEnc.NPDU)

		log.Printf("decoded SACK reply:\n\tService: %d\n", sACKEnc.Service)

		sentRequests++

		if sentRequests == wpN {
			break
		}

		time.Sleep(time.Duration(wpPeriod) * time.Second)
	}
}

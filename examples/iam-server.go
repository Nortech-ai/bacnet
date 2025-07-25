package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Nortech-ai/bacnet"
	"github.com/Nortech-ai/bacnet/common"
	"github.com/Nortech-ai/bacnet/objects"
	"github.com/Nortech-ai/bacnet/services"
	"github.com/spf13/cobra"
)

var (
	IAmCmd = &cobra.Command{
		Use:   "iam",
		Short: "Send IAm requests.",
		Long: "This example will wait until it receives a WhoIs request. Upon reception\n" +
			"it'll just reply with the configured IAm fields",
		Args: argValidation,
		Run:  IAmExample,
	}
)

func IAmExample(cmd *cobra.Command, args []string) {
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

	mIAm, err := bacnet.NewIAm(321, 31)
	if err != nil {
		log.Fatalf("error generating initial IAm: %v\n", err)
	}

	reqRaw := make([]byte, 1024)

	var nBytes int
	var remoteAddr net.Addr
	for {
		nBytes, remoteAddr, err = listenConn.ReadFrom(reqRaw)
		if err != nil {
			log.Fatalf("error reading incoming packet: %v\n", err)
		}

		if common.IsLocalAddr(ifaceAddrs, remoteAddr) {
			log.Printf("got our own broadcast, back to listening...\n")
			continue
		}

		log.Printf("read %d bytes from %s: %x\n", nBytes, remoteAddr, reqRaw[:nBytes])

		serviceMsg, err := bacnet.Parse(reqRaw[:nBytes])
		if err != nil {
			log.Fatalf("error parsing the received message: %v\n", err)
		}
		// switch between recieved messages
		s := serviceMsg.GetService()
		switch s {
		case services.ServiceUnconfirmedWhoIs:
			whoIsMessage, ok := serviceMsg.(*services.UnconfirmedWhoIs)
			if !ok {
				log.Printf("we didn't receive a WhoIs reply...\n")
			}

			log.Printf("received a WhoIs request!\n")
			_, err := whoIsMessage.Decode()
			if err != nil {
				log.Fatalf("couldn't decode the WhoIs request: %v\n", err)
			}

			if _, err := listenConn.WriteTo(mIAm, remoteUDPAddr); err != nil {
				log.Fatalf("error sending our IAm response: %v\n", err)
			}
		case services.ServiceConfirmedReadPropMultiple:
			readPropertyMessage, ok := serviceMsg.(*services.ConfirmedReadProperty)
			if !ok {
				log.Printf("we didn't receive a ReadPropertyMultiple request! Back to listening...\n")
				continue
			}
			log.Printf("received a ReadPropertyMultiple request!\n")
			decodedReadPropertyMessage, err := readPropertyMessage.DecodeRPM()
			if err != nil {
				log.Fatalf("error decoding the ReadPropertyMultiple message: %v\n", err)
			}
			out := "Decoded ReadPropertyMultiple message:\n"
			out += fmt.Sprintf(
				"\n\tObject Type: %d\n\tInstance Id: %d\n",
				decodedReadPropertyMessage.ObjectType, decodedReadPropertyMessage.InstanceNum,
			)
			for i, t := range decodedReadPropertyMessage.Tags {
				out += fmt.Sprintf(
					"\tTag %d:\n\t\tAppTag Type: %s\n\t\tValue: %+v\n\t\tData Length: %d\n",
					i, objects.TagMap[t.TagNumber], t.Value, t.Length,
				)
			}
			log.Print(out)
		}
	}
}

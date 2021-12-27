package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

var MAX_COURSE_UPLOADS uint32 = 0x64 // Nintendo has this set to 100, but we can change it if we want

func getApplicationConfig(err error, client *nex.Client, callID uint32, applicationID uint32) {

	config := make([]uint32, 0)

	switch applicationID {
	case 0: // Unknown?
		config = getApplicationConfig_Unknown0(client, callID, applicationID)
	case 1: // PIDs?
		config = getApplicationConfig_PID(client, callID, applicationID)
	case 2: // Unknown?
		config = getApplicationConfig_Unknown2(client, callID, applicationID)
	case 10: // Unknown?
		config = getApplicationConfig_Unknown10(client, callID, applicationID)
	default:
		fmt.Printf("[Warning] DataStoreSMMProtocol::GetApplicationConfig Unsupported applicationID: %v\n", applicationID)
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteListUInt32LE(config)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodGetApplicationConfig, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}

func getApplicationConfig_Unknown0(client *nex.Client, callID uint32, applicationID uint32) []uint32 {
	// I have no idea what this is
	// Just replaying data sent from the real server
	return []uint32{
		0x01000000, 0x32000000, 0x96000000, 0x2c010000, 0xf4010000,
		0x20030000, 0x14050000, 0xd0070000, 0xb80b0000, 0x88130000,
		MAX_COURSE_UPLOADS, 0x14000000, 0x1e000000, 0x28000000, 0x32000000,
		0x3c000000, 0x46000000, 0x50000000, 0x5a000000, 0x64000000,
		0x23000000, 0x4b000000, 0x23000000, 0x4b000000, 0x32000000,
		0x00000000, 0x03000000, 0x03000000, 0x64000000, 0x06000000,
		0x01000000, 0x60000000, 0x05000000, 0x60000000, 0x00000000,
		0xe4070000, 0x01000000, 0x01000000, 0x0c000000, 0x00000000,
	}
}

func getApplicationConfig_PID(client *nex.Client, callID uint32, applicationID uint32) []uint32 {
	// This looks like user PIDs?
	// Sending an empty list here crashes the game
	return []uint32{
		0x02000000, 0x70cc8269, 0x50cc8269,
		0x38cc8269, 0xdbd08269, 0xa9d08269,
		0x89d08269, 0x59c48269, 0x36c48269,
	}
}

func getApplicationConfig_Unknown2(client *nex.Client, callID uint32, applicationID uint32) []uint32 {
	// I have no idea what this is
	// Just replaying data sent from the real server
	return []uint32{0xdf070000, 0x0c000000, 0x16000000, 0x05000000, 0x00000000}
}

func getApplicationConfig_Unknown10(client *nex.Client, callID uint32, applicationID uint32) []uint32 {
	// I have no idea what this is
	// Just replaying data sent from the real server
	// Only seen on the 3DS
	return []uint32{35, 75, 96, 40, 5, 6}
}

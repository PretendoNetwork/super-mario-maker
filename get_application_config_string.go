package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getApplicationConfigString(err error, client *nex.Client, callID uint32, applicationID uint32) {

	switch applicationID {
	case 128: // Word blacklist?
		getApplicationConfigString_WordBlacklist(client, callID, applicationID)
	case 129: // PIDs?
		getApplicationConfigString_Unknown129(client, callID, applicationID)
	case 130: // Unknown?
		getApplicationConfigString_Unknown130(client, callID, applicationID)
	default:
		fmt.Printf("[Warning] DataStoreSMMProtocol::GetApplicationConfig Unsupported applicationID: %v\n", applicationID)
	}
}

func getApplicationConfigString_WordBlacklist(client *nex.Client, callID uint32, applicationID uint32) {
	rmcResponseStream := nex.NewStreamOut(nexServer)

	// TODO complete this

	// Looks like a word blacklist
	// Just replaying data sent from the real server
	//config, _ := hex.DecodeString("420000000a00e38191e38195e3828c000a00e6b688e38195e3828c000d00e5898ae999a4e38195e3828c001300e383aae382bbe38383e38388e38195e3828c000a0042414ee38195e3828c001000efbca2efbca1efbcaee38195e3828c001300e382ade3839fe381aee382b3e383bce382b9001000e5909be381aee382b3e383bce382b9001300e3818de381bfe381aee382b3e383bce382b9000c00e3818420e3818420e381ad000d00e9818ae381b3e381bee38199000700e59cb0e99c87000700e99c87e781bd000700e8a2abe781bd000700e6b4a5e6b3a2000d00e38390e383b3e38195e3828c000800e381847ee381ad000700e99c87e5baa6000a00e38198e38197e38293000a0062616ee38195e3828c001000e3818fe3828fe38197e3818fe381af000d00e8a9b3e38197e3818fe381af000d00e381a1e38293e381a1e38293000800e381a130e381930006006269637468000e00e381842ee38184efbc8ee381ad000d00e3838ae382a4efbd9ee382b9000800e3818426e38184000b00e381842de38184e381ad000a00e38184e38183e381ad0007006e69676765720006006e67676572000a00737461722069662075000a00537461722069662075000c005374617220696620796f75000c007374617220696620796f7500060050454e6c53000a00e3839ee383b3e382b300090062757474686f6c650005004c494c49000700766167696e61000700766167796e61000a00e38186e38293e381a1000a00e38186e38293e38193000a00e382a6e383b3e382b3000d00efbca9efbd89efbd8eefbd8500050045454e45000a00e381bee38293e38193000a00e382a6e383b3e383810007006e69676c65740008006e6967676c6574000c00706c65617365206c696b65000d00e3818de38293e3819fe381be00090042757474686f6c650006006c6ce381ad00080069e38184e381ad000c006769766520612073746172000a00e381a1e38293e381bd000700e4ba80e9a0ad00060070656e6973000a00efbdb3efbe9defbdba000f00706c7a206d6f72652073746172730009007374617220706c7a000900e381842829e381ad000c00504c454153452073746172000d00426974746520537465726e6500")
	//rmcResponseStream.Grow(int64(len(config)))
	//rmcResponseStream.WriteBytesNext(config)

	config := []string{
		"けされ", "消され", "削除され", "リセットされ",
		"BANされ", "ＢＡＮされ", "キミのコース", "君のコース",
		"きみのコース", "い い ね", "遊びます", "地震",
		"震災", "被災", "津波", "バンされ",
		"い~ね", "震度", "じしん", "banされ",
		"くわしくは", "詳しくは", "ちんちん", "ち0こ",
		"bicth", "い.い．ね", "ナイ～ス", "い&い",
		"い-いね", "いぃね", "nigger", "ngger",
		"star if u", "Star if u", "Star if you", "star if you",
		"PENlS", "マンコ", "butthole", "LILI",
		"vagina", "vagyna", "うんち", "うんこ",
		"ウンコ", "Ｉｉｎｅ", "EENE", "まんこ",
		"ウンチ", "niglet", "nigglet", "please like",
		"きんたま", "Butthole", "llね", "iいね",
		"give a star", "ちんぽ", "亀頭", "penis",
		"ｳﾝｺ", "plz more stars", "star plz", "い()ね",
		"PLEASE star", "Bitte Sterne",
	}

	rmcResponseStream.WriteListString(config)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodGetApplicationConfigString, rmcResponseBody)

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

func getApplicationConfigString_Unknown129(client *nex.Client, callID uint32, applicationID uint32) {
	rmcResponseStream := nex.NewStreamOut(nexServer)

	// TODO complete this

	// No idea what this is
	// Just replaying data sent from the real server
	//config, _ := hex.DecodeString("070000000d00e382bce383ade3818be38289000a00efbc90e3818be3828900080030e3818be38289001600e38184e38080e38080e38184e38080e38080e381ad000700e38184e38184000a00e69db1e697a5e69cac000700e5a4a7e99c8700")
	//rmcResponseStream.Grow(int64(len(config)))
	//rmcResponseStream.WriteBytesNext(config)

	config := []string{
		"ゼロから", "０から", "0から", "い　　い　　ね", "いい", "東日本", "大震",
	}

	rmcResponseStream.WriteListString(config)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodGetApplicationConfigString, rmcResponseBody)

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

func getApplicationConfigString_Unknown130(client *nex.Client, callID uint32, applicationID uint32) {
	rmcResponseStream := nex.NewStreamOut(nexServer)

	// TODO complete this

	// No idea what this is
	// Just replaying data sent from the real server
	//config, _ := hex.DecodeString("3d0000000a00e38184e38184e381ad000a00e4b88be38195e38184000d00e3818fe381a0e38195e38184000a00e68abce38197e381a6000a00e3818ae38197e381a6000700e8bf94e38199000a00e3818be38188e38199000400e6989f000d00e38197e381a6e3818fe3828c000a00e38199e3828be38288001000e29886e3818fe3828ce3819fe38289001000e29886e38182e38192e381bee38199001000e29885e3818fe3828ce3819fe38289001000e29885e38182e38192e381bee38199000700e38197e381ad000a00e38193e3828de38199001000e38193e3828de38195e3828ce3819f000a00e382a2e3838ae383ab000d00e38395e382a1e38383e382af000d00e382ade383b3e382bfe3839e000700e2978be381ad000d00e382ade38381e382ace382a4000a00e38186e38293e381930008004b495449474149000700e98791e78e89000d00e3818ae381a3e381b1e38184000a00e29886e3818ae38199000a00e29886e68abce38199000a00e29885e3818ae38199000a00e29885e68abce38199000d00e38184e38184e38199e3828b000a00e38184e38184e38288000a00e382a4e382a4e3838d000700e382b1e38384000a00e38186e38293e381a1001300e3818be3818fe3819be38184e38196e38184000d00e8a69ae3819be38184e589a4000a00e382b7e383a3e38396000d00e3818de38293e3819fe381be000d00e381a1e38293e381a1e38293000d00e3818ae38197e381a3e38193000d00e381a1e38293e381bde38193000d00e38193e3828de38197e381a6000a00e382b0e38383e38389000a00e382b0e38383e38388000a00e383ace2978fe38397000a00e38390e383bce382ab000d00e3818de381a1e3818ce38184000a00e381a1e38293e38192000a00e3839ee383b3e382b3000a00e381bee38293e38193000a00e38381e383b3e3839d000700e382afe382ba000a00e382a6e383b3e382b3001f00e3838ae382a4e382b9e3818ae381ade3818ce38184e38197e381bee3819900060070656e6973000a00e382a4e382a4e381ad000a00e29886e38288e3828d001100e3838ae382a4e382b921e38197e381a6000b00e381be2fe38293e38193000b00e381bee382932fe3819300")
	//rmcResponseStream.Grow(int64(len(config)))
	//rmcResponseStream.WriteBytesNext(config)

	config := []string{
		"いいね", "下さい", "ください",
		"押して", "おして", "返す",
		"かえす", "星", "してくれ",
		"するよ", "☆くれたら", "☆あげます",
		"★くれたら", "★あげます", "しね",
		"ころす", "ころされた", "アナル",
		"ファック", "キンタマ", "○ね",
		"キチガイ", "うんこ", "KITIGAI",
		"金玉", "おっぱい", "☆おす",
		"☆押す", "★おす", "★押す",
		"いいする", "いいよ", "イイネ",
		"ケツ", "うんち", "かくせいざい",
		"覚せい剤", "シャブ", "きんたま",
		"ちんちん", "おしっこ", "ちんぽこ",
		"ころして", "グッド", "グット",
		"レ●プ", "バーカ", "きちがい",
		"ちんげ", "マンコ", "まんこ",
		"チンポ", "クズ", "ウンコ",
		"ナイスおねがいします", "penis", "イイね",
		"☆よろ", "ナイス!して", "ま/んこ",
		"まん/こ",
	}

	rmcResponseStream.WriteListString(config)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodGetApplicationConfigString, rmcResponseBody)

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

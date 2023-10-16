package nex_datastore_super_mario_maker

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetApplicationConfigString(err error, client *nex.Client, callID uint32, applicationID uint32) uint32 {
	// Word blacklists?
	config := make([]string, 0)

	switch applicationID {
	case 128:
		config = getApplicationConfigString_WordBlacklist1()
	case 129:
		config = getApplicationConfigString_WordBlacklist2()
	case 130:
		config = getApplicationConfigString_WordBlacklist3()
	default:
		fmt.Printf("[Warning] DataStoreSMMProtocol::GetApplicationConfigString Unsupported applicationID: %v\n", applicationID)
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteListString(config)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodGetApplicationConfigString, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.SecureServer.Send(responsePacket)

	return 0
}

func getApplicationConfigString_WordBlacklist1() []string {
	// Just replaying data sent from Nintendo's servers
	return []string{
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
}

func getApplicationConfigString_WordBlacklist2() []string {
	// Just replaying data sent from Nintendo's servers
	return []string{
		"ゼロから", "０から", "0から", "い　　い　　ね", "いい", "東日本", "大震",
	}
}

func getApplicationConfigString_WordBlacklist3() []string {
	// Just replaying data sent from Nintendo's servers
	return []string{
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
}

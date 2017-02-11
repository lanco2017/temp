
package main

import (
	"fmt"
	"log"

	"net/http"
	"os"
	//"regexp"

	"strconv"
	"strings"
	
	"github.com/line/line-bot-sdk-go/linebot"

	// "bytes"

	"io/ioutil"

	"image/jpeg"

    "crypto/md5"
    "encoding/hex"

    // "encoding/json"
    	// "github.com/bitly/go-simplejson"

)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

//https://gist.github.com/synr/d3d68d42b12204d981b39203a0b16762
func GetMD5Hash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

func real_num(text string) string {
	text = strings.Replace(text, "１", "1", -1)
	text = strings.Replace(text, "２", "2", -1)
	text = strings.Replace(text, "３", "3", -1)
	text = strings.Replace(text, "４", "4", -1)
	text = strings.Replace(text, "５", "5", -1)
	text = strings.Replace(text, "６", "6", -1)
	text = strings.Replace(text, "７", "7", -1)
	text = strings.Replace(text, "８", "8", -1)
	text = strings.Replace(text, "９", "9", -1)
	text = strings.Replace(text, "０", "0", -1)
	return text
}



//http://qiita.com/koki_cheese/items/66980888d7e8755d01ec
// func handleTask(w http.ResponseWriter, r *http.Request) {
// }

	//修改時主要參考官方文件以及：
	// https://github.com/line/line-bot-sdk-go/blob/master/examples/kitchensink/server.go
		// KEY = handleText 等
	// https://github.com/dongri/line-bot-sdk-go
		// KEY = linebot.NewAudioMessage(originalContentURL, duration)
func callbackHandler(w http.ResponseWriter, r *http.Request) {
		// allow cross domain AJAX requests
		// http://stackoverflow.com/questions/12830095/setting-http-headers-in-golang/
		//	https://developer.mozilla.org/ja/docs/Web/HTTP/HTTP_access_control
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set( "Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS" )
		//http://qiita.com/futosu/items/b49f7d9e28101daaa99e
		//https://play.golang.org/p/xHp44c_pJm
	// w.Header().Set("Access-Control-Allow-Headers","Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		// https://groups.google.com/forum/#!topic/golang-nuts/-Sh616lXNRE

	//-----------------------------------------------

	// log.Print("r")
	// log.Print(r)

	events, err := bot.ParseRequest(r)
												//c := appengine.NewContext(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	
	for _, event := range events {

		//2016.12.23+ 統一基本資訊集中
		//2016.12.24+ 嘗試抓使用者資訊 https://github.com/line/line-bot-sdk-go/blob/master/examples/kitchensink/server.go
		target_id_code := event.Source.UserID + event.Source.GroupID + event.Source.RoomID//target_id_code := ""
 		log.Print("event.Source.UserID = " + event.Source.UserID)
		log.Print("event.Source.GroupID = " + event.Source.GroupID)
		log.Print("event.Source.RoomID = " + event.Source.RoomID)
		log.Print("target_id_code = " + target_id_code)
		target_item := ""
		if event.Source.UserID!="" {
			target_item = "好友"
		}
		if event.Source.GroupID!="" {
			target_item = "群組對話"
		}
		if event.Source.RoomID!="" {
			target_item = "房間"
		}
		log.Print("target_item = " + target_item)

		username := ""
		userStatus := ""
		userImageUrl := ""
																				//userLogo_url := ""
		switch target_id_code{
			case "U77931c3187565323bcd8838498dcf651":
				username = "LL"
			case "U0a8152d2cea8c981aa2436a0ab596bca":
				username = "♪"
			case "C259ba6d5ace1c4932361612bf55394b3":
				username = "群組封測人員"
			case "Ceb12b46125861074c3b84106a5dcf750":
				username = "會友群組"
		}
		log.Print("username = " + username)

		//如果是群組會出錯，只能 1 對 1的時候。
		//if target_item == "好友"{
		if event.Source.UserID!="" {
			//2016.12.24+ 嘗試抓使用者資訊 https://github.com/line/line-bot-sdk-go/blob/master/examples/kitchensink/server.go
			profile, err := bot.GetProfile(event.Source.UserID).Do()
			if err != nil {
				log.Print(1162)
			    log.Print(err)
			}
			log.Print("profile.DisplayName = " + profile.DisplayName)			// println(res.Displayname)
			log.Print("profile.StatusMessage " + profile.StatusMessage)			// println(res.StatusMessage)
			log.Print("profile.PictureURL = " + profile.PictureURL)

														// log.Print("userLogo_url = " +  userLogo_url)
			//如果不是認識的 ID，就取得對方的名
			if username == ""{
				username = profile.DisplayName
			}
			userStatus = profile.StatusMessage
			userImageUrl = profile.PictureURL

			log.Print("username = " + username)
			log.Print("userStatus = " + userStatus)
			log.Print("userImageUrl = " + userImageUrl)

		}

		user_talk := ""
		if username == ""{
			user_talk = "【" + target_item + "】 " + target_id_code
		}else{
			user_talk = username
		}
		log.Print("※ user_talk = " + user_talk)

		//2016.12.27+

		// https://devdocs.line.me/en/#template-messages
		// HTTPS
		// JPEG or PNG
		// Aspect ratio: 1:1.51
		// Max width: 1024px
		// Max: 1 MB

		// 1024 = 1.51x
		// X = 678.145

		// 300 = 1.51x
		// x = 300 / 1.51 = 長的 / 1.51 = 198
		// 300 * 1.51 = 453 (用 300 當短)

		//LineTemplate_chat := linebot.NewURITemplateAction("線上與開發者聊天", "http://www.smartsuppchat.com/widget?key=77b943aeaffa11a51bb483a816f552c70e322417&vid=" + target_id_code + "&lang=tw&pageTitle=%E9%80%99%E6%98%AF%E4%BE%86%E8%87%AA%20LINE%40%20%E9%80%B2%E4%BE%86%E7%9A%84%E5%8D%B3%E6%99%82%E9%80%9A%E8%A8%8A")

		imageURL := "https://trello-attachments.s3.amazonaws.com/52ff05f27a3c676c046c37f9/585e3fb981c1240b4df88d73/34fba56ed5cbb5d7f0a9d7d2543ff238/C02018kong-koan_8481.JPG" //單位圖
		SystemImageURL := "https://trello-attachments.s3.amazonaws.com/52ff05f27a3c676c046c37f9/5831e5e304f9fac88ac50a23/d390ae079971c82074b5174c98899e9e/2017.png"

		LineTemplate_CarouselColumn_feedback := linebot.NewCarouselColumn(
			SystemImageURL, "意見回饋 feedback", "你可以透過此功能\n對 LINE 機器人的 開發者 提出建議\n或錯誤回報、其他提案。",
			LineTemplate_addme,
			LineTemplate_chat,
			//linebot.NewMessageTemplateAction("聯絡 LINE 機器人開發者", "開發者"),
			linebot.NewPostbackTemplateAction("發訊息給教會", "取得發訊息給教會的提示",""),
		)

		//正題

		//只會抓到透過按鈕按下去的東西。方便做新的觸發點。(缺點是沒有 UI 介面的時候會無法使用)
		if event.Type == linebot.EventTypePostback {
				//這裡用來設計按下某按鈕後要做什麼事情
				log.Print("觸發 Postback 功能（不讓使用者察覺的程式利用）")
				log.Print("event.Postback.Data = " + event.Postback.Data)
				HttpPost_JANDI(target_item + " " + "[" + user_talk + "](" + userImageUrl + ") 觸發了按鈕並呼了 event.Postback.Data = " + event.Postback.Data + `\n` + userStatus, "brown" , "LINE 程式觀察",target_id_code)
				HttpPost_IFTTT(target_item + " " + user_talk + " 觸發了按鈕並呼了 event.Postback.Data = " + event.Postback.Data + `\n<br>` + userImageUrl + `\n<br>` + userStatus , "LINE 程式觀察" ,target_id_code)
				HttpPost_Zapier(target_item + " " + "[" + user_talk + "](" + userImageUrl + ") 觸發了按鈕並呼了 event.Postback.Data = " + event.Postback.Data + `\n` + userStatus, "LINE 程式觀察" ,target_id_code,user_talk)

				// if event.Postback.Data == "測試"{

				// }

				if event.Postback.Data == "test"{


					// https://devdocs.line.me/en/#imagemap-message
					// "x": 0,
     				//	"y": 0,
		   			// "width": 520,
		   			// "height": 1040

		   			//log.Print("test MD5 = " + GetMD5Hash(event.Postback.Data))

		   			//測試圖片地圖
					obj_message := linebot.NewImagemapMessage(
							"https://synr.github.io/test",
							"Imagemap alt text",
							linebot.ImagemapBaseSize{1040, 1040},
							linebot.NewURIImagemapAction("https://store.line.me/family/manga/en", linebot.ImagemapArea{0, 0, 520, 520}),
							linebot.NewURIImagemapAction("https://store.line.me/family/music/en", linebot.ImagemapArea{520, 0, 520, 520}),
							linebot.NewURIImagemapAction("https://store.line.me/family/play/en", linebot.ImagemapArea{0, 520, 520, 520}),
							linebot.NewMessageImagemapAction("URANAI!", linebot.ImagemapArea{520, 520, 520, 520}),	//上限 400 字
					)

					if _, err := bot.ReplyMessage(event.ReplyToken,obj_message).Do(); err != nil {
						log.Print(1586)
						log.Print(err)
					}
				}

				if event.Postback.Data == "開啟管理者選單"{
					switch target_id_code {
						case "U6f738a70b63c5900aa2c0cbbe0af91c4":
							imageURL = SystemImageURL
							LineTemplate_test := linebot.NewCarouselTemplate(
								linebot.NewCarouselColumn(
									imageURL, "管理模式", "測試中",
									linebot.NewURITemplateAction("本季行事曆","https://docs.google.com/spreadsheets/d/1RYchaiPFyPNzCC7paUzg4tfcJ6Y2XlsASer4V5K4_eU/pubhtml"),
									linebot.NewPostbackTemplateAction("無效選項","admin", ""),
									linebot.NewPostbackTemplateAction("登出","登出管理者", ""),
								),
								// LineTemplate_other_example,
								// LineTemplate_other,
								LineTemplate_CarouselColumn_feedback,
							)
							no_temp_msg := "請更新使用最新版本的 LINE APP 才能查看內容 。"
							obj_message := linebot.NewTemplateMessage(no_temp_msg, LineTemplate_test)
							if _, err = bot.ReplyMessage(event.ReplyToken, obj_message).Do(); err != nil {
									log.Print(572)
									log.Print(err)
							}
						default:
					}
				}

				if event.Postback.Data == "passcheck"{
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("請輸入暗號驗證管理者身分")).Do(); err != nil {
							log.Print(1929)
							log.Print(err)
					}
				}

				//2017.01.03+
				if event.Postback.Data == "admin"{
					switch target_id_code {
						case "U6f738a70b63c5900aa2c0cbbe0af91c4":
							imageURL = SystemImageURL
							LineTemplate_test := linebot.NewCarouselTemplate(
								linebot.NewCarouselColumn(
									imageURL, "管理模式", "For ADMIN mode.",
									linebot.NewPostbackTemplateAction("核對「暗號」","passcheck", ""),
									linebot.NewPostbackTemplateAction("管理模式","admin", ""),
									linebot.NewPostbackTemplateAction("測試圖片地圖","test", ""),
								),
								// LineTemplate_other_example,
								// LineTemplate_other,
								LineTemplate_CarouselColumn_feedback,
							)
							no_temp_msg := "請更新使用最新版本的 LINE APP 才能查看內容 。"
							obj_message := linebot.NewTemplateMessage(no_temp_msg, LineTemplate_test)
							if _, err = bot.ReplyMessage(event.ReplyToken, obj_message).Do(); err != nil {
									log.Print(605)
									log.Print(err)
							}
						default:
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("您無法使用此功能。")).Do(); err != nil {
									log.Print(1955)
									log.Print(err)
							}
					}
				}

				if event.Postback.Data == "登出管理者"{
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("你已登出管理模式")).Do(); err != nil {
						log.Print(1965)
						log.Print(err)
					}
				}






				if event.Postback.Data == "取消離開群組"{
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("你已經取消請我離開 :)")).Do(); err != nil {
						log.Print(1207)
						log.Print(err)
					}
				}

				//2016.12.26+
				if event.Postback.Data == "按下確定離開群組對話"{
					template := linebot.NewCarouselTemplate(
						linebot.NewCarouselColumn(
							SystemImageURL, "請機器人離開群組", "你確定要請我離開嗎QAQ？\n如果確定請按下方按鈕 QQ",
							linebot.NewPostbackTemplateAction("請機器人離開群組","離開群組", "機器人已經自動離開。\n如要加回來請找：\nhttps://line.me/R/ti/p/@rgp6918s\n如要聯絡開發者請找：\nhttps://line.me/R/ti/p/@uwk0684z"),
							//linebot.NewPostbackTemplateAction("請機器人離開群組","離開群組", "機器人已經自動離開。\n如要加回來請找：\nhttps://line.me/R/ti/p/@sjk2434l\n如要聯絡開發者請找：\nhttps://line.me/R/ti/p/@uwk0684z"),
							LineTemplate_addme,
							LineTemplate_chat,
						),
					)
					obj_message := linebot.NewTemplateMessage("這是命令機器人自己離開群組的方法。\n這功能只支援 APP 使用。\n請用 APP 端查看下一步。", template)
					if _, err = bot.ReplyMessage(event.ReplyToken, obj_message).Do(); err != nil {
						log.Print(1225)
						log.Print(err)
					}
				}

				if event.Postback.Data == "離開群組"{
					if target_item == "群組對話" {
						if _, err := bot.LeaveGroup(target_id_code).Do(); err != nil {
							log.Print(1233)
						    log.Print(err)
						}
						//HttpPost_JANDI("自動離開 "  + user_talk , "gray" , "LINE 離開群組",target_id_code)
						//HttpPost_IFTTT("自動離開 "  + user_talk , "LINE 離開群組",target_id_code)
						HttpPost_Zapier("自動離開 "  + user_talk , "LINE 離開群組",target_id_code,user_talk)
						log.Print("觸發自動離開 " + user_talk +  " 群組")
					}
				}
		}
		//觸發加入好友
		if event.Type == linebot.EventTypeFollow {
				HttpPost_JANDI("有新的好朋友：["  + user_talk + "](" + userImageUrl  + ")" + `\n` + userStatus, "blue" , "LINE 新好友",target_id_code)
				HttpPost_IFTTT("有新的好朋友："  + user_talk  + `\n<br>` + userImageUrl + `\n<br>` + userStatus, "LINE 新好友" ,target_id_code)
				HttpPost_Zapier("有新的好朋友：["  + user_talk + "](" + userImageUrl  + ")" + `\n` + userStatus, "LINE 新好友" ,target_id_code,user_talk)
				//target_id_code := event.Source.UserID + event.Source.GroupID + event.Source.RoomID	//target_id_code := ""
				log.Print("觸發與 " + user_talk + " 加入好友")

			    imageURL = SystemImageURL
				//template := LineTemplate_firstinfo
				t_msg := "請用最新版本的 LINE APP 觀看效果。\n如有其他建議或想討論，請對我輸入「開發者」進行聯絡。"
				obj_message := linebot.NewTemplateMessage(t_msg, LineTemplate_firstinfo)

				// username := ""
				// if target_id_code == "U6f738a70b63c5900aa2c0cbbe0af91c4"{
				// 	username = "LL"
				// }
				// if target_id_code == "Uf150a9f2763f5c6e18ce4d706681af7f"{
				// 	username = "包包"
				// }
				//reply 的寫法
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("你好啊！" + username + "～\n想知道我的功能，可以說：「簡介」\n單獨輸入「聖經」可以知道查詢方法。\n\nPS：最新版的 LINE APP 上可以看到不一樣的內容喔！"),obj_message).Do(); err != nil {
						log.Print(1288)
						log.Print(err)
				}
		}
		//觸發解除好友
		if event.Type == linebot.EventTypeUnfollow {
				HttpPost_JANDI("與 ["  + user_talk + "](" + userImageUrl + ") 解除好友" + `\n` + userStatus, "gray" , "LINE 被解除好友",target_id_code)
				HttpPost_IFTTT("與 "  + user_talk + " 解除好友" + `\n<br>` + userImageUrl + `\n<br>` + userStatus , "LINE 被解除好友" ,target_id_code)
				HttpPost_Zapier("與 ["  + user_talk + "](" + userImageUrl + ") 解除好友" + `\n` + userStatus , "LINE 被解除好友" ,target_id_code,user_talk)
				log.Print("觸發與 " + user_talk + " 解除好友")
		}
		//觸發加入群組聊天
		if event.Type == linebot.EventTypeJoin {
				HttpPost_JANDI("加入了 "  + user_talk , "blue" , "LINE 已加入群組",target_id_code)
				HttpPost_IFTTT("加入了 "  + user_talk , "LINE 已加入群組" ,target_id_code)
				HttpPost_Zapier("加入了 "  + user_talk , "LINE 已加入群組" ,target_id_code,user_talk)
				log.Print("觸發加入" + user_talk)
 				//source := event.Source
 				//log.Print("觸發加入群組聊天事件 = " + source.GroupID)
 				push_string := "很高興你邀請我進來這裡聊天！"

				//if source.GroupID == "Ca78bf89fa33b777e54b4c13695818f81"{
				if target_id_code == "Ca78bf89fa33b777e54b4c13695818f81"{
					push_string += "\n你好，" + user_talk + "。"
				}
					//push 的寫法
					// 				if _, err = bot.PushMessage(source.GroupID, linebot.NewTextMessage(push_string)).Do(); err != nil {
					// 					log.Print(err)
					// 				}
					// 				if _, err = bot.PushMessage("Ca78bf89fa33b777e54b4c13695818f81", linebot.NewTextMessage("這裡純測試對嗎？\n只發於測試聊天室「test」")).Do(); err != nil {
					// 					log.Print(err)
					// 				}
					//target_id_code := event.Source.UserID + event.Source.GroupID + event.Source.RoomID	//target_id_code := ""
			    imageURL = SystemImageURL
				//template := LineTemplate_firstinfo
				t_msg := "請用最新版本的 LINE APP 觀看效果。\n如有其他建議或想討論，請對我輸入「開發者」進行聯絡。"
				obj_message := linebot.NewTemplateMessage(t_msg, LineTemplate_firstinfo)

				//reply 的寫法
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("群組聊天的各位大家好哇～！\n" + push_string + "\n\n想知道我的功能，請說：「簡介」\n單獨輸入「聖經」可以知道查詢方法。"),linebot.NewTextMessage("這是一種資訊整合的便捷應用，效果類似於自動回話小助理。\n\n概念上最基本的應用類似於遊戲 NPC 或 0800 電話總機，會根據指示自動回覆相關基本資訊。\n也可做其他延伸應用，像是聖經查詢 或 留言給意見...等等。\n\n目前除了教會相關資訊外，還可查詢 24 本聖經。\n支援 10 種語言、24 種聖經版本的精準經節查詢機能。\n並支援範圍查詢的寫法。（例如：聖經 創世紀 1:1-10）\n\n詳細說明可輸入「聖經」，有完整的使用說明介紹。"),obj_message).Do(); err != nil {
						log.Print(1351)
						log.Print(err)
				}
		}
		//觸發離開群組聊天
		if event.Type == linebot.EventTypeLeave {
				HttpPost_JANDI("被請離開 "  + user_talk , "gray" , "LINE 離開群組",target_id_code)
				HttpPost_IFTTT("被請離開 "  + user_talk , "LINE 離開群組",target_id_code)
				HttpPost_Zapier("被請離開 "  + user_talk , "LINE 離開群組",target_id_code,user_talk)
				log.Print("觸發被踢出 " + user_talk +  " 群組")
		}
		//？？？？？
			//https://admin-official.line.me/beacon/register
			//https://devdocs.line.me/en/#line-beacon
			//https://devdocs.line.me/ja/#line-beacon
			//這個應用要有硬體
		if event.Type == linebot.EventTypeBeacon {
			HttpPost_JANDI(target_item + " " + "[" + user_talk + "](" + userImageUrl + ") 觸發 Beacon（啥鬼）" + `\n` + userStatus, "yellow" , "LINE 對話同步",target_id_code)
			HttpPost_IFTTT(target_item + " " + user_talk + " 觸發 Beacon（啥鬼）" + `\n<br>` + userImageUrl + `\n<br>` + userStatus, "LINE 對話同步",target_id_code)
			HttpPost_Zapier(target_item + " " + "[" + user_talk + "](" + userImageUrl + ") 觸發 Beacon（啥鬼）" + `\n` + userStatus, "LINE 對話同步",target_id_code,user_talk)
			log.Print(user_talk + " 觸發 Beacon（啥鬼）")
		}
		//觸發收到訊息
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:

			case *linebot.ImageMessage:
				log.Print("對方丟圖片 message.ID = " + message.ID)

				//log.Print("對方丟圖片 linebot.EventSource = " + linebot.EventSource

				//----------------------------------------------------------------取得使用者資訊的寫法
				// source := event.Source

				// userID := event.Source.UserID
				// groupID := event.Source.GroupID
				// RoomID := event.Source.RoomID
				// markID := userID + groupID + RoomID
				
				// log.Print(source.UserID)
				//----------------------------------------------------------------取得使用者資訊的寫法

				// username := ""
				// if markID == "U6f738a70b63c5900aa2c0cbbe0af91c4"{//if source.UserID == "U6f738a70b63c5900aa2c0cbbe0af91c4"{
				// 	username = "LL = " + userID + groupID + RoomID //2016.12.20+
				// }
				// if markID == "Uf150a9f2763f5c6e18ce4d706681af7f"{
				// 	username = "包包"
				// }

				//https://devdocs.line.me/en/#get-content
				//[GAE/GoでLineBotをつくったよ〜 - ベーコンの裏](http://sun-bacon.hatenablog.com/entry/2016/10/10/233520)
				content, err := bot.GetMessageContent(message.ID).Do()
				if err != nil {
					log.Print(2141)
					log.Print(err)
				}
				defer content.Content.Close()
				log.Print("content.ContentType = " + content.ContentType)
				log.Print("content.ContentLength = ")
				log.Print(content.ContentLength) //檔案大小??
				log.Print("content.Content = ")
				log.Print(content.Content)

				//https://github.com/line/line-bot-sdk-go/blob/master/linebot/get_content_test.go
				//ContentLength
				//https://golang.org/pkg/image/jpeg/

				//目標是把 content.Content 存起來

                image, err := jpeg.Decode(content.Content)
                if err != nil {
                	log.Print(2167)
                    log.Print(err)
                }
                log.Printf("image %v", image.Bounds())
                //http://ithelp.ithome.com.tw/articles/10161612
                //https://webcache.googleusercontent.com/search?q=cache:cLTwZS5RNmMJ:https://libraries.io/go/github.com%252Fline%252Fline-bot-sdk-go%252Flinebot+&cd=6&hl=zh-TW&ct=clnk&gl=tw



				var imgByte []byte
				imgByte, err = ioutil.ReadAll(content.Content)
				if err != nil {
					log.Print(err)
				}

				log.Print(imgByte)

                //暫時放棄 = =

									// file, err := ioutil.TempFile("temp.jpg", "")
									// if err != nil {
									// 	log.Print(2175)
									// 	log.Print(err)
									// }
									// defer file.Close()
									
									// _, err = ioutil.WriteFile("temp.jpg", []byte(image.Bounds()), 0600)//io.Copy(file, content.Content)
									// if err != nil {
									// 	log.Print(2182)
									// 	log.Print(err)
									// }
									// log.Printf("Saved %s", file.Name())


                //可以
				// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("這圖片是？\n\n" + username + "你丟給我圖片幹嘛！\n我眼睛還沒長好看不懂XD")).Do(); err != nil {
				// 	log.Print(1845)
				// 	log.Print(err)
				// }
			case *linebot.VideoMessage:
				// //https://github.com/dongri/line-bot-sdk-go
			 //    originalContentURL := "https://dl.dropboxusercontent.com/u/358152/linebot/resource/video-original.mp4"
			 //    previewImageURL := "https://dl.dropboxusercontent.com/u/358152/linebot/resource/video-preview.png"
			 //    obj_message := linebot.NewVideoMessage(originalContentURL, previewImageURL)
 			// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("這影片是？\n我也給你影片吧！\n\n這只是測試功能"),obj_message).Do(); err != nil {
 			// 		log.Print(1854)
 			// 		log.Print(err)
 			// 	}
			case *linebot.AudioMessage:
				// //下面都是 OK 的寫法，但是還是沒辦法取得...........
				// //另外因為現在這個專案不適合這樣玩
				// originalContentURL := "https://dl.dropboxusercontent.com/u/358152/linebot/resource/ok.m4a"
				// duration := 1000
				// obj_message := linebot.NewAudioMessage(originalContentURL, duration)
 			// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("這是什麼聲音？"),obj_message).Do(); err != nil {
 			// 		log.Print(1862)
 			// 		log.Print(err)
 			// 	}
			case *linebot.LocationMessage:
				log.Print("message.Title = " + message.Title)
				log.Print("message.Address = " + message.Address)
				log.Print("message.Latitude = ")
				log.Print(message.Latitude)
				log.Print("message.Longitude = ")
				log.Print(message.Longitude)
				//obj_message := linebot.NewLocationMessage(message.Title, message.Address, message.Latitude, message.Longitude)
				obj_message_map := linebot.NewLocationMessage("台北公館教會", "11677 台北市汀州路四段85巷2號", 25.007408,121.537688) //台北市信義區富陽街46號

				//case 1
				//obj_message_1 := linebot.NewLocationMessage("歡迎光臨", "地球", 25.022413, 121.556427) //台北市信義區富陽街46號
					//obj_message_2 := linebot.NewLocationMessage("歡迎光臨", "哪個近", 25.022463, 121.556454) //這個遠

				//if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("你好！小天使來幫你帶路！\n你在這裡？\n要看看附近的教會嗎？\nhttps://www.google.com/maps/search/%E6%95%99%E6%9C%83/@" + fmt.Sprintf("%f",message.Latitude) + "%2C" + fmt.Sprintf("%f",message.Longitude) + ",15z\n以下是我剛剛收到的 GPS 定位地址！"),obj_message,linebot.NewTextMessage("我們教會在這裡～\n為您預備導航路線圖：\nhttp://maps.google.com/maps?f=d&saddr=" + fmt.Sprintf("%f",message.Latitude) + "," + fmt.Sprintf("%f",message.Longitude) + "&daddr=25.007408,121.537688&hl=zh-tw&dirflg=&sort=num&mrsp=0&doflg=ptk&ttype=now\n以下是我們教會的確切地址！"),obj_message_map,linebot.NewStickerMessage("2", "514")).Do(); err != nil {
				if target_item == "群組對話"{
					LineTemplate_address := linebot.NewCarouselTemplate(
						linebot.NewCarouselColumn(
							imageURL, "小天使來幫你帶路！", "你在「" + message.Address + "」？",
							linebot.NewURITemplateAction("看看附近的教會？", "https://www.google.com/maps/search/%E6%95%99%E6%9C%83/@" + fmt.Sprintf("%f",message.Latitude) + "%2C" + fmt.Sprintf("%f",message.Longitude) + ",15z"),
							linebot.NewURITemplateAction("帶你去公館教會","http://maps.google.com/maps?f=d&saddr=" + fmt.Sprintf("%f",message.Latitude) + "," + fmt.Sprintf("%f",message.Longitude) + "&daddr=25.007408,121.537688&hl=zh-tw&dirflg=&sort=num&mrsp=0&doflg=ptk&ttype=now"),
							linebot.NewPostbackTemplateAction("公館教會地圖","公館教會地圖 POST", "教會地圖"),
						),
						linebot.NewCarouselColumn(
							imageURL, "教會資訊", "我可以幫大家取得教會資訊。",
							linebot.NewPostbackTemplateAction("本週週報 & 聚會時間", "週報 POST","週報"),
							linebot.NewPostbackTemplateAction("網站資訊","網站資訊 POST", "網站資訊"),
							linebot.NewPostbackTemplateAction("聯絡資訊","聯絡資訊 POST", "聯絡資訊"),
						),
						LineTemplate_CarouselColumn_bible_one,
						// LineTemplate_other_example,
						// LineTemplate_other,
						LineTemplate_CarouselColumn_feedback,
					)
					t_address := "帶你去公館教會！\nhttp://maps.google.com/maps?f=d&saddr=" + fmt.Sprintf("%f",message.Latitude) + "," + fmt.Sprintf("%f",message.Longitude) + "&daddr=25.007408,121.537688&hl=zh-tw&dirflg=&sort=num&mrsp=0&doflg=ptk&ttype=now"
					obj_message_address := linebot.NewTemplateMessage(t_address, LineTemplate_address)
					if _, err = bot.ReplyMessage(
						event.ReplyToken, 
						//linebot.NewStickerMessage("2", "514"),
						//linebot.NewTextMessage("你好！小天使來幫你帶路！\n你在「" + message.Address + "」？\n要看看附近的教會嗎？\nhttps://www.google.com/maps/search/%E6%95%99%E6%9C%83/@" + fmt.Sprintf("%f",message.Latitude) + "%2C" + fmt.Sprintf("%f",message.Longitude) + ",15z"),
						// obj_message,
						//linebot.NewTextMessage("我們教會在這裡～\n為您預備導航路線圖：\nhttp://maps.google.com/maps?f=d&saddr=" + fmt.Sprintf("%f",message.Latitude) + "," + fmt.Sprintf("%f",message.Longitude) + "&daddr=25.007408,121.537688&hl=zh-tw&dirflg=&sort=num&mrsp=0&doflg=ptk&ttype=now\n下面是我們教會的地圖！"),
						//obj_message_map,
						obj_message_address,
					).Do(); err != nil {
						log.Print(18766)
						log.Print(err)
					}
				}else{
					if _, err = bot.ReplyMessage(
						event.ReplyToken, 
						linebot.NewStickerMessage("2", "514"),
						linebot.NewTextMessage("你好！小天使來幫你帶路！\n你在「" + message.Address + "」？\n要看看附近的教會嗎？\nhttps://www.google.com/maps/search/%E6%95%99%E6%9C%83/@" + fmt.Sprintf("%f",message.Latitude) + "%2C" + fmt.Sprintf("%f",message.Longitude) + ",15z"),
						// obj_message,
						linebot.NewTextMessage("我們教會在這裡～\n為您預備導航路線圖：\nhttp://maps.google.com/maps?f=d&saddr=" + fmt.Sprintf("%f",message.Latitude) + "," + fmt.Sprintf("%f",message.Longitude) + "&daddr=25.007408,121.537688&hl=zh-tw&dirflg=&sort=num&mrsp=0&doflg=ptk&ttype=now\n下面是我們教會的地圖！"),
						obj_message_map,
					).Do(); err != nil {
						log.Print(1876)
						log.Print(err)
					}
				}
			case *linebot.StickerMessage:
				log.Print("message.PackageID = " + message.PackageID)
				log.Print("message.StickerID = " + message.StickerID)
					//https://github.com/line/line-bot-sdk-go/blob/master/examples/kitchensink/server.go handleSticker
					//message.PackageID, message.StickerID
				//丟跟對方一樣的貼圖回他
			// obj_message_moto := linebot.NewStickerMessage(message.PackageID, message.StickerID)
					//https://github.com/line/line-bot-sdk-go/blob/master/examples/kitchensink/server.go
					//2016.12.20+ 多次框框的方式成功！（最多可以五個）
					//.NewStickerMessage 發貼貼圖成功	 //https://devdocs.line.me/files/sticker_list.pdf			
				obj_message := linebot.NewStickerMessage("2", "514") //https://devdocs.line.me/en/?go#send-message-object
 				//if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("OU<"),linebot.NewTextMessage("0.0"),linebot.NewTextMessage("．ω．"),linebot.NewTextMessage("．ω．")).Do(); err != nil {

				PackageID_int := 0
				StickerID_int := 0
				if PackageID_int, err = strconv.Atoi(message.PackageID); err != nil {
					log.Print("7793 字串轉整數失敗")
					log.Print(PackageID_int)
					log.Print(err.Error())
				}

				if StickerID_int, err = strconv.Atoi(message.StickerID); err != nil {
					log.Print("7798 字串轉整數失敗")
					log.Print(StickerID_int)
					log.Print(err.Error())
				}

				//特別處理過貼圖範圍外的貼圖
				if (PackageID_int!=0) && (PackageID_int<=4){
					// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("．ω．"),obj_message_moto,obj_message).Do(); err != nil {
					// 	log.Print(7806)
					// 	log.Print(err)
					// }
					if _, err = bot.ReplyMessage(event.ReplyToken, obj_message).Do(); err != nil {
						log.Print(7806)
						log.Print(err)
					}
				}else{
					// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("．ω．"),obj_message).Do(); err != nil {
					// 	log.Print(7811)
					// 	log.Print(err)
					// }
					if(username=="LL"){
						//https://store.line.me/stickershop/product/1021884/zh-Hant
						//if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("PackageID_int = " + strconv.Itoa(PackageID_int) + "\nStickerID_int = " + strconv.Itoa(StickerID_int)), obj_message).Do(); err != nil {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("https://store.line.me/stickershop/product/" + strconv.Itoa(PackageID_int) + "/zh-Hant"), obj_message).Do(); err != nil {	
							log.Print(7806)
							log.Print(err)
						}
					}else{

					}
				}

					// if _, err = bot.ReplyMessage(event.ReplyToken, obj_message).Do(); err != nil {
					// 	log.Print(7806)
					// 	log.Print(err)
					// }
			}
		}
	}
}

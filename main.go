package main

import (
    "log"
    "math/rand"
    "time"
    "net/http"
    "os"
    "fmt"
    // "flag"
	"github.com/gin-gonic/gin"
    "github.com/line/line-bot-sdk-go/linebot"  // ① SDKを追加

    "database/sql"
		_ "github.com/go-sql-driver/mysql"
)

// 可愛い系構造体
type Cute struct {
    id        int
    name      string
    url       string
    pic_path  string
    picture_id int
}

// 綺麗系構造体
type Beautiful struct {
    id        int
    name      string
    url       string
    pic_path  string
    picture_id int
}

// ユーザーログ構造体
type userlogs struct {
    id        int
    userID    string
    state     int
    log1      int
    finish    int
}




func main() {
	// 環境変数取得　heroku config:setで設定済み
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    db, err := sql.Open("mysql", "")
	if err != nil {
        log.Fatal("DB is not valid")
    }
    defer db.Close()

    fmt.Println("DB接続完了")
    // ② LINE bot instanceの作成
    bot, err := linebot.New(
        os.Getenv("CHANNEL_SECRET"),
        os.Getenv("CHANNEL_TOKEN"),
    )
    if err != nil {
        log.Fatal(err)
    }
    router := gin.New()
    router.Use(gin.Logger())
    router.LoadHTMLGlob("templates/*.tmpl.html")
    router.Static("/static", "static")
    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl.html", nil)
    })

    // ③ LINE Messaging API用の Routing設定
    router.POST("/callback", func(c *gin.Context) {
        fmt.Println("got a message")
        events, err := bot.ParseRequest(c.Request)
        if err != nil {
            if err == linebot.ErrInvalidSignature {
                log.Print(err)
            }
            return
        }

        // todo  リッチメニューあとで
        // R := setRichmenu()
        // fmt.Println(R)

        // var filePath = flag.String("./top.png", "", "path to image, used in upload/download mode")
        // var rid  = flag.String(R, "", "richmenu id")

        // // fmt.Println(U)
        // if _, err = bot.UploadRichMenuImage(*rid, *filePath).Do(); err != nil {
		// 	log.Fatal(err)
		// }

        // if _, err = bot.LinkUserRichMenu("", *rid).Do(); err != nil {
		// 	log.Fatal(err)
        // }




        // event にはidexが入るので、値は_に詰める

        for _, event := range events {
            // user_id　取得

            userID := event.Source.UserID
            state :=1

            stmtIns, err := db.Prepare(fmt.Sprintf("INSERT INTO userlogs(userID,state) VALUES (?, ?)"))
            if err != nil {
                log.Fatal(err)
            }
            defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

            fmt.Println(userID)

            _, err = stmtIns.Exec(userID, state)

            fmt.Println("insert fin")

            // event setting
            if event.Type == linebot.EventTypeMessage {
                switch message := event.Message.(type){
                // メッセージがText stringで入ってくる場合
                case *linebot.TextMessage:

                    // DO. で　http.Response型のポインタ（とerror）が返ってくる
                    // ReplyMessage関数呼ぶ
                    if message.Text == "こんにちは" {
                        beautyImage := "https://www.vivi.tv/uploads/images/20171004141324_glmjrivv6crr95qhalek0usn45_main_yoshioka.jpg"
                        cuteImage :="https://www.cinemacafe.net/imgs/thumb_h1/58888.jpg"

                        template := linebot.NewCarouselTemplate(
                            linebot.NewCarouselColumn(
                                cuteImage, "可愛い系", "好きならタップ↓↓↓!!",
                                linebot.NewMessageAction("可愛い系", "可愛い系"),
                            ),

                            linebot.NewCarouselColumn(
                                beautyImage, "綺麗系", "好きならタップ↓↓↓!!",
                                linebot.NewMessageAction("綺麗系", "綺麗系"),
                            ),
                        )

                        if _, err := bot.ReplyMessage(
                            event.ReplyToken,
                            linebot.NewTemplateMessage("Carousel alt text", template),
                        ).Do(); err != nil {
                            log.Print(err)
                        }
                    }

                    if message.Text == "可愛い系"{
                        rows, err := db.Query("SELECT * FROM cute")
                        if err != nil {
                            log.Fatal(err)
                        }
                        fmt.Println("クエリ取ってくる")
                        cute := make([]Cute, 0)

                        //  rowsを文字列に成形
                        for rows.Next() {
                            c := Cute{}
                            if err := rows.Scan(&c.id, &c.name, &c.url, &c.pic_path, &c.picture_id); err != nil {
                                log.Fatal("Data is not correct")
                                log.Fatal(err)
                            }
                            cute = append(cute, c)
                        }
                        rand.Seed(time.Now().UnixNano())
                        index := rand.Intn(len(cute))

                        fmt.Println(index)

                        theCute := cute[index]

                        fmt.Println(theCute.url)
                        imagePic := theCute.pic_path
                        template := linebot.NewCarouselTemplate(
                            linebot.NewCarouselColumn(
                                imagePic, theCute.name + "さん" , "インスタも見ましょう！",
                                linebot.NewURIAction("インスタを見る", theCute.url),
                                linebot.NewMessageAction("他の子も見たい", "可愛い系"),
                                linebot.NewMessageAction("この子がタイプ", "この子がタイプ"),
                            ),
                        )

                        var pictureID int = theCute.picture_id
                        stmtIns, err := db.Prepare(fmt.Sprintf("UPDATE userlogs SET log1 = ? WHERE (userID = ?)"))
                        _, err = stmtIns.Exec(pictureID, userID)
                        if err != nil {
                            log.Print(err)
                        }

			            defer stmtIns.Close()

                        if _, err := bot.ReplyMessage(
                            event.ReplyToken,
                            linebot.NewTemplateMessage("Carousel alt text", template),
                        ).Do(); err != nil {
                            log.Print(err)
                        }
                        // rows, err := db.Query("SELECT * FROM sample")
                        // stmtIns, err := db.Prepare(fmt.Sprintf("INSERT INTO userlogs(userID,state) VALUES (?, ?)"))

                    }




                    if message.Text == "綺麗系"{
                        rows, err := db.Query("SELECT * FROM beautiful")
                        if err != nil {
                            fmt.Println("綺麗系テーブル")
                            log.Fatal(err)
                        }

                        beautiful := make([]Beautiful, 0)

                        for rows.Next() {
                            b := Beautiful{}
                            if err := rows.Scan(&b.id, &b.name, &b.url, &b.pic_path, &b.picture_id); err != nil {
                                log.Fatal("Data is not correct")
                            }
                            beautiful = append(beautiful, b)
                        }

                        rand.Seed(time.Now().UnixNano())
                        index := rand.Intn(len(beautiful))

                        theBeautiful := beautiful[index]

                        fmt.Println(theBeautiful)
                        imagePic := theBeautiful.pic_path
                        template := linebot.NewCarouselTemplate(
                            linebot.NewCarouselColumn(
                                imagePic, theBeautiful.name + "さん" , "インスタも見ましょう！",
                                linebot.NewURIAction("インスタを見る", theBeautiful.url),
                                linebot.NewMessageAction("他の子も見たい", "綺麗系"),
                                linebot.NewMessageAction("この子がタイプ", "この子がタイプ"),
                            ),
                        )

                        if _, err := bot.ReplyMessage(
                            event.ReplyToken,
                            linebot.NewTemplateMessage("Carousel alt text", template),
                        ).Do(); err != nil {
                            log.Print(err)
                        }
                    }

                    if message.Text == "この子がタイプ"  {
                        if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("また遊びに来てくださいね！")).Do(); err != nil {
                            log.Print(err)
                        }

                        stmtIns, err := db.Prepare(fmt.Sprintf("UPDATE userlogs SET finish = 1 WHERE (userID = ?)"))
                        _, err = stmtIns.Exec(userID)
                        if err != nil {
                            log.Print(err)
                        }

			            defer stmtIns.Close()
                        break
                    }

                    if message.Text != "こんにちは" || message.Text != "綺麗系" || message.Text != "可愛い系" || message.Text != "この子がタイプ"{
                        if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("可愛い系or綺麗系を選んでね")).Do(); err != nil {
                            log.Print(err)
                        }
                    }
                }
            }
        }
    })
    router.Run(":" + port)
}



// // 綺麗系ランダム生成
// func choice(beatiful map[int]string) int {
//     // 乱数作成　ここ課題　乱数の作り方
//     rand.Seed(time.Now().UnixNano())
//     i := rand.Intn(len(beatiful))
//     // 乱数確認ログ
//     fmt.Println(i)
//     return i
// }




// リッチメニュー
func setRichmenu() string {

    bot, err := linebot.New(
        os.Getenv("CHANNEL_SECRET"),
        os.Getenv("CHANNEL_TOKEN"),
    )
    if err != nil {
        log.Fatal(err)
    }



    richMenu := linebot.RichMenu{

        Size: linebot.RichMenuSize{Width: 2500, Height: 1686},
        Selected: false,
        Name: "Nice richmenu",
        ChatBarText: "Tap here",
        Areas: []linebot.AreaDetail{
            {
                Bounds: linebot.RichMenuBounds{X: 0, Y: 0, Width: 1250, Height: 843},
                Action: linebot.RichMenuAction{
                    Type: linebot.RichMenuActionTypePostback,
                    Data: "action=buy&itemid=123",
                },
            },
            // {
            //     Bounds: linebot.RichMenuBounds{X: 1250, Y: 843, Width: 1250, Height: 843},
            //     Action: linebot.RichMenuAction{
            //         Type: linebot.RichMenuActionTypeDatetimePicker,
            //         Data: "datetime picker!",
            //     },
            // },
        },
    }
    res, err := bot.CreateRichMenu(richMenu).Do()
	if err != nil {
		log.Fatal(err)
    }
    log.Println(res.RichMenuID)
    fmt.Println("リッチ")
    fmt.Println(res.RichMenuID)
    richId := res.RichMenuID

    return richId
}

//  紹介ログインデータインサート　のちにリファクタリング

// func getUserFirstState(userID string state int)　(result bool , err error){

//     stmtIns, err := db.Prepare(fmt.Sprintf("INSERT INTO userlogs(userID,state) VALUES (?, ?)"))
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

//     fmt.Println(userID)

//     _, err = stmtIns.Exec(userID, state)

//     fmt.Println("insert fin")

//     return result
// }





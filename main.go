package main

import (
    "log"
    "math/rand"
    "time"
    "net/http"
    "os"
    "fmt"
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

// 綺麗系構造体
type Cool struct {
    id        int
    name      string
    url       string
    pic_path  string
    picture_id int
}

// 綺麗系構造体
type Lovely struct {
    id        int
    name      string
    url       string
    pic_path  string
    picture_id int
}

// ユーザーログ構造体
type userLogs struct {
    id        int
    gender    int
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
    log.Print(db)

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

        for _, event := range events {
        // todo  リッチメニューあとで
        userID := event.Source.UserID
        var R string = setRichmenu()
        fmt.Println(R)

        var filePath string = "./top.png"
        var rid  string = R

        // /U := "/Users/kuriokatomomi/go/src/github.com/heroku/go-getting-started/top.png"
        // fmt.Println(U)
        fmt.Println(rid)
        fmt.Println("パース通過")
        fmt.Println(filePath)

        if _, err = bot.UploadRichMenuImage(rid, filePath).Do(); err != nil {
			log.Fatal(err)
        }

        fmt.Println("filepath通過")

        fmt.Println("getRichMenu")
        if _, err = bot.GetRichMenu(rid).Do(); err != nil {
			log.Fatal(err)
        }
        fmt.Println("gotRichMenu")


        fmt.Println(userID)
        fmt.Println(rid)

        if _, err = bot.LinkUserRichMenu(userID, rid).Do(); err != nil {
			log.Fatal(err)
        }
        fmt.Println("PreviewRichMenu")
        // if _, err = bot.LinkUserRichMenu("Ucbc84552b616bccf726859cd416548d9", *rid).Do(); err != nil {
		// 	log.Fatal(err)
        // }

        // event にはidexが入るので、値は_に詰める

        // for _, event := range events {
            // user_id　取得

            // userID := event.Source.UserID
            var state int =1

            sum :=0
            sum += 0
            fmt.Println(sum)

            rows, err := db.Query("SELECT COUNT(*) as count  FROM userlogs WHERE userID = ?", userID)
            if err != nil {
                log.Fatal(err)
            }
            defer rows.Close()
            var checkCount int = checkCount(rows)
            fmt.Println(checkCount)

            if checkCount < 1 {
                stmtIns, err := db.Prepare(fmt.Sprintf("INSERT INTO userlogs(userID,state) VALUES (?, ?)"))
                if err != nil {
                    log.Fatal(err)
                }
                defer stmtIns.Close() // Close the statement when we leave main() / the program terminates
                _, err = stmtIns.Exec(userID, state)

                fmt.Println("insert fin")
            }




            // event setting
            if event.Type == linebot.EventTypeMessage {
                switch message := event.Message.(type){
                // メッセージがText stringで入ってくる場合
                case *linebot.TextMessage:
                    if message.Text == "START" && checkCount < 1  {
                        template := linebot.NewConfirmTemplate(
                                "性別を選択",
                                linebot.NewMessageAction("女性", "女性"),
                                linebot.NewMessageAction("男性", "男性"),
                        )

                        if _, err := bot.ReplyMessage(
                            event.ReplyToken,
                            linebot.NewTemplateMessage("Confirm alt text", template),
                        ).Do(); err != nil {
                            log.Print(err)
                        }
                    }

                    // DO. で　http.Response型のポインタ（とerror）が返ってくる
                    // ReplyMessage関数呼ぶ
                    if message.Text == "女性" || message.Text == "男性" || message.Text == "START" {

                        if (message.Text == "女性") {
                            stmtIns, err := db.Prepare(fmt.Sprintf("UPDATE userlogs SET gender = 1 WHERE (userID = ?)"))
                            _, err = stmtIns.Exec(userID)
                            if err != nil {
                                log.Print(err)
                            }
                            defer stmtIns.Close()

                        }else if(message.Text == "男性"){
                            stmtIns, err := db.Prepare(fmt.Sprintf("UPDATE userlogs SET gender = 2 WHERE (userID = ?)"))
                            _, err = stmtIns.Exec(userID)
                            if err != nil {
                                log.Print(err)
                            }
                            defer stmtIns.Close()
                        }else{
                            fmt.Println("start")
                            stmtIns, err := db.Prepare(fmt.Sprintf("UPDATE userlogs SET log1 = 0, finish = 0 WHERE (userID = ?)"))
                            _, err = stmtIns.Exec(userID)
                            if err != nil {
                                log.Print(err)
                            }
                            defer stmtIns.Close()
                        }

                        // 凛
                        beautyImage := "https://www.yutori528.com/wp-content/uploads/2017/11/AS20150720002647_comm.jpg"
                        // 艶
                        cuteImage :="https://www.yutori528.com/wp-content/uploads/2018/05/1-3.jpg"
                        // 萌
                        lovelyImage :="https://www.asahicom.jp/articles/images/AS20171226001781_comm.jpg"
                        // 清
                        coolImage :="https://www.lespros.co.jp/files/talent/4/profile.jpg"

                        template := linebot.NewCarouselTemplate(
                            linebot.NewCarouselColumn(
                                beautyImage, "凛", "好きならタップ↓↓↓!!",
                                linebot.NewMessageAction("凛", "凛"),
                            ),

                            linebot.NewCarouselColumn(
                                cuteImage, "艶", "好きならタップ↓↓↓!!",
                                linebot.NewMessageAction("艶", "艶"),
                            ),

                            linebot.NewCarouselColumn(
                                lovelyImage, "萌", "好きならタップ↓↓↓!!",
                                linebot.NewMessageAction("萌", "萌"),
                            ),

                            linebot.NewCarouselColumn(
                                coolImage, "清", "好きならタップ↓↓↓!!",
                                linebot.NewMessageAction("清", "清"),
                            ),
                        )

                        if _, err := bot.ReplyMessage(
                            event.ReplyToken,
                            linebot.NewTemplateMessage("Carousel alt text", template),
                        ).Do(); err != nil {
                            log.Print(err)
                        }
                    }

                    // 凛

                    if message.Text == "凛"{
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
                        // var theCuteId int = theCute.id

                        // fmt.Println(theCuteId)

                        // var cuteCount []int = make([]int, 0)
                        // cuteCount = append(cuteCount, theCuteId)

                        // fmt.Println(cuteCount)



                        fmt.Println(theCute.url)
                        imagePic := theCute.pic_path
                        template := linebot.NewCarouselTemplate(
                            linebot.NewCarouselColumn(
                                imagePic, theCute.name + "さん" , "インスタも見ましょう！",
                                linebot.NewURIAction("インスタを見る", theCute.url),
                                linebot.NewMessageAction("他の子も見たい", "凛"),
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
                    }



                    // 艶
                    if message.Text == "艶"{
                        rows, err := db.Query("SELECT * FROM beautiful")
                        if err != nil {
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
                                linebot.NewMessageAction("他の子も見たい", "艶"),
                                linebot.NewMessageAction("この子がタイプ", "この子がタイプ"),
                            ),
                        )

                        var pictureID int = theBeautiful.picture_id
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
                        log.Print(err)
                        break
                    }



                    // 萌
                    if message.Text == "萌"{
                        rows, err := db.Query("SELECT * FROM lovely")
                        if err != nil {
                            log.Fatal(err)
                        }

                        lovely := make([]Lovely, 0)

                        for rows.Next() {
                            l := Lovely{}
                            if err := rows.Scan(&l.id, &l.name, &l.url, &l.pic_path, &l.picture_id); err != nil {
                                log.Fatal("Data is not correct")
                            }
                            lovely = append(lovely, l)
                        }

                        rand.Seed(time.Now().UnixNano())
                        index := rand.Intn(len(lovely))

                        theLovely := lovely[index]

                        fmt.Println(theLovely)
                        imagePic := theLovely.pic_path
                        template := linebot.NewCarouselTemplate(
                            linebot.NewCarouselColumn(
                                imagePic, theLovely.name + "さん" , "インスタも見ましょう！",
                                linebot.NewURIAction("インスタを見る", theLovely.url),
                                linebot.NewMessageAction("他の子も見たい", "萌"),
                                linebot.NewMessageAction("この子がタイプ", "この子がタイプ"),
                            ),
                        )

                        var pictureID int = theLovely.picture_id
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
                    }


                    // 清
                    if message.Text == "清"{
                        rows, err := db.Query("SELECT * FROM cool")
                        if err != nil {
                            log.Fatal(err)
                        }

                        cool := make([]Cool, 0)

                        for rows.Next() {
                            c := Cool{}
                            if err := rows.Scan(&c.id, &c.name, &c.url, &c.pic_path, &c.picture_id); err != nil {
                                log.Fatal("Data is not correct")
                            }
                            cool = append(cool, c)
                        }

                        rand.Seed(time.Now().UnixNano())
                        index := rand.Intn(len(cool))

                        theCool := cool[index]

                        fmt.Println(theCool)
                        imagePic := theCool.pic_path
                        template := linebot.NewCarouselTemplate(
                            linebot.NewCarouselColumn(
                                imagePic, theCool.name + "さん" , "インスタも見ましょう！",
                                linebot.NewURIAction("インスタを見る", theCool.url),
                                linebot.NewMessageAction("他の子も見たい", "清"),
                                linebot.NewMessageAction("この子がタイプ", "この子がタイプ"),
                            ),
                        )

                        var pictureID int = theCool.picture_id
                        stmtIns, err := db.Prepare(fmt.Sprintf("UPDATE userlogs SET log1 = ? WHERE (userID = ?)"))
                        _, err = stmtIns.Exec(pictureID, userID)
                        if err != nil {
                            log.Print(err)
                        }

                        if _, err := bot.ReplyMessage(
                            event.ReplyToken,
                            linebot.NewTemplateMessage("Carousel alt text", template),
                        ).Do(); err != nil {
                            log.Print(err)
                        }
                    }


                    if message.Text != "START" || message.Text != "艶" || message.Text != "清" || message.Text != "この子がタイプ" || message.Text != "女性" || message.Text != "男性" || message.Text != "萌" || message.Text != "凛" {
                        if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("凛・艶・萌・清を選んでね")).Do(); err != nil {
                            log.Print(err)
                        }
                    }
                }
            }
        }
    })
    router.Run(":" + port)
}


// リッチメニュー　
func setRichmenu() string {

    bot, err := linebot.New(
        os.Getenv("CHANNEL_SECRET"),
        os.Getenv("CHANNEL_TOKEN"),
    )
    if err != nil {
        log.Fatal(err)
    }


// リッチメニューコンテンツ
    richMenu := linebot.RichMenu{

        Size: linebot.RichMenuSize{Width: 2500, Height: 1686},
        Selected: false,
        Name: "Nice richmenu",
        ChatBarText: "Tap here",
        Areas: []linebot.AreaDetail{
            {
                Bounds: linebot.RichMenuBounds{X: 0, Y: 0, Width: 1250, Height: 843},
                Action: linebot.RichMenuAction{
                    // Type: linebot.RichMenuActionTypePostback,
                    // Data: "action=buy&itemid=123",
                    Type:  "message",
                    Text: "女性",
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

// func createDB(message, userID){
//     if (message.Text == 女性) {
//         stmtIns, err := db.Prepare(fmt.Sprintf("UPDATE userlogs SET log1 = 1 WHERE (userID = ?)"))
//         _, err = stmtIns.Exec(userID)
//         if err != nil {
//             log.Print(err)
//         }
//         defer stmtIns.Close()

//     }else{
//         stmtIns, err := db.Prepare(fmt.Sprintf("UPDATE userlogs SET log1 = 2 WHERE (userID = ?)"))
//         _, err = stmtIns.Exec(userID)
//         if err != nil {
//             log.Print(err)
//         }
//         defer stmtIns.Close()
//     }

// }

func checkCount(rows *sql.Rows) (count int) {
    for rows.Next() {
       err:= rows.Scan(&count)
       if err != nil {
        log.Println("num of count")
        log.Fatal(err)
    }
   }
   return count
}




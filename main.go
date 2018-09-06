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
}

// 綺麗系構造体
type Beautiful struct {
    id        int
    name      string
    url       string
    pic_path  string
}


func main() {
	// 環境変数取得　heroku config:setで設定済み
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    db, err := sql.Open("mysql", "username:pass@tcp(us-cdbr-iron-east-01.cleardb.net:3306)/heroku_2178104d727ee3e")
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

        // event にはidexが入るので、値は_に詰める
        for _, event := range events {
            // event setting
            if event.Type == linebot.EventTypeMessage {
                switch message := event.Message.(type){
                // メッセージがText stringで入ってくる場合
                case *linebot.TextMessage:
                    // DO. で　http.Response型のポインタ（とerror）が返ってくる
                    // ReplyMessage関数呼ぶ
                    if message.Text == "こんにちは" {
                        beautyImage := "https://www.vivi.tv/uploads/images/20171004141324_glmjrivv6crr95qhalek0usn45_main_yoshioka.jpg"
                        cuteImage :="https://gorilla.clinic/cms/wp-content/uploads/2014/11/kasumi.jpg"

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
                            fmt.Println("DBquery")
                            log.Fatal(err)
                        }
                        fmt.Println("クエリ取ってくる")
                        cute := make([]Cute, 0)

                        //  rowsを文字列に成形
                        for rows.Next() {
                            c := Cute{}
                            if err := rows.Scan(&c.id, &c.name, &c.url, &c.pic_path); err != nil {
                                log.Fatal("Data is not correct")
                            }
                            cute = append(cute, c)
                        }
                        rand.Seed(time.Now().UnixNano())
                        index := rand.Intn(len(cute))

                        theCute := cute[index]

                        fmt.Println(theCute.url)
                        imagePic := theCute.pic_path
                        template := linebot.NewCarouselTemplate(
                            linebot.NewCarouselColumn(
                                imagePic, theCute.name + "さん" , "インスタも見ましょう！",
                                linebot.NewURIAction("インスタを見る", theCute.url),
                                linebot.NewMessageAction("他の子も見たい", "可愛い系"),
                            ),
                        )

                        if _, err := bot.ReplyMessage(
                            event.ReplyToken,
                            linebot.NewTemplateMessage("Carousel alt text", template),
                        ).Do(); err != nil {
                            log.Print(err)
                        }
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
                            if err := rows.Scan(&b.id, &b.name, &b.url, &b.pic_path); err != nil {
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
                            ),
                        )

                        if _, err := bot.ReplyMessage(
                            event.ReplyToken,
                            linebot.NewTemplateMessage("Carousel alt text", template),
                        ).Do(); err != nil {
                            log.Print(err)
                        }
                    }
                    if message.Text != "こんにちは" || message.Text != "綺麗系" || message.Text != "可愛い系" {
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
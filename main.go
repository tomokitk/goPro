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
    // "database/sql"
	// 	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 環境変数取得　heroku config:setで設定済み
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    // db, err := sql.Open("mysql", "root:J02M05A004@tcp(127.0.0.1:13306)/gosample" )
	// if err != nil {
    //     log.Fatal("DB is not valid")
    // }
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
                    if message.Text == "綺麗系が好き"{

                        beatifulMap := map[int]string{
                            0:"\n 貴島明日香ちゃんはどうですか？\n\n"+"https://www.instagram.com/asuka_kijima/?hl=ja",
                            1: "\n おゆみちゃんはどうですか？\n\n"+"https://www.instagram.com/youme_mlk/?hl=ja",
                            2: "\n あやプーさんはどうですか？\n\n"+"https://www.instagram.com/ayapooh_22/?hl=ja",
                            3: "\n yuu__aaaちゃんはどうですか？\n\n"+"https://www.instagram.com/yuu__aaa/?hl=ja",
                        }

                        index := choice(beatifulMap)

                        theBeautiful := beatifulMap[index]


                        // db, err := sql.Open("mysql", "root:J02M05A004@tcp(127.0.0.1:13306)/gosample" )

                        // if err != nil {
                        //     fmt.Println("DB接続エラー")
                        //     log.Fatal(err)
                        // }
                        // rows, err := db.Query("SELECT id, name FROM sample")
                        // if err != nil {
                        //     fmt.Println("DBquery")
                        //     log.Fatal(err)
                        // }

                        // for rows.Next() {
                        //     var id int
                        //     var name string
                        //     if err := rows.Scan(&id, &name); err != nil {
                        //         log.Fatal("Data is not correct")
                        //     }
                        //     fmt.Println(id, name)
                        // }
                        // fmt.Println(rows)
                        // db.Close()

                        if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(theBeautiful)).Do(); err != nil {
                            log.Print(err)
                        }
                    }else{
                        if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("https://www.instagram.com/kuriokan/")).Do(); err != nil {
                            log.Print(err)
                        }
                    }
                }
            }
        }
    })
    router.Run(":" + port)
}

// 綺麗系ランダム生成
func choice(beatiful map[int]string) int {
    // 乱数作成　ここ課題　乱数の作り方
    rand.Seed(time.Now().UnixNano())
    i := rand.Intn(len(beatiful))
    // 乱数確認ログ
    fmt.Println(i)
    return i

}
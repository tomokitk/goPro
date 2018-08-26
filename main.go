package main

import (
    "log"
    "net/http"
    "os"
    "fmt"
	"github.com/gin-gonic/gin"
    "github.com/line/line-bot-sdk-go/linebot"  // ① SDKを追加
)

func main() {
	// 環境変数取得　heroku config:setで設定済み
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }
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
                switch message := event.Message.(type) {
                // メッセージがText stringで入ってくる場合
                case *linebot.TextMessage:
                    // DO. で　http.Response型のポインタ（とerror）が返ってくる
                    if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
                        log.Print(err)
                    }
                }
            }
        }
    })

    router.Run(":" + port)
}
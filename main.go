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

// Cute構造体
type Cute struct {
    id        int
    name      string
    url       string
}

func main() {
	// 環境変数取得　heroku config:setで設定済み
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    db, err := sql.Open("mysql", "bfccbf8ad2f3fd:da1ad3db@tcp(us-cdbr-iron-east-01.cleardb.net:3306)/heroku_2178104d727ee3e")
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

                    if message.Text == "かわいい"{
                        if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("性別を教えてね！(男・女　を入力！)")).Do(); err != nil {
                            log.Print(err)
                        }
                    }



                    if message.Text == "綺麗系が好き"{
                        rows, err := db.Query("SELECT * FROM cute")
                        if err != nil {
                            fmt.Println("DBquery")
                            log.Fatal(err)
                        }
                        fmt.Println("クエリ取ってくる")

                        //  rowsを文字列に成形
                        cute := make([]Cute, 0)

                        for rows.Next() {
                            c := Cute{}
                            // var id int
                            // var name string
                            // var url string
                            if err := rows.Scan(&c.id, &c.name, &c.url); err != nil {
                                log.Fatal("Data is not correct")
                            }
                            cute = append(cute, c)
                            // fmt.Println(id, name)
                        }
                        // fmt.Println(len(cute))
                        rand.Seed(time.Now().UnixNano())
                        index := rand.Intn(len(cute))

                        theCute := cute[index]

                        fmt.Println(index)
                        fmt.Println(theCute.url)
                        fmt.Println(theCute.name)



                        // db.Close()

                            // /var cuteName string = name
                            // fmt.Println(cuteName)
                            if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(theCute.name + "さんはどうですか？\n" + theCute.url)).Do(); err != nil {
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
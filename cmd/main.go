package main

import (
	"fmt"
	"github.com/ginkt/newsman/config"
	api_parser "github.com/ginkt/newsman/internal/api-parser"
	"github.com/ginkt/newsman/internal/db"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/mailru/dbr"
)

var (
	dbConn *dbr.Connection
	cfg *config.Config
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		ForceQuote:      true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetLevel(log.TraceLevel)
	cfg = config.InitConfig()

	var err error
	dbConn, err = db.CreatePostgresClient(cfg)
	if err != nil {
		log.Fatalln(err)
	}
}

var subscribeKeyboardMarkup = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Бизнес", "biz"),
		tgbotapi.NewInlineKeyboardButtonData("Развлечение", "fun"),
		tgbotapi.NewInlineKeyboardButtonData("Общее", "gen"),
		tgbotapi.NewInlineKeyboardButtonData("Здоровье", "health"),
		tgbotapi.NewInlineKeyboardButtonData("Наука", "sci"),
		tgbotapi.NewInlineKeyboardButtonData("Спорт", "sport"),
		tgbotapi.NewInlineKeyboardButtonData("Технологии", "tech"),
	),
)

func main() {
	log.Println(cfg.TgBotApiToken)
	bot, err := tgbotapi.NewBotAPI(cfg.TgBotApiToken)
	if err != nil {
		log.Fatalln("bot creation error", err)
	}
	log.Infof("%s bot was initialized", bot.Self.FirstName)

	bot.Debug = true

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(cfg.TgWebhookUrl))
	if err != nil {
		log.Fatalln("webhook error", err)
	}

	parserClient := api_parser.CreateParserClient(&http.Client{})

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe(":8080", nil)
	log.Println("Starting to listen and serve :8080")

	for update := range updates {
		if update.CallbackQuery != nil{
			log.Println("callback query", update)

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID,update.CallbackQuery.Data))

			var callbackMsg tgbotapi.MessageConfig
			switch update.CallbackQuery.Data {
			case "Nudes":
				callbackMsg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Okay buddy, you chose %s", update.CallbackQuery.Data))
			case "Japanese Girls":
				callbackMsg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Nice choice, you chose %s", update.CallbackQuery.Data))
			}

			_, err = bot.Send(callbackMsg)
			if err != nil {
				log.Errorln("error sending message:", err)
				continue
			}
		}
		if update.Message == nil {
			continue
		}

		log.Printf("Got a new message from %d:%s\n[%s]", update.Message.Chat.ID, update.Message.Chat.FirstName, update.Message.Text)

		var msgToSend tgbotapi.MessageConfig
		switch update.Message.Text {
		case "/start":
			msgToSend = tgbotapi.NewMessage(update.Message.Chat.ID, "Хэй, факинг слэйв, выбери одну из тем для подписки, быстрееенько.")
			msgToSend.ReplyMarkup = subscribeKeyboardMarkup
		case "/subscribe":
			msgToSend = tgbotapi.NewMessage(update.Message.Chat.ID, "Окей братан ты выбрал подписаться")
		case "/debug":
			articles := parserClient.GetArticlesByTags([]string{"sport", "tech"})
			log.Printf("Loaded %d tags\n%+v", len(articles), articles)
			for tag, response := range articles {
				if len(response.Articles) == 0 {
					continue
				}

				msgToSend = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("New %s response only for you, bro ^-^\n%s\n%s\n%s", tag,
					response.Articles[0].Description, response.Articles[0].Url, response.Articles[0].UrlToImage))
				_, err = bot.Send(msgToSend)
				if err != nil {
					log.Errorln("error sending message from debug:", err)
					continue
				}
			}
		default:
			msgToSend = tgbotapi.NewMessage(update.Message.Chat.ID, "Что то на богатом")
		}

		//_, err = bot.Send(msgToSend)
		//if err != nil {
		//	log.Errorln("error sending message:", err)
		//	continue
		//}
	}

}

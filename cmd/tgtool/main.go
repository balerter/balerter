package main

import (
	"context"
	"flag"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	token = flag.String("token", "", "Telegram Bot Token")
)

func main() {
	flag.Parse()

	*token = strings.TrimSpace(*token)

	if *token == "" {
		log.Printf("error: empty token")
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		log.Printf("error: error create bot: %v", err)
		os.Exit(1)
	}

	updateConfig := tgbotapi.UpdateConfig{
		Offset:  0,
		Limit:   0,
		Timeout: 0,
	}

	ch, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Printf("error: get updates channel: %v", err)
		os.Exit(1)
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	go listen(ctx, ch)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	<-stop
	ctxCancel()

	log.Printf("bye")
}

func listen(ctx context.Context, ch tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case mes := <-ch:
			if mes.Message == nil {
				continue
			}

			log.Printf("Text '%s' in ChatID %d", mes.Message.Text, mes.Message.Chat.ID)
		}
	}
}

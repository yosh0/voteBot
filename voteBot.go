package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"flag"
	"time"
	"errors"
	"regexp"
	"runtime"
	"strconv"
	"encoding/json"
	"golang.org/x/net/context"
	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
)

var(
	C 		= Config{}
)

func main() {
	token := flag.String("token", C.Serv.Token, "telegram bot token")
	debug := flag.Bool("debug", C.Serv.Debug, "show debug information")
	flag.Parse()
	if *token == "" {
		log.Fatal("token flag is required")
	}

	api := telegram.New(*token)
	api.Debug(*debug)
	bot := telebot.NewWithAPI(api)
	bot.Use(telebot.Recover())
	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.HandleFunc(func(ctx context.Context) error {
		update := telebot.GetUpdate(ctx)
		if update.CallbackQuery != nil {
			if SU, ok := SavedUserLoadCache(update.CallbackQuery.From.ID); ok {
				if SU.EndVote == true {
					err := errors.New("VoteEndTrue1")
					return err
				}
			}
			rex, err := regexp.Compile(`^(variant)_(\d+)$`)
			if err != nil {
				LogFuncStr(fName(), err.Error())
			}
			res := rex.FindStringSubmatch(update.CallbackQuery.Data)
			if res != nil {
				if SU, ok := SavedUserLoadCache(update.CallbackQuery.From.ID); ok {
					LogFuncStr(fName(), "EndVote")
					variant, err := strconv.Atoi(res[2])
					if err != nil {
						LogFuncStr(fName(), err.Error())
					}
					SU.EndVote = true
					SU.VoteVariant = variant
					UpdateSavedUser(SU)
					SU.dbInsert()
				}
				msg := MessageHTML(update, C.BotMessage.VoteEnd, "", false)
				_, err = api.Send(ctx, msg)
				if err != nil {
					LogFuncStr(fName(), err.Error())
					return err
				}
				return nil
			}
		}
/*
		if update.Message == nil {
			return nil
		}
*/
		err := initCustomKeyboard(update, api, ctx)
		if err != nil {
			LogFuncStr(fName(), err.Error())
			return err
		}
		return nil
	})

	bot.Use(telebot.Commands(map[string]telebot.Commander{
		"start": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {
				api := telebot.GetAPI(ctx)
				update := telebot.GetUpdate(ctx)
				if SU, ok := SavedUserLoadCache(update.Message.From.ID); ok {
					if SU.EndVote == true {
						err := errors.New("VoteEndTrue2")
						msg := MessageHTML(update, C.BotMessage.VoteLimit, "", false)
						_, err = api.Send(ctx, msg)
						if err != nil {
							LogFuncStr(fName(), err.Error())
							return err
						}
						return err
					}
				}
				msg := MessageHTML(update, C.StartMsg.Command, C.KbBtnText.Start, false)
				_, err := api.Send(ctx, msg)
				return err
			}),
	}))
	err := bot.Serve(netCtx)
	if err != nil {
		log.Fatal(err)
	}
}

func NewSavedUser(update *telegram.Update) SavedUser {
	SU := SavedUser{
		TgID: update.Message.From.ID,
		UserName: fmt.Sprintf("%s %s %s", update.Message.From.FirstName, update.Message.From.Username, update.Message.From.LastName),
		StartVote: false,
		EndVote: false,
		VoteVariant: 0,
		Category: C.Vote.Category,
		UpdatedAt: time.Now().Unix(),
	}
	SavedUserStoreCache(SU)
	LogSavedUsers(fName(), SU)
	return SU
}

func UpdateSavedUser(SU SavedUser) {
	SavedUserStoreCache(SU)
	LogSavedUsers(fName(), SU)
}

func initCustomKeyboard(update *telegram.Update, api *telegram.API, ctx context.Context) error {
	if SU, ok := SavedUserLoadCache(update.Message.From.ID); ok {
		if SU.EndVote == true {
			err := errors.New("VoteEndTrue2")
			msg := MessageHTML(update, C.BotMessage.VoteLimit, "", false)
			_, err = api.Send(ctx, msg)
			if err != nil {
				LogFuncStr(fName(), err.Error())
				return err
			}
			return err
		}
	}
	if update.Message.Text == C.StartMsg.Start {
		if SU, ok := SavedUserLoadCache(update.Message.From.ID); ok {
			if SU.StartVote == true {
				err := errors.New("VoteStartTrue1")
				msg := MessageHTML(update, C.BotMessage.VoteStarted, "", false)
				_, err = api.Send(ctx, msg)
				if err != nil {
					LogFuncStr(fName(), err.Error())
					return err
				}
				return err
			}
		}
		for elem, photo := range C.Vote.Photos {
			votePictureSet := NewVotePictureSet(elem, photo)
			NewSendPhoto(update, api, ctx, votePictureSet)
		}
		LogFuncStr(fName(), update.Message.Text)
		msg := Keyb(update)
		_, err := api.SendMessage(ctx, msg)
		if err != nil {
			LogFuncStr(fName(), err.Error())
			return err
		}
		SU := NewSavedUser(update)
		SU.StartVote = true
		UpdateSavedUser(SU)
		return nil
	}
	err := errors.New("error1")
	return err
}

func NewVotePictureSet(elem int, photo string) VotePictureSet {
	return VotePictureSet{
		Element: elem + 1,
		MimeType: C.Vote.MimeType,
		PictureUrl: photo,
	}
}

func NewSendPhoto(update *telegram.Update, api *telegram.API, ctx context.Context, votePictureSet VotePictureSet) {
	var r io.Reader
	api.SendPhoto(ctx, telegram.PhotoCfg{
		Caption: strconv.Itoa(votePictureSet.Element),
		BaseFile: telegram.BaseFile{
			BaseMessage: telegram.BaseMessage{
				BaseChat : telegram.BaseChat{
					ID: update.Chat().ID,
					ChannelUsername: strconv.Itoa(int(update.Chat().ID)),
				},
			},
			FileID: votePictureSet.PictureUrl,
			MimeType: votePictureSet.MimeType,
			InputFile: telegram.NewInputFile(votePictureSet.PictureUrl, r),
		},
	})
}

func InitVoteButtons() ([]string, []string) {
	var array1 []string
	var array2 []string
	for elem, _ := range C.Vote.Photos {
		array1 = append(array1, strconv.Itoa(elem + 1))
		array2 = append(array2, fmt.Sprintf("%s%d", "variant_", elem + 1))
	}
	return array1, array2
}

func Keyb(update *telegram.Update) telegram.MessageCfg {
	text := C.BotMessage.VoteVariantSelect
	buttonArray1, buttonArray2 := InitVoteButtons()
	msg := telegram.NewMessage(update.Chat().ID, text)
	msg.ReplyMarkup = telegram.InlineKeyboardMarkup{
		InlineKeyboard: telegram.NewVInlineKeyboard(
			"",
			buttonArray1,
			buttonArray2,
		),
	}
	return msg
}

func MessageHTML(update *telegram.Update, text string, kbtext string, rc bool) telegram.MessageCfg {
	textHtml := TextHTMLBold(text)
	msg := telegram.NewMessage(update.Chat().ID, textHtml)
	msg.ParseMode = "HTML"
	if kbtext == "" {
		return msg
	}
	msg.ReplyMarkup = telegram.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: NewCustomKeyboard(kbtext, false),
	}
	return msg
}

func TextHTMLBold(text string) string {
	return fmt.Sprintf("<b>%s</b>", text)
}

func NewCustomKeyboard(text string, request_contact bool) [][]telegram.KeyboardButton {
	r := make([][]telegram.KeyboardButton, 1)
	r[0] = []telegram.KeyboardButton{
		{
			Text: text,
			RequestContact: request_contact,
			RequestLocation: false,
		},
	}
	return r
}

func init() {
	c, err := os.Open(".config.json")
	if err != nil {
		LogFuncStr(fName(), err.Error())
	}
	decoder := json.NewDecoder(c)
	conf := Config{}
	err = decoder.Decode(&conf)
	C = conf
}

func fName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

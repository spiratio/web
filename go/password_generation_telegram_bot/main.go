package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	Vowels string = "aeiouy"
	Consonants string = "bcdfghjklmnpqrstvwxz"
	SpecialSymbols string = "1234567890!$@%&()_-+*"
	PasswordLength = 10
	PasswordGenerationCommand string = "пароль"
	StartCommand string = "/start"
	AboutBot string = "Бот генератор паролей приветствует Вас! Чтобы сгенерировать пароль, напишите 'Пароль'"
	UndefinedCommandError = "Команда не распознана. Введите новую команду."
)

//---------------------------------------------------------------------//
func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	botToken := "secret key"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0

	for ;; {
		updates, err := getUpdates(botUrl, offset)
		if err != nil{
			log.Println("Smth went wrong: ", err.Error())
		}
		for _, update := range updates {
			err = respond(botUrl, update)
			offset = update.UpdateId + 1
		}
		fmt.Println(updates)
	}
}

//---------------------------------------------------------------------//
func getUpdates(botUrl string ,offset int )([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

//---------------------------------------------------------------------//
func respond(botUrl string ,update Update) error {
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	if checkPasswordCommand(update.Message.Text) {
		botMessage.Text = generatePassword()
	} else if update.Message.Text == StartCommand {
		botMessage.Text = AboutBot
	} else {
		botMessage.Text = UndefinedCommandError
	}
	buf, err := json.Marshal(botMessage)
	if err != nil{
		return err
	}
	_, err = http.Post(botUrl + "/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	return nil
}

//---------------------------------------------------------------------//
func checkPasswordCommand (message string) bool {
	command := strings.ToLower(message)
	return strings.Contains(command, PasswordGenerationCommand)
}

//---------------------------------------------------------------------//
func generatePassword() string {
	var result string
	specialSymbolsCount := rand.Int() % 3 + 1
	var isVowel bool = rand.Int() % 2 == 0

	for i := 0; i < (PasswordLength - specialSymbolsCount); i++ {
		var newLetter string
		if isVowel {
			newLetter = string(Vowels[rand.Int() % len(Vowels)])
		} else {
			newLetter = string(Consonants[rand.Int() % len(Consonants)])
		}
		isVowel = !isVowel
		if rand.Int() % 2 == 0 {
			newLetter = strings.ToUpper(newLetter)
		}
		result = result + newLetter
	}
	for i := 0; i < specialSymbolsCount; i++ {
		result = result + string(SpecialSymbols[rand.Int() % len(SpecialSymbols)])
	}

	return result
}

package main

import (
    "log"
    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "strings"
)

var (
    messages = map[string]string{
        "welcome": "Welcome to the Kdrama Recommendation Bot!\n What's your name?",
        "genre":   "Which genre do you prefer?\n (Romance, Melodrama, Historical, Thriller, Action, Comedy)",
    }
    dramas = map[string][]string{
        "Romance": {"Twinkling Watermelon✨🍉\nA high school student with a passion for music and playing guitar in a rock band walks into a mysterious musical instrument store and is transported from 2023 to 1995.\n https://doramy.club/40627-mercayushhij-arbuz.html"},
        "Melodrama": {"Heirs😏💸\n The series tells the story of young heirs of wealthy families who study at an exclusive school for the elite\n https://kinogo.inc/drama/13744-nasledniki-dorama-hd-dolby2-v13-vs12.html"},
        "Historical": {"Moon Lovers🌑🤍: Scarlet Heart Ryeo\n During a solar eclipse, a girl from the 21st century, Ko Ha-jin, finds herself in the past, during the time of the state of Goryeo. She awakens in the year 941 in the body of a 16-year-old aristocrat, Hae-su. She remembers everything, but cannot go back and is forced to live according to the laws of the old world.\n https://doramy.club/639-dorama-lunnye-vlyublyonnye-alye-serdca-koryo.html"},
        "Thriller": {"Mouse🐁🗡️\n A gripping story set in a world where people can identify psychopaths through DNA testing of fetuses in the womb\n https://doramy.club/24814-7-mysh.html"},
        "Action": {"K2⚖️💪🏼\n Having survived betrayal, a disgraced mercenary soldier codenamed K2 finds himself at the center of the intrigues of dishonest politicians. He is trying to avenge his collapse and save the unfortunate daughter of a presidential candidate.\n https://doramy.club/498-k2-teloxranitel.html"},
        "Comedy": {"Mr. Queen👸🏻🥴\n Jung Bong Hwan works as a chef at the Blue House, but suddenly he finds himself in the Joseon era, and now his soul is trapped in the body of Queen Chorin. King Cheol Jong looks like a puppet in the hands of those who really rule the country, but in fact he hides quite strong traits of his character.\n https://doramy.club/23665-62-koroleva-chorin.html"},
    }
)

type UserState struct {
    Name string
    Step string
}

func capitalizeFirst(s string) string {
    if s == "" {
        return ""
    }
    r := []rune(s)
    return strings.ToUpper(string(r[0])) + string(r[1:])
}

func main() {
    botToken := "6673045880:AAFyEN4jHRjqETu1yIY3HimwXfFIORA5xjE" // Your actual bot token
    bot, err := tgbotapi.NewBotAPI(botToken)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    userStates := make(map[int64]*UserState)

    for update := range updates {
        if update.Message != nil {
            chatID := update.Message.Chat.ID

            if _, ok := userStates[chatID]; !ok {
                userStates[chatID] = &UserState{Step: "name"}
                msg := tgbotapi.NewMessage(chatID, messages["welcome"])
                bot.Send(msg)
                continue
            }

            userState := userStates[chatID]

            switch userState.Step {
            case "name":
                userState.Name = update.Message.Text
                userState.Step = "genre"
                msgText := "Hello, " + update.Message.Text + "! " + messages["genre"]
                msg := tgbotapi.NewMessage(chatID, msgText)
                bot.Send(msg)

            case "genre":
                genreInput := strings.ToLower(update.Message.Text)
                genre := capitalizeFirst(genreInput)
                dramaList, exists := dramas[genre]
                if exists {
                    for _, drama := range dramaList {
                        msg := tgbotapi.NewMessage(chatID, drama)
                        bot.Send(msg)
                    }
                } else {
                    errorMsg := "Sorry, I couldn't find dramas for this genre."
                    msg := tgbotapi.NewMessage(chatID, errorMsg)
                    bot.Send(msg)
                }
                delete(userStates, chatID)
            }
        }
    }
}

package main

import (
    "log"

    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
    messages = map[string]string{
        "welcome": "Welcome to the Kdrama Recommendation Bot! What's your name?",
        "choice":  "What do you prefer: Kdramas or Anime?\nPlease type 'Kdrama' or 'Anime'.",
        "genre":   "Which genre do you prefer?\n(Romance, Melodrama, Historical, Thriller, Action, Comedy)",
        "more":    "Do you want to choose another genre?",
        "thanks":  "Thank you for using the Kdrama Recommendation Bot! Enjoy your viewing!🤍",
    }
    recommendations = map[string]map[string][]string{
        "Kdrama": {
            "Romance":   []string{"Twinkling Watermelon✨🍉\nA high school student with a passion for music and playing guitar in a rock band walks into a mysterious musical instrument store and is transported from 2023 to 1995.\n https://doramy.club/40627-mercayushhij-arbuz.html"},
            "Melodrama": []string{"Heirs😏💸\nThe series tells the story of young heirs of wealthy families who study at an exclusive school for the elite\n https://kinogo.inc/drama/13744-nasledniki-dorama-hd-dolby2-v13-vs12.html"},
            "Historical": []string{"Moon Lovers🌑🤍: Scarlet Heart Ryeo\n During a solar eclipse, a girl from the 21st century, Ko Ha-jin, finds herself in the past, during the time of the state of Goryeo. She awakens in the year 941 in the body of a 16-year-old aristocrat, Hae-su. She remembers everything, but cannot go back and is forced to live according to the laws of the old world.\n https://doramy.club/639-dorama-lunnye-vlyublyonnye-alye-serdca-koryo.html"},
            "Thriller":   []string{"Mouse🐁🗡️\nA gripping story set in a world where people can identify psychopaths through DNA testing of fetuses in the womb\n https://doramy.club/24814-7-mysh.html"},
            "Action":     []string{"K2⚖️💪🏼\nHaving survived betrayal, a disgraced mercenary soldier codenamed K2 finds himself at the center of the intrigues of dishonest politicians. He is trying to avenge his collapse and save the unfortunate daughter of a presidential candidate.\n https://doramy.club/498-k2-teloxranitel.html"},
            "Comedy":     []string{"Mr. Queen👸🏻🥴\nJung Bong Hwan works as a chef at the Blue House, but suddenly he finds himself in the Joseon era, and now his soul is trapped in the body of Queen Chorin. King Cheol Jong looks like a puppet in the hands of those who really rule the country, but in fact he hides quite strong traits of his character.\n https://doramy.club/23665-62-koroleva-chorin.html"},
        },
        "Anime": {
            "Romance":   []string{"Anime Romance Recommendation(soon)..."},
            "Melodrama": []string{"Anime Melodrama Recommendation(soon)..."},
            "Historical": []string{"Anime Historical Recommendation(soon)..."},
            "Thriller":   []string{"Anime Thriller Recommendation(soon)..."},
            "Action":     []string{"Anime Action Recommendation(soon)..."},
            "Comedy":     []string{"Anime Comedy Recommendation(soon)..."},
        },
    }
)

type UserState struct {
    Name       string
    Preference string
    Step       string
    More       bool
}

func main() {
    botToken := "6673045880:AAFyEN4jHRjqETu1yIY3HimwXfFIORA5xjE" // Replace with your actual bot token
    bot, err := tgbotapi.NewBotAPI(botToken)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 90

    updates := bot.GetUpdatesChan(u)

    userStates := make(map[int64]*UserState)

    for update := range updates {
        if update.Message != nil {
            chatID := update.Message.Chat.ID
            text := update.Message.Text

            if _, ok := userStates[chatID]; !ok {
                userStates[chatID] = &UserState{Step: "welcome"}
                msg := tgbotapi.NewMessage(chatID, messages["welcome"])
                bot.Send(msg)
                continue
            }

            userState := userStates[chatID]

            switch userState.Step {
            case "welcome":
                userState.Name = text
                userState.Step = "choice"
                msg := tgbotapi.NewMessage(chatID, "Welcome to the Kdrama Recommendation Bot, "+userState.Name+"!\nWhat's your preference: Kdramas or Anime?\nPlease type 'Kdrama' or 'Anime'.")
                bot.Send(msg)

            case "choice":
                if text == "Kdrama" || text == "Anime" {
                    userState.Preference = text
                    userState.Step = "genre"
                    msg := tgbotapi.NewMessage(chatID, messages["genre"])
                    bot.Send(msg)
                } else {
                    msg := tgbotapi.NewMessage(chatID, "Please type 'Kdrama' or 'Anime'.")
                    bot.Send(msg)
                }

            case "genre":
                genreRecommendations, exists := recommendations[userState.Preference][text]
                if exists {
                    for _, rec := range genreRecommendations {
                        msg := tgbotapi.NewMessage(chatID, rec)
                        bot.Send(msg)
                    }
                    userState.Step = "more"
                    msg := tgbotapi.NewMessage(chatID, messages["more"])
                    bot.Send(msg)
                } else {
                    msg := tgbotapi.NewMessage(chatID, "Sorry, I couldn't find recommendations for this genre.")
                    bot.Send(msg)
                }

            case "more":
                if text == "yes" {
                    userState.Step = "genre"
                    msg := tgbotapi.NewMessage(chatID, messages["genre"])
                    bot.Send(msg)
                } else if text == "no" {
                    msg := tgbotapi.NewMessage(chatID, messages["thanks"])
                    bot.Send(msg)
                    delete(userStates, chatID) // Reset user state after thanking
                } else {
                    msg := tgbotapi.NewMessage(chatID, "Please type 'yes' or 'no'.")
                    bot.Send(msg)
                }
            }
        }
    }
}

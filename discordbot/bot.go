package discordBot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var Discord *discordgo.Session

func Connect() {
	var err error
	Discord, err = discordgo.New("Bot " + "OTYyNjg0OTE1NjMwMTA0NjI2.YlLIMQ.Mr8A8xO02cyjNyGupTsrn7H_6vE")
	if err != nil {
		log.Panic("err: ", err)
	}

	Discord.AddHandler(func(discord *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", discord.State.User.Username, discord.State.User.Discriminator)
	})
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ricksanchez",
			Description: "warum nicht",
		},
	}
)

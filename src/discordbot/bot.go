package discordBot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func discordBot() {
	discord, err := discordgo.New("Bot " + "OTYyNjg0OTE1NjMwMTA0NjI2.YlLIMQ.Mr8A8xO02cyjNyGupTsrn7H_6vE")

	if err != nil {
		log.Panic("err: ", err)
	}

	println(discord)
}

package discordBot

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/Oddernumse/akiya/db"
	"github.com/bwmarrin/discordgo"
	http "github.com/useflyent/fhttp"
)

type Payment struct {
	Title      string `json:"title"`
	Product_id string `json:"product_id"`
	Gateway    string `json:"gateway"`
	Quantity   int    `json:"quantity"`
	Email      string `json:"email"`
	WhiteLabel bool   `json:"white_label"`
	ReturnUrl  string `json:"return_url"`
}

type InvoiceObject struct {
	Status int `json:"status"`
	Data   struct {
		Url string `json:"url"`
	} `json:"data"`
}

var (
	GuildID        = flag.String("guild", "962702002444443738", "weed")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var Discord *discordgo.Session

func Connect() *discordgo.Session {
	var err error
	Discord, err = discordgo.New("Bot " + "OTYyNjg0OTE1NjMwMTA0NjI2.YlLIMQ.Mr8A8xO02cyjNyGupTsrn7H_6vE")
	if err != nil {
		log.Panic("err: ", err)
	}

	Discord.AddHandler(func(discord *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", discord.State.User.Username, discord.State.User.Discriminator)
	})

	err = Discord.Open()
	if err != nil {
		log.Panic("error opening connection,", err)
	}

	return Discord
}

func CreateCommands(dg *discordgo.Session) {
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "account",
			Description: "Add token you wish gifts to be redeemed on",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "token",
					Description: "Your discord user token",
					Required:    true,
				},
			},
		},
		{
			Name:        "buy",
			Description: "Buy snipes with your chosen payment method",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "gateway",
					Description: "Choose what currency you wish to pay with",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Bitcoin",
							Value: "bitcoin",
						},
						{
							Name:  "Litecoin",
							Value: "litecoin",
						},
						{
							Name:  "Monero",
							Value: "monero",
						},
						{
							Name:  "Ethereum",
							Value: "ethereum",
						},
						{
							Name:  "Paypal",
							Value: "paypal",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "quantity",
					Description: "Decide how many snipes you want",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "email",
					Description: "The mail to which you want the order sent",
					Required:    true,
				},
			},
		},
		{
			Name:        "overview",
			Description: "Quick overview for your order",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"account": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: uint64(discordgo.MessageFlagsEphemeral),
				},
			})
			options := i.ApplicationCommandData().Options

			if db.GetUser(i.Interaction.User.ID).Userid == i.Interaction.User.ID {
				db.UpdateCustomer(i.Interaction.User.ID, options[0].StringValue())
			} else {
				db.CreateCustomer(i.Interaction.Member.User.ID, options[0].StringValue())
			}

			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{Content: "Added Account"})
		},
		"buy": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: uint64(discordgo.MessageFlagsEphemeral),
				},
			})

			//dbClient := db.Table{}
			s.GuildMemberEdit("962702002444443738", i.Interaction.Member.User.ID, []string{"965234270266339349"})
			options := i.ApplicationCommandData().Options

			client := &http.Client{}
			reqe, _ := json.Marshal(Payment{Title: "Snipe", Product_id: "625c70089c528", Gateway: options[0].StringValue(), Quantity: int(options[1].IntValue()), Email: options[2].StringValue(), ReturnUrl: "https://demo.sellix.io/return"})
			req, _ := http.NewRequest("POST", "https://dev.sellix.io/v1/payments", bytes.NewBuffer(reqe))
			req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
			req.Header.Set("authorization", "Bearer FxvWdiW1o0CVVe71WDueKowPuyPZ4CFJqhMfvCQy7809ud6IgUcJsSYt8SQHEzyB")
			resp, err := client.Do(req)

			if err != nil {
				log.Fatal(err)
			}

			BuyUrl := InvoiceObject{}
			processed, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal([]byte(processed), &BuyUrl)

			if options[0].StringValue() == "bitcoin" {
				s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{Content: "Bitcoin orders need to be more than 3.5$"})

			} else {
				s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{Content: "Payment Created! " + BuyUrl.Data.Url})
			}

		},
		"overview": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		},
	}
)

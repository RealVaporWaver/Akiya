package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	discordBot "github.com/Oddernumse/akiya/discordbot"
	"github.com/gorilla/websocket"
	http "github.com/useflyent/fhttp"
)

// Structs
type HeartbeatInterval struct {
	HeartbeatInterval int32 `json:"heartbeat_interval"`
}

type WSPayload struct {
	Op   int32       `json:"op"`
	T    string      `json:"t"`
	Data interface{} `json:"d"`
}

type MessagePayload struct {
	Content string `json:"content"`
}

type IdentifyPayload struct {
	Token           string                 `json:"token"`
	SuperProperties map[string]interface{} `json:"properties"`
}

type Client struct {
	Lock sync.Mutex
	Ws   *websocket.Conn
}

// Class Methods
func (c *Client) heartbeat(ms int32, ws *websocket.Conn) {
	payload := WSPayload{Op: 1, Data: nil}

	for {
		ws.WriteJSON(payload)
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}

func convert(in interface{}, out interface{}) {
	data, err := json.Marshal(in)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &out)

	if err != nil {
		panic(err)
	}
}

func (c *Client) dial() {
	ws, _, err := websocket.DefaultDialer.Dial("wss://gateway.discord.gg?v=9&encoding=json", nil)
	if err != nil {
		log.Panic("dial:", err)
	}

	c.Ws = ws
}

func (c *Client) readMessage() WSPayload {
	_, message, err := c.Ws.ReadMessage()
	if err != nil {
		log.Println("read:", err)
	}

	payload := WSPayload{}
	json.Unmarshal(message, &payload)

	return payload
}

func (c *Client) write(op int32, data interface{}) {
	c.Lock.Lock()
	c.Ws.WriteJSON(&WSPayload{
		Op:   op,
		Data: data,
	})
	c.Lock.Unlock()
}

//------------------------------------
// MAIN FUNC
//------------------------------------

func main() {
	discordBot.Connect()

	f, _ := os.Open("bot_tokens.txt")

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		go func() {
			client := Client{}
			client.dial()

			for {
				payload := client.readMessage()

				switch payload.Op {
				case 10:
					var heart HeartbeatInterval
					convert(payload.Data, &heart)

					client.write(2, IdentifyPayload{Token: line, SuperProperties: map[string]interface{}{
						"$os":      "windows",
						"$browser": "brave",
						"$device":  "brave",
					}})

					go client.heartbeat(heart.HeartbeatInterval, client.Ws)
					break
				case 0:
					switch payload.T {
					case "MESSAGE_CREATE":
						var msg MessagePayload
						convert(payload.Data, &msg)

						// This if statements entire job is to redeem codes once its found a possible code
						if strings.Contains(msg.Content, "discord.gift/") {
							code := strings.Split(msg.Content, "discord.gift/")
							if len(code[1]) >= 16 {
								httpClient := &http.Client{}

								req, _ := http.NewRequest("POST", "https://discord.com/api/v9/entitlements/gift-codes/"+code[1][0:16]+"/redeem", nil)
								req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
								req.Header.Set("authorization", "Mjg4Njk1MzQ0ODYxMDIwMTYw.YlGRdQ.-nHHh8rnAed6Wo_LgjYAp9gJNKw")

								resp, err := httpClient.Do(req)

								// Error handling for if the redeem goes bad
								if err != nil {
									log.Println("err: ", err)
								}
								if resp.StatusCode != 200 {
									log.Println("Failed Redeem")
								} else {
									//insert queing system here
									log.Println("Successful redeem")
								}
							}
						}
						break
					case "READY":
						println("READY: ", line)
						break
					}
				}
			}
		}()
	}

	select {}
}

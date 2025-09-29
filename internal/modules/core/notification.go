package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	"github.com/yuin/gopher-lua"
)

// NotificationModule provides notification services
type NotificationModule struct {
	info   CoreModuleInfo
	client *http.Client
}

// NewNotificationModule creates a new notification module
func NewNotificationModule() *NotificationModule {
	info := CoreModuleInfo{
		Name:         "notify",
		Version:      "1.0.0",
		Description:  "Notification services including Slack, Discord, Email, Webhook",
		Author:       "Sloth Runner Team",
		Category:     "core",
		Dependencies: []string{},
	}

	return &NotificationModule{
		info: info,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Info returns module information
func (n *NotificationModule) Info() CoreModuleInfo {
	return n.info
}

// Loader loads the notification module into Lua
func (n *NotificationModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"slack":    n.luaSlack,
		"discord":  n.luaDiscord,
		"email":    n.luaEmail,
		"webhook":  n.luaWebhook,
		"teams":    n.luaTeams,
		"telegram": n.luaTelegram,
	})

	L.Push(mod)
	return 1
}

// SlackMessage represents a Slack message
type SlackMessage struct {
	Text        string       `json:"text,omitempty"`
	Username    string       `json:"username,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment represents a Slack attachment
type Attachment struct {
	Color     string  `json:"color,omitempty"`
	Title     string  `json:"title,omitempty"`
	Text      string  `json:"text,omitempty"`
	Footer    string  `json:"footer,omitempty"`
	Timestamp int64   `json:"ts,omitempty"`
	Fields    []Field `json:"fields,omitempty"`
}

// Field represents a Slack field
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// luaSlack sends a Slack notification
func (n *NotificationModule) luaSlack(L *lua.LState) int {
	webhookURL := L.CheckString(1)
	messageTable := L.CheckTable(2)

	message := SlackMessage{}
	
	// Parse message from Lua table
	messageTable.ForEach(func(key, value lua.LValue) {
		switch key.String() {
		case "text":
			message.Text = value.String()
		case "username":
			message.Username = value.String()
		case "channel":
			message.Channel = value.String()
		case "icon_emoji":
			message.IconEmoji = value.String()
		case "attachments":
			if attachTable, ok := value.(*lua.LTable); ok {
				n.parseAttachments(attachTable, &message)
			}
		}
	})

	jsonData, err := json.Marshal(message)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	resp, err := n.client.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	L.Push(lua.LBool(success))
	if !success {
		L.Push(lua.LString(fmt.Sprintf("HTTP %d", resp.StatusCode)))
		return 2
	}
	return 1
}

// parseAttachments parses Slack attachments from Lua table
func (n *NotificationModule) parseAttachments(table *lua.LTable, message *SlackMessage) {
	table.ForEach(func(_, value lua.LValue) {
		if attachTable, ok := value.(*lua.LTable); ok {
			attachment := Attachment{
				Timestamp: time.Now().Unix(),
			}
			
			attachTable.ForEach(func(key, val lua.LValue) {
				switch key.String() {
				case "color":
					attachment.Color = val.String()
				case "title":
					attachment.Title = val.String()
				case "text":
					attachment.Text = val.String()
				case "footer":
					attachment.Footer = val.String()
				case "fields":
					if fieldsTable, ok := val.(*lua.LTable); ok {
						n.parseFields(fieldsTable, &attachment)
					}
				}
			})
			
			message.Attachments = append(message.Attachments, attachment)
		}
	})
}

// parseFields parses Slack fields from Lua table
func (n *NotificationModule) parseFields(table *lua.LTable, attachment *Attachment) {
	table.ForEach(func(_, value lua.LValue) {
		if fieldTable, ok := value.(*lua.LTable); ok {
			field := Field{}
			
			fieldTable.ForEach(func(key, val lua.LValue) {
				switch key.String() {
				case "title":
					field.Title = val.String()
				case "value":
					field.Value = val.String()
				case "short":
					field.Short = lua.LVAsBool(val)
				}
			})
			
			attachment.Fields = append(attachment.Fields, field)
		}
	})
}

// DiscordMessage represents a Discord message
type DiscordMessage struct {
	Content   string          `json:"content,omitempty"`
	Username  string          `json:"username,omitempty"`
	AvatarURL string          `json:"avatar_url,omitempty"`
	Embeds    []DiscordEmbed  `json:"embeds,omitempty"`
}

// DiscordEmbed represents a Discord embed
type DiscordEmbed struct {
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	Color       int                 `json:"color,omitempty"`
	Footer      *DiscordEmbedFooter `json:"footer,omitempty"`
	Timestamp   string              `json:"timestamp,omitempty"`
	Fields      []DiscordField      `json:"fields,omitempty"`
}

// DiscordEmbedFooter represents a Discord embed footer
type DiscordEmbedFooter struct {
	Text string `json:"text"`
}

// DiscordField represents a Discord field
type DiscordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// luaDiscord sends a Discord notification
func (n *NotificationModule) luaDiscord(L *lua.LState) int {
	webhookURL := L.CheckString(1)
	messageTable := L.CheckTable(2)

	message := DiscordMessage{}
	
	// Parse message from Lua table
	messageTable.ForEach(func(key, value lua.LValue) {
		switch key.String() {
		case "content":
			message.Content = value.String()
		case "username":
			message.Username = value.String()
		case "avatar_url":
			message.AvatarURL = value.String()
		case "embeds":
			if embedTable, ok := value.(*lua.LTable); ok {
				n.parseDiscordEmbeds(embedTable, &message)
			}
		}
	})

	jsonData, err := json.Marshal(message)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	resp, err := n.client.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	L.Push(lua.LBool(success))
	if !success {
		L.Push(lua.LString(fmt.Sprintf("HTTP %d", resp.StatusCode)))
		return 2
	}
	return 1
}

// parseDiscordEmbeds parses Discord embeds from Lua table
func (n *NotificationModule) parseDiscordEmbeds(table *lua.LTable, message *DiscordMessage) {
	table.ForEach(func(_, value lua.LValue) {
		if embedTable, ok := value.(*lua.LTable); ok {
			embed := DiscordEmbed{
				Timestamp: time.Now().Format(time.RFC3339),
			}
			
			embedTable.ForEach(func(key, val lua.LValue) {
				switch key.String() {
				case "title":
					embed.Title = val.String()
				case "description":
					embed.Description = val.String()
				case "color":
					if num, ok := val.(lua.LNumber); ok {
						embed.Color = int(num)
					}
				case "footer":
					embed.Footer = &DiscordEmbedFooter{Text: val.String()}
				case "fields":
					if fieldsTable, ok := val.(*lua.LTable); ok {
						n.parseDiscordFields(fieldsTable, &embed)
					}
				}
			})
			
			message.Embeds = append(message.Embeds, embed)
		}
	})
}

// parseDiscordFields parses Discord fields from Lua table
func (n *NotificationModule) parseDiscordFields(table *lua.LTable, embed *DiscordEmbed) {
	table.ForEach(func(_, value lua.LValue) {
		if fieldTable, ok := value.(*lua.LTable); ok {
			field := DiscordField{}
			
			fieldTable.ForEach(func(key, val lua.LValue) {
				switch key.String() {
				case "name":
					field.Name = val.String()
				case "value":
					field.Value = val.String()
				case "inline":
					field.Inline = lua.LVAsBool(val)
				}
			})
			
			embed.Fields = append(embed.Fields, field)
		}
	})
}

// luaEmail sends an email notification
func (n *NotificationModule) luaEmail(L *lua.LState) int {
	configTable := L.CheckTable(1)
	
	var smtpHost, smtpPort, username, password, from, to, subject, body string
	
	configTable.ForEach(func(key, value lua.LValue) {
		switch key.String() {
		case "smtp_host":
			smtpHost = value.String()
		case "smtp_port":
			smtpPort = value.String()
		case "username":
			username = value.String()
		case "password":
			password = value.String()
		case "from":
			from = value.String()
		case "to":
			to = value.String()
		case "subject":
			subject = value.String()
		case "body":
			body = value.String()
		}
	})

	// Create message
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)
	
	// Setup authentication
	auth := smtp.PlainAuth("", username, password, smtpHost)
	
	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	
	L.Push(lua.LBool(err == nil))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// luaWebhook sends a generic webhook notification
func (n *NotificationModule) luaWebhook(L *lua.LState) int {
	url := L.CheckString(1)
	payloadTable := L.CheckTable(2)
	
	// Convert Lua table to map
	payload := make(map[string]interface{})
	payloadTable.ForEach(func(key, value lua.LValue) {
		payload[key.String()] = n.luaValueToGo(value)
	})
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	resp, err := n.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer resp.Body.Close()
	
	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	L.Push(lua.LBool(success))
	if !success {
		L.Push(lua.LString(fmt.Sprintf("HTTP %d", resp.StatusCode)))
		return 2
	}
	return 1
}

// luaTeams sends a Microsoft Teams notification
func (n *NotificationModule) luaTeams(L *lua.LState) int {
	webhookURL := L.CheckString(1)
	messageTable := L.CheckTable(2)
	
	message := make(map[string]interface{})
	messageTable.ForEach(func(key, value lua.LValue) {
		message[key.String()] = n.luaValueToGo(value)
	})
	
	jsonData, err := json.Marshal(message)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	resp, err := n.client.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer resp.Body.Close()
	
	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	L.Push(lua.LBool(success))
	if !success {
		L.Push(lua.LString(fmt.Sprintf("HTTP %d", resp.StatusCode)))
		return 2
	}
	return 1
}

// luaTelegram sends a Telegram notification
func (n *NotificationModule) luaTelegram(L *lua.LState) int {
	botToken := L.CheckString(1)
	chatID := L.CheckString(2)
	message := L.CheckString(3)
	
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    message,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	resp, err := n.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer resp.Body.Close()
	
	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	L.Push(lua.LBool(success))
	if !success {
		L.Push(lua.LString(fmt.Sprintf("HTTP %d", resp.StatusCode)))
		return 2
	}
	return 1
}

// luaValueToGo converts Lua value to Go interface{}
func (n *NotificationModule) luaValueToGo(value lua.LValue) interface{} {
	switch v := value.(type) {
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case lua.LBool:
		return bool(v)
	case *lua.LTable:
		result := make(map[string]interface{})
		v.ForEach(func(key, val lua.LValue) {
			result[key.String()] = n.luaValueToGo(val)
		})
		return result
	default:
		return value.String()
	}
}
package bot

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/porty/chrome-screens/chrome"

	"github.com/nlopes/slack"
)

type Bot struct {
	name     string
	id       string
	rtm      *slack.RTM
	chrome   *chrome.Chrome
	sleeping bool
	onNewURL func(newURL string)
}

func New(rtm *slack.RTM, chrome *chrome.Chrome, onNewURL func(newURL string)) *Bot {
	if onNewURL == nil {
		onNewURL = func(_ string) {}
	}
	return &Bot{
		rtm:      rtm,
		chrome:   chrome,
		onNewURL: onNewURL,
	}
}

func (b *Bot) Run() {
	for msg := range b.rtm.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.MessageEvent:
			b.handleMessage(event)
		case *slack.ConnectedEvent:
			b.name = event.Info.User.Name
			b.id = event.Info.User.ID
			log.Printf("Connected to slack, my name is %s", b.name)
		case *slack.ConnectingEvent:
			// ignored
		case *slack.HelloEvent:
			// ignored
		case *slack.UnmarshallingErrorEvent:
			// log.Print("Unmarshalling error: " + event.Error())
		case *slack.AckErrorEvent:
			log.Print("ACK error: " + event.Error())
		case *slack.PresenceChangeEvent:
			// ignored
		case *slack.ReconnectUrlEvent:
			// ignored
		case *slack.LatencyReport:
			// ignored
		case *slack.UserTypingEvent:
			// ignored
		case *slack.PrefChangeEvent, *slack.ReactionAddedEvent:
			// ignored
		default:
			log.Printf("Received event of type " + msg.Type)
		}
	}
}

func (b *Bot) handleMessage(message *slack.MessageEvent) {
	if !b.isForMe(message) {
		return
	}

	log.Printf("Received message %q from %s", message.Text, message.User)

	// TODO support direct messaging without having to @tv1 each time
	parts := strings.Split(message.Text, " ")
	if len(parts) == 3 && parts[1] == "set" {
		b.set(parts[2], message)
		return
	}
	if len(parts) == 2 && parts[1] == "help" {
		b.help(message)
		return
	}
	if len(parts) == 2 && parts[1] == "wake" {
		b.wake(message)
		return
	}
	if len(parts) == 2 && parts[1] == "sleep" {
		b.sleep(message)
		return
	}
	if len(parts) == 2 && b.getSlackURL(parts[1]) != "" {
		b.set(parts[1], message)
		return
	}

	log.Printf("I don't know what to do with %q", message.Text)
	b.addReaction("badpokerface", message)
}

func (b *Bot) isForMe(message *slack.MessageEvent) bool {
	if message.User == b.id {
		return false
	}
	//log.Printf("message.User = %q, message.Username = %q", message.User, message.Username)
	upcase := strings.ToUpper(strings.TrimSpace(message.Text))
	identifier := "<@" + strings.ToUpper(b.id) + ">"

	return strings.HasPrefix(upcase, strings.ToUpper(b.name)) ||
		strings.Contains(upcase, identifier) ||
		strings.HasPrefix(message.Channel, "D")
}

func (b *Bot) addReaction(emoji string, message *slack.MessageEvent) {
	b.rtm.AddReaction(emoji, slack.ItemRef{
		Channel:   message.Channel,
		Timestamp: message.Timestamp,
	})
}

func (b *Bot) getSlackURL(url string) string {
	if !strings.HasPrefix(url, "<") && !strings.HasSuffix(url, ">") {
		return ""
	}
	url = strings.TrimPrefix(url, "<")
	url = strings.TrimSuffix(url, ">")

	if pipe := strings.Index(url, "|"); pipe != -1 {
		return url[:pipe]
	}
	return url
}

func (b *Bot) help(message *slack.MessageEvent) {
	text := "Say _<@" + b.id + "> www.disney.com_ to set to disney.com\n"
	text += "Say _<@" + b.id + "> wake_ to wake me up\n"
	text += "Say _<@" + b.id + "> sleep_ to put me to sleep"
	b.rtm.SendMessage(&slack.OutgoingMessage{
		Channel: message.Channel,
		Text:    text,
		Type:    "message",
	})
}

func (b *Bot) sleep(message *slack.MessageEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "xset", "dpms", "force", "off")
	env := os.Environ()
	env = append(env, "DISPLAY=:0")
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		log.Print("Failed to run sleep command: " + err.Error())
		b.addReaction("boom", message)
		return
	}
	b.addReaction("ok_hand", message)
}

func (b *Bot) wake(message *slack.MessageEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "xset", "dpms", "force", "on")
	env := os.Environ()
	env = append(env, "DISPLAY=:0")
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		log.Print("Failed to run sleep command: " + err.Error())
		b.addReaction("boom", message)
		return
	}
	if message != nil {
		b.addReaction("ok_hand", message)
	}
}

func (b *Bot) set(dirtyURL string, message *slack.MessageEvent) {
	url := b.getSlackURL(dirtyURL)
	b.chrome.SetURL(url)
	b.onNewURL(url)

	if message != nil {
		b.addReaction("ok_hand", message)
	}
}

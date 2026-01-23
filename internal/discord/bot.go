package discord

import (
	"fmt"
	"go-notify-hub/internal/database"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	validCommand = regexp.MustCompile(`^/\w+`)
	validTicker  = regexp.MustCompile(`^\$?[A-Z]{1,5}$`)
)

var slashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "gex",
		Description: "Get GEX data for a ticker",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "ticker",
				Description: "Stock ticker symbol (e.g. AAPL, TSLA)",
				Required:    true,
			},
		},
	},
	{
		Name:        "help",
		Description: "Show available commands",
	},
}

type Bot struct {
	session  *discordgo.Session
	commands []*discordgo.ApplicationCommand
}

type Message struct {
	UserID    string
	ChannelID string
	GuildID   string
	Cmd       string
	Params    []string
	Content   string
}

func NewBot(token string) (*Bot, error) {
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("create discord session: %w", err)
	}

	bot := &Bot{session: session}

	session.AddHandler(bot.handleMessageCreate)
	session.AddHandler(bot.handleInteraction) // 新增
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentMessageContent

	if err := session.Open(); err != nil {
		return nil, fmt.Errorf("open websocket connection: %w", err)
	}

	// 新增：註冊 Slash Commands
	if err := bot.registerCommands(); err != nil {
		session.Close()
		return nil, fmt.Errorf("register commands: %w", err)
	}

	slog.Info("bot is running", slog.String("user", session.State.User.Username))

	return bot, nil
}

// 新增
func (b *Bot) registerCommands() error {
	appID := b.session.State.User.ID

	for _, cmd := range slashCommands {
		registered, err := b.session.ApplicationCommandCreate(appID, "", cmd)
		if err != nil {
			return fmt.Errorf("register command %s: %w", cmd.Name, err)
		}
		b.commands = append(b.commands, registered)
		slog.Info("registered command", slog.String("name", cmd.Name))
	}

	return nil
}

func (b *Bot) Close() error {
	slog.Info("shutting down")
	if b.session == nil {
		return nil
	}
	return b.session.Close()
}

func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()

	var userID, username string
	if i.Member != nil {
		userID = i.Member.User.ID
		username = i.Member.User.Username
	} else if i.User != nil {
		userID = i.User.ID
		username = i.User.Username
	}

	slog.Info("slash command received",
		slog.String("command", data.Name),
		slog.String("user_id", userID),
		slog.String("user", username),
	)

	switch data.Name {
	case "gex":
		b.slashGex(s, i)
	case "help":
		b.slashHelp(s, i)
	}
}

func (b *Bot) slashGex(s *discordgo.Session, i *discordgo.InteractionCreate) {
	const fn = "DiscordHandler/slashGex"

	options := i.ApplicationCommandData().Options
	ticker := strings.ToUpper(options[0].StringValue())

	if !validTicker.MatchString(ticker) {
		b.respond(s, i, "Usage: /gex <ticker>")
		return
	}

	ticker = strings.TrimPrefix(ticker, "$")

	result, err := database.SelectTicker(ticker)
	if err != nil {
		slog.Error(fn+": failed to get ticker data", "ticker", ticker, "error", err)
		b.respond(s, i, "Failed to fetch ticker data")
		return
	}

	text := fmt.Sprintf(`親愛的投資人您好，今天是 %s，
盤前觀察 $%s
%s

%s

%s`,
		time.Now().Format("2006-01-02"),
		ticker,
		result.RationaleCompact,
		result.ExpiryNearTerm,
		result.ExpiryFarTerm,
	)

	b.respondWithImageURL(s, i, text, result.ImageURL)
}

func (b *Bot) slashHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	help := `**Available Commands**
/gex <ticker> - Get GEX data for a ticker
/help - Show this message`

	b.respond(s, i, help)
}

func (b *Bot) respond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		slog.Error("respond: failed to send", slog.String("error", err.Error()))
	}
}

func (b *Bot) respondWithImageURL(s *discordgo.Session, i *discordgo.InteractionCreate, content, imageURL string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Embeds: []*discordgo.MessageEmbed{
				{
					Image: &discordgo.MessageEmbedImage{
						URL: imageURL,
					},
				},
			},
		},
	})
	if err != nil {
		slog.Error("respondWithImageURL: failed to send", slog.String("error", err.Error()))
	}
}

func (b *Bot) handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	slog.Info("message received",
		slog.String("guild_id", m.GuildID),
		slog.String("channel_id", m.ChannelID),
		slog.String("author_id", m.Author.ID),
		slog.String("author", m.Author.Username),
		slog.String("content", m.Content),
	)

	msg, err := parseMessage(m)
	if err != nil {
		return
	}

	if err := verifyMessage(msg); err != nil {
		b.replyText(s, m, err.Error())
		return
	}

	switch msg.Cmd {
	case "/gex":
		b.commandGex(s, m, msg)
	case "/help":
		b.commandHelp(s, m)
	default:
		b.replyText(s, m, fmt.Sprintf("Unknown command: %s", msg.Cmd))
	}
}

func parseMessage(m *discordgo.MessageCreate) (*Message, error) {
	content := strings.TrimSpace(m.Content)

	if !validCommand.MatchString(content) {
		return nil, fmt.Errorf("not a command")
	}

	fields := strings.Fields(content)

	return &Message{
		UserID:    m.Author.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
		Cmd:       strings.ToLower(fields[0]),
		Params:    fields[1:],
		Content:   content,
	}, nil
}

func verifyMessage(msg *Message) error {
	switch msg.Cmd {
	case "/gex":
		if len(msg.Params) < 1 {
			return fmt.Errorf("Usage: /gex $<ticker>")
		}
	}
	return nil
}

func (b *Bot) commandGex(s *discordgo.Session, m *discordgo.MessageCreate, msg *Message) {
	const fn = "DiscordHandler/commandGex"

	ticker := strings.ToUpper(msg.Params[0])
	if !validTicker.MatchString(ticker) {
		b.replyText(s, m, "Usage: /gex $<ticker>")
		return
	}

	ticker = strings.TrimPrefix(ticker, "$")

	result, err := database.SelectTicker(ticker)
	if err != nil {
		slog.Error(fn+": failed to get ticker data", "ticker", ticker, "error", err)
		return
	}

	text := fmt.Sprintf(`親愛的投資人您好，今天是 %s，
盤前觀察 $%s
%s

%s

%s`,
		time.Now().Format("2006-01-02"),
		ticker,
		result.RationaleCompact,
		result.ExpiryNearTerm,
		result.ExpiryFarTerm,
	)

	b.replyWithImageURL(s, m, text, result.ImageURL)
}

func (b *Bot) commandHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	help := `**Available Commands**
/gex $<ticker> - Get GEX data for a ticker
/help - Show this message`

	b.replyText(s, m, help)
}

func (b *Bot) replyText(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	_, err := s.ChannelMessageSendReply(m.ChannelID, content, m.Reference())
	if err != nil {
		slog.Error("replyText: failed to send",
			slog.String("channel_id", m.ChannelID),
			slog.String("error", err.Error()),
		)
	}
}

func (b *Bot) replyWithImage(s *discordgo.Session, m *discordgo.MessageCreate, content, imagePath string) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("open image file: %w", err)
	}
	defer file.Close()

	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content:   content,
		Reference: m.Reference(),
		Files: []*discordgo.File{
			{
				Name:   filepath.Base(imagePath),
				Reader: file,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (b *Bot) replyWithImageURL(s *discordgo.Session, m *discordgo.MessageCreate, content, imageURL string) error {
	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content:   content,
		Reference: m.Reference(),
		Embeds: []*discordgo.MessageEmbed{
			{
				Image: &discordgo.MessageEmbedImage{
					URL: imageURL,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (b *Bot) replyWithImageReader(s *discordgo.Session, m *discordgo.MessageCreate, content, filename string, reader io.Reader) error {
	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content:   content,
		Reference: m.Reference(),
		Files: []*discordgo.File{
			{
				Name:   filename,
				Reader: reader,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

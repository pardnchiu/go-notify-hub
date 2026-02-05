package handler

type Message struct {
	UserID    string
	ChannelID string
	GuildID   string
	Cmd       string
	Params    []string
	Content   string
}

type Reply struct {
	Content    string
	ImageURL   string
	PreviewURL string
}

func Handler(msg *Message) []Reply {
	var replies []Reply

	if msg.Cmd != "" {
		switch msg.Cmd {
		case "help", "/help":
			return []Reply{{Content: "How to use"}}
		default:
			return []Reply{{Content: "No such command: " + msg.Cmd}}
		}
	} else {
		replies = append(replies, Reply{Content: "Receive message：" + msg.Content})
	}

	return replies
}

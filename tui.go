package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/celestiaorg/celestia-node/api/rpc/client"
	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-node/header"
	"github.com/celestiaorg/celestia-node/share"
	"github.com/celestiaorg/celestia-node/state"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding = 2
)

var (
	pad         = strings.Repeat(" ", padding)
	docStyle    = lipgloss.NewStyle().Margin(1, 1, 1, 1)
	senderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	tiCmd       tea.Cmd
	vpCmd       tea.Cmd
)

type tickMsg time.Time

type model struct {
	ctx           context.Context
	client        *client.Client
	namespace     share.Namespace
	addr          state.Address
	roomID        string
	encryptionKey string
	username      string
	public        bool
	headerSub     <-chan *header.ExtendedHeader

	viewport    viewport.Model
	messages    []Message
	textarea    textarea.Model
	senderStyle lipgloss.Style

	width, height int
}

func NewModel(ctx context.Context, celestiaClient *client.Client, key string, public bool, username string) (*model, error) {
	h, v := docStyle.GetFrameSize()

	addr, err := celestiaClient.State.AccountAddress(ctx)
	if err != nil {
		return nil, err
	}

	namespace, err := share.NewBlobNamespaceV0([]byte(chatNamespaceStr))
	if err != nil {
		return nil, err
	}

	headers, err := celestiaClient.Header.Subscribe(ctx)
	if err != nil {
		return nil, err
	}

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(300)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(300, 50)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	m := &model{
		ctx:         ctx,
		client:      celestiaClient,
		namespace:   namespace,
		addr:        addr,
		headerSub:   headers,
		public:      public,
		username:    username,
		textarea:    ta,
		messages:    []Message{},
		viewport:    vp,
		height:      v,
		width:       h,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
	}
	if m.public {
		m.roomID = key
	} else {
		m.encryptionKey = key
	}
	return m, nil
}

const historyDisplayRange = 7200 // ~day

// TODO: Make sure we call this before we start listening(or displaying) new header
func (m *model) DisplayHistory() error {
	head, err := m.client.Header.NetworkHead(m.ctx)
	if err != nil {
		return err
	}

	start, end := head.Height()-syncPeriod, head.Height()
	msgCh := GetMessagesBackwardsAsync(m.ctx, m.client, m.namespace, start, end)
	for msg := range msgCh {
		// TODO: Messages has to be synchronized as here we write and View reads async
		if msg.ID == m.roomID {
			// TODO: needs to be added backwards, not to the end of slice
			m.messages = append(m.messages, msg)
		}
	}

	return nil
}

func (m *model) handleIncomingHeader(h *header.ExtendedHeader) tea.Cmd {
	height := h.Height()
	blobs, err := m.client.Blob.GetAll(m.ctx, height, []share.Namespace{m.namespace})
	if err != nil {
		blobs = make([]*blob.Blob, 0)
	}

	for _, b := range blobs {
		var newMessage Message
		err = json.Unmarshal(b.Data, &newMessage)
		if err != nil {
			// TODO LOG error
		}

		if newMessage.ID == m.roomID {
			m.messages = append(m.messages, newMessage)
		}
	}

	m.viewport.SetContent(displayMessages(m.messages))
	m.viewport.GotoBottom()
	return waitForActivity(m.headerSub)
}

func displayMessages(messages []Message) string {
	str := ""
	for _, m := range messages {
		ansiColor := strconv.Itoa(int(m.Username[len(m.Username)-1]) % 15)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(ansiColor))
		str += style.Render(m.Username) + ": " + m.Content + "\n"
	}
	return str
}

func waitForActivity(sub <-chan *header.ExtendedHeader) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func (m *model) sendMessage(ctx context.Context, content string) error {
	text := content
	var err error
	if !m.public {
		text, err = encrypt(content, m.encryptionKey)
		if err != nil {
			return err
		}
	}

	msg := &Message{
		Username: m.username,
		ID:       m.roomID,
		Public:   m.public,
		Content:  text,
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	msgBlob, err := blob.NewBlobV0(m.namespace, msgBytes)
	if err != nil {
		return err
	}

	_, err = m.client.Blob.Submit(ctx, []*blob.Blob{msgBlob}, nil)
	return err
}

func (m model) Init() tea.Cmd {
	return tea.Batch(waitForActivity(m.headerSub), textarea.Blink)
}
func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*250, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case *header.ExtendedHeader:
		return m, m.handleIncomingHeader(msg)
	case tea.WindowSizeMsg:
		m.height, m.width = msg.Height, msg.Width
		return m, nil
	case tickMsg:
		m.viewport.SetContent(displayMessages(m.messages))
		m.viewport.GotoBottom()
		return m, tickCmd()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			// TODO: loading circle
			m.textarea.Reset()
			m.sendMessage(m.ctx, m.textarea.Value())
			// m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	// h, _ := docStyle.GetFrameSize()

	// m.textarea.SetWidth(h)

	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func (m *model) getPanelDimensions(scaleW, scaleH float64) (w, h int) {
	h, v := docStyle.GetFrameSize()
	return int(float64(m.width)*scaleW) - h, int(float64(m.height)*scaleH) - v
}

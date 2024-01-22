package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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

type model struct {
	ctx       context.Context
	client    *client.Client
	namespace share.Namespace
	addr      state.Address
	headerSub <-chan *header.ExtendedHeader

	viewport    viewport.Model
	messages    []Message
	textarea    textarea.Model
	senderStyle lipgloss.Style
}

type Message struct {
	Username string
	Content  string
	Time     time.Time
}

func main() {
	celestiaClient, err := client.NewClient(
		context.TODO(), // TODO
		os.Args[1],
		os.Args[2],
	)
	if err != nil {
		// TODO
		panic(err)
	}

	addr, err := celestiaClient.State.AccountAddress(context.TODO())
	if err != nil {
		// TODO
		panic(err)
	}

	namespace, err := share.NewBlobNamespaceV0([]byte{5, 5, 5, 5})
	if err != nil {
		// TODO
		panic(err)
	}

	// TODO
	headers, err := celestiaClient.Header.Subscribe(context.TODO())
	if err != nil {
		// TODO
		panic(err)
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

	vp := viewport.New(300, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	m := model{
		ctx:       context.Background(),
		client:    celestiaClient,
		namespace: namespace,
		addr:      addr,
		headerSub: headers,

		textarea:    ta,
		messages:    []Message{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
	}

	if _, err := tea.NewProgram(&m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

func (m *model) handleIncomingHeader(h *header.ExtendedHeader) tea.Cmd {
	height := h.Height()
	blobs, err := m.client.Blob.GetAll(m.ctx, height, []share.Namespace{m.namespace})
	if err != nil {
		//TODO
		panic(err)
	}
	for _, b := range blobs {
		var newMessage Message
		err = json.Unmarshal(b.Data, &newMessage)
		if err != nil {
			// TODO LOG error
		}

		m.messages = append(m.messages, newMessage)
	}

	m.viewport.SetContent(displayMessages(m.messages))
	m.viewport.GotoBottom()
	return waitForActivity(m.headerSub)
}

func displayMessages(messages []Message) string {
	str := ""
	for _, m := range messages {
		str += m.Time.Local().String() + "\t" + m.Username + ": " + m.Content + "\n"
	}
	return str
}

func waitForActivity(sub <-chan *header.ExtendedHeader) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func (m *model) sendMessage(ctx context.Context, content string) error {
	msg := &Message{
		Username: m.addr.String(),
		Content:  content,
		Time:     time.Now(),
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case *header.ExtendedHeader:
		return m, m.handleIncomingHeader(msg)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.sendMessage(m.ctx, m.textarea.Value())
			// m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			m.textarea.Reset()
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

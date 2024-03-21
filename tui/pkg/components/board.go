package components

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enkdress/go-todo/pkg/model"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

var baseUrl = "http://127.0.0.1:3000/v1"

var (
	bgColor          = lipgloss.AdaptiveColor{Light: "#faf4ed", Dark: "#232136"}
	highlightColor   = lipgloss.AdaptiveColor{Light: "#d7827e", Dark: "#ea9a97"}
	docStyle         = lipgloss.NewStyle().Padding(1, 2, 1, 2).Align(lipgloss.Left)
	boardStyle       = lipgloss.NewStyle().Width(50).Align(lipgloss.Center)
	activeBoardStyle = boardStyle.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(highlightColor)
)

type Board struct {
	list   list.Model
	cursor model.Task
}

func NewBoard(title string, tasks []list.Item) *Board {
	model := Board{
		list: list.New(tasks, list.NewDefaultDelegate(), 200, 20),
	}
	model.list.Title = title
	return &model
}

func (b *Board) SetCursor(cursor model.Task) {
	b.cursor = cursor
}

type Kanban struct {
	boards      []Board
	activeBoard int
}

func (b *Kanban) SetCursor(cursor int) {
	b.activeBoard = cursor
}

type Response struct {
	Data []model.Task `json:"data"`
}

func InitialModel() *Kanban {
	res, err := httpClient.Get(fmt.Sprintf("%s/tasks", baseUrl))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	todoItems := []list.Item{}
	doneItems := []list.Item{}

	for _, t := range data.Data {
		if t.IsFinished == 1 {
			doneItems = append(doneItems, t)
		} else {
			todoItems = append(todoItems, t)
		}
	}

	todoBoard := NewBoard("To Do", todoItems)
	doneBoard := NewBoard("Done", doneItems)
	boards := make([]Board, 0, 0)
	boards = append(boards, *todoBoard, *doneBoard)

	return &Kanban{
		boards:      boards,
		activeBoard: 0,
	}
}

func (m Kanban) Init() tea.Cmd {
	return nil
}

func (m Kanban) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	activeBoard := &m.boards[m.activeBoard]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := activeBoardStyle.GetFrameSize()
		activeBoard.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h", "left":
			if m.activeBoard != 0 {
				m.SetCursor(m.activeBoard - 1)
			}
			return m, nil
		case "l", "right":
			if m.activeBoard != len(m.boards)-1 {
				m.SetCursor(m.activeBoard + 1)
			}
			return m, nil
		case "d":
			selectedIndex := activeBoard.list.Cursor()
			selectedItem := activeBoard.list.Items()[selectedIndex].(model.Task)

			selectedItem.IsFinished = 1
			jsonData, err := json.Marshal(selectedItem)
			if err != nil {
				log.Fatal(err)
			}
			jsonBuf := bytes.NewBuffer(jsonData)
			req, err := http.NewRequest("PUT", fmt.Sprintf("%s/tasks", baseUrl), jsonBuf)

			if err != nil {
				log.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			_, err = httpClient.Do(req)

			if err != nil {
				log.Fatal(err)
				return m, nil
			}

			activeBoard.list.RemoveItem(selectedIndex)
			return m, m.boards[1].list.InsertItem(0, selectedItem)
		case "u":
			selectedIndex := activeBoard.list.Cursor()
			selectedItem := activeBoard.list.Items()[selectedIndex].(model.Task)

			selectedItem.IsFinished = 0
			jsonData, err := json.Marshal(selectedItem)
			if err != nil {
				log.Fatal(err)
			}
			jsonBuf := bytes.NewBuffer(jsonData)
			req, err := http.NewRequest("PUT", fmt.Sprintf("%s/tasks", baseUrl), jsonBuf)

			if err != nil {
				log.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			_, err = httpClient.Do(req)

			if err != nil {
				log.Fatal(err)
				return m, nil
			}

			activeBoard.list.RemoveItem(selectedIndex)
			return m, m.boards[0].list.InsertItem(0, selectedItem)
		}
	}

	var cmd tea.Cmd
	activeBoard.list, cmd = activeBoard.list.Update(msg)

	return m, cmd
}

func (k Kanban) View() string {
	var renderedBoards []string
	docStyle := docStyle

	for i, board := range k.boards {
		var style lipgloss.Style
		isActive := i == k.activeBoard

		if isActive {
			style = activeBoardStyle
		} else {
			style = boardStyle
		}

		board.list.SetHeight(30)
		board.list.Styles.TitleBar.Margin(5, 0, 0, 0)
		board.list.Styles.Title.Background(highlightColor)
		renderedBoards = append(renderedBoards, style.Render(board.list.View()))
	}

	// return docStyle.Render(k.boards[0].list.View())
	return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, renderedBoards...))
}

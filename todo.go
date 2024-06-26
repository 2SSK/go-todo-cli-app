package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

// Structre of todo item
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

// Slice of todo items
type Todos []item

// Function add a new todo item
func (t *Todos) Add(task string) {
	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)
}

// Function to mark a todo item as done
func (t *Todos) Complete(index int) error {
	ls := *t
	if index <= 0 || index > len(ls) {
		return errors.New("Invalid index")
	}

	ls[index-1].CompletedAt = time.Now()
	ls[index-1].Done = true

	return nil
}

// Function to delete a todo item
func (t *Todos) Delete(index int) error {
	ls := *t
	if index <= 0 || index > len(ls) {
		return errors.New("Invalid index")
	}

	*t = append(ls[:index-1], ls[index:]...)

	return nil
}

// Function to list all todo items
func (t *Todos) Load(filename string) error {
	file, err := os.ReadFile(filename)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

// Function to store all todo items
func (t *Todos) Store(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// Function to print all todo items
func (t *Todos) Print() {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignCenter, Text: "CreatedAt"},
			{Align: simpletable.AlignCenter, Text: "CompletedAt"},
		},
	}

	var cells [][]*simpletable.Cell

	for idx, item := range *t {
		idx++
		var done string
		if item.Done {
			done = "\u2705"
		} else {
			done = "\u274C"
		}

		task := blue(item.Task)
		if item.Done {
			task = green(item.Task)
		}

		cells = append(cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", idx)},
			{Text: task},
			{Text: done},
			//{Text: fmt.Sprintf("%t", item.Done)},
			{Text: item.CreatedAt.Format(time.RFC822)},
			{Text: item.CompletedAt.Format(time.RFC822)},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: red(fmt.Sprintf("You have %d pending todos", t.CountPending()))},
	}}
	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func (t *Todos) CountPending() int {
	total := 0
	for _, item := range *t {
		if !item.Done {
			total++
		}
	}

	return total
}

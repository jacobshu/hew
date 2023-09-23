package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/term"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func BuildCmdTree() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "hew",
		Short: "A handy haversack with tools ready to hand",
		Args:  cobra.NoArgs,
		RunE:  hewRoot,
	}

	var taskCmd = &cobra.Command{
		Use:   "task",
		Short: "A CLI task management tool for ~slaying~ your to do list.",
		Args:  cobra.NoArgs,
		Run:   taskRoot,
	}
	rootCmd.AddCommand(taskCmd)

	var addCmd = &cobra.Command{
		Use:   "add NAME",
		Short: "Add a new task with an optional project name",
		Args:  cobra.ExactArgs(1),
		RunE:  taskAdd,
	}
	addCmd.Flags().StringP("project", "p", "", "specify a project for your task")
	taskCmd.AddCommand(addCmd)

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all your tasks",
		Args:  cobra.NoArgs,
		RunE:  taskList,
	}
	taskCmd.AddCommand(listCmd)

	var updateCmd = &cobra.Command{
		Use:   "update ID",
		Short: "Update a task by ID",
		Args:  cobra.ExactArgs(1),
		RunE:  taskUpdate,
	}
	updateCmd.Flags().StringP("name", "n", "", "specify a name for your task")
	updateCmd.Flags().StringP("project", "p", "", "specify a project for your task")
	updateCmd.Flags().IntP("status", "s", int(todo), "specify a status for your task")
	taskCmd.AddCommand(updateCmd)

	// task delete command
	var deleteCmd = &cobra.Command{
		Use:   "delete ID",
		Short: "Delete a task by ID",
		Args:  cobra.ExactArgs(1),
		RunE:  taskDelete,
	}
	taskCmd.AddCommand(deleteCmd)

	return rootCmd
}

func hewRoot(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func taskRoot(cmd *cobra.Command, args []string) {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}

func taskAdd(cmd *cobra.Command, args []string) error {
	project, err := cmd.Flags().GetString("project")
	if err != nil {
		return err
	}
	if err := devDb.insertTask(args[0], project); err != nil {
		return err
	}
	return nil
}

func taskList(cmd *cobra.Command, args []string) error {
	tasks, err := devDb.getTasks()
	if err != nil {
		return err
	}
	table := createListTable(tasks)
	fmt.Print(table.View())
	return nil
}

func taskUpdate(cmd *cobra.Command, args []string) error {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	project, err := cmd.Flags().GetString("project")
	if err != nil {
		return err
	}
	prog, err := cmd.Flags().GetInt("status")
	if err != nil {
		return err
	}

	var status string
	switch prog {
	case int(inProgress):
		status = inProgress.String()
	case int(done):
		status = done.String()
	default:
		status = todo.String()
	}

	newTask := task{Name: name, Project: project, Status: status, Created: time.Time{}}
	if err := devDb.updateTask(args[0], newTask); err != nil {
		return err
	}
	return nil
}

func taskDelete(cmd *cobra.Command, args []string) error {
	if err := devDb.deleteTaskById(args[0]); err != nil {
		return err
	}
	return nil
}

func calculateWidth(min, width int) int {
	p := width / 10
	switch min {
	case XS:
		if p < XS {
			return XS
		}
		return p / 2

	case SM:
		if p < SM {
			return SM
		}
		return p
	case MD:
		if p < MD {
			return MD
		}
		return p * 2
	case LG:
		if p < LG {
			return LG
		}
		return p * 3
	default:
		return p
	}
}

const (
	XS int = 5
	SM int = 10
	MD int = 15
	LG int = 20
)

func createListTable(tasks []task) table.Model {
	// get term size
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	log.Printf("terminal size: %+v x %+v", w, h)
	if err != nil {
		// we don't really want to fail it...
		log.Println("unable to calculate height and width of terminal")
	}

	columns := []table.Column{
		{Title: "ID", Width: 24},
		{Title: "Name", Width: calculateWidth(MD, w)},
		{Title: "Project", Width: calculateWidth(SM, w)},
		{Title: "Status", Width: calculateWidth(XS, w)},
		{Title: "Created At", Width: calculateWidth(XS, w)},
	}
	var rows []table.Row
	for _, task := range tasks {
		rows = append(rows, table.Row{
			fmt.Sprintf("%s", task.ID.Hex()),
			task.Name,
			task.Project,
			task.Status,
			task.Created.Format("2006-01-02"),
		})
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		//table.WithFocused(false),
		table.WithHeight(len(tasks)),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		Border(lipgloss.ThickBorder(), true, false, true, false).
		Foreground(lipgloss.Color(ff["cyan"])).
		BorderForeground(lipgloss.Color(ff["yellow"]))
	s.Selected = s.Selected.Bold(false).Foreground(lipgloss.Color(ff["white"]))
	t.SetStyles(s)
	return t
}

// convert tasks to items for a list
func tasksToItems(tasks []task) []list.Item {
	var items []list.Item
	for _, t := range tasks {
		items = append(items, t)
	}
	return items
}

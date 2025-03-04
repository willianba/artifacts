package cmd

import (
	"artifacts/internal"
	"os"

	"github.com/0xN0x/go-artifactsmmo"
	"github.com/0xN0x/go-artifactsmmo/models"
)

func AcceptNewTaskAction(client *artifactsmmo.ArtifactsMMO) *models.TaskData {
	task, err := client.AcceptNewTask()
	if err != nil {
		internal.Logger.Error("Error accepting new task", "err", err)
		os.Exit(1)
	}

	return task
}

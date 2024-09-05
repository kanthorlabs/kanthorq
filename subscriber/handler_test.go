package subscriber

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/stretchr/testify/require"
)

func TestPrinterHandler(t *testing.T) {
	// Save the original stdout
	stdout := os.Stdout

	// Create a pipe to capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	event := entities.NewEvent("system.say_hello", []byte("{\"msg\": \"Hello, World!\"}"))
	task := entities.TaskFromEvent(event)
	// Call the function that uses fmt.Printf
	err := PrinterHandler()(context.Background(), &Message{Event: event, Task: task})
	require.NoError(t, err)

	// Close the writer and restore stdout
	w.Close()
	os.Stdout = stdout

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Check if the output matches the expected value
	require.Contains(t, buf.String(), "PRINTER:")
}

func TestRandomErrorHandler(t *testing.T) {
	event := entities.NewEvent("system.say_hello", []byte("{\"msg\": \"Hello, World!\"}"))
	task := entities.TaskFromEvent(event)

	err := RandomErrorHandler(event.CreatedAt+2)(context.Background(), &Message{Event: event, Task: task})
	require.NoError(t, err)
}

func TestRandomErrorHandler_Error(t *testing.T) {
	event := entities.NewEvent("system.say_hello", []byte("{\"msg\": \"Hello, World!\"}"))
	task := entities.TaskFromEvent(event)
	// Call the function that uses fmt.Printf
	err := RandomErrorHandler(event.CreatedAt)(context.Background(), &Message{Event: event, Task: task})
	require.ErrorContains(t, err, "random error because")
}

func TestRandomErrorHandler_Panic(t *testing.T) {
	event := entities.NewEvent("system.say_hello", []byte("{\"msg\": \"Hello, World!\"}"))
	event.CreatedAt = event.CreatedAt * 2
	task := entities.TaskFromEvent(event)

	defer func() {
		if r := recover(); r != nil {
			require.ErrorContains(t, r.(error), "random error because")
		}
	}()
	// Call the function that uses fmt.Printf
	RandomErrorHandler(event.CreatedAt)(context.Background(), &Message{Event: event, Task: task})
}

func TestRandomErrorHandler_ContextTimeout(t *testing.T) {
	event := entities.NewEvent("system.say_hello", []byte("{\"msg\": \"Hello, World!\"}"))
	task := entities.TaskFromEvent(event)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Call  function that uses fmt.Printf
		err := RandomErrorHandler(event.CreatedAt/100)(ctx, &Message{Event: event, Task: task})
		require.ErrorIs(t, err, context.Canceled)
	}()
	cancel()
}

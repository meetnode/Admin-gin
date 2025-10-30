package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/rivo/tview"
)

// streamLogs streams stdout + stderr from a command into a TextView.
// streamLogs streams stdout + stderr from a command into a TextView.
func streamLogs(cmd *exec.Cmd, view *tview.TextView, label string, color string) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := stripansi.Strip(scanner.Text()) // remove escape codes
			fmt.Fprintf(view, "[%s]%s | %s[-]\n", color, label, line)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := stripansi.Strip(scanner.Text())
			fmt.Fprintf(view, "[%s]%s | %s[-]\n", color, label, line)
		}
	}()

	return nil
}

// killProcessTree attempts to terminate the given command and its children.
// On Windows we use taskkill /T /F. On UNIX-like systems we kill the process group.
func killProcessTree(cmd *exec.Cmd, timeout time.Duration) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	pid := cmd.Process.Pid

	// Best-effort termination
	if runtime.GOOS == "windows" {
		// taskkill will kill the process tree
		_ = exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T", "/F").Run()
	}
	// else {
	// 	// send TERM to the process group (requires Setpgid on start)
	// 	_ = syscall.Kill(-pid, syscall.SIGTERM)
	// }

	// Wait for process to exit with timeout, otherwise force kill
	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	select {
	case <-time.After(timeout):
		// force kill
		if runtime.GOOS == "windows" {
			_ = exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T", "/F").Run()
		}
		//  else {
		// 	_ = syscall.Kill(-pid, syscall.SIGKILL)
		// }
	case <-done:
		// exited
	}
}

func main() {
	// Parse command line flags
	columnLayout := flag.Bool("column", false, "Use column layout instead of row layout")
	flag.Parse()

	app := tview.NewApplication()

	// Explicitly create TextViews (not Boxes!)
	backendView := tview.NewTextView()
	backendView.SetDynamicColors(true)
	backendView.SetBorder(true)
	backendView.SetTitle(" Backend Logs ")
	backendView.SetScrollable(true)
	backendView.SetWordWrap(true)
	backendView.ScrollToEnd()
	backendView.SetChangedFunc(func() { app.Draw() })

	frontendView := tview.NewTextView()
	frontendView.SetDynamicColors(true)
	frontendView.SetBorder(true)
	frontendView.SetTitle(" Frontend Logs ")
	frontendView.SetScrollable(true)
	frontendView.SetWordWrap(true)
	frontendView.ScrollToEnd()
	frontendView.SetChangedFunc(func() { app.Draw() })

	// Set layout direction based on command line flag
	direction := tview.FlexRow
	if *columnLayout {
		direction = tview.FlexColumn
	}

	// Create flexible layout
	flex := tview.NewFlex().
		SetDirection(direction).
		AddItem(backendView, 0, 1, false).
		AddItem(frontendView, 0, 1, false)

	// Prepare commands. On UNIX set a new process group so we can kill children.
	backendCmd := exec.Command("air")
	frontendCmd := exec.Command("npm", "run", "dev", "--prefix", "frontend")

	// if runtime.GOOS != "windows" {
	// 	backendCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	// 	frontendCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	// }

	if err := streamLogs(backendCmd, backendView, "Backend", "green"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start backend: %v\n", err)
	}
	if err := streamLogs(frontendCmd, frontendView, "Frontend", "cyan"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start frontend: %v\n", err)
	}

	// Setup signal handling to ensure subprocesses are killed on exit.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Run UI in a goroutine so we can listen for signals concurrently.
	runErrCh := make(chan error, 1)
	go func() {
		runErrCh <- app.SetRoot(flex, true).EnableMouse(true).Run()
	}()

	select {
	case sig := <-sigCh:
		// Received interrupt/terminate
		fmt.Fprintf(os.Stderr, "received signal: %v, shutting down...\n", sig)
		// stop the UI (in case it's still running)
		app.Stop()
	case err := <-runErrCh:
		if err != nil {
			fmt.Fprintf(os.Stderr, "application error: %v\n", err)
		}
	}

	// Give UI some time to teardown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
	default:
	}

	// Ensure subprocesses are cleaned up
	killProcessTree(backendCmd, 3*time.Second)
	killProcessTree(frontendCmd, 3*time.Second)

	// Small pause to let termination complete
	time.Sleep(200 * time.Millisecond)
}

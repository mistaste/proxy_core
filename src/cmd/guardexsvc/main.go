// guardexsvc is the privileged Windows service that owns the wintun
// adapter and the routing table edits. The Flutter desktop app talks to
// it over a named pipe so it never needs to run elevated itself.
//
// Usage:
//   guardexsvc install      Register as an auto-start Windows service
//   guardexsvc uninstall    Stop and remove the service
//   guardexsvc start        Start the installed service
//   guardexsvc stop         Stop the installed service
//   guardexsvc debug        Run in the foreground (no SCM), handy for dev
//
// When started by the Service Control Manager with no args, the binary
// runs as a service.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"segment/service_win"
)

const serviceName = "GuardexVPN"
const serviceDisplayName = "Guardex VPN Service"
const serviceDescription = "Manages the Guardex wintun adapter and routing so the UI can run without elevation."

type guardexService struct{}

func (guardexService) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (ssec bool, errno uint32) {
	const accepted = svc.AcceptStop | svc.AcceptShutdown
	s <- svc.Status{State: svc.StartPending}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := service_win.NewServer()
	errCh := make(chan error, 1)
	go func() { errCh <- server.Run(ctx) }()

	s <- svc.Status{State: svc.Running, Accepts: accepted}

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				s <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			}
		case err := <-errCh:
			if err != nil {
				log.Printf("pipe server exited: %v", err)
			}
			break loop
		}
	}

	s <- svc.Status{State: svc.StopPending}
	cancel()
	service_win.StopActiveSession()
	return false, 0
}

func setupFileLog() {
	dir := os.Getenv("ProgramData")
	if dir == "" {
		dir = `C:\ProgramData`
	}
	dir = filepath.Join(dir, "Guardex")
	_ = os.MkdirAll(dir, 0755)
	f, err := os.OpenFile(filepath.Join(dir, "service.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	// Redirect os.Stdout / os.Stderr to the log file too so any panic
	// traces, tun2socks slog output, or log.Fatalf from third-party code
	// is visible to us — otherwise a service process has no stderr sink.
	os.Stdout = f
	os.Stderr = f
}

func main() {
	if len(os.Args) < 2 {
		// Launched by SCM.
		setupFileLog()
		log.Println("guardexsvc: starting (SCM mode)")
		defer func() {
			if r := recover(); r != nil {
				log.Printf("guardexsvc: top-level panic: %v", r)
			}
		}()
		if err := svc.Run(serviceName, guardexService{}); err != nil {
			log.Fatalf("service run: %v", err)
		}
		return
	}

	switch os.Args[1] {
	case "install":
		if err := installService(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("installed")
	case "uninstall":
		if err := uninstallService(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("uninstalled")
	case "start":
		if err := startService(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("started")
	case "stop":
		if err := controlService(svc.Stop); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("stopped")
	case "debug":
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		fmt.Println("guardexsvc: running in foreground on", service_win.PipeName)
		if err := service_win.NewServer().Run(ctx); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Fprintln(os.Stderr, "unknown subcommand:", os.Args[1])
		os.Exit(2)
	}
}

func installService() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe, err = filepath.Abs(exe)
	if err != nil {
		return err
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	if s, err := m.OpenService(serviceName); err == nil {
		s.Close()
		return fmt.Errorf("service %s already installed", serviceName)
	}

	s, err := m.CreateService(serviceName, exe, mgr.Config{
		DisplayName:      serviceDisplayName,
		Description:      serviceDescription,
		StartType:        mgr.StartAutomatic,
		ErrorControl:     mgr.ErrorNormal,
	})
	if err != nil {
		return err
	}
	defer s.Close()
	return nil
}

func uninstallService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("open service: %w", err)
	}
	defer s.Close()

	_, _ = s.Control(svc.Stop)
	if err := s.Delete(); err != nil {
		return err
	}
	return nil
}

func startService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer s.Close()
	return s.Start()
}

func controlService(c svc.Cmd) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer s.Close()
	_, err = s.Control(c)
	return err
}

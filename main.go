package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"

	"github.com/manishbadgotra/vehicle-details/controllers"
	"github.com/manishbadgotra/vehicle-details/database"
	"github.com/manishbadgotra/vehicle-details/utils"
)

// ─── SERVICE NAME ─────────────────────────────────────────────────────────────
const svcName = "VehicleGetAPI"

func initDBPath() error {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	godotenv.Load()

	if err := database.CreateDB(); err != nil {
		return fmt.Errorf("unable to create tables in database table: %v", err.Error())
	}

	return nil
}

// ─── WINDOWS SERVICE HANDLER ─────────────────────────────────────────────────
// windowsService implements svc.Handler — this is what makes the .exe
// talk the SCM (Service Control Manager) protocol natively.
// SCM sends control signals (Stop, Shutdown, Interrogate) via the channel.
type windowsService struct{}

func (ws *windowsService) Execute(args []string, req <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	status <- svc.Status{State: svc.StartPending}

	// ✅ Move initDBPath here — after SCM handshake begins
	if err := initDBPath(); err != nil {
		log.Printf("DB init failed: %v", err)
		return false, 1 // clean exit with error code
	}

	serverErr := make(chan error, 1)

	var PORT = os.Getenv("PORT")
	if PORT == "" {
		serverErr <- fmt.Errorf("PORT is not defined")
	}

	srv := buildServer(PORT)

	// Start HTTP server in background
	go func() {
		log.Println("HTTP server starting on :5898")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	log.Printf("%s service running", svcName)

loop:
	for {
		select {
		case err := <-serverErr:
			log.Printf("HTTP server crashed: %v", err)
			break loop
		case c := <-req:
			switch c.Cmd {
			case svc.Interrogate:
				// SCM health-poll — echo current status back
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				log.Printf("%s service stopping", svcName)
				status <- svc.Status{State: svc.StopPending}

				// Graceful HTTP shutdown with a deadline
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				if err := srv.Shutdown(ctx); err != nil {
					log.Printf("graceful shutdown error: %v", err)
				}
				break loop
			default:
				log.Printf("unexpected SCM control request: %d", c.Cmd)
			}
		}
	}

	return false, 0
}

func buildServer(port string) *http.Server {
	mux := http.NewServeMux()

	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowedHeaders:   []string{"Content-Type"},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           100,
	}

	mux.Handle("GET /v1/vehicles", http.HandlerFunc(controllers.GetAllVehicles))

	return &http.Server{
		Addr:         "0.0.0.0:" + port,
		Handler:      cors.New(corsOptions).Handler(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func main() {

	go utils.GetVehiclesFromList()

	isService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to detect service context: %v", err)
	}

	if isService {
		// ── Running as a Windows Service ──────────────────────────────────
		// Wire up the Windows Event Log so logs appear in Event Viewer
		elog, err := eventlog.Open(svcName)
		if err == nil {
			defer elog.Close()
			elog.Info(1, fmt.Sprintf("%s service starting", svcName))
		}

		// Hand control to SCM — this blocks until the service is stopped
		if err := svc.Run(svcName, &windowsService{}); err != nil {
			log.Fatalf("%s service failed: %v", svcName, err)
		}
	} else {
		// ── Running interactively (dev / debug mode) ───────────────────────
		// debug.Run mimics the SCM loop in the console so you can Ctrl+C to stop
		log.Printf("Running in debug (console) mode — use Ctrl+C to stop")
		if err := debug.Run(svcName, &windowsService{}); err != nil {
			log.Fatalf("debug run failed: %v", err)
		}
	}
}

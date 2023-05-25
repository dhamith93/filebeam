package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dhamith93/share_core/internal/api"
	"github.com/dhamith93/share_core/internal/database"
	"github.com/dhamith93/share_core/internal/system"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// App struct
type App struct {
	ctx           context.Context
	db            database.Database
	apiServer     api.Server
	listeningPort string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.listeningPort = "9292"
	a.db = database.Database{}
	a.db.CreateDB(filepath.Join(os.TempDir(), "filebeam.db"))
	// defer a.db.Db.Close()

	a.apiServer = api.Server{Database: &a.db, Key: generateKey(6)}

	// var wg sync.WaitGroup
	// wg.Add(1)

	go func() {
		// defer wg.Done()
		lis, err := net.Listen("tcp", ":"+a.listeningPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		api.RegisterFileServiceServer(grpcServer, &a.apiServer)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	ch := make(chan string)
	go collectLocalDevicesWithServiceRunning(a.listeningPort, ch)

	for {
		resp := <-ch
		if resp == "done" {
			break
		}
		err := a.db.AddDevice(resp)
		if err != nil {
			log.Println(err)
		}
	}
}

func (a *App) shutdown(ctx context.Context) {
	a.db.Db.Close()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetDevices() []string {
	devices, err := a.db.GetDevices()
	if err != nil {
		log.Fatal(err.Error())
	}
	return devices
}

func (a *App) GetKey() string {
	return a.apiServer.Key
}

func createClient(endpoint string) (*grpc.ClientConn, api.FileServiceClient, context.Context, context.CancelFunc) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("connection error: " + err.Error())
		return nil, nil, nil, nil
	}
	c := api.NewFileServiceClient(conn)
	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), nil), time.Second*10)
	return conn, c, ctx, cancel
}

func collectLocalDevicesWithServiceRunning(port string, ch chan string) {
	ips := system.GetLocalIPs()
	var wg sync.WaitGroup
	wg.Add(len(ips) - 1)
	for _, ip := range ips {
		if ip == system.GetIp() {
			continue
		}
		go func(ip string, ch chan string) {
			defer wg.Done()
			host := ip + ":" + port
			conn, c, ctx, cancel := createClient(host)
			if conn == nil {
				return
			}
			defer conn.Close()
			defer cancel()
			_, err := c.Hello(ctx, &api.Void{})
			if err != nil {
				return
			}
			ch <- host
		}(ip, ch)
	}

	wg.Wait()
	ch <- "done"
}

func generateKey(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset)-1)]
	}
	return string(b)
}

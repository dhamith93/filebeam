package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/dhamith93/filebeam/internal/api"
	"github.com/dhamith93/filebeam/internal/database"
	"github.com/dhamith93/filebeam/internal/file"
	"github.com/dhamith93/filebeam/internal/system"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// App struct
type App struct {
	ctx           context.Context
	db            database.MemDatabase
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
	a.db = database.MemDatabase{}
	a.db.CreateDB()
	a.apiServer = api.Server{Database: &a.db, Key: generateKey(6)}

	go func() {
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

	a.refreshDevices()

	ticker := time.NewTicker(time.Second * 5)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				a.refreshDevices()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	ticker2 := time.NewTicker(time.Millisecond * 500)
	quit2 := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker2.C:
				handlePendingTransfers(&a.db, a.listeningPort)
			case <-quit2:
				ticker2.Stop()
				return
			}
		}
	}()
}

func (a *App) refreshDevices() {
	ch := make(chan string)
	go collectLocalDevicesWithServiceRunning(a.listeningPort, ch)

	for {
		resp := <-ch
		if resp == "done" {
			break
		}
		existing, _ := a.db.GetDevices()
		exists := false
		for _, d := range existing {
			if d == resp {
				exists = true
			}
		}
		if !exists {
			err := a.db.AddDevice(resp)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func handlePendingTransfers(db *database.MemDatabase, listeningPort string) {
	if db.FileTransfersInProgress(1) {
		return
	}
	files, err := db.GetPendingTransfers()
	if err != nil {
		log.Fatalf("failed to load transfers: %s", err)
	}
	for _, f := range files {
		if !f.IsFile() {
			db.UpdateTransferStatus(f.Dest, f.Path, "cannot_read_file")
			continue
		}
		conn, c, ctx, cancel := createClient(f.Dest + ":" + listeningPort)
		if conn == nil {
			log.Printf("error creating connection")
			return
		}
		defer conn.Close()
		defer cancel()
		_, err := c.FilePush(ctx, &api.FilePushRequest{File: getAPIFile(f), Key: f.Key, Port: listeningPort})
		if err != nil {
			log.Printf("error sending data: %s", err.Error())
			db.UpdateTransferStatus(f.Dest, f.Path, "error")
		} else {
			err = db.UpdateTransferStatus(f.Dest, f.Path, "processing")
			if err == nil { // successfully set the file to processing
				break
			}
		}
	}
}

func (a *App) shutdown(ctx context.Context) {

}

func (a *App) domready(ctx context.Context) {

}

func (a *App) GetDevices() []string {
	devices, err := a.db.GetDevices()
	if err != nil {
		log.Fatal(err.Error())
	}
	return devices
}

func (a *App) GetDirectoryContent(path string) ([]File, error) {
	return getDirectoryContent(path)
}

func (a *App) GetHomeDir() string {
	return getHomeDir()
}

func (a *App) GetKey() string {
	return a.apiServer.Key
}

func (a *App) GetIp() string {
	return system.GetIp()
}

func (a *App) AddToQueue(files []File, host string, key string) error {
	for _, f := range files {
		file := file.CreateFile(f.Path)
		err := a.db.AddTransfer(host, key, file.Name, file.Type, file.Extension, file.Path, file.Size)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) GetTransfers() ([]database.Transfer, error) {
	return a.db.GetAllTransfers()
}

func (a *App) CancelTransfer(ip string, filename string, isDownload bool) error {
	if isDownload {
		return a.db.UpdateIncomingTransferStatus(ip, filename, "cancelled")
	}
	return a.db.UpdateTransferStatus(ip, filename, "cancelled")
}

func getAPIFile(in file.File) *api.File {
	return &api.File{
		Name:      in.Name,
		Size:      in.Size,
		Type:      in.Type,
		Extension: in.Extension,
	}
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
			ch <- ip
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

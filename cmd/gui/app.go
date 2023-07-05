package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/dhamith93/filebeam/internal/api"
	"github.com/dhamith93/filebeam/internal/file"
	"github.com/dhamith93/filebeam/internal/queue"
	"github.com/dhamith93/filebeam/internal/system"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// App struct
type App struct {
	ctx           context.Context
	devices       []string
	apiServer     api.Server
	uploadQueue   queue.Queue
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
	a.apiServer = api.CreateServer()
	a.apiServer.Port = a.listeningPort
	a.apiServer.Key = generateKey(16)

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
				a.handlePendingTransfers(a.apiServer.UploadQueue, a.listeningPort)
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
		if !slices.Contains(a.devices, resp) {
			a.devices = append(a.devices, resp)
		}
	}
}

func (a *App) handlePendingTransfers(q *queue.Queue, listeningPort string) {
	if a.uploadQueue.FileTransfersInProgress(1) {
		return
	}
	files := a.uploadQueue.GetPendingTransfers()

	for _, f := range files {
		if !f.File.IsFile() {
			a.uploadQueue.UpdateTransferStatus(f.Ip+":"+f.FilePort, f.File, "cannot_read_file")
			continue
		}

		err := a.apiServer.PushFile(f.Ip+":"+listeningPort, f.File)

		if err != nil {
			log.Println(err.Error())
		}
		a.uploadQueue.Remove(f)
		break
	}
}

func (a *App) shutdown(ctx context.Context) {

}

func (a *App) domready(ctx context.Context) {

}

func (a *App) GetDevices() []string {
	return a.devices
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

func (a *App) GetPendingDownloads() []queue.Transfer {
	return a.apiServer.DownloadQueue.GetPendingTransfers()
}

func (a *App) AddToQueue(files []File, host string, key string) {
	for _, f := range files {
		file := file.CreateFile(f.Path)
		file.Key = key
		a.uploadQueue.AddToQueue(host+":xxxx", key, file, false)
	}
}

func (a *App) GetTransfers() []queue.Transfer {
	return append(a.apiServer.DownloadQueue.Items, a.apiServer.UploadQueue.Items...)
}

func (a *App) CancelTransfer(host string, filename string, isDownload bool) {
	if isDownload {
		a.apiServer.StopDownload(host, filename)
	} else {
		a.apiServer.StopUpload(host, filename)
	}
}

func (a *App) DownloadTransfer(host string, filename string) {
	a.apiServer.StartDownload(host, filename)
}

func (a *App) AmIRunningOnMacos() bool {
	return runtime.GOOS == "darwin"
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

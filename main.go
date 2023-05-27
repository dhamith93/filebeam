package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/dhamith93/share_core/internal/api"
	"github.com/dhamith93/share_core/internal/database"
	"github.com/dhamith93/share_core/internal/file"
	"github.com/dhamith93/share_core/internal/system"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	log.Println("Hello, World!")
	var listeningPort, dbPath string
	flag.StringVar(&listeningPort, "listen", "9292", "Port to listen")
	flag.StringVar(&dbPath, "db-path", "file_process.db", "Path to the sqlite DB")
	flag.Parse()

	db := database.MemDatabase{}
	db.CreateDB()

	var wg sync.WaitGroup
	wg.Add(2)
	s := api.Server{Database: &db, Key: generateKey(6)}
	db.SetKey(s.Key)
	log.Println(s.Key)

	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", ":"+listeningPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		api.RegisterFileServiceServer(grpcServer, &s)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	ticker := time.NewTicker(time.Millisecond * 500)
	quit := make(chan struct{})
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				handlePendingTransfers(&db, listeningPort)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	ch := make(chan string)
	go collectLocalDevicesWithServiceRunning(listeningPort, ch)

	for {
		resp := <-ch
		if resp == "done" {
			break
		}
		err := db.AddDevice(resp)
		if err != nil {
			log.Println(err)
		}
	}

	wg.Wait()
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
		conn, c, ctx, cancel := createClient(f.Dest)
		if conn == nil {
			log.Printf("error creating connection")
			return
		}
		defer conn.Close()
		defer cancel()
		_, err := c.FilePush(ctx, &api.FilePushRequest{File: getAPIFile(f), Key: f.Key, Port: listeningPort})
		if err != nil {
			log.Printf("error sending data: %s", err.Error())
			db.UpdateTransferStatus(f.Dest, f.Path, err.Error())
		}
		err = db.UpdateTransferStatus(f.Dest, f.Path, "processing")
		if err == nil { // successfully set the file to processing
			break
		}
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

func getAPIFile(in file.File) *api.File {
	return &api.File{
		Name:      in.Name,
		Size:      in.Size,
		Type:      in.Type,
		Extension: in.Extension,
	}
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

package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/dhamith93/share_core/internal/api"
	"github.com/dhamith93/share_core/internal/file"
	"github.com/dhamith93/share_core/internal/system"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	log.Println("Hello, World!")
	var mode, listeningPort, host, path string

	flag.StringVar(&mode, "mode", "rx", "tx/rx")
	flag.StringVar(&listeningPort, "listen", "9292", "Port to listen")
	flag.StringVar(&host, "host", "192.168.123.201:8080", "host")
	flag.StringVar(&path, "path", "", "File path")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)
	s := api.Server{}
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

	if mode == "tx" {
		conn, c, ctx, cancel := createClient(host)
		if conn == nil {
			log.Printf("error creating connection")
			return
		}
		defer conn.Close()
		defer cancel()
		s.PendingFile = path
		_, err := c.FilePush(ctx, &api.FilePushRequest{File: getAPIFile(file.CreateFile(path)), Port: listeningPort})
		if err != nil {
			log.Printf("error sending data")
			os.Exit(1)
		}
	}

	ch := make(chan string)
	go collectLocalDevicesWithServiceRunning(listeningPort, ch)

	for {
		resp := <-ch
		if resp == "done" {
			break
		}
		log.Println(resp)
	}

	wg.Wait()
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
			conn, c, ctx, cancel := createClient(ip + ":" + port)
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

package api

import (
	context "context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dhamith93/filebeam/internal/database"
	"github.com/dhamith93/filebeam/internal/file"
	fileservice "github.com/dhamith93/filebeam/internal/file_service"
	"github.com/dhamith93/filebeam/internal/queue"
	"github.com/dhamith93/filebeam/internal/system"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Server struct {
	Key         string
	FileService fileservice.FileService
	PendingFile string
	Database    *database.MemDatabase
	Queue       *queue.Queue
	UnimplementedFileServiceServer
}

func (s *Server) Init() {
	s.FileService = fileservice.FileService{}
	queue := queue.CreateQueue()
	s.Queue = &queue
}

func (s *Server) FilePush(ctx context.Context, fileRequest *FilePushRequest) (*FilePushResponse, error) {
	queue := queue.CreateQueue()
	s.Queue = &queue
	if fileRequest.Key != s.Key {
		return nil, fmt.Errorf("key does not match")
	}
	p, _ := peer.FromContext(ctx)
	ip := strings.Split(p.Addr.String(), ":")[0]
	s.FileService.Database = s.Database
	s.FileService.Queue = s.Queue
	s.Database.AddIncomingTransfer(ip, fileRequest.File.Name, fileRequest.File.Type, fileRequest.File.Extension, fileRequest.File.Size)
	go s.FileService.Receive(s.getFileStruct("", fileRequest.File))
	s.sendClearToSend(ip+":"+fileRequest.Port, fileRequest.File)
	return &FilePushResponse{Accepted: true}, nil
}

func (s *Server) ClearToSend(ctx context.Context, fileResponse *FilePushResponse) (*Void, error) {
	s.FileService.Database = s.Database
	s.FileService.Queue = s.Queue
	go s.FileService.Send(fileResponse.Host+":"+fileResponse.Port, s.getFileStruct(fileResponse.Host, fileResponse.File))
	return &Void{}, nil
}

func (s *Server) Hello(ctx context.Context, void *Void) (*Void, error) {
	return &Void{}, nil
}

func (s *Server) getFileStruct(dest string, in *File) file.File {
	path := s.Database.GetFilePath(dest, in.Name)
	return file.File{
		Id:        in.Id,
		Name:      in.Name,
		Type:      in.Type,
		Extension: in.Extension,
		Size:      in.Size,
		Path:      path,
	}
}

func (s *Server) sendClearToSend(host string, file *File) {
	conn, c, ctx, cancel := createClient(host)
	if conn == nil {
		log.Printf("error creating connection")
		return
	}
	defer conn.Close()
	defer cancel()
	_, err := c.ClearToSend(ctx, &FilePushResponse{File: file, Host: system.GetIp(), Port: s.FileService.Port})
	if err != nil {
		log.Println(err)
	}
}

func createClient(endpoint string) (*grpc.ClientConn, FileServiceClient, context.Context, context.CancelFunc) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("connection error: " + err.Error())
		return nil, nil, nil, nil
	}
	c := NewFileServiceClient(conn)
	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), nil), time.Second*10)
	return conn, c, ctx, cancel
}

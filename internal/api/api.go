package api

import (
	context "context"
	"fmt"
	"log"
	"strings"
	"time"

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
	Key           string
	FileService   fileservice.FileService
	PendingFile   file.File
	DownloadQueue *queue.Queue
	UploadQueue   *queue.Queue
	UnimplementedFileServiceServer
}

func CreateServer() Server {
	downloadQueue := queue.CreateQueue()
	uploadQueue := queue.CreateQueue()

	fileService := fileservice.FileService{}
	fileService.UploadQueue = &uploadQueue
	fileService.DownloadQueue = &downloadQueue

	return Server{
		UploadQueue:   &uploadQueue,
		DownloadQueue: &downloadQueue,
		FileService:   fileservice.FileService{},
	}
}

func (s *Server) FilePush(ctx context.Context, fileRequest *FilePushRequest) (*FilePushResponse, error) {
	if fileRequest.Key != s.Key {
		return nil, fmt.Errorf("key does not match")
	}
	p, _ := peer.FromContext(ctx)
	ip := strings.Split(p.Addr.String(), ":")[0]
	s.FileService.DownloadQueue = s.DownloadQueue
	go s.FileService.Receive(s.getFileStruct(fileRequest.File))
	s.sendClearToSend(ip+":"+fileRequest.Port, fileRequest.File)
	return &FilePushResponse{Accepted: true}, nil
}

func (s *Server) ClearToSend(ctx context.Context, fileResponse *FilePushResponse) (*Void, error) {
	s.FileService.UploadQueue = s.UploadQueue
	s.UploadQueue.AddToQueue(fileResponse.Host+":"+fileResponse.Port, "", s.PendingFile)
	go s.FileService.Send(fileResponse.Host+":"+fileResponse.Port, s.PendingFile)
	// s.UploadQueue.UpdateFilePortOfTransfer(fileResponse.Host, "xxxx", fileResponse.Port, file)
	return &Void{}, nil
}

func (s *Server) Hello(ctx context.Context, void *Void) (*Void, error) {
	return &Void{}, nil
}

func (s *Server) PushFile(host string, f file.File) error {
	conn, c, ctx, cancel := createClient(host)
	if conn == nil {
		return fmt.Errorf("error creating connection")
	}
	defer conn.Close()
	defer cancel()
	s.PendingFile = f
	_, err := c.FilePush(ctx, &FilePushRequest{File: s.getAPIFile(f), Key: f.Key, Port: strings.Split(host, ":")[1]})
	if err != nil {
		s.PendingFile = file.File{}
	}
	return err
}

func (s *Server) getFileStruct(in *File) file.File {
	return file.File{
		Id:        in.Id,
		Name:      in.Name,
		Type:      in.Type,
		Extension: in.Extension,
		Size:      in.Size,
	}
}

func (s *Server) getAPIFile(in file.File) *File {
	return &File{
		Name:      in.Name,
		Size:      in.Size,
		Type:      in.Type,
		Extension: in.Extension,
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

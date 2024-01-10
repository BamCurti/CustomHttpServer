package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/http-server-starter-go/app/server"
)

func main() {
	ip := "0.0.0.0"
	port := "4221"

	// enable directory flag
	dirFlag := flag.String("directory", ".", "Directory to serve files from")
	flag.Parse()
	server.DirFlag = dirFlag

	s, err := server.New(ip, port)

	if err != nil {
		log.Fatalln(err)
	}

	s.Get("/", helloServer)
	s.Get("/echo/{word}", echoHandler)
	s.Get("/user-agent", userAgent)
	s.Get("/files/{filename}", getFile)
	s.Post("/files/{filename}", postFile)
	s.Serve()
}

func helloServer(req *server.HttpRequest, res *server.HttpResponse) {
	res.ContentType = server.TEXT_PLAIN
	res.StatusCode = server.OK_MSG
	res.Write("")
}

func echoHandler(req *server.HttpRequest, res *server.HttpResponse) {
	payload, _ := req.UrlParam("word")
	res.ContentType = server.TEXT_PLAIN
	res.StatusCode = server.OK_MSG
	res.Write(payload)
}

func userAgent(req *server.HttpRequest, res *server.HttpResponse) {
	payload := req.Headers["User-Agent"]
	res.ContentType = server.TEXT_PLAIN
	res.StatusCode = server.OK_MSG
	res.Write(payload)
}

func getFile(req *server.HttpRequest, res *server.HttpResponse) {
	fileName, _ := req.UrlParam("filename")
	filePath := filepath.Join(*server.DirFlag, fileName)

	fileContent, err := os.ReadFile(filePath)
	var payload string

	if err != nil {
		payload = ""
		res.StatusCode = server.NOT_FOUND_MSG
		res.ContentType = server.TEXT_PLAIN
	} else {
		payload = string(fileContent)
		res.StatusCode = server.OK_MSG
		res.ContentType = server.APP_OCTET_STREAM
	}
	res.Write(payload)
}

func postFile(req *server.HttpRequest, res *server.HttpResponse) {
	fileName, _ := req.UrlParam("filename")
	filepath := filepath.Join(*server.DirFlag, fileName)

	err := createFile(filepath, req.Content)
	if err != nil {
		res.StatusCode = server.INTERNAL_ERROR_MSG
		res.ContentType = server.TEXT_PLAIN
		res.Write("")
	}

	res.StatusCode = server.CREATED_MSG
	res.ContentType = server.APP_OCTET_STREAM
	res.Write(req.Content)
}

func createFile(filePath, body string) error {
	f, err := os.Create(filePath)
	if err != nil {
		log.Println("File not created", filePath)
		return err
	}

	defer f.Close()

	_, err = f.Write([]byte(body))
	if err != nil {
		log.Println("File not written", err.Error())
		os.Remove(filePath)
		return err
	}
	return nil
}

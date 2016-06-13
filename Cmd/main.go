package main

import (
	"io"
	"os"
	"log"
	"fmt"
	"errors"
)

const (
	LUN1 = "lun1.bin"
	LUN2 = "lun2.bin"

	logfile	= "file.log"
)

type ServerConfig struct{
//	Login, Pass string
	Address		string
	TargetName	string
	LUNs		[]string
}

// Device is an interface for block device.
type Devicer interface{
	ReadAt(p []byte, off int64) (n int, err error)
	WriteAt(p []byte, off int64) (n int, err error)
	Discard(size int64, off int64) error
	Flush() error
	Close() error
}

type Device io.ReadWriteCloser

type Server struct{
	TargetName	string
	LUNs		[]*Device
}

var _ interface{
//	ListenAndServe() error
	Close() error
} = (*Server)(nil)

var Logger *log.Logger

func main() {
	var dev Device

	Logger = createLogger()

	srv := NewServer(dev, &ServerConfig{	"172.24.1.3",
						"iqn.2016-04.npp.sit-1920:storage",
						[]string{LUN1, LUN2}})
	fmt.Println(srv)

	return
}

func NewServer(dev Device, conf *ServerConfig) (srv *Server) {
	srv.TargetName = conf.TargetName
	srv.LUNs = make([]*Device, len(conf.LUNs))
	for i, fileName := range conf.LUNs {
		fi, err := os.OpenFile(fileName, os.O_RDWR, os.ModePerm)
		if err != nil {
			Logger.Println("Can't open " + fileName + " file. ", err)
			continue
		}
		srv.LUNs[i] = NewDevice(fi)
	}

	return srv
}

func createLogger() (l *log.Logger) {
	w, err := os.OpenFile(logfile, os.O_WRONLY, os.ModePerm)
	if err != nil {
		w = os.Stdout
	}
	l = log.New(w, "", log.LstdFlags)
	return l
}

func NewDevice(fi *os.File) *Device {
	d := Device(fi)
	return &d
}

func (s *Server)Close() (err error) {
	errStr := ""
	for _, lun := range s.LUNs {
		if err := lun.Close(); err != nil {
			Logger.Println(err)
			errStr += fmt.Sprintln(err)
		}
	}
	if errStr == "" {
		err = nil
	} else {
		err = errors.New(errStr)
	}
	return err
}

func (d *Device)Close() error {
	return io.Closer()
}
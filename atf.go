package atf

import (
	"os"
	"fmt"
	"errors"
	"log"
)

const (
	logfile	= "file.log"
)

var Logger *log.Logger

type ServerConfig struct{
	//	Login, Pass string
	Address		string
	TargetName	string
	LUNs		[]string
	Log		*log.Logger
}

// Device is an interface for block device.
type Devicer interface{
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	ReadAt(p []byte, off int64) (n int, err error)
	WriteAt(p []byte, off int64) (n int, err error)
	Discard(size int64, off int64) error
	Flush() error
	Close() error
}

type Device struct {
	I	Devicer
	F	*os.File
}

type Server struct{
	TargetName	string
	LUNs		[]Device
	Log		log.Logger
}

var _ interface{
	//	ListenAndServe() error
	Close() error
} = (*Server)(nil)

func NewServer(dev Device, conf *ServerConfig) *Server {
	var srv Server

	srv.TargetName = conf.TargetName
	srv.LUNs = make([]Device, len(conf.LUNs))
	for i, fileName := range conf.LUNs {
		fi, err := os.OpenFile(fileName, os.O_RDWR, os.ModePerm)
		if err != nil {
			Logger.Println(err)
			continue
		}
		srv.LUNs[i] = *NewDevice(fi)
		//srv.LUNs[i] = fi
	}

	return &srv
}

func CreateLogger() *log.Logger {
	w, err := os.OpenFile(logfile, os.O_WRONLY, os.ModePerm)
	if err != nil {
		w = os.Stdout
	}
	return log.New(w, "", log.LstdFlags)
}

func NewDevice(fi *os.File) *Device {
	//d := Device(fi)
	return &Device{F:fi}
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

func (d *Device)Read(p []byte) (n int, err error) {
	return d.F.Read(p)
}

func (d *Device)Write(p []byte) (n int, err error) {
	return d.F.Write(p)
}
func (d *Device)ReadAt(p []byte, off int64) (n int, err error) {
	return d.F.ReadAt(p, off)
}

func (d *Device)WriteAt(p []byte, off int64) (n int, err error) {
	return d.F.WriteAt(p, off)
}

func (d *Device)Discard(size int64, off int64) error {
	return nil
}
func (d *Device)Flush() error {
	return nil
}
func (d *Device)Close() error {
	return d.F.Close()
}

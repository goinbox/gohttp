package httpserver

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/goinbox/golog"
	"github.com/goinbox/pcontext"
)

const (
	defaultReadTimeout  = time.Second * 60
	defaultWriteTimeout = time.Second * 60

	environKeyListenerFd = "LISTENER_FD"
)

type listener struct {
	net.Listener

	File *os.File
}

type Server struct {
	*http.Server

	ln         *listener
	shutdownCh chan struct{}
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		Server: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadHeaderTimeout: defaultReadTimeout,
			WriteTimeout:      defaultWriteTimeout,
		},
	}
}

func (s *Server) ListenAndServe(ctx pcontext.Context) error {
	ln, err := s.createListener(ctx)
	if err != nil {
		return fmt.Errorf("createListener error: %w", err)
	}

	ctx.Logger().Notice("ListenAndServe", []*golog.Field{
		{
			Key:   "addr",
			Value: s.Addr,
		},
		{
			Key:   "pid",
			Value: os.Getpid(),
		},
	}...)
	err = s.serve(ctx, ln)
	if err != nil {
		return fmt.Errorf("server.Serve error: %s", err)
	}

	return nil
}

func (s *Server) serve(ctx pcontext.Context, ln *listener) error {
	s.ln = ln
	s.shutdownCh = make(chan struct{})

	go s.signalHandler(ctx)

	err := s.Server.Serve(ln)
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("Server.Serve error: %w", err)
		}
	}

	ctx.Logger().Notice("server shutdown, wait for all connections complete")
	<-s.shutdownCh

	return nil
}

func (s *Server) createListener(ctx pcontext.Context) (*listener, error) {
	ln, err := s.createListenerFromEnviron()
	if err != nil {
		return nil, fmt.Errorf("createListenerFromEnviron error: %w", err)
	}
	if ln != nil {
		ctx.Logger().Debug("createListenerFromEnviron")
		return ln, nil
	}

	var lc net.ListenConfig
	tln, err := lc.Listen(ctx, "tcp", s.Addr)
	if err != nil {
		return nil, fmt.Errorf("net.Listen error: %w", err)
	}
	f, err := tln.(*net.TCPListener).File()
	if err != nil {
		return nil, fmt.Errorf("TCPListener.File error: %w", err)
	}

	ctx.Logger().Debug("createListenerFromNetAddr")
	return &listener{
		Listener: tln,
		File:     f,
	}, nil
}

func (s *Server) createListenerFromEnviron() (*listener, error) {
	fdStr := os.Getenv(environKeyListenerFd)
	if fdStr == "" {
		return nil, nil
	}

	fd, err := strconv.Atoi(fdStr)
	if err != nil {
		return nil, fmt.Errorf("strconv.Atoi error: %w", err)
	}

	f := os.NewFile(uintptr(fd), "")
	ln, err := net.FileListener(f)
	if err != nil {
		return nil, fmt.Errorf("net.FileListener error: %w", err)
	}

	return &listener{
		Listener: ln,
		File:     f,
	}, nil
}

func (s *Server) signalHandler(ctx pcontext.Context) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR2)
	logger := ctx.Logger()

	for {
		sig := <-ch
		logger.Notice(fmt.Sprintf("receive signal %s", sig.String()))

		switch sig {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP:
			s.shutdown(ctx, ch)
		case syscall.SIGUSR2:
			err := s.startNewProcess(ctx)
			if err != nil {
				logger.Error("server.startNewProcess error", golog.ErrorField(err))
			} else {
				s.shutdown(ctx, ch)
			}
		default:
			continue
		}
	}
}

func (s *Server) shutdown(ctx pcontext.Context, ch chan os.Signal) {
	ctx.Logger().Notice("graceful shutdown")

	signal.Stop(ch)

	err := s.Shutdown(ctx)
	if err != nil {
		ctx.Logger().Error("server.Shutdown error", golog.ErrorField(err))
	} else {
		close(s.shutdownCh)
	}
}

func (s *Server) startNewProcess(ctx pcontext.Context) error {
	files := []*os.File{os.Stdin, os.Stdout, os.Stderr, s.ln.File}
	p, err := os.StartProcess(os.Args[0], os.Args, &os.ProcAttr{
		Env:   append(os.Environ(), fmt.Sprintf("%s=%d", environKeyListenerFd, len(files)-1)),
		Files: files,
	})
	if err != nil {
		return fmt.Errorf("os.StartProcess error: %w", err)
	}

	ctx.Logger().Notice(fmt.Sprintf("new process pid is %d", p.Pid))

	return nil
}

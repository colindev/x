package cli

import (
	"fmt"
	"log"
	"os"
)

// Handler handle cli command
type Handler interface {
	Usage()
	Exec([]string)
}

// Service serve cli commands
type Service interface {
	Register(string, Handler) Service
	Handler
}

type service struct {
	m map[string]Handler
}

func (s *service) Register(name string, handler Handler) Service {

	if _, exists := s.m[name]; exists {
		log.Fatalf("[%s] already exists", name)
	}

	s.m[name] = handler

	return s
}

func (s *service) Usage() {
	fmt.Println("Usage: COMMAND [SUBCOMMAND] [OPTIONS]")
	fmt.Println("SUBCOMMANDS:")
	for name, fn := range s.m {
		fmt.Printf("\t%s\n", name)
		fn.Usage()
	}
}

func (s *service) Exec(args []string) {

	if len(args) == 0 {
		s.Usage()
		os.Exit(2)
	}

}

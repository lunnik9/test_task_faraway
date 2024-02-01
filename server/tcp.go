package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"test_task_faraway/repository"
	"test_task_faraway/services"
)

type Server struct {
	challenger    services.Challenger
	cache         repository.Cache
	log           *log.Logger
	statistic     repository.StatModule
	wordsOfWisdom *repository.WordsOfWisdom
}

type WordOfWisdomResponse struct {
	WordOfWisdom string `json:"word_of_wisdom"`
	Error        string `json:"error"`
}

func NewServer(
	challenger services.Challenger,
	cache repository.Cache,
	statistic repository.StatModule,
	wordsOfWisdom *repository.WordsOfWisdom,
	log *log.Logger,
) *Server {
	return &Server{
		challenger:    challenger,
		cache:         cache,
		log:           log,
		statistic:     statistic,
		wordsOfWisdom: wordsOfWisdom,
	}
}

func (s *Server) Listen(ctx context.Context, port string) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	go func(ctx context.Context) {
		for {
			var conn net.Conn

			select {
			case <-ctx.Done():
				return
			default:
				conn, err = ln.Accept()
				if err != nil {
					s.log.Printf("Error accepting connection: %v", err)
					continue
				}
			}

			go func() {
				var resp WordOfWisdomResponse

				ok, err := s.DDoSMiddleware(ctx, conn)
				if err != nil {
					s.log.Printf("Error accepting connection: %v", err)
					conn.Close()
					return
				}

				if !ok {
					defer conn.Close()

					resp.Error = "challenge is unsolved"

					respJson, err := json.Marshal(resp)
					if err != nil {
						s.log.Printf("Error writing response: %v", err)
					}

					_, err = fmt.Fprintln(conn, respJson)
					if err != nil {
						s.log.Printf("Error writing response: %v", err)
					}

					return
				}

				resp.WordOfWisdom = s.wordsOfWisdom.GetQuote()

				respJson, err := json.Marshal(resp)
				if err != nil {
					s.log.Printf("Error writing response: %v", err)
					conn.Close()
				}

				_, err = conn.Write(respJson)
				if err != nil {
					s.log.Printf("Error writing response: %v", err)
					conn.Close()
				}
			}()

		}
	}(ctx)

	<-ctx.Done()

	ln.Close()

	return nil
}

package client

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	connDeadline = 10
)

type Client struct {
	host                     string
	personalDataShareAllowed bool
	clientID                 int64
	challenge                string

	log *log.Logger

	// can be extracted once
	os     string
	arch   string
	numCpu int
}

func NewClient(log *log.Logger, host string) *Client {
	return &Client{log: log, host: host}
}

func (c *Client) Run(ctx context.Context) {
	var (
		reader = bufio.NewReader(os.Stdin)
	)

	fmt.Println("Do you allow sharing personal data with server? Insert 'y' for yes and whatever else for no")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "y" {
		c.personalDataShareAllowed = true

		c.os = runtime.GOOS
		c.arch = runtime.GOARCH
		c.numCpu = runtime.NumCPU()
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("Press any key to connect to the server or 'q' to quit:")

			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "q" {
				break
			}

			quote, err := c.getQuote(ctx)
			if err != nil {
				c.log.Printf("err getting quote: %s", err)
			}

			fmt.Println(quote)
		}
	}
}

func (c *Client) getQuote(ctx context.Context) (string, error) {
	var (
		challengeResponse    ChallengeDataResponse
		wordOfWisdomResponse WordOfWisdomResponse

		accessReq = AccessRequest{
			UserID:    c.clientID,
			Challenge: c.challenge,
			OS:        c.os,
			Arch:      c.arch,
			NumCPU:    c.numCpu,
		}
	)

	conn, err := net.Dial("tcp", c.host)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(connDeadline * time.Second))
	if err != nil {
		return "", err
	}

	accessReqJson, err := json.Marshal(accessReq)
	if err != nil {
		return "", err
	}

	_, err = conn.Write(accessReqJson)
	if err != nil {
		return "", err
	}

	challengeResponseDecoder := json.NewDecoder(conn)

	err = challengeResponseDecoder.Decode(&challengeResponse)
	if err != nil {
		return "", err
	}

	if challengeResponse.Error != "" {
		return "", errors.New(challengeResponse.Error)
	}

	if c.clientID == 0 {
		c.clientID = challengeResponse.UserID
	}

	nonce, ok := findNonce(ctx, challengeResponse.Challenge, challengeResponse.Difficulty)
	if !ok {
		return "", errors.New("cannot calculate nonce")
	}

	challengeRequestJson, err := json.Marshal(ChallengeDataRequest{
		Challenge: challengeResponse.Challenge,
		Answer:    nonce,
	})

	_, err = conn.Write(challengeRequestJson)
	if err != nil {
		return "", err
	}

	wordOfWisdomResponseDecoder := json.NewDecoder(conn)

	err = wordOfWisdomResponseDecoder.Decode(&wordOfWisdomResponse)
	if err != nil {
		return "", err
	}

	if wordOfWisdomResponse.Error != "" {
		return "", fmt.Errorf("error from server: %w", errors.New(wordOfWisdomResponse.Error))
	}

	return wordOfWisdomResponse.WordOfWisdom, nil
}

func findNonce(ctx context.Context, challenge string, difficulty int32) (int64, bool) {
	var nonce int64
	targetPrefix := strings.Repeat("0", int(difficulty))

	for {
		select {
		case <-ctx.Done():
			return 0, false
		default:
			data := fmt.Sprintf("%s:%d", challenge, nonce)
			hash := sha256.Sum256([]byte(data))
			hexHash := hex.EncodeToString(hash[:])

			if strings.HasPrefix(hexHash, targetPrefix) {
				return nonce, true
			}

			nonce++
		}
	}
}

type AccessRequest struct {
	UserID              int64  `json:"user_id"`
	Challenge           string `json:"challenge"`
	OS                  string `json:"os"`
	Arch                string `json:"arch"`
	NumCPU              int    `json:"num_cpu"`
	PersonalDataAllowed bool   `json:"personal_data_allowed"`
}

type ChallengeDataResponse struct {
	Challenge  string `json:"challenge"`
	Difficulty int32  `json:"difficulty"`
	UserID     int64  `json:"user_id"`
	Error      string `json:"error"`
}

type ChallengeDataRequest struct {
	Challenge string `json:"challenge"`
	Answer    int64  `json:"answer"`
}

type WordOfWisdomResponse struct {
	WordOfWisdom string `json:"word_of_wisdom"`
	Error        string `json:"error"`
}

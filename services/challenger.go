package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	// difficulty defines max length of the number client supposed to find
	defaultDifficulty = 3
	MaxDifficulty     = 7

	// 1 stands for the version, second value is difficulty,
	// third value is for time stamp to ensure the token is time-sensitive,
	// fourth value is for the challenge itself
	// fifth value is user id
	// sixth value stands for test token
	challengeLayout = "1:%d:%d:%s:%d:%s"

	// we make these random for better obfuscation
	noTestToken = "Dij93Nm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type UserData struct {
	OS                  string
	Arch                string
	NumCPU              int
	UserID              int64
	PersonalDataAllowed bool
}

type Challenger interface {
	GenerateChallenge(data UserData) (string, int32, error)
	ValidateChallenge(challenge string, answer int64) (bool, error)
}

type ChallengerStruct struct {
	loader Loader
	logger *log.Logger
}

func NewChallengerStruct(loader Loader, logger *log.Logger) *ChallengerStruct {
	return &ChallengerStruct{
		loader: loader,
		logger: logger,
	}
}

func (c *ChallengerStruct) GenerateChallenge(data UserData) (string, int32, error) {
	var (
		difficulty, testToken = c.getDifficulty(data)
		randomBytes           = make([]byte, 16)
	)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", 0, fmt.Errorf("failed to generate random data: %w", err)
	}

	challenge := fmt.Sprintf(
		challengeLayout,
		difficulty,
		time.Now().UnixNano(),
		hex.EncodeToString(randomBytes),
		data.UserID,
		testToken,
	)

	return challenge, difficulty, nil
}

func (c *ChallengerStruct) ValidateChallenge(challenge string, answer int64) (bool, error) {
	challengeData := strings.Split(challenge, ":")

	if len(challengeData) != 6 {
		return false, errors.New("invalid challenge")
	}

	difficulty, err := strconv.Atoi(challengeData[1])
	if err != nil {
		return false, errors.New("invalid challenge")
	}

	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%d", challenge, answer)))
	hexHash := hex.EncodeToString(hash[:])
	
	if strings.HasPrefix(hexHash, strings.Repeat("0", difficulty)) {
		return true, nil
	}

	return false, nil
}

func (c *ChallengerStruct) getDifficulty(data UserData) (int32, string) {
	usage, err := c.loader.GetCpuUsage()
	if err != nil {
		c.logger.Printf("error retrieving cpu data: %s", err)
		return defaultDifficulty, noTestToken
	}

	// if cpu is actively used, we should not conduct tests
	if usage > 80 {
		return MaxDifficulty, noTestToken
	}

	if data.PersonalDataAllowed {
		// we randomly with 33% chance try out different difficulties based on
		// user's characteristics for further research and optimization.
		// an example below sets different difficulty based on their OS or Arch.
		// ideally we should store these tests in persistent storage, but this should be
		// discussed with analysts because
		// 1. we need to keep track on these tests - codenames, amount of tests we need to run, ...
		// 2. storage that suits us, can be very hard for data analyst who'll work with this data
		// so we'll keep this in memory until we hire an analyst :)
		test := rand.Intn(3)

		if test == 1 {
			switch data.OS {
			case "linux":
				return 2, "32mso2MC"
			case "windows":
				return 4, "ioc9dZl2"
			case "darwin":
				return 5, "Loapl020"
			}
		}

		if test == 2 {
			switch data.Arch {
			case "amd64":
				return 1, "923icmAl"
			case "arm64":
				return 2, "Ompa29m"
			case "386":
				return 4, "mlpoGh6"
			}
		}
	}
	return defaultDifficulty, noTestToken
}

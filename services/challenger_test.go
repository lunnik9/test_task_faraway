package services_test

import (
	"github.com/golang/mock/gomock"
	"log"
	"os"
	"strings"
	"test_task_faraway/mocks"
	"test_task_faraway/services"
	"testing"
)

func TestGenerateChallenge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoader := mocks.NewMockLoader(ctrl)
	mockLoader.EXPECT().GetCpuUsage().Return(float64(50), nil).AnyTimes()

	logger := log.New(os.Stdout, "testLogger: ", log.Lshortfile)
	challenger := services.NewChallengerStruct(mockLoader, logger)

	userData := services.UserData{
		OS:                  "linux",
		Arch:                "amd64",
		NumCPU:              4,
		UserID:              12345,
		PersonalDataAllowed: true,
	}

	challenge, difficulty, err := challenger.GenerateChallenge(userData)
	if err != nil {
		t.Fatalf("GenerateChallenge returned an error: %v", err)
	}

	if !strings.HasPrefix(challenge, "1:") || difficulty < 1 || difficulty > services.MaxDifficulty {
		t.Errorf("Generated challenge has incorrect format or difficulty")
	}
}

func TestValidateChallenge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoader := mocks.NewMockLoader(ctrl)
	mockLoader.EXPECT().GetCpuUsage().Return(float64(50), nil).AnyTimes()

	logger := log.New(os.Stdout, "testLogger: ", log.Lshortfile)
	challenger := services.NewChallengerStruct(mockLoader, logger)

	validChallenge := "1:3:1706731837614662000:76b04aaa7f9ee725cd780574be6d4e81:15:Dij93Nm"
	isValid, err := challenger.ValidateChallenge(validChallenge, 8821)
	if err != nil || !isValid {
		t.Errorf("ValidateChallenge failed to validate a valid challenge")
	}

	invalidChallenge := "1:2:123456789:abcdef:12345:Dij93Nm"
	isValid, err = challenger.ValidateChallenge(invalidChallenge, 456)
	if err == nil && isValid {
		t.Errorf("ValidateChallenge incorrectly validated an invalid challenge")
	}
}

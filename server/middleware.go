package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"test_task_faraway/repository"
	"test_task_faraway/services"
	"time"
)

const (
	middlewareTimeout = 20
	challengeTtl      = 60 * 60 * 24
	susUserTtl        = 60
)

type AccessRequest struct {
	UserID              int64  `json:"user_id"`
	Challenge           string `json:"challenge"`
	OS                  string `json:"os"`
	Arch                string `json:"arch"`
	NumCPU              int    `json:"num_cpu"`
	PersonalDataAllowed bool   `json:"personal_data_allowed"`
}

type ChallengeDataRequest struct {
	Challenge string `json:"challenge"`
	Answer    int64  `json:"answer"`
}

type ChallengeDataResponse struct {
	Challenge  string `json:"challenge"`
	Difficulty int32  `json:"difficulty"`
	UserID     int64  `json:"user_id"`
	Error      string `json:"error"`
}

func (s *Server) DDoSMiddleware(ctx context.Context, conn net.Conn) (bool, error) {
	var (
		accessRequest    AccessRequest
		challengeRequest ChallengeDataRequest
		resp             ChallengeDataResponse
	)

	middlewareCtx, cancelFunc := context.WithTimeout(ctx, middlewareTimeout*time.Second)
	defer cancelFunc()

	accessRequestDecoder := json.NewDecoder(conn)
	if err := accessRequestDecoder.Decode(&accessRequest); err != nil {
		s.log.Printf("err decoding access request: %s", err)
		return false, err
	}

	if accessRequest.UserID == 0 {
		userID, err := s.cache.GetUserID(middlewareCtx)
		if err != nil {
			return false, err
		}

		resp.UserID = userID
	} else {
		ok, err := s.cache.PresentSusUser(ctx, accessRequest.UserID)
		if err != nil {
			s.log.Printf("err retrieving sus user: %s", err)
		}

		if ok {
			resp.Error = "you marked as a suspicious user, try again later"

			respJson, err := json.Marshal(resp)
			if err != nil {
				s.log.Printf("marshalling err: %s", err)
				return false, err
			}

			_, err = conn.Write(respJson)
			if err != nil {
				s.log.Printf("err writing to conn: %s", err)
				return false, err
			}

			return false, errors.New("got a suspicious user")
		}

		resp.UserID = accessRequest.UserID
	}

	challenge, err := s.cache.GetChallenge(middlewareCtx, resp.UserID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		s.log.Printf("err retrieving challenge: %s", err)
		return false, err
	}

	if challenge == accessRequest.Challenge && challenge != "" {
		return true, nil
	}

	resp.Challenge, resp.Difficulty, err = s.challenger.GenerateChallenge(services.UserData{
		OS:                  accessRequest.OS,
		Arch:                accessRequest.Arch,
		NumCPU:              accessRequest.NumCPU,
		UserID:              accessRequest.UserID,
		PersonalDataAllowed: accessRequest.PersonalDataAllowed,
	})
	if err != nil {
		s.log.Printf("err generating challenge: %s", err)
		return false, err
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		s.log.Printf("marshalling err: %s", err)
		return false, err
	}

	_, err = conn.Write(respJson)
	if err != nil {
		s.log.Printf("err writing to conn: %s", err)
		return false, err
	}

	challengeSendTime := time.Now()

	deadline, _ := middlewareCtx.Deadline()

	err = conn.SetDeadline(deadline)
	if err != nil {
		s.log.Printf("err setting context deadline: %s", err)
		return false, err
	}

	challengeRequestDecoder := json.NewDecoder(conn)

	err = challengeRequestDecoder.Decode(&challengeRequest)
	if err != nil {
		s.log.Printf("unmarshalling error: %s", err)
		return false, err
	}

	ok, err := s.challenger.ValidateChallenge(challengeRequest.Challenge, challengeRequest.Answer)
	if err != nil {
		s.log.Printf("err validating challenge: %s", err)
		return false, err
	}

	if ok {
		err = s.cache.InsertChallenge(
			ctx,
			fmt.Sprintf("%s:%d", challengeRequest.Challenge, challengeRequest.Answer),
			resp.UserID,
			challengeTtl,
		)
		if err != nil {
			s.log.Printf("error saving challenge to cache: %s", err)
		}

		go func(
			statisticCtx context.Context,
			accessRequest AccessRequest,
			challenge string,
			answer, userID, timeTaken int64,
		) {
			err = s.statistic.SaveUserStat(statisticCtx, repository.UserStat{
				OS:                  accessRequest.OS,
				Arch:                accessRequest.Arch,
				NumCPU:              accessRequest.NumCPU,
				UserID:              resp.UserID,
				PersonalDataAllowed: accessRequest.PersonalDataAllowed,
				TimeSpent:           timeTaken,
				Challenge:           challenge,
				Answer:              answer,
			})
			if err != nil {
				s.log.Printf("error saving statistics: %s", err)
			}
		}(ctx, accessRequest, challenge, challengeRequest.Answer, resp.UserID, int64(time.Now().Sub(challengeSendTime)))
	} else {

		// if challenge is unsolved, we mark this user as suspicious and restrict theirs access for a short time
		err = s.cache.InsertSusUser(ctx, resp.UserID, susUserTtl)
		if err != nil {
			s.log.Printf("error inserting sus user: %s", err)
		}
	}

	return ok, nil
}

package config

import (
	"Statistic/internal/storage"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	Adr         = "127.0.0.1"
	ListenPort  = "8383"
	SelfPort    = "8282"
	PageMMS     = "/mms"
	PageSMS     = "/sms"
	PageSupport = "/support"
	PageCrush   = "/accendent"
)

type Service map[int]*any

func (s *Service) ConnectionToHost(page string) (ar []any, err error) {
	resp, err := http.Get("http://" + Adr + ":" + storage.ListenPort + page)
	if err != nil {
		log.Println("resp err")
	}
	if resp == nil {
		log.Println("resp nil")
		return nil, err

	} else {
		defer resp.Body.Close()
		if resp.Status != "200 OK" {
			log.Println("break!")
			return nil, err
		} else {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("boddy err")
			}
			// если убрать до return функции  MMS, VoiceCall, CrushReport - "отстёгиваются"
			err = json.Unmarshal(body, &ar)
			if err != nil {
				log.Println("decoder err")
			}

			return ar, nil
		}
	}
}

func NewSomeStorageMMS() *Service {
	var service Service = make(map[int]*any)
	return &service
}

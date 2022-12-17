package config

//package main

import (
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

//var fu HiddenFunctions
/*
type HiddenFunctions interface {
	ConnectionToHost(page string) (mmsked []*model.MMSData, err error)
	ForCreatePlaceMMS(mmsked []*model.MMSData) (err error)
	ReadOnServiceSMS(fullContent []byte) (err error)
	ForPrintAll()
}
*/

//type DataService map[int]*any

//var ddb *DataService

func (s *Service) ConnectionToHost(page string) (ar []any, err error) { //ar []any, err error) {
	resp, err := http.Get("http://" + Adr + ":" + ListenPort + page)
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
			//todo если убрать код ниже до return функции  MMS, VoiceCall, CrushReport - "отстёгиваются"
			err = json.Unmarshal(body, &ar)
			if err != nil {
				log.Println("decoder err")
			}

			return ar, nil
		}
	}
}

/*
	func (s *ServiceMMS) ForCreatePlaceMMS(mmsked []*model.MMSData) (err error) {
		id := 1
		for i := 0; i < len(mmsked); i++ {
			for {
				_, ok := (*s)[id]
				if ok {
					id++
					continue
				}
				break
			}
			(*s)[id] = mmsked[i]
		}
		return nil

}

func (s *ServiceMMS) ForPrintAll() (mms *model.MMSData) {

		for _, mms = range *s {
			//	log.Println(mms)
		}
		return mms
	}
*/
func NewSomeStorageMMS() *Service {
	var service Service = make(map[int]*any)
	return &service
}

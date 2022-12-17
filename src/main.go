package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

type ServiceMMS map[int]*MMSData

func NewStorage() *ServiceMMS {
	var MMS ServiceMMS = make(map[int]*MMSData)
	return &MMS
}

const (
	adr  = "127.0.0.1"
	port = "8383"
	page = "/mms"
)

func main() {
	MMS1 := NewStorage()

	mmsked, err := MMS1.ConnectionToHost()
	if err != nil {
		fmt.Println("err connectionToHost")
	}

	err = MMS1.ForCreatePlace(mmsked)
	if err != nil {
		fmt.Println("ForCreatePlace err")
	}

	MMS1.ForPrintAll()
}

func (s *ServiceMMS) ForCreatePlace(mmsked []*MMSData) (err error) {
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

func (s *ServiceMMS) ForPrintAll() {

	for _, mms := range *s {
		fmt.Println(mms)
	}
	return
}

func (s *ServiceMMS) ConnectionToHost() (mmsked []*MMSData, err error) {

	resp, err := http.Get("http://" + adr + ":" + port + page)
	if err != nil {
		fmt.Println("resp err")
	}
	if resp == nil {
		fmt.Println("resp nil")
		return nil, err

	} else {

		defer resp.Body.Close()
		if resp.Status != "200 OK" {
			fmt.Println("break!")
			return nil, err
		} else {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("boddy err")
			}
			err = json.Unmarshal(body, &mmsked)
			//	if err := json.NewDecoder(resp.Body).Decode(&mms);
			if err != nil {
				fmt.Println("decoder err")
			}

			//fmt.Println(string(body))

			return mmsked, nil
		}
	}
}

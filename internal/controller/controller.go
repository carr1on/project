package controller

import (
	"Statistic/internal/config"
	"Statistic/internal/controller/tools"
	"Statistic/internal/model"
	"Statistic/internal/storage"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Server struct {
	Router *chi.Mux
}

var Resp model.ResultSetT

func CreateNewServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()

	return s
}

func (s *Server) MountHandlers() {
	r := s.Router
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.CleanPath)

	r.Route("/", func(r chi.Router) {
		r.Get("/", DefaultRoute)

		// sub-route 'Statistic'
		r.Route("/stat", func(r chi.Router) {
			r.Get(config.PageSMS, CreateReportSMS)
			r.Get(config.PageMMS, CreateReportMMS)
			r.Get("/email", CreateReportEmail)
			r.Get("/billing", CreateReportBilling)
			r.Get(config.PageSupport, CreateReportSupport)
			r.Get("/voice_call", CreateReportVoiceCall)
			r.Get("/history_of_crush", CreateReportHistoryOfCrush)

		})
	})
}

func DefaultRoute(w http.ResponseWriter, _ *http.Request) {

	CreateReportSMS(w, nil)
	CreateReportEmail(w, nil)
	CreateReportVoiceCall(w, nil)
	CreateReportBilling(w, nil)
	CreateReportMMS(w, nil)
	CreateReportSupport(w, nil)
	CreateReportHistoryOfCrush(w, nil)

}

func ServiceSMS() (SMSbyCountry, SMSbyProvider []model.SMSData, err error) {
	smsdata := "sms.data"

	_, record, err := tools.OpenFile(smsdata)
	if err != nil {
		fmt.Println(err)
	}

	contentText, err := sortedSMSFirstStep(record)
	if err != nil {
		log.Println("err on First Step Sort SMS")
	}

	SMSbyCountry, err = sortedSMSbyCountry(contentText)
	if err != nil {
		log.Println("err on Country Step Sort SMS")
	}
	SMSbyProvider, err = sortedSMSbyProvider(contentText)
	if err != nil {
		log.Println("err on Provider Step Sort SMS")
	}

	return SMSbyCountry, SMSbyProvider, nil
}

func sortedSMSFirstStep(fullContent [][]string) (ContentText string, err error) {
	for i := range fullContent {
		smsslice := fullContent[i]

		if len(smsslice) == 4 {
			if smsslice[3] == "Topolo" || smsslice[3] == "Rond" || smsslice[3] == "Kildy" {
				forWriteString := smsslice[0] + " " + smsslice[1] + " " + smsslice[2] + " " + smsslice[3] + "\n"
				ContentText += forWriteString
			} else {
				break
			}
		} else {
			smsslice = nil
		}
	}
	return ContentText, nil
}

func sortedSMSbyCountry(ContentText string) (allSMS []model.SMSData, err error) {
	var (
		temp    string
		country []storage.ISO3886
		oneSMS  model.SMSData
	)

	country = storage.Countres
	spltn := strings.SplitN(ContentText, "\n", -1)

	tools.SortByAlgorithmABC(spltn)

	for i := range spltn {
		spltn1 := strings.SplitN(spltn[i], " ", -1)

		smsslice := spltn1
		if len(smsslice) == 4 {

			for l, _ := range country {
				comparison := strings.Compare(country[l].Alpa2Code, smsslice[0])
				if comparison == 0 {
					smsslice[0] = "(" + smsslice[0] + ") " + country[l].CountryName

					forWriteString := smsslice[0] + " " + smsslice[1] + " " + smsslice[2] + " " + smsslice[3] + "\n"
					temp += forWriteString

					oneSMS.Country = smsslice[0]
					oneSMS.Bandwidth = smsslice[1]
					oneSMS.ResponseTime = smsslice[2]
					oneSMS.Provider = smsslice[3]
					allSMS = append(allSMS, oneSMS)
				}
			}
		} else {
			smsslice = nil
		}
	}

	return allSMS, nil
}

func sortedSMSbyProvider(ContentText string) (allSMS []model.SMSData, err error) {
	var (
		temp   string
		oneSMS model.SMSData
	)

	ContentText, _ = ColumnsForServiceSMS(ContentText)

	spltn := strings.SplitN(ContentText, "\n", -1)

	tools.SortByAlgorithmABC(spltn)

	for i := range spltn {
		spltn1 := strings.SplitN(spltn[i], " ", -1)

		smsslice := spltn1
		if len(smsslice) == 4 {

			forWriteString := smsslice[0] + " " + smsslice[1] + " " + smsslice[2] + " " + smsslice[3] + "\n"
			temp += forWriteString

			oneSMS.Country = smsslice[1]
			oneSMS.Bandwidth = smsslice[2]
			oneSMS.ResponseTime = smsslice[3]
			oneSMS.Provider = smsslice[0]
			allSMS = append(allSMS, oneSMS)

		} else {
			smsslice = nil
		}
	}

	return allSMS, nil
}

// ColumnsForServiceSMS используется для смены очерёдности колонок. Упрощает сортировку.
func ColumnsForServiceSMS(ContentTextIn string) (ContentTextOut string, err error) {
	spltn := strings.SplitN(ContentTextIn, "\n", -1)
	for i := range spltn {
		spltn1 := strings.SplitN(spltn[i], " ", -1)

		smsslice := spltn1
		if len(smsslice) == 4 {

			forWriteString := tools.Swapper(smsslice)
			ContentTextOut += forWriteString

		} else {
			smsslice = nil
		}
	}
	return ContentTextOut, nil
}

func CreateReportSMS(w http.ResponseWriter, _ *http.Request) {

	SMSbyCountry, SMSbyProvider, err := ServiceSMS()

	dataCountry, _ := json.Marshal(SMSbyCountry)
	if _, err = w.Write(dataCountry); err != nil {
		log.Printf("err write by Country")
	}

	json.Unmarshal(dataCountry, &Resp.SMS)
	dataPrrovider, _ := json.Marshal(SMSbyProvider)
	if _, err = w.Write(dataPrrovider); err != nil {
		log.Printf("err write by Country")
	}

	json.Unmarshal(dataCountry, &Resp.SMS)
	w.WriteHeader(http.StatusOK)
}

func sortedEmail(fullContent [][]string) (forWrite, jsonDataLow, jsonDataHigh []byte, err error) {
	var (
		eData                          model.EmailData
		eDataStore, eDataStoreForWrite []model.EmailData
	)

	County := storage.CountyOnService
	Provider := storage.EmailProvider

	for k := range fullContent {
		emailslice := fullContent[k]

		if len(emailslice) == 3 {

			for i := 0; i < len(County); i++ {
				for j := 0; j < len(Provider); j++ {
					if emailslice[0] == County[i] {
						if emailslice[1] == Provider[j] {

							eData.Country = emailslice[0]
							eData.Provider = emailslice[1]
							eData.DeliveryTime, err = strconv.Atoi(emailslice[2])
							if err != nil {
								log.Println("AtoI err")
							}
							eDataStore = append(eDataStore, eData)

						}
					}
				}
			}
		} else {
			forWrite, err = json.Marshal(eDataStore)
			if err != nil {
				log.Printf("Marshal err in Sorted Email")
			}
			eDataStoreForWrite = nil

			sort.Slice(eDataStore, func(i, j int) bool { return eDataStore[i].DeliveryTime < eDataStore[j].DeliveryTime })
			eDataStoreForWrite = append(eDataStoreForWrite, eDataStore[0], eDataStore[1], eDataStore[2])

			jsonDataLow, err = json.Marshal(eDataStoreForWrite)
			if err != nil {
				log.Printf("Marshal err in Sorted Email")
			}
			eDataStoreForWrite = nil

			sort.Slice(eDataStore, func(i, j int) bool { return eDataStore[i].DeliveryTime > eDataStore[j].DeliveryTime })
			eDataStoreForWrite = append(eDataStoreForWrite, eDataStore[0], eDataStore[1], eDataStore[2])
			jsonDataHigh, err = json.Marshal(eDataStoreForWrite)
			if err != nil {
				log.Printf("Marshal err in Sorted Email")
			}
		}
	}
	return forWrite, jsonDataLow, jsonDataHigh, nil

}

func ServiceEmail() (eData, eDataLow, eDataHigh []model.EmailData, err error) {
	emaildata := "email.data"

	_, record, err := tools.OpenFile(emaildata)
	if err != nil {
		fmt.Println(err)
	}

	forWriteAll, forWriteLow, forWriteHigh, err := sortedEmail(record)
	if err != nil {
		log.Println("err")
	}

	err = json.Unmarshal(forWriteAll, &eData)
	if err != nil {
		log.Println("json Unmarshal err")
		return
	}

	err = json.Unmarshal(forWriteLow, &eDataLow)
	if err != nil {
		log.Println("json Unmarshal err")
		return
	}

	err = json.Unmarshal(forWriteHigh, &eDataHigh)
	if err != nil {
		log.Println("json Unmarshal err")
		return
	}

	return eData, eDataLow, eDataHigh, nil
}

func CreateReportEmail(w http.ResponseWriter, _ *http.Request) {

	eData, eDataLow, eDataHigh, err := ServiceEmail()

	Resp.Email["LOW PING ON EMAIL SERVICE"] = eDataLow
	Resp.Email["HiGH PING ON EMAIL SERVICE"] = eDataHigh
	Resp.Email["ALL RESPONSE ON EMAIL SERVICE"] = eData

	dataCountry, _ := json.Marshal(Resp.Email)
	if _, err = w.Write(dataCountry); err != nil {
		log.Printf("err write by Email Service")
	}

	w.WriteHeader(http.StatusOK)
}

func sortedVoiceCall(fullContent [][]string) (voiceCall []model.VoiceCallData, err error) {
	var (
		vCall    model.VoiceCallData
		County   = storage.CountyOnService
		Provider = storage.ProviderVoiceCall
	)

	for k := range fullContent {
		VCSlc := fullContent[k]

		if len(VCSlc) == 8 {

			for i := 0; i < len(County); i++ {
				for j := 0; j < len(Provider); j++ {

					if VCSlc[0] == County[i] {
						if VCSlc[3] == Provider[j] {

							vCall.Country = VCSlc[0]
							vCall.ResponseTime = VCSlc[1]
							vCall.Bandwidth = VCSlc[2]
							vCall.Provider = VCSlc[3]

							if temp, err := strconv.ParseFloat(VCSlc[4], 32); err == nil {
								vCall.ConnectionStability = float32(temp)
							}

							vCall.TTFB, err = strconv.Atoi(VCSlc[5])
							vCall.VoicePurity, err = strconv.Atoi(VCSlc[6])
							vCall.MediaOfCallsTime, err = strconv.Atoi(VCSlc[7])
							if err != nil {
								log.Print("err in convert strings to value for Voice Call")
							}

							voiceCall = append(voiceCall, vCall)
						}
					}
				}
			}
		} else {
			return nil, err
		}
	}

	return voiceCall, nil

}

func ServiceVoiceCall() (voiceCall []model.VoiceCallData, err error) {
	voicedata := "voice.data"

	_, record, err := tools.OpenFile(voicedata)
	if err != nil {
		fmt.Println(err)
	}

	voiceCall, err = sortedVoiceCall(record)
	if err != nil {
		log.Println("err")
	}

	return voiceCall, nil
}

func CreateReportVoiceCall(w http.ResponseWriter, _ *http.Request) {

	forWrite, err := ServiceVoiceCall()

	if err != nil {
		log.Printf("err ServiceVoiceCall")
	}

	w.WriteHeader(http.StatusOK)

	jsonData, err := json.Marshal(forWrite)
	if err != nil {
		log.Printf("err in json marshal on Voice Call ")
	}
	json.Unmarshal(jsonData, &Resp.VoiceCall)
	if _, err = w.Write(jsonData); err != nil {
		log.Printf("err write")
	}

}

// billing

func CreateReportBilling(w http.ResponseWriter, _ *http.Request) {
	forWrite, billingData, err := Billing()

	if err != nil {
		log.Printf("err ServiceBilling")
	}
	Resp.Billing = billingData
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(forWrite); err != nil {
		log.Printf("err write")
	}

}

func Billing() (forWrite []byte, billingData model.BiliingData, err error) {

	billingdata := "billing.data"

	file, _, err := tools.OpenFile(billingdata)
	if err != nil {
		fmt.Println(err)
	}
	if _, err = file.Seek(0, 0); err != nil {
		log.Printf("err seek")
	}

	fullContent, err := io.ReadAll(file)

	_, err = strconv.Atoi(string(fullContent))
	if err != nil {
		fmt.Print("err")
	}

	summ, err := strconv.ParseInt(string(fullContent), 2, 10)
	if err != nil {
		fmt.Println("err ParseInt")
	}
	var sliceForSepPrimary []int
	var sliceForSep_Secondary []int
	for i := 0; i < len(fullContent); i++ {

		sliceForSepPrimary = append(sliceForSepPrimary, i+1)
		if fullContent[i] == 49 {
			sliceForSep_Secondary = append(sliceForSep_Secondary, 1)
		} else if fullContent[i] == 48 {
			sliceForSep_Secondary = append(sliceForSep_Secondary, 0)
		}

	}
	var slicestring []bool
	for i := 0; i < len(sliceForSep_Secondary); i++ {

		if sliceForSep_Secondary[i] == 1 {
			slicestring = append(slicestring, true)
		} else {
			slicestring = append(slicestring, false)
		}
	}

	billingData.CreateCustomer = slicestring[0]
	billingData.Purchase = slicestring[1]
	billingData.Payout = slicestring[2]
	billingData.Recurring = slicestring[3]
	billingData.FraudControl = slicestring[4]
	billingData.CheckoutPage = slicestring[5]

	var (
		stringsForPrint []string
		count, amount   int
	)
	if billingData.CreateCustomer == true {
		stringsForPrint = append(stringsForPrint, "CreateCustomer")
		amount++
	} else {
		count++
	}
	if billingData.Purchase == true {
		stringsForPrint = append(stringsForPrint, "Purchase")
		amount++
	} else {
		count++
	}
	if billingData.Payout == true {
		stringsForPrint = append(stringsForPrint, "Payout")
		amount++
	} else {
		count++
	}
	if billingData.Recurring == true {
		stringsForPrint = append(stringsForPrint, "Recurring")
		amount++
	} else {
		count++
	}
	if billingData.FraudControl == true {
		stringsForPrint = append(stringsForPrint, "FraudControl")
		amount++
	} else {
		count++
	}
	if billingData.CheckoutPage == true {
		stringsForPrint = append(stringsForPrint, "CheckoutPage")
		amount++
	} else {
		count++
	}
	var mask string
	mask = string(fullContent) + " = " + strconv.Itoa(int(summ))

	var output string
	if count == 0 {
		output = fmt.Sprintf("%s %s, должны значения True\n", mask, stringsForPrint)
	} else {
		if amount == 0 {
			output = fmt.Sprintf("%s Все значения False", mask)
		} else if amount == 1 {
			output = fmt.Sprintf("%s %s, должно быть True\n Остальное False", mask, stringsForPrint)
		} else if count > 0 {
			output = fmt.Sprintf("%s, %s, должны быть True\n Остальное False", mask, stringsForPrint)
		}
	}

	temp := output + "\n "
	forWrite = []byte(temp)
	return forWrite, billingData, nil
}

// out http requests
var store *config.Service

func init() {
	//TODO
	Resp.Email = make(map[string][]model.EmailData)
	store = config.NewSomeStorageMMS()
}

func CreateReportMMS(w http.ResponseWriter, _ *http.Request) {

	var (
		page      = config.PageMMS
		sortedMMS []model.MMSData
	)

	AllMMS, err := store.ConnectionToHost(page)
	if err != nil {
		log.Println("err connectionToHost")
	}
	jsonData, err := json.Marshal(AllMMS)
	if err != nil {
		log.Printf("[SERVER] can't prepare response: %v\n", err)
		return
	}
	err = json.Unmarshal(jsonData, &sortedMMS)
	if err != nil {
		log.Println("json Unmarshal err")
		return
	}

	if len(sortedMMS) != 0 {
		byCountry, byProvider, err := sortedMMSbyCountry(sortedMMS)
		if err != nil {
			log.Println("err in sort func")
		}

		jsonDataP, err := json.Marshal(byProvider)
		if _, err = w.Write(jsonDataP); err != nil {
			log.Printf("err write by Country")
		}
		json.Unmarshal(jsonDataP, &Resp.MMS)

		jsonDataC, err := json.Marshal(byCountry)
		if _, err = w.Write(jsonDataC); err != nil {
			log.Printf("err write by Country")
		}
		json.Unmarshal(jsonDataC, &Resp.MMS)
		w.WriteHeader(http.StatusOK)

	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func sortedMMSbyCountry(mms []model.MMSData) (mmsCountry, mmProvider []model.MMSData, err error) {
	var (
		temp   string
		mmsked model.MMSData
		mmsTmp []model.MMSData
	)
	var country []storage.ISO3886
	country = storage.Countres

	for i := range mms {
		temp += mms[i].Country + " " + mms[i].Bandwidth + " " + mms[i].ResponseTime + " " + mms[i].Provider + "\n"
	}
	spltn := strings.SplitN(temp, "\n", -1)
	tools.SortByAlgorithmABC(spltn)

	for i := range spltn {
		spltn1 := strings.SplitN(spltn[i], " ", -1)
		smsslice := spltn1

		if len(smsslice) == 4 {
			for l := range country {
				comparison := strings.Compare(country[l].Alpa2Code, smsslice[0])
				if comparison == 0 {
					smsslice[0] = country[l].CountryName + "(" + smsslice[0] + ") "

					mmsked.Country = smsslice[0]
					mmsked.Provider = smsslice[3]
					mmsked.Bandwidth = smsslice[1]
					mmsked.ResponseTime = smsslice[2]
					mmsTmp = append(mmsTmp, mmsked)
				}
			}

		}
	}
	mms = nil

	sort.Slice(mmsTmp, func(i, j int) bool { return mmsTmp[i].Provider < mmsTmp[j].Provider })
	_, err = json.Marshal(mmsTmp)
	if err != nil {
		log.Printf("Marshal err in Sorted Email")
	}
	mmsCountry = mmsTmp

	sort.Slice(mmsTmp, func(i, j int) bool { return mmsTmp[i].Country < mmsTmp[j].Country })
	_, err = json.Marshal(mmsTmp)
	if err != nil {
		log.Printf("Marshal err in Sorted Email")
	}
	mmProvider = mmsTmp

	return mmsCountry, mmProvider, nil
}

func CreateReportSupport(w http.ResponseWriter, _ *http.Request) {
	page := config.PageSupport

	AllRead, err := store.ConnectionToHost(page)
	if err != nil {
		log.Println("err connectionToHost")
	}

	jsonData, err := json.Marshal(AllRead)
	if err != nil {
		log.Printf("[SERVER] can't prepare response: %v\n", err)
		return
	}

	var support []model.SupportData
	err = json.Unmarshal(jsonData, &support)
	if err != nil {
		log.Println("json Unmarshal err")
		return
	}

	if len(support) != 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		sendResponse(w, http.StatusInternalServerError, support, nil)
	}

	TicketsInt := ActiveTickets(support)
	if TicketsInt != 0 {
		forWrite := SupportTimeTickets(TicketsInt)

		dataSupport, err := json.Marshal(forWrite)
		if _, err = w.Write(dataSupport); err != nil {
			log.Printf("err write by Country")
		}

		json.Unmarshal(dataSupport, Resp.Support)

		if _, err = w.Write([]byte(dataSupport)); err != nil {
			log.Printf("err write")
			return
		}
	}
}

func ActiveTickets(all []model.SupportData) (ticketsValue int) {

	for _, j := range all {

		i := j.ActiveTickets
		ticketsValue += i
	}
	return ticketsValue
}

func SupportTimeTickets(allInt int) (timeForAccept string) {

	var (
		baseTicketInHour float32 = 18
		timePerBase              = 60
		percentHundred   float32 = 100
		percentRequired  float32
		degreeStatus     string
	)

	someInteger := float32(allInt) - baseTicketInHour
	//someInteger := 32 - baseTicketInHour //изначально тестировалось на тикетах = 32

	if int(someInteger) < 9 {
		degreeStatus = fmt.Sprintf("level.1")
	} else if 9 <= int(someInteger) && int(someInteger) < 16 {
		degreeStatus = fmt.Sprintf("level.2")
	} else if 16 <= int(someInteger) {
		degreeStatus = fmt.Sprintf("level.3")
	}

	percentOne := baseTicketInHour / percentHundred

	for someInteger < baseTicketInHour {
		baseTicketInHour = baseTicketInHour - percentOne
		percentRequired++

		if someInteger == baseTicketInHour {
			break
		} else if someInteger > baseTicketInHour {
			break
		}
	}

	percentRequired = percentHundred - percentRequired
	timeToAccept := float32(timePerBase) * (percentRequired / 100)
	if timeToAccept == 60 {
		timeForAccept = fmt.Sprintf("%s время ожидания до ответа на новый запрос в Support более %d минут", degreeStatus, int(timeToAccept))

		return timeForAccept
	}

	timeForAccept = fmt.Sprintf("%s время ожидания до ответа на новый запрос в Support: %d минут", degreeStatus, int(timeToAccept))
	return timeForAccept
}

func CreateReportHistoryOfCrush(w http.ResponseWriter, _ *http.Request) {
	page := config.PageCrush

	body, err := store.ConnectionToHost(page)
	if err != nil {
		log.Println("err connectionToHost")
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Printf("[SERVER] can't prepare response: %v\n", err)
		return
	}
	var historyOfCrush []*model.IncidentData
	err = json.Unmarshal(jsonData, &historyOfCrush)
	if err != nil {
		log.Println("json Unmarshal err")
		return
	}

	if len(historyOfCrush) != 0 {
		sort.Slice(historyOfCrush, func(i, j int) bool { return historyOfCrush[i].Status < historyOfCrush[j].Status })

		w.WriteHeader(http.StatusOK)

		jsonData, err := json.Marshal(historyOfCrush)
		if err != nil {
			log.Print("err in history of crush")
		}
		json.Unmarshal(jsonData, &Resp.Incidents)
		//sendResponse(w, http.StatusOK, historyOfCrush, nil)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		sendResponse(w, http.StatusInternalServerError, historyOfCrush, nil)
	}

}

type ResultT struct {
	Status bool             `json:"status"`
	Data   model.ResultSetT `json:"data"`
	Error  error            `json:"error,omitempty"`
}

type UserResponse struct {
	Status   int         `json:"status"`
	Response interface{} `json:"response,omitempty"`
	Error    error       `json:"error,omitempty"`
}

func sendResponse(w http.ResponseWriter, status int, forWrite interface{}, err error) {

	response := UserResponse{
		Status:   status,
		Response: forWrite,
		Error:    err,
	}
	jsonData, err := json.Marshal(response)

	if err != nil {
		log.Printf("[SERVER] can't prepare response: %v\n", err)
		return
	}
	if _, err = w.Write(jsonData); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
		return
	}
}

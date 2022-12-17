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

func CreateNewServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	return s
}

func (s *Server) MountHandlers() {
	r := s.Router
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	//	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.CleanPath)

	r.Route("/", func(r chi.Router) {
		r.Get("/", DefaultRoute)

		// sub-route 'users'
		r.Route("/stat", func(r chi.Router) {
			r.Get(config.PageSMS, CreateReportSMS)
			r.Get(config.PageMMS, CreateReportMMS)
			r.Get("/email", CreateReportEmail)
			r.Get("/billing", CreateReportBilling)
			r.Get(config.PageSupport, CreateReportSupport)
			r.Get("/voice_call", CreateReportVoiceCall)
			r.Get("/history_of_crush", CreateReportHistoryOfCrush)

			/*				r.Put("/update", Put)
							r.Post("/make_friend", MakeFriend)
							r.Delete("/delete", DeleteUser)
							r.Route("/this_user", func(r chi.Router) {
							r.Get("/friends", GetFriends)
			*/
		})
	})
}

func DefaultRoute(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("SMS\n\n")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
	CreateReportSMS(w, nil)
	if _, err := w.Write([]byte("-------------------------------------\n----------------------------------\n")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}

	//
	if _, err := w.Write([]byte("Email\n\n")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
	CreateReportEmail(w, nil)
	if _, err := w.Write([]byte("\n---------------------------------------\n--------------------------------\n")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}

	//
	if _, err := w.Write([]byte("VoiceCall\n\n")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
	voiceString := "|Country|Bandwidth|ResponseTime|      Provider     |TTFB|VoicePurity|MediaOfCallsTime\n"
	if _, err := w.Write([]byte(voiceString)); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
	CreateReportVoiceCall(w, nil)

	//
	if _, err := w.Write([]byte("---------------------------------------\n--------------------------------\n")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
	if _, err := w.Write([]byte("Billing\n\n")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
	//
	CreateReportBilling(w, nil)
	if _, err := w.Write([]byte("---------------------------------------\n--------------------------------\n")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
	if _, err := w.Write([]byte("MMS\n\n")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
	CreateReportMMS(w, nil)
	if _, err := w.Write([]byte("\n---------------------------------------\n--------------------------------\n")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
	//
	if _, err := w.Write([]byte("Support\n\n")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
	CreateReportSupport(w, nil)
	if _, err := w.Write([]byte("\n---------------------------------------\n--------------------------------\n")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
	//
	if _, err := w.Write([]byte("HistoryOfCrush\n\n")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
	CreateReportHistoryOfCrush(w, nil)
	if _, err := w.Write([]byte("\n---------------------------------------\n--------------------------------\n")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
	if _, err := w.Write([]byte("server alive")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
}

func ServiceSMS() (forWrite, forWriteByCountry, forWriteByProvider []byte, err error) {
	smsdata := "sms.data"

	_, record, err := tools.OpenFile(smsdata)
	if err != nil {
		fmt.Println(err)
	}

	forWrite, contentText, err := sortedSMS(record)
	if err != nil {
		log.Println("err")
	}

	forWriteByCountry, err = sortedSMSbyCountry(contentText)
	if err != nil {
		log.Println("err")
	}
	forWriteByProvider, err = sortedSMSbyProvider(contentText)
	if err != nil {
		log.Println("err")
	}

	return forWrite, forWriteByCountry, forWriteByProvider, nil
}

func sortedSMS(fullContent [][]string) (forWrite []byte, ContentText string, err error) {
	for i := range fullContent {
		smsslice := fullContent[i]

		if len(smsslice) == 4 {
			if smsslice[3] == "Topolo" || smsslice[3] == "Rond" || smsslice[3] == "Kildy" {

				forWriteString := smsslice[0] + " " + smsslice[1] + " " + smsslice[2] + " " + smsslice[3] + "\n"
				ContentText += forWriteString

				forWrite = []byte(ContentText)
			} else {
				break
			}
		} else {
			smsslice = nil
		}
	}
	return forWrite, ContentText, nil
}

func sortedSMSbyCountry(ContentText string) (forWrite []byte, err error) {
	var (
		temp    string
		country []storage.ISO3886
	)

	country = storage.Countres
	spltn := strings.SplitN(ContentText, "\n", -1)

	tools.SortByAlgorithmNarayana(spltn)

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

					forWrite = []byte(temp)
				}
			}
		} else {
			smsslice = nil
		}
	}

	return forWrite, nil
}

func sortedSMSbyProvider(ContentText string) (forWrite []byte, err error) {
	var temp string

	ContentText, _ = ColumnsForServiceSMS(ContentText)

	spltn := strings.SplitN(ContentText, "\n", -1)

	tools.SortByAlgorithmNarayana(spltn)

	for i := range spltn {
		spltn1 := strings.SplitN(spltn[i], " ", -1)

		smsslice := spltn1
		if len(smsslice) == 4 {

			forWriteString := "Provider: " + smsslice[0] + " " + smsslice[1] + " " + smsslice[2] + " " + smsslice[3] + "\n"
			temp += forWriteString

			forWrite = []byte(temp)

		} else {
			smsslice = nil
		}
	}

	return forWrite, nil
}

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

	_, forWriteByCountry, forWriteByProvider, err := ServiceSMS()

	if err != nil {
		log.Printf("err ServiceSMS")
	}

	w.WriteHeader(http.StatusOK)

	/*
		if _, err = w.Write(forWrite); err != nil {
			log.Printf("err write")
		}
	*/
	if _, err = w.Write(forWriteByCountry); err != nil {
		log.Printf("err write by Country")
	}

	if _, err = w.Write([]byte("\n--------\n")); err != nil {
		log.Printf("err write")
	}
	if _, err = w.Write(forWriteByProvider); err != nil {
		log.Printf("err write by Provider")
	}

	/*var smsData []*model.SMSData
	err = json.Unmarshal(forWrite, &smsData)
	if err != nil {
		log.Println("json Unmarshal err")
		return
	}*/
	//	w.WriteHeader(http.StatusOK)
	//	sendResponse(w, http.StatusOK, forWrite, nil) // u, nil)
}

func sortedEmail(fullContent [][]string) ([]byte, []byte, []byte, error) {
	var (
		eData                               model.EmailData
		eDataStore, eDataStoreForWrite      []model.EmailData
		forWrite, jsonDataLow, jsonDataHigh []byte
		err                                 error
	)

	County := []string{"RU", "US", "GB", "FR", "BL", "AT", "BG", "DK", "CA", "ES", "CH", "TR", "PE", "NZ", "MC"}
	Provider := []string{"Gmail", "Yahoo", "Hotmail", "MSN", "Orange", "Comcast", "AOL", "Live", "RediffMail", "GMX", "Protomail", "Yandex", "Mail.ru"}

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

	//}
	//return nil, nil, nil, err
}

func ServiceEmail() (forWriteAll, forWriteLow, forWriteHigh []byte, err error) {
	emaildata := "email.data"

	_, record, err := tools.OpenFile(emaildata)
	if err != nil {
		fmt.Println(err)
	}

	forWriteAll, forWriteLow, forWriteHigh, err = sortedEmail(record)
	if err != nil {
		log.Println("err")
	}

	return forWriteAll, forWriteLow, forWriteHigh, nil
}

func CreateReportEmail(w http.ResponseWriter, _ *http.Request) {

	forWrite, forWriteLow, forWriteHigh, err := ServiceEmail()

	if err != nil {
		log.Printf("err ServiceSMS")
	}

	w.WriteHeader(http.StatusOK)

	if _, err = w.Write([]byte("\n-low Ping----\n")); err != nil {
		log.Printf("err write")
	}

	if _, err = w.Write(forWriteLow); err != nil {
		log.Printf("err write")
	}
	if _, err = w.Write([]byte("\n-high ping----\n")); err != nil {
		log.Printf("err write")
	}

	if _, err = w.Write(forWriteHigh); err != nil {
		log.Printf("err write")
	}

	if _, err = w.Write([]byte("\n-All----\n")); err != nil {
		log.Printf("err write")
	}
	if _, err = w.Write(forWrite); err != nil {
		log.Printf("err write")
	}

}

func sortedVoiceCall(fullContent [][]string) ([]byte, error) {
	var (
		temp     string
		forWrite []byte
		err      error
	)

	County := []string{"RU", "US", "GB", "FR", "BL", "AT", "BG", "DK", "CA", "ES", "CH", "TR", "PE", "NZ", "MC"}
	Provider := []string{"  TransparentCalls   ", "      E-Voice       ", "     JustPhone      ", "JustPhone"}

	for k := range fullContent {
		VCSlc := fullContent[k]

		if len(VCSlc) == 8 {

			for i := 0; i < len(County); i++ {
				for j := 0; j < len(Provider); j++ {

					if VCSlc[3] == "E-Voice" {
						VCSlc[3] = "      E-Voice       "
					} else if VCSlc[3] == "JustPhone" {
						VCSlc[3] = "     JustPhone      "
					} else if VCSlc[3] == "TransparentCalls" {
						VCSlc[3] = "  TransparentCalls   "
					}
					if VCSlc[0] == County[i] {
						if VCSlc[3] == Provider[j] {

							forWriteString := VCSlc[0] + " " + VCSlc[1] + " " + VCSlc[2] + " " + VCSlc[3] + " " + VCSlc[4] + VCSlc[5] + " " + VCSlc[6] + VCSlc[7] + " " + "\n"
							temp += forWriteString

							forWrite = []byte(temp)
						}
					}
				}
			}
		} else {
			return nil, err
		}
	}

	//	slice = nil
	//return nil, err

	return forWrite, nil

}

func ServiceVoiceCall() (forWrite []byte, err error) {
	voicedata := "voice.data"

	_, record, err := tools.OpenFile(voicedata)
	if err != nil {
		fmt.Println(err)
	}

	forWrite, err = sortedVoiceCall(record)
	if err != nil {
		log.Println("err")
	}

	return forWrite, nil
}

func CreateReportVoiceCall(w http.ResponseWriter, _ *http.Request) {

	forWrite, err := ServiceVoiceCall()

	if err != nil {
		log.Printf("err ServiceVoiceCall")
	}

	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(forWrite); err != nil {
		log.Printf("err write")
	}

	//sendResponseSlice(w, http.StatusOK, _, nil) // u, nil)
}

// billing

func CreateReportBilling(w http.ResponseWriter, _ *http.Request) {

	forWrite, err := Billing()

	if err != nil {
		log.Printf("err ServiceBilling")
	}

	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(forWrite); err != nil {
		log.Printf("err write")
	}

	//sendResponseSlice(w, http.StatusOK, _, nil) // u, nil)
}

func Billing() (forWrite []byte, err error) {

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
	var sliceforsep []int
	var sliceforsep1 []int
	for i := 0; i < len(fullContent); i++ {

		sliceforsep = (append(sliceforsep, i+1))
		if fullContent[i] == 49 {
			sliceforsep1 = (append(sliceforsep1, 1))
		} else if fullContent[i] == 48 {
			sliceforsep1 = (append(sliceforsep1, 0))
		}

	}
	var slicestring []bool
	for i := 0; i < len(sliceforsep1); i++ {

		if sliceforsep1[i] == 1 {
			slicestring = append(slicestring, true)
		} else {
			slicestring = append(slicestring, false)
		}
	}

	var d model.BiliingData
	d.CreateCustomer = slicestring[0]
	d.Purchase = slicestring[1]
	d.Payout = slicestring[2]
	d.Recurring = slicestring[3]
	d.FraudControl = slicestring[4]
	d.CheckoutPage = slicestring[5]

	var stringgForPrint []string
	var count, amount int
	if d.CreateCustomer == true {
		stringgForPrint = append(stringgForPrint, "CreateCustomer")
		amount++
	} else {
		count++
	}
	if d.Purchase == true {
		stringgForPrint = append(stringgForPrint, "Purchase")
		amount++
	} else {
		count++
	}
	if d.Payout == true {
		stringgForPrint = append(stringgForPrint, "Payout")
		amount++
	} else {
		count++
	}
	if d.Recurring == true {
		stringgForPrint = append(stringgForPrint, "Recurring")
		amount++
	} else {
		count++
	}
	if d.FraudControl == true {
		stringgForPrint = append(stringgForPrint, "FraudControl")
		amount++
	} else {
		count++
	}
	if d.CheckoutPage == true {
		stringgForPrint = append(stringgForPrint, "CheckoutPage")
		amount++
	} else {
		count++
	}
	var mask string
	mask = string(fullContent) + " = " + strconv.Itoa(int(summ))

	var output string
	if count == 0 {
		output = fmt.Sprintf("%s %s, должны значения True\n", mask, stringgForPrint)
	} else {
		if amount == 0 {
			output = fmt.Sprintf("%s Все значения False", mask)
		} else if amount == 1 {
			output = fmt.Sprintf("%s %s, должно быть True\n Остальное False", mask, stringgForPrint)
		} else if count > 0 {
			output = fmt.Sprintf("%s, %s, должны быть True\n Остальное False", mask, stringgForPrint)
		}
	}
	//	log.Println(output)

	temp := output + "\n "
	forWrite = []byte(temp)
	return forWrite, nil
}

// http request

type ServiceMMS map[int]*model.MMSData

var store *config.Service

func init1() {
	//TODO
	store = config.NewSomeStorageMMS()
}

func CreateReportMMS(w http.ResponseWriter, _ *http.Request) {
	page := config.PageMMS
	init1()
	//var mmsked []model.MMSData
	mmsked, err := store.ConnectionToHost(page)
	if err != nil {
		log.Println("err connectionToHost")
	}

	jsonData, err := json.Marshal(mmsked)
	if err != nil {
		log.Printf("[SERVER] can't prepare response: %v\n", err)
		return
	}

	var (
		byProvider, byCountry []byte
		mms                   []model.MMSData
	)
	err = json.Unmarshal(jsonData, &mms)
	if err != nil {
		log.Println("json Unmarshal err")
		return
	}

	if len(mms) != 0 {
		_, byProvider, byCountry, err = sortedMMSbyCountry(mms)
		if err != nil {
			log.Println("err in sort func")
		}

		//todo если убрать объединение срезов выводится только 1 вариант
		for i := range byProvider {
			byCountry = append(byCountry, byProvider[i])
		}

		//todo если оставить незакоментированными строки ниже:  "json Unmarshal err"
		//todo  выводит вариант 1
		/*err = json.Unmarshal(byCountry, &mms)
		if err != nil {
			log.Println("json Unmarshal err")
			return
		}*/
		//todo вариант 2
		/*err = json.Unmarshal(byProvider, &mms)
		if err != nil {
			log.Println("json Unmarshal err")
			return
		}*/

		sendResponse(w, http.StatusOK, string(byCountry), nil) //todo при выводе тут было "...(w, http.StatusOK, mms).."
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		sendResponse(w, http.StatusInternalServerError, mms, nil)
	}
}

func sortedMMSbyCountry(mms []model.MMSData) ([]byte, []byte, []byte, error) {
	var temp string
	var mmsked model.MMSData
	var mmsTmp []model.MMSData
	var forWrite, jsonCountry, jsonProvider []byte
	var err error

	var country []storage.ISO3886
	country = storage.Countres

	for i := range mms {

		temp += mms[i].Country + " " + mms[i].Bandwidth + " " + mms[i].ResponseTime + " " + mms[i].Provider + "\n"
	}
	spltn := strings.SplitN(temp, "\n", -1)
	tools.SortByAlgorithmNarayana(spltn)

	for i := range spltn {
		spltn1 := strings.SplitN(spltn[i], " ", -1)

		smsslice := spltn1

		if len(smsslice) == 4 {

			for l, _ := range country {
				comparison := strings.Compare(country[l].Alpa2Code, smsslice[0])
				if comparison == 0 {
					smsslice[0] = "(" + smsslice[0] + ") " + country[l].CountryName

					forWriteString := smsslice[0] + " " + smsslice[1] + " " + smsslice[2] + " " + smsslice[3] + "\n"

					mmsked.Country = smsslice[0]
					mmsked.Provider = smsslice[3]
					mmsked.Bandwidth = smsslice[1]
					mmsked.ResponseTime = smsslice[2]
					mmsTmp = append(mmsTmp, mmsked)

					temp += forWriteString

					forWrite = []byte(temp)

				}
			}

		}
	}

	forWrite, err = json.Marshal(mms)
	if err != nil {
		log.Printf("Marshal err in Sorted Email")
	}
	sort.Slice(mmsTmp, func(i, j int) bool { return mmsTmp[i].Provider < mmsTmp[j].Provider })

	jsonProvider, err = json.Marshal(mmsTmp)
	if err != nil {
		log.Printf("Marshal err in Sorted Email")
	}
	sort.Slice(mmsTmp, func(i, j int) bool { return mmsTmp[i].Country < mmsTmp[j].Country })

	jsonCountry, err = json.Marshal(mmsTmp)
	if err != nil {
		log.Printf("Marshal err in Sorted Email")
	}

	return forWrite, jsonCountry, jsonProvider, nil
}

func CreateReportSupport(w http.ResponseWriter, _ *http.Request) {
	page := config.PageSupport
	init1()

	mmsked, err := store.ConnectionToHost(page)
	if err != nil {
		log.Println("err connectionToHost")
	}

	jsonData, err := json.Marshal(mmsked)

	if err != nil {
		log.Printf("[SERVER] can't prepare response: %v\n", err)
		return
	}

	var support []*model.SupportData
	err = json.Unmarshal(jsonData, &support)
	if err != nil {
		log.Println("json Unmarshal err")
		return
	}

	if len(support) != 0 {
		w.WriteHeader(http.StatusOK)
		sendResponse(w, http.StatusOK, support, nil)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		sendResponse(w, http.StatusInternalServerError, support, nil)
	}

	TicketsInt := ActiveTickets(support)
	if TicketsInt != 0 {
		tmp := SupportTimeTickets(TicketsInt)
		if _, err = w.Write([]byte("\n" + tmp + "\n")); err != nil {
			log.Printf("err write")
			return
		}
	}
}

func ActiveTickets(all []*model.SupportData) (mmskedInt int) {

	for _, j := range all {

		i := j.ActiveTickets
		mmskedInt += i
	}
	return mmskedInt
}

func SupportTimeTickets(allInt int) (timeForAccept string) {

	var (
		base               float32 = 18
		timeHourPerTickets         = 60
		percentHundred     float32 = 100
		percentRequired    float32
		degreeStatus       string
	)

	someInteger := float32(allInt) - base
	//someInteger := 32 - base

	if int(someInteger) < 9 {
		degreeStatus = fmt.Sprintf("\nlevel.1\n")
	} else if 9 <= int(someInteger) && int(someInteger) < 16 {
		degreeStatus = fmt.Sprintf("\nlevel.2\n")
	} else if 16 <= int(someInteger) {
		degreeStatus = fmt.Sprintf("\nlevel.3\n")
	}

	percentOne := base / percentHundred

	for someInteger < base {
		base = base - percentOne
		percentRequired++

		if someInteger == base {
			break
		} else if someInteger > base {
			break
		}
	}

	percentRequired = percentHundred - percentRequired
	timeToAccept := float32(timeHourPerTickets) * (percentRequired / 100)
	if timeToAccept == 60 {
		timeForAccept = fmt.Sprintf("%s время ожидания до ответа на новый запрос в Support более %d минут", degreeStatus, int(timeToAccept))
		//fmt.Println(timeForAccept)
		return timeForAccept
	}

	timeForAccept = fmt.Sprintf("%s время ожидания до ответа на новый запрос в Support: %d минут", degreeStatus, int(timeToAccept))
	//fmt.Println(timeForAccept)
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
		w.WriteHeader(http.StatusOK)
		sendResponse(w, http.StatusOK, historyOfCrush, nil)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		sendResponse(w, http.StatusInternalServerError, historyOfCrush, nil)
	}

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

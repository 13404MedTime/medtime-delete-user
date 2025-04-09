package function

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	ucodesdk "github.com/golanguzb70/ucode-sdk"
	"github.com/spf13/cast"
)

const (
	IsHTTP = true // if this is true banchmark test works.
	appId  = "P-JV2nVIRUtgyPO5xRNeYll2mT4F5QG4bS"
)

/*
Answer below questions before starting the function.

When the function invoked?
 - CARGO -> HTTP -> CREATE
What does it do?
- Explain the purpose of the function.(O'zbekcha yozilsa ham bo'ladi.)
- Driver tomonidan bitta zakazga otklik qilsh uchun.
*/

// Request structures
type (
	// Handle request body
	NewRequestBody struct {
		RequestData HttpRequest `json:"request_data"`
		Auth        AuthData    `json:"auth"`
		Data        Data        `json:"data"`
	}

	HttpRequest struct {
		Method  string      `json:"method"`
		Path    string      `json:"path"`
		Headers http.Header `json:"headers"`
		Params  url.Values  `json:"params"`
		Body    []byte      `json:"body"`
	}

	AuthData struct {
		Type string                 `json:"type"`
		Data map[string]interface{} `json:"data"`
	}

	// Function request body >>>>> GET_LIST, GET_LIST_SLIM, CREATE, UPDATE
	Request struct {
		Data map[string]interface{} `json:"data"`
	}

	// most common request structure -> UPDATE, MULTIPLE_UPDATE, CREATE, DELETE
	Data struct {
		AppId      string                 `json:"app_id"`
		Method     string                 `json:"method"`
		ObjectData map[string]interface{} `json:"object_data"`
		ObjectIds  []string               `json:"object_ids"`
		TableSlug  string                 `json:"table_slug"`
		UserId     string                 `json:"user_id"`
	}

	FunctionRequest struct {
		BaseUrl     string  `json:"base_url"`
		TableSlug   string  `json:"table_slug"`
		AppId       string  `json:"app_id"`
		Request     Request `json:"request"`
		DisableFaas bool    `json:"disable_faas"`
	}
)

// Response structures
type (
	// Create function response body >>>>> CREATE
	Datas struct {
		Data struct {
			Data struct {
				Data map[string]interface{} `json:"data"`
			} `json:"data"`
		} `json:"data"`
	}

	// ClientApiResponse This is get single api response >>>>> GET_SINGLE_BY_ID, GET_SLIM_BY_ID
	ClientApiResponse struct {
		Data ClientApiData `json:"data"`
	}

	ClientApiData struct {
		Data ClientApiResp `json:"data"`
	}

	ClientApiResp struct {
		Response map[string]interface{} `json:"response"`
	}

	Response struct {
		Status string                 `json:"status"`
		Error  string                 `json:"error"`
		Data   map[string]interface{} `json:"data"`
	}

	// GetListClientApiResponse This is get list api response >>>>> GET_LIST, GET_LIST_SLIM
	GetListClientApiResponse struct {
		Data GetListClientApiData `json:"data"`
	}

	GetListClientApiData struct {
		Data GetListClientApiResp `json:"data"`
	}

	GetListClientApiResp struct {
		Response []map[string]interface{} `json:"response"`
	}

	// ClientApiUpdateResponse This is single update api response >>>>> UPDATE
	ClientApiUpdateResponse struct {
		Status      string `json:"status"`
		Description string `json:"description"`
		Data        struct {
			TableSlug string                 `json:"table_slug"`
			Data      map[string]interface{} `json:"data"`
		} `json:"data"`
	}

	// ClientApiMultipleUpdateResponse This is multiple update api response >>>>> MULTIPLE_UPDATE
	ClientApiMultipleUpdateResponse struct {
		Status      string `json:"status"`
		Description string `json:"description"`
		Data        struct {
			Data struct {
				Objects []map[string]interface{} `json:"objects"`
			} `json:"data"`
		} `json:"data"`
	}

	ResponseStatus struct {
		Status string `json:"status"`
	}
)

type MultipleDeleteStruct struct {
	Ids []string `json:"ids"`
}

/*
Answer below questions before starting the function.

When the function invoked?
 - table_slug -> AFTER | BEFORE | HTTP -> CREATE | UPDATE | MULTIPLE_UPDATE | DELETE | APPEND_MANY2MANY | DELETE_MANY2MANY
What does it do?
- Explain the purpose of the function.(O'zbekcha yozilsa ham bo'ladi.)
*/

// Testing types
type (
	Asserts struct {
		Request  ucodesdk.Request
		Response ucodesdk.Response
	}

	FunctionAssert struct{}
)

func (f FunctionAssert) GetAsserts() []Asserts {
	return []Asserts{
		{
			Request: ucodesdk.Request{
				Data: ucodesdk.Data{
					ObjectData: map[string]interface{}{
						"guid": "e06494f0-18dc-4b90-9adc-4de811d846a9",
					},
				},
			},
			Response: ucodesdk.Response{
				Status: "done",
			},
		},
		{
			Request: ucodesdk.Request{
				Data: ucodesdk.Data{
					ObjectData: map[string]interface{}{
						"guid": "e06494f0-18dc-4b90-9adc-4de811d846a9",
					},
				},
			},
			Response: ucodesdk.Response{Status: "error"},
		},
	}
}

func (f FunctionAssert) GetBenchmarkRequest() Asserts {
	return Asserts{
		Request: ucodesdk.Request{
			Data: ucodesdk.Data{
				ObjectData: map[string]interface{}{
					"guid": "ded64958-8a89-4587-9263-426d0605c054",
				},
			},
		},
		Response: ucodesdk.Response{
			Status: "done",
		},
	}
}

type GetListWithCount struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		TableSlug string `json:"table_slug"`
		Data      struct {
			Count    int                      `json:"count"`
			Response []map[string]interface{} `json:"response"`
		} `json:"data"`
		IsCached bool `json:"is_cached"`
	} `json:"data"`
	CustomMessage string `json:"custom_message"`
}

var tableSlugs = []string{
	"patient_card",
	"patient_cards",
	"client_files",
	"patient_visits",
	"selected_doctors",
	"puls",
	"blood_sugar",
	"imt",
	"patient_medication",
	"naznachenie",
	"subscription",
	"subscription_report",
	"transactions",
	"medicine_taking",
	"medicine_taking_test",
	"patient_medication_test",
	"walk",
}

// Handle a serverless request
func Handle(req []byte) string {

	var (
		request  NewRequestBody
		response = Response{
			Status: "done",
		}
	)

	//Send2Bot("BEGIN >>>>>>>>" + string(req))

	err := json.Unmarshal(req, &request)
	if err != nil {
		response.Status = "error"
		response.Data = map[string]interface{}{
			"message": "Error while unmarshaling request",
		}
		//Send2Bot("Error while unmarshaling request" + " " + err.Error())
		resp, _ := json.Marshal(response)
		return string(resp)
	}

	//clientId := cast.ToString(request.Data.ObjectIds[0])

	// timeNow := time.Now()

	//var wg sync.WaitGroup

	//for _, v := range tableSlugs {
	//	wg.Add(1)
	//	go func(v string) {
	//		DeleteWithRelations(v, clientId)
	//		defer wg.Done()
	//
	//	}(v)
	//}

	//wg.Wait()

	//DeleteWithRelationsForNotification(clientId)

	// for _, v := range tableSlugs {

	// 	DeleteWithRelations(v, clientId)

	// }

	// fmt.Println(time.Since(timeNow))

	//Send2Bot("ENDING FUNC >>>>>")

	response.Status = "done"
	dataByte, _ := json.Marshal(response)

	return string(dataByte)
}

//func Send2Bot(text string) {
//	bot, _ := tgbotapi.NewBotAPI("6877626035:AAFO2hwSqDnUTHGqqQvsLe9ps5vh1QIiK3Y")
//
//	msg := tgbotapi.NewMessage(-412867319, text)
//
//	bot.Send(msg)
//}

func DeleteWithRelations(slug, id string) {

	// fmt.Println(slug)

	getRequest := Request{
		Data: map[string]interface{}{
			"cleints_id": id,
		},
	}

	url := "https://api.admin.u-code.io/v2/object/get-list/" + slug

	getListResponseInByte, err := DoRequest(url, "POST", getRequest, appId)

	if err != nil {
		// fmt.Println("error 1")
		//panic(err)
		return
	}

	var getListObject GetListWithCount

	err = json.Unmarshal(getListResponseInByte, &getListObject)

	if err != nil {
		// fmt.Println("log 123: ", string(getListResponseInByte))
		// fmt.Println("error 2")
		//panic(err)
		return
	}

	deleteArr := []string{}

	for i := range getListObject.Data.Data.Response {

		deleteArr = append(deleteArr, cast.ToString(getListObject.Data.Data.Response[i]["guid"]))

	}

	// count := cast.ToInt(getListObject.Data.Data.Count)

	// fmt.Println(deleteArr)

	requestForDelete := MultipleDeleteStruct{
		Ids: deleteArr,
	}

	urlForDelete := "https://api.admin.u-code.io/v1/object/" + slug

	_, err = DoRequest(urlForDelete, "DELETE", requestForDelete, appId)

	if err != nil {
		// fmt.Println("error 3")
		//panic(err)
		return

	}

}

func DeleteWithRelationsForNotification(id string) {

	getRequest := Request{
		Data: map[string]interface{}{
			"client_id": id,
		},
	}

	url := "https://api.admin.u-code.io/v2/object/get-list/notifications"

	getListResponseInByte, err := DoRequest(url, "POST", getRequest, appId)

	if err != nil {
		// fmt.Println("error 1")
		//panic(err)
		return
	}

	var getListObject GetListWithCount

	err = json.Unmarshal(getListResponseInByte, &getListObject)

	if err != nil {
		// fmt.Println("log 123: ", string(getListResponseInByte))
		// fmt.Println("error 2")
		//panic(err)
		return
	}

	deleteArr := []string{}

	for i := range getListObject.Data.Data.Response {

		deleteArr = append(deleteArr, cast.ToString(getListObject.Data.Data.Response[i]["guid"]))

	}

	// count := cast.ToInt(getListObject.Data.Data.Count)

	// fmt.Println("nitification:", deleteArr)

	requestForDelete := MultipleDeleteStruct{
		Ids: deleteArr,
	}

	urlForDelete := "https://api.admin.u-code.io/v1/object/notifications"

	_, err = DoRequest(urlForDelete, "DELETE", requestForDelete, appId)

	if err != nil {
		// fmt.Println("error 3")
		//panic(err)
		return

	}

}

func DoRequest(url string, method string, body interface{}, appId string) ([]byte, error) {
	data, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	request.Header.Add("authorization", "API-KEY")
	request.Header.Add("X-API-KEY", appId)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respByte, nil
}

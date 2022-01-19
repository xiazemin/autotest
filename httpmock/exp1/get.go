package exp1

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getAPIResponse(url string) (string, error) {
	var err error
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	myClient := http.Client{}
	fmt.Println("xiazemin1", myClient.Transport, http.DefaultClient)
	response, err := myClient.Do(request)
	if err != nil {
		return "", err
	}
	fmt.Println("xiazemin2")
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", errors.New("response not 200!")
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

//https://www.jianshu.com/p/545963b593de
//https://github.com/jarcoal/httpmock

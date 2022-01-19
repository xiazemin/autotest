package exp1

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
)

func Test_getAPIResponse(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		prepare func(args *args)
	}{
		// TODO: Add test cases.
		{
			name: "case1",

			args: args{
				url: "https://api.mybiz.com/articles", //string
			},
			want:    `[{"id": 1, "name": "My Great Article"}]`, //string,
			wantErr: false,
			prepare: func(args *args) {
				fmt.Println("xiazemin0", http.DefaultTransport)
				fmt.Println(httpmock.DefaultTransport)
				fmt.Println("xiazemin01", http.DefaultTransport)

				// Exact URL match
				httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles", httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Article"}]`))

				// Regexp match (could use httpmock.RegisterRegexpResponder instead)
				// httpmock.RegisterResponder("GET", `=~^https://api\.mybiz\.com/articles/id/\d+\z`,
				// 	httpmock.NewStringResponder(200, `{"id": 1, "name": "My Great Article"}`))

				// do stuff that makes a request to articles

				// get count info
				//httpmock.GetTotalCallCount()

				// get the amount of calls for the registered responder
				//info := httpmock.GetCallCountInfo()
				//fmt.Println(info["GET https://api.mybiz.com/articles"])               // number of GET calls made to https://api.mybiz.com/articles
				//fmt.Println(info["GET https://api.mybiz.com/articles/id/12"])         // number of GET calls made to https://api.mybiz.com/articles/id/12
				//fmt.Println(info[`GET =~^https://api\.mybiz\.com/articles/id/\d+\z`]) // number of GET calls made to https://api.mybiz.com/articles/id/<any-number>
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		httpmock.Activate()
		if tt.prepare != nil {
			tt.prepare(&tt.args)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAPIResponse(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAPIResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getAPIResponse() = %v, want %v", got, tt.want)
			}
		})
		httpmock.DeactivateAndReset()
	}
}

//% go test -v -timeout 3s -tags dynamic -run ^Test_getAPIResponse$ github.com/xiazemin/autotest/httpmock/exp1

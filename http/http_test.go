package toolshttp

import (
	"bytes"
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNewHttpClient(t *testing.T) {
	type fields struct {
		client  *http.Client
		method  string
		url     string
		headers []map[string]string
		body    io.Reader
	}
	data, _ := json.Marshal(map[string]string{"hello": "golang"})
	testPost := bytes.NewBuffer(data)
	tests := []struct {
		name   string
		fields fields

		want  int
		want1 *bytes.Buffer
		want2 error
	}{
		{
			name: "test-GET",
			fields: fields{
				client: &http.Client{
					Timeout: 10 * time.Second,
				},
				method: "GET",
				url:    "http://www.bigfool.cn/test",
			},
			want:  200,
			want1: bytes.NewBuffer([]byte("\"ok\"")),
			want2: nil,
		},
		{
			name: "test-POST",
			fields: fields{
				client: &http.Client{
					Timeout: 10 * time.Second,
				},
				method:  "POST",
				url:     "http://www.bigfool.cn/test",
				headers: []map[string]string{{"Content-Type": "application/json"}},
				body:    testPost,
			},
			want:  200,
			want1: bytes.NewBuffer([]byte("{\"hello\":\"golang\"}")),
			want2: nil,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.RegisterResponder(tt.fields.method, tt.fields.url, httpmock.NewStringResponder(200, tt.want1.String()))
			gp := &GoHttpClient{
				client:  tt.fields.client,
				method:  tt.fields.method,
				url:     tt.fields.url,
				headers: tt.fields.headers,
				body:    tt.fields.body,
			}
			got, got1, got2 := gp.Do()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("do() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("do() got2 = %v, want %v", got2, tt.want2)
			}
			if got2 != tt.want2 {
				t.Errorf("do() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

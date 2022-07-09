package service

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/TadahTech/foodlogiq-demo/pkg/model"
)

func Test_userFromBearer(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userFromBearer(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("userFromBearer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userFromBearer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

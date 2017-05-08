package encrypt

import (
	"reflect"
	"testing"
)

func TestNewSaultEncryptor(t *testing.T) {
	type args struct {
		salt1 string
		salt2 string
	}
	tests := []struct {
		name string
		args args
		want *saultEncryptor
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSaultEncryptor(tt.args.salt1, tt.args.salt2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSaultEncryptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaultEncryptor_Encrypt(t *testing.T) {
	type args struct {
		content string
		param   interface{}
	}
	tests := []struct {
		name      string
		encryptor *saultEncryptor
		args      args
		want      string
		wantErr   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.encryptor.Encrypt(tt.args.content, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("saultEncryptor.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("saultEncryptor.Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

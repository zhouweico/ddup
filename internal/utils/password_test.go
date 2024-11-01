package utils

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "有效密码",
			password: "validPassword123",
			wantErr:  false,
		},
		{
			name:     "空密码",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !strings.HasPrefix(got, "$2a$") {
				t.Errorf("HashPassword() = %v, want bcrypt hash", got)
			}
		})
	}
}

func TestComparePasswords(t *testing.T) {
	password := "testPassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() failed: %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		plainPassword  string
		want           bool
	}{
		{
			name:           "正确密码",
			hashedPassword: hash,
			plainPassword:  password,
			want:           true,
		},
		{
			name:           "错误密码",
			hashedPassword: hash,
			plainPassword:  "wrongPassword",
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComparePasswords(tt.hashedPassword, tt.plainPassword); got != tt.want {
				t.Errorf("ComparePasswords() = %v, want %v", got, tt.want)
			}
		})
	}
}

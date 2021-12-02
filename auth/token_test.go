package auth

import (
	"testing"
)

func TestJWTSettings_validate(t *testing.T) {
	type fields struct {
		Issuer     string
		Lifespan   int
		SigningKey string
		Algorithm  signingAlgorithm
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid settings",
			fields: fields{
				Issuer:     "Heimdall",
				Lifespan:   60,
				SigningKey: "secretkey",
				Algorithm:  HMAC256Algorithm,
			},
			wantErr: false,
		},
		{
			name: "Invalid settings - missing issuer",
			fields: fields{
				Lifespan:   60,
				SigningKey: "secretkey",
				Algorithm:  HMAC256Algorithm,
			},
			wantErr: true,
		},
		{
			name: "Invalid settings - missing lifespan",
			fields: fields{
				Issuer:     "Heimdall",
				SigningKey: "secretkey",
				Algorithm:  HMAC256Algorithm,
			},
			wantErr: true,
		},
		{
			name: "Invalid settings - missing signing key",
			fields: fields{
				Issuer:    "Heimdall",
				Lifespan:  60,
				Algorithm: HMAC256Algorithm,
			},
			wantErr: true,
		},
		{
			name: "Invalid settings - missing algorithm",
			fields: fields{
				Issuer:     "Heimdall",
				Lifespan:   60,
				SigningKey: "secretkey",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := JWTSettings{
				Issuer:     tt.fields.Issuer,
				Lifespan:   tt.fields.Lifespan,
				SigningKey: tt.fields.SigningKey,
				Algorithm:  tt.fields.Algorithm,
			}
			if err := s.validate(); (err != nil) != tt.wantErr {
				t.Errorf("JWTSettings.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_generateJWT(t *testing.T) {
	tests := []struct {
		name     string
		settings JWTSettings
		want     Token
		wantErr  bool
	}{
		{
			name: "Success",
			settings: JWTSettings{
				Issuer:     "Heimdall",
				Lifespan:   60,
				SigningKey: "secretkey",
				Algorithm:  HMAC256Algorithm,
			},
			want: Token{
				Lifespan: 60,
			},
			wantErr: false,
		},
		{
			name:     "Invalid settings returns error",
			settings: JWTSettings{},
			want:     Token{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateJWT(tt.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && (got.Lifespan != tt.want.Lifespan || got.AccessToken == "") {
				t.Errorf("generateJWT() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

package crypto

import (
	"bytes"
	"encoding/base64"
	"reflect"
	"testing"
)

// This test ensures that the default params don't accidentally change. Changing these
// could result in the invalidation of all passwords if a deployed intance of Heimdall
// is using the defaults.
func Test_defaultParamsAreCorrect(t *testing.T) {
	expected := ArgonParams{
		Time:    1,
		Memory:  2 * 1024 * 1024,
		Threads: 4,
		KeyLen:  32,
		SaltLen: 16,
	}

	if !reflect.DeepEqual(expected, DefaultParams) {
		t.Errorf("Default params have been changed: expected %v, got %v", expected, DefaultParams)
	}
}

func Test_hashesAreEqual(t *testing.T) {
	tests := []struct {
		name string
		a    []byte
		b    []byte
		want bool
	}{
		{
			name: "Equal",
			a:    []byte{'a', 'b', 'c'},
			b:    []byte{'a', 'b', 'c'},
			want: true,
		},
		{
			name: "Not equal - same length",
			a:    []byte{'a', 'b', 'c'},
			b:    []byte{'a', 'b', 'd'},
			want: false,
		},
		{
			name: "Not equal - different length",
			a:    []byte{'b', 'c'},
			b:    []byte{'a', 'b', 'd'},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hashesAreEqual(tt.a, tt.b); got != tt.want {
				t.Errorf("hashesAreEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashPassword(t *testing.T) {
	password := "password123"
	params := ArgonParams{
		Time:    3,
		Memory:  64 * 1024,
		Threads: 4,
		KeyLen:  32,
		SaltLen: 16,
	}
	salt, err := base64.RawStdEncoding.DecodeString("sxCtsSYtbBo4tnUj6v7sCw")
	if err != nil {
		t.Fatalf("Unexpected error decoding salt: %v", err)
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString("Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY")
	if err != nil {
		t.Fatalf("Unexpected error decoding hash: %v", err)
	}

	hash, err := hashPassword(password, salt, params)
	if err != nil {
		t.Fatalf("Unexpected error hashing password: %v", err)
	}

	if !bytes.Equal(expectedHash, hash) {
		t.Errorf("Incorrect password hash: expected %v, got %v", expectedHash, hash)
	}
}

func Test_encodePassword(t *testing.T) {
	params := ArgonParams{
		Time:    3,
		Memory:  64 * 1024,
		Threads: 4,
		KeyLen:  32,
		SaltLen: 16,
	}
	salt, err := base64.RawStdEncoding.DecodeString("sxCtsSYtbBo4tnUj6v7sCw")
	if err != nil {
		t.Fatalf("Unexpected error decoding salt: %v", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString("Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY")
	if err != nil {
		t.Fatalf("Unexpected error decoding hash: %v", err)
	}

	expectedString := "$argon2id$v=19$m=65536,t=3,p=4$sxCtsSYtbBo4tnUj6v7sCw$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY"

	encodedHash := encodeHash(hash, salt, params)

	if encodedHash != expectedString {
		t.Errorf("Incorrect encoded hash: expected %v, got %v", expectedString, encodedHash)
	}
}

func Test_decodeHash(t *testing.T) {
	tests := []struct {
		name        string
		encodedHash string
		hash        []byte
		salt        []byte
		params      ArgonParams
		wantErr     bool
	}{
		{
			name:        "Success",
			encodedHash: "$argon2id$v=19$m=65536,t=3,p=4$sxCtsSYtbBo4tnUj6v7sCw$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY",
			hash:        []byte{86, 41, 169, 218, 143, 172, 94, 234, 42, 16, 235, 244, 244, 84, 58, 153, 97, 137, 44, 7, 92, 211, 138, 238, 122, 57, 13, 203, 33, 70, 72, 246},
			salt:        []byte{179, 16, 173, 177, 38, 45, 108, 26, 56, 182, 117, 35, 234, 254, 236, 11},
			params: ArgonParams{
				Time:    3,
				Memory:  64 * 1024,
				Threads: 4,
				KeyLen:  32,
				SaltLen: 16,
			},
			wantErr: false,
		},
		{
			name:        "Malformed hash - incorrect number of $'s",
			encodedHash: "argon2id$v=19$m=65536,t=3,p=4$sxCtsSYtbBo4tnUj6v7sCw$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY",
			wantErr:     true,
		},
		{
			name:        "Malformed hash - unexpected prefix",
			encodedHash: "abc$argon2id$v=19$m=65536,t=3,p=4$sxCtsSYtbBo4tnUj6v7sCw$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY",
			wantErr:     true,
		},
		{
			name:        "Malformed hash - unsupported algorithm",
			encodedHash: "$unsupported$v=19$m=65536,t=3,p=4$sxCtsSYtbBo4tnUj6v7sCw$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY",
			wantErr:     true,
		},
		{
			name:        "Malformed hash - unsupported version",
			encodedHash: "$argon2id$v=-1$m=65536,t=3,p=4$sxCtsSYtbBo4tnUj6v7sCw$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY",
			wantErr:     true,
		},
		{
			name:        "Malformed hash - missing memory param",
			encodedHash: "$argon2id$v=19$m=,t=3,p=4$sxCtsSYtbBo4tnUj6v7sCw$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY",
			wantErr:     true,
		},
		{
			name:        "Malformed hash - missing time param",
			encodedHash: "$argon2id$v=19$m=65536,t=,p=4$sxCtsSYtbBo4tnUj6v7sCw$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY",
			wantErr:     true,
		},
		{
			name:        "Malformed hash - missing threads param",
			encodedHash: "$argon2id$v=19$m=65536,t=3,p=$sxCtsSYtbBo4tnUj6v7sCw$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY",
			wantErr:     true,
		},
		{
			name:        "Malformed hash - invalid salt encoding",
			encodedHash: "$argon2id$v=19$m=65536,t=3,p=4$sxCtsSYtbBo4tnUj6v7sCw&$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY",
			wantErr:     true,
		},
		{
			name:        "Malformed hash - empty salt",
			encodedHash: "$argon2id$v=19$m=65536,t=3,p=4$$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY",
			wantErr:     true,
		},
		{
			name:        "Malformed hash - invalid hash encoding",
			encodedHash: "$argon2id$v=19$m=65536,t=3,p=4$sxCtsSYtbBo4tnUj6v7sCw$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY&",
			wantErr:     true,
		},
		{
			name:        "Malformed hash - empty hash",
			encodedHash: "$argon2id$v=19$m=65536,t=3,p=4$sxCtsSYtbBo4tnUj6v7sCw$",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := decodeHash(tt.encodedHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.hash) {
				t.Errorf("decodeHash() got = %v, want %v", got, tt.hash)
			}
			if !reflect.DeepEqual(got1, tt.salt) {
				t.Errorf("decodeHash() got1 = %v, want %v", got1, tt.salt)
			}
			if !reflect.DeepEqual(got2, tt.params) {
				t.Errorf("decodeHash() got2 = %v, want %v", got2, tt.params)
			}
		})
	}
}

func Test_generateSaltGeneratesSaltOfCorrectSize(t *testing.T) {
	l := 16
	salt, err := generateSalt(uint32(l))
	if err != nil {
		t.Fatalf("Unexpected error generating salt: %v", err)
	}

	if len(salt) != l {
		t.Errorf("Incorrect salt length: expected %d, got %d", l, len(salt))
	}
}

func Test_generateSaltGeneratesRandomSalts(t *testing.T) {
	l := uint32(16)
	salt1, err := generateSalt(l)
	if err != nil {
		t.Fatalf("Unexpected error generating salt: %v", err)
	}

	salt2, err := generateSalt(l)
	if err != nil {
		t.Fatalf("Unexpected error generating salt: %v", err)
	}

	if bytes.Equal(salt1, salt2) {
		t.Errorf("Successive salts are equal")
	}
}

func TestValidatePassword(t *testing.T) {
	type args struct {
		password    string
		encodedHash string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Valid password",
			args: args{
				password:    "password123",
				encodedHash: "$argon2id$v=19$m=65536,t=3,p=4$sxCtsSYtbBo4tnUj6v7sCw$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Invalid password",
			args: args{
				password:    "foobar",
				encodedHash: "$argon2id$v=19$m=65536,t=3,p=4$sxCtsSYtbBo4tnUj6v7sCw$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Invalid hash - malformed encoding",
			args: args{
				password:    "foobar",
				encodedHash: "argon2id$v=19$m=65536,t=3,p=4$sxCtsSYtbBo4tnUj6v7sCw$Vimp2o+sXuoqEOv09FQ6mWGJLAdc04ruejkNyyFGSPY",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidatePassword(tt.args.password, tt.args.encodedHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPasswordHash(t *testing.T) {
	type args struct {
		password string
		p        ArgonParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				password: "password123",
				p: ArgonParams{
					Time:    1,
					Memory:  1024,
					Threads: 4,
					KeyLen:  32,
					SaltLen: 16,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPasswordHash(tt.args.password, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("GetPasswordHash() = %v, want \"\"", got)
			}
		})
	}
}

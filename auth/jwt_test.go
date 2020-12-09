package auth

import (
	"context"
	"reflect"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	type args struct {
		ctx context.Context
		jt  JWTToken
		au  Auth
	}
	tests := []struct {
		name    string
		args    args
		want    *SignedToken
		wantErr bool
	}{
		{
			name: "generatetoken - success",
			args: args{
				ctx: context.Background(),
				au:  Auth{Username: "stefan", Password: "password"},
				jt:  JWTToken{SigningKey: "SuperSecureKey"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.jt.GenerateToken(tt.args.ctx, tt.args.au)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJWTToken_ValidateToken(t *testing.T) {
	// Replace SignedToken with token from above
	type fields struct {
		SigningKey string
	}
	type args struct {
		ctx context.Context
		st  SignedToken
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "validtoken - success",
			args: args{
				ctx: context.Background(),
				st:  SignedToken{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGFpbXMiOnsiZXhwIjoxNjA3NDgwOTY0fSwidXNlcm5hbWUiOiJzdGVmYW4ifQ.GqNsyBLfENUHNFhjmdVi4nbzbv3md1DSlCqiuMz-seU"},
			},
			fields:  fields{SigningKey: "SuperSecureKey"},
			want:    "stefan",
			wantErr: false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jt := JWTToken{
				SigningKey: tt.fields.SigningKey,
			}
			got, err := jt.ValidateToken(tt.args.ctx, tt.args.st)
			if (err != nil) != tt.wantErr {
				t.Errorf("JWTToken.ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JWTToken.ValidateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

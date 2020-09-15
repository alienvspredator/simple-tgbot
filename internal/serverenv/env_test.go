package serverenv_test

import (
	"context"
	"testing"

	"github.com/alienvspredator/simple-tgbot/internal/serverenv"
)

func TestServerEnv_Close(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		s       *serverenv.ServerEnv
		args    args
		wantErr bool
	}{
		{
			name: "Close nil env",
			s:    nil,
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Close non-nil env without all-nil fields",
			s:    serverenv.New(context.Background()),
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if err := tt.s.Close(tt.args.ctx); (err != nil) != tt.wantErr {
					t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

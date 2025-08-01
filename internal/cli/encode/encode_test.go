package encode

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

func TestEncoder_Encode(t *testing.T) {
	type testcase struct {
		name string
		encoder Encoder
		message proto.Message
		want string
	}
	tests := []testcase{
		{
			name: "encode to JSON",
			encoder: json,
			message: &v1alpha1.LongRunningOperation{
				Name: "operation-1234",
			},
			want: `{"name":"operation-1234"}`,
		},
		{
			name: "encode to YAML",
			encoder: yaml,
			message: &v1alpha1.LongRunningOperation{
				Name: "operation-1234",
			},
			want: "name: operation-1234\n",
		},
		{
			name: "encode to text",
			encoder: text,
			message: &v1alpha1.LongRunningOperation{
				Name: "operation-1234",
			},
			want: `name:"operation-1234"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.Buffer
			tt.encoder.Encode(&got, tt.message)
			if diff := cmp.Diff(tt.want, got.String()); diff != "" {
				t.Errorf("Encoder.Encode() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEncoder_VarP(t *testing.T) {
	type testcase struct {
		name string
		args []string
		want Encoder
	}
	tests := []testcase{
		{
			name: "set to JSON",
			args: []string{"--output", "json"},
			want: json,
		},
		{
			name: "set to YAML",
			args: []string{"--output", "yaml"},
			want: yaml,
		},
		{
			name: "set to text",
			args: []string{"--output", "text"},
			want: text,
		},
		{
			name: "set to JSON with alias",
			args: []string{"-o", "json"},
			want: json,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var encoder Encoder
			cmd := &cobra.Command{}
			encoder.VarP(cmd)
			cmd.SetArgs(tt.args)
			if err := cmd.Execute(); err != nil {
				t.Fatalf("cmd.Execute() error = %v", err)
			}
			if diff := cmp.Diff(tt.want, encoder); diff != "" {
				t.Errorf("Encoder mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

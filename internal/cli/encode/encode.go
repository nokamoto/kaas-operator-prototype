package encode

import (
	"fmt"
	"io"

	"buf.build/go/protoyaml"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

// Encoder defines an interface for encoding protobuf messages into different formats.
// It supports JSON, YAML, and text formats.
//
// It implements pflag.Value interface to allow it to be used as a flag value in Cobra commands.
type Encoder string

const (
	json Encoder = "json"
	yaml Encoder = "yaml"
	text Encoder = "text"
)

func suppressError(cmd *cobra.Command, err error) {
	// no-op for now, can be used to handle errors gracefully
}

// VarP registers the Encoder as a flag with the given command.
func (e *Encoder) VarP(cmd *cobra.Command) {
	cmd.Flags().VarP(e, "output", "o", "Output format (json, yaml, text)")
}

// Print encodes the given protobuf message and writes it to the command's output.
func (e *Encoder) Print(cmd *cobra.Command, v proto.Message) {
	e.Encode(cmd.OutOrStdout(), v)
}

// Encode encodes the given protobuf message into the specified format and writes it to the provided writer.
// It supports JSON, YAML, and text formats based on the value of the Encoder.
func (e *Encoder) Encode(w io.Writer, v proto.Message) error {
	var bytes []byte
	var err error
	switch *e {
	case json:
		bytes, err = protojson.Marshal(v)
	case text:
		bytes, err = prototext.Marshal(v)
	default:
		bytes, err = protoyaml.Marshal(v)
	}
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}

func (e *Encoder) String() string {
	return string(*e)
}

func (e *Encoder) Set(v string) error {
	switch v {
	case "json":
		*e = json
	case "yaml":
		*e = yaml
	case "text":
		*e = text
	default:
		return fmt.Errorf("unknown output format: %s, must be one of json, yaml, text", v)
	}
	return nil
}

func (e *Encoder) Type() string {
	return "Encoder"
}

package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/frawleyskid/w3s-upload/parse"
	"github.com/ipfs/go-cid"
	"github.com/spf13/viper"
	"reflect"
)

const prettyJsonIndent = "  "

var ErrInvalidOutput = errors.New("invalid output parameter")

type Output struct {
	Artifacts []parse.Artifact `json:"artifacts"`
}

// initializeNilSlicesOfValue accepts a struct and initializes any nil slices in
// its fields to slices of len = 0 and cap = 0. If one of its fields is a
// struct or slice of structs, it does the same for those struct's fields
// recursively.
func initializeNilSlicesOfValue(value reflect.Value) {
	for fieldIndex := 0; fieldIndex < value.NumField(); fieldIndex++ {
		switch field := value.Field(fieldIndex); field.Kind() {
		case reflect.Struct:
			initializeNilSlicesOfValue(field)
		case reflect.Slice:
			if field.IsNil() {
				field.Set(reflect.MakeSlice(field.Type(), 0, 0))
			} else if field.Type().Elem().Kind() == reflect.Struct {
				for sliceIndex := 0; sliceIndex < field.Len(); sliceIndex++ {
					initializeNilSlicesOfValue(field.Index(sliceIndex))
				}
			}
		}
	}
}
func initializeNilSlices(value interface{}) {
	initializeNilSlicesOfValue(reflect.ValueOf(value).Elem())
}

func marshalArtifact(artifacts []parse.Artifact, pretty bool) (string, error) {
	output := Output{Artifacts: artifacts}

	// To make the output more consistent and easier to parse, we want to
	// normalize any nil slices to empty slices before we serialize so that
	// they're serialized as `[]` and not `null`.
	initializeNilSlices(&output)

	var (
		marshalledOutput []byte
		err              error
	)

	if pretty {
		marshalledOutput, err = json.MarshalIndent(output, "", prettyJsonIndent)
	} else {
		marshalledOutput, err = json.Marshal(output)
	}

	if err != nil {
		return "", err
	}

	return string(marshalledOutput), nil
}

func marshalCid(cids []cid.Cid, pretty bool) (string, error) {
	var (
		marshalledOutput []byte
		err              error
	)

	marshalledCids := make([]string, len(cids))

	for i, id := range cids {
		marshalledCids[i] = id.String()
	}

	if pretty {
		marshalledOutput, err = json.MarshalIndent(marshalledCids, "", prettyJsonIndent)
	} else {
		marshalledOutput, err = json.Marshal(marshalledCids)
	}

	if err != nil {
		return "", err
	}

	return string(marshalledOutput), nil
}

func Print(entries []parse.Artifact, cidList []cid.Cid) error {
	if viper.GetBool("action") {
		artifactOutput, err := marshalArtifact(entries, false)
		if err != nil {
			return err
		}

		fmt.Printf("::set-output name=artifacts::%s\n", artifactOutput)

		cidOutput, err := marshalCid(cidList, false)
		if err != nil {
			return err
		}

		fmt.Printf("::set-output name=cids::%s\n", cidOutput)

		return nil
	}

	switch outputMode := viper.GetString("output"); outputMode {
	case "artifacts":
		artifactOutput, err := marshalArtifact(entries, true)
		if err != nil {
			return err
		}

		fmt.Println(artifactOutput)
	case "cids":
		cidOutput, err := marshalCid(cidList, true)
		if err != nil {
			return err
		}

		fmt.Println(cidOutput)
	case "summary":
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOutput, outputMode)
	}

	return nil
}

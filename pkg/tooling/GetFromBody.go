package tooling

import (
	"encoding/json"
	"io"
)

func GetFromBody(closer io.ReadCloser, any2 any) error {
	body, err := io.ReadAll(closer)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &any2)
	if err != nil {
		return err
	}
	return nil
}

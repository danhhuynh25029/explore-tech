package iplocate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type HTTPGetter interface {
	Get(url string) (*http.Response, error)
}

type Activities struct {
	HTTPClient HTTPGetter
}

func (i *Activities) ActivityA(ctx context.Context) (string, error) {
	fmt.Println("ActivityA")
	return "", errors.New("error")
}

func (i *Activities) ActivityB(ctx context.Context) (string, error) {
	fmt.Println("ActivityB")
	return "", nil
}

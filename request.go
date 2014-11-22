package goblin

import (
    "net/http"
)

type Request struct {
    *http.Request
}

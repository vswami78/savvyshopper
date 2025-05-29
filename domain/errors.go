package domain

import "errors"

var (
    ErrNetwork   = errors.New("network error")
    ErrAuth      = errors.New("authentication error")
    ErrNoResults = errors.New("no offers found")
)

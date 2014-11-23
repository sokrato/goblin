package goblin

import (
    "net/url"
    "strconv"
)

type Form url.Values

func (f Form) String(key string) (string, error) {
    val, ok := f[key]
    if !ok {
        return "", ErrKeyIsNotSet
    }
    return val[0], nil
}

func (f Form) Int(key string) (int, error) {
    val, ok := f[key]
    if !ok {
        return 0, ErrKeyIsNotSet
    }
    return strconv.Atoi(val[0])
}

func (f Form) Uint(key string) (uint, error) {
    val, ok := f[key]
    if !ok {
        return 0, ErrKeyIsNotSet
    }
    num, err := strconv.ParseUint(val[0], 10, strconv.IntSize)
    if err != nil {
        return 0, err
    }
    return uint(num), nil
}

func (f Form) Int64(key string) (int64, error) {
    val, ok := f[key]
    if !ok {
        return 0, ErrKeyIsNotSet
    }
    return strconv.ParseInt(val[0], 10, 64)
}

func (f Form) Uint64(key string) (uint64, error) {
    val, ok := f[key]
    if !ok {
        return 0, ErrKeyIsNotSet
    }
    num, err := strconv.ParseInt(val[0], 10, 64)
    if err != nil {
        return 0, err
    }
    return uint64(num), nil
}

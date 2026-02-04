package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var ErrEnvFileNotFound = errors.New("env file not found in default search paths")

type Config struct {
	loaded bool
}

func NewConfig(path ...string) (*Config, error) {
	if envPath := os.Getenv("ENV_PATH"); envPath != "" {
		return &Config{loaded: true}, nil
	}

	if len(path) > 0 {
		return loadFromPath(path[0])
	}

	return searchAndLoadEnv()
}

func loadFromPath(path string) (*Config, error) {
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if err := godotenv.Load(fullPath); err != nil {
		return nil, errors.New("failed to load env file: " + fullPath)
	}
	os.Setenv("ENV_PATH", fullPath)
	return &Config{loaded: true}, nil
}

func searchAndLoadEnv() (*Config, error) {
	execPath, _ := os.Executable()
	fileName := filepath.Base(execPath)
	wd, _ := os.Getwd()

	// 로딩 순위 - ./.파일명.env -> ./.config.env -> ../.config.env
	searchPaths := []string{
		filepath.Join(wd, "."+fileName+".env"),
		filepath.Join(wd, ".config.env"),
		filepath.Join(filepath.Dir(wd), ".config.env"),
	}

	for _, path := range searchPaths {
		if err := godotenv.Load(path); err == nil {
			os.Setenv("ENV_PATH", path)
			return &Config{loaded: true}, nil
		}
	}
	return nil, ErrEnvFileNotFound
}

func (e *Config) Get(key string, defaultVal ...string) string {
	val := os.Getenv(key)
	if val == "" && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return val
}

func (e *Config) String(key string, defaultVal ...string) string {
	return e.Get(key, defaultVal...)
}

func (e *Config) Int(key string, defaultVal ...int) int {
	val := e.Get(key)
	if data, err := strconv.Atoi(val); err == nil {
		return data
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return 0
}

func (e *Config) Int64(key string, defaultVal ...int64) int64 {
	val := e.Get(key)
	if data, err := strconv.ParseInt(val, 10, 64); err == nil {
		return data
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return 0
}

func (e *Config) Float64(key string, defaultVal ...float64) float64 {
	val := e.Get(key)
	if data, err := strconv.ParseFloat(val, 64); err == nil {
		return data
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return 0
}

func (e *Config) Bool(key string, defaultVal ...bool) bool {
	val := e.Get(key)
	if data, err := strconv.ParseBool(val); err == nil {
		return data
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return false
}

func (e *Config) SliceString(key string, defaultVal ...[]string) []string {
	val := e.Get(key)
	if val == "" {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return []string{}
	}
	return strings.Split(val, ",")
}

func (e *Config) SliceInt(key string, defaultVal ...[]int) []int {
	return parseSlice(e.Get(key), strconv.Atoi, defaultVal...)
}

func (e *Config) SliceInt64(key string, defaultVal ...[]int64) []int64 {
	return parseSlice(e.Get(key), func(s string) (int64, error) { return strconv.ParseInt(s, 10, 64) }, defaultVal...)
}

func (e *Config) SliceFloat64(key string, defaultVal ...[]float64) []float64 {
	return parseSlice(e.Get(key), func(s string) (float64, error) { return strconv.ParseFloat(s, 64) }, defaultVal...)
}

func parseSlice[T any](raw string, parser func(string) (T, error), defaultVal ...[]T) []T {
	if raw == "" {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return []T{}
	}
	parts := strings.Split(raw, ",")
	res := make([]T, 0, len(parts))
	for _, p := range parts {
		if val, err := parser(strings.TrimSpace(p)); err == nil {
			res = append(res, val)
		} else {
			return []T{} // 파싱 실패 시 빈 슬라이스 반환 (기존 로직 유지)
		}
	}
	return res
}

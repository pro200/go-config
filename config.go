package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const ERROR_NOT_FOUND = "not found .fileName.env or .config.env"

type Config struct {
	loaded bool
}

func NewConfig(path ...string) (*Config, error) {
	// 이미 로드된 경우
	if envPath := os.Getenv("ENV_PATH"); envPath != "" {
		return &Config{loaded: true}, nil
	}

	// from path
	if len(path) > 0 {
		fullPath, err := filepath.Abs(path[0])
		if err != nil {
			return nil, err
		}

		if err := godotenv.Load(fullPath); err == nil {
			os.Setenv("ENV_PATH", fullPath)
			return &Config{loaded: true}, nil
		}

		return nil, errors.New("not found env file in " + fullPath)
	}

	// 로딩 순위 - ./.파일명.env -> ./.config.env -> ../.config.env
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	fileName := filepath.Base(execPath)
	//if strings.HasPrefix(fileName, "___go_") {
	//	return nil, fmt.Errorf("this file name \"%s\" not supported", fileName)
	//}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	paths := strings.Split(wd, string(os.PathSeparator))

	envFiles := []string{
		filepath.Join(wd, "."+fileName+".env"),
		filepath.Join(wd, ".config.env"),
	}

	if len(paths) > 1 {
		envFiles = append(envFiles,
			filepath.Join(strings.Join(paths[:len(paths)-1], string(os.PathSeparator)), ".config.env"),
		)
	}

	for _, file := range envFiles {
		if err := godotenv.Load(file); err == nil {
			os.Setenv("ENV_PATH", file)
			return &Config{loaded: true}, nil
		}
	}

	return nil, errors.New(ERROR_NOT_FOUND)
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
	result := e.Get(key)
	if result == "" {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}
	data, err := strconv.Atoi(result)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return data
}

func (e *Config) Int64(key string, defaultVal ...int64) int64 {
	result := e.Get(key)
	if result == "" {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}
	data, err := strconv.ParseInt(result, 10, 64)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return data
}

func (e *Config) Float64(key string, defaultVal ...float64) float64 {
	result := e.Get(key)
	if result == "" && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	data, err := strconv.ParseFloat(result, 64)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return data
}

func (e *Config) Bool(key string, defaultVal ...bool) bool {
	result := e.Get(key)
	if result == "" {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return false
	}
	data, err := strconv.ParseBool(result)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return data
}

func (e *Config) SliceString(key string, defaultVal ...[]string) []string {
	var s []string

	result := e.Get(key)
	if result == "" {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return []string{}
	}

	parts := strings.Split(result, ",")
	s = append(s, parts...)

	return s
}

func (e *Config) SliceInt(key string, defaultVal ...[]int) []int {
	var s []int

	result := e.Get(key)
	if result == "" {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return []int{}
	}

	parts := strings.Split(result, ",")
	for _, part := range parts {
		val, err := strconv.Atoi(part)
		if err != nil {
			return []int{}
		}
		s = append(s, val)
	}

	return s
}

func (e *Config) SliceInt64(key string, defaultVal ...[]int64) []int64 {
	var s []int64

	result := e.Get(key)
	if result == "" {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return []int64{}
	}

	parts := strings.Split(result, ",")
	for _, part := range parts {
		val, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return []int64{}
		}
		s = append(s, val)
	}

	return s
}

func (e *Config) SliceFloat64(key string, defaultVal ...[]float64) []float64 {
	var s []float64

	result := e.Get(key)
	if result == "" {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return []float64{}
	}

	parts := strings.Split(result, ",")
	for _, part := range parts {
		val, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return []float64{}
		}
		s = append(s, val)
	}

	return s
}

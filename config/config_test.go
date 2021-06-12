package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestConfig_LoadSecrets(t *testing.T) {
	// test scenarios
	const (
		defaultEnv = iota
		envFile
	)

	tests := []struct {
		name     string
		scenario int
		expected Secrets
	}{
		{
			name:     "Test default variables",
			scenario: defaultEnv,
			expected: Secrets{
				Port:        defaultPort,
				CoindeskURL: defaultUrl,
			},
		},
		{
			name:     "Test .env file",
			scenario: envFile,
			expected: Secrets{
				Port:        "1234",
				CoindeskURL: "url",
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			s := &Secrets{}
			switch testCase.scenario {
			case defaultEnv:
				s = LoadSecrets()
			case envFile:
				// create a temp .env file in the current working directory
				path, err := os.Getwd()
				assert.NoError(t, err)

				file, err := ioutil.TempFile(path, "*.env")
				assert.NoError(t, err)

				// add sample env data to temp file
				_, err = file.Write([]byte(fmt.Sprintf(
					"PORT=%v\nCOINDESK_URL=%v",
					testCase.expected.Port,
					testCase.expected.CoindeskURL,
				)))
				assert.NoError(t, err)

				// cleanup temp file
				defer func() {
					err := os.Remove(file.Name())
					assert.NoError(t, err)
				}()

				s = LoadSecrets(file.Name())
			}

			assert.Equal(t, testCase.expected, *s)
		})
	}
}

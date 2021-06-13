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
				Port:      defaultPort,
				JWTSecret: defaultSecret,
				DBName:    defaultDbName,
				DBURL:     defaultDbUrl,
			},
		},
		{
			name:     "Test .env file",
			scenario: envFile,
			expected: Secrets{
				Port:      "1234",
				JWTSecret: "secret",
				DBName:    "track",
				DBURL:     "someUrl",
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
					"PORT=%v\nDATABASE_URL=%v\nDATABASE_NAME=%v\nJWT_SECRET=%v",
					testCase.expected.Port,
					testCase.expected.DBURL,
					testCase.expected.DBName,
					testCase.expected.JWTSecret,
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

func TestNew(t *testing.T) {
	fmt.Println("testing new")
}

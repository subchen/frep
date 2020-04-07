package secretsmanager

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

type mockedSecret struct {
	secretsmanageriface.SecretsManagerAPI
	Resp secretsmanager.GetSecretValueOutput
}

// GetSecret return mocked secret value
func (m mockedSecret) GetSecretValue(in *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	// for now just return empty resp
	return &m.Resp, nil
}

func TestGetSecret(t *testing.T) {
	testCase := struct {
		Resp     secretsmanager.GetSecretValueOutput
		Expected string
	}{
		Resp: secretsmanager.GetSecretValueOutput{
			Name:         aws.String("test/secret"),
			SecretString: aws.String(`{"key": "test"}`),
		},
		Expected: "test",
	}

	s := Secret{
		Client: mockedSecret{Resp: testCase.Resp},
	}
	resp, err := s.GetSecret("test/secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != testCase.Expected {
		t.Fatalf("expected secret %v, got %v", resp, testCase.Expected)
	}
}

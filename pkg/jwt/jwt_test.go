package jwt

import (
	"os"
	"testing"

	"github.com/1ef7yy/medods_test_task/internal/errors"
	"github.com/1ef7yy/medods_test_task/models"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "secret")
	JWTSecret = []byte(os.Getenv("JWT_SECRET"))

	code := m.Run()

	os.Unsetenv("JWT_SECRET")
	os.Exit(code)
}

func TestGenerateTokenPair(t *testing.T) {
	t.Run("successful token generation", func(t *testing.T) {
		req := models.GenerateTokenRequest{
			Guid:       "test-guid",
			IP:         "192.168.1.1",
			Generation: 1,
		}

		tokenPair, err := GenerateTokenPair(req)
		if err != nil {
			t.Fatalf("GenerateTokenPair failed: %s", err.Error())
		}

		if tokenPair.AccessToken == "" {
			t.Error("access token is empty")
		}
		if tokenPair.RefreshToken == "" {
			t.Error("refresh token is empty")
		}
	})

	t.Run("missing JWT secret", func(t *testing.T) {
		originalSecret := JWTSecret
		JWTSecret = nil
		// if test fails, return to original secret
		defer func() {
			JWTSecret = originalSecret
		}()

		req := models.GenerateTokenRequest{
			Guid:       "test-guid",
			IP:         "192.168.1.1",
			Generation: 1,
		}

		_, err := GenerateTokenPair(req)
		if err == nil {
			t.Error("expected error when JWT_SECRET is missing")
		}

		if err != errors.CouldNotFindSecretErr {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})
}

func TestDecodeRefresh(t *testing.T) {
	t.Run("valid refresh token", func(t *testing.T) {
		req := models.GenerateTokenRequest{
			Guid:       "test-guid",
			IP:         "192.168.1.1",
			Generation: 1,
		}

		tokenPair, err := GenerateTokenPair(req)
		if err != nil {
			t.Fatalf("GenerateTokenPair failed: %s", err.Error())
		}

		refresh, err := DecodeRefresh(tokenPair.RefreshToken)
		if err != nil {
			t.Fatalf("DecodeRefresh failed: %s", err.Error())
		}

		if refresh.Guid != req.Guid {
			t.Errorf("expected guid '%s', instead got '%s'", req.Guid, refresh.Guid)
		}
		if refresh.IP != req.IP {
			t.Errorf("expected ip '%s', instead got '%s'", req.IP, refresh.IP)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := DecodeRefresh("invalid.token.string")
		if err == nil {
			t.Error("expected error for invalid token")
		}
		if err != errors.TokenInvalidErr {
			t.Errorf("expected TokenInvalidErr, got %v", err)
		}
	})

	t.Run("malformed token", func(t *testing.T) {
		req := models.GenerateTokenRequest{
			Guid:       "test-guid",
			IP:         "192.168.1.1",
			Generation: 1,
		}

		tokenPair, _ := GenerateTokenPair(req)

		malformedToken := tokenPair.RefreshToken + "f"

		_, err := DecodeRefresh(malformedToken)

		if err == nil {
			t.Error("expected error for malformed token, got nothing")
		}
	})

	t.Run("missing JWT secret", func(t *testing.T) {
		originalSecret := JWTSecret
		JWTSecret = nil

		defer func() {
			JWTSecret = originalSecret
		}()

		_, err := DecodeRefresh("any.token")
		if err == nil {
			t.Error("expected missing JWT secret error")
		}

		if err != errors.CouldNotFindSecretErr {
			t.Errorf("expected %s, got %s", errors.CouldNotFindSecretErr.Error(), err.Error())
		}
	})
}

func TestDecodeAccess(t *testing.T) {
	t.Run("valid access token", func(t *testing.T) {
		req := models.GenerateTokenRequest{
			Guid:       "test-guid",
			IP:         "192.168.1.1",
			Generation: 1,
		}

		tokenPair, err := GenerateTokenPair(req)
		if err != nil {
			t.Fatalf("GenerateTokenPair failed: %s", err.Error())
		}

		access, err := DecodeAccess(tokenPair.AccessToken)
		if err != nil {
			t.Fatalf("DecodeAccess failed: %s", err.Error())
		}

		if access.Guid != req.Guid {
			t.Errorf("expected guid '%s', instead got '%s'", req.Guid, access.Guid)
		}
		if access.Generation != req.Generation {
			t.Errorf("expected generation '%d', instead got '%d'", req.Generation, access.Generation)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := DecodeAccess("invalid.token.string")
		if err == nil {
			t.Error("expected error for invalid token")
		}
		if err != errors.TokenInvalidErr {
			t.Errorf("expected TokenInvalidErr, got %v", err)
		}
	})

	t.Run("malformed token", func(t *testing.T) {
		req := models.GenerateTokenRequest{
			Guid:       "test-guid",
			IP:         "192.168.1.1",
			Generation: 1,
		}

		tokenPair, _ := GenerateTokenPair(req)

		malformedToken := tokenPair.AccessToken + "f"

		_, err := DecodeAccess(malformedToken)

		if err == nil {
			t.Error("expected error for malformed token, got nothing")
		}
	})

	t.Run("missing JWT secret", func(t *testing.T) {
		var originalSecret []byte
		copy(originalSecret, JWTSecret)
		JWTSecret = nil

		defer func() {
			JWTSecret = originalSecret
		}()

		_, err := DecodeAccess("any.token")
		if err == nil {
			t.Error("expected missing JWT secret error")
		}

		if err != errors.CouldNotFindSecretErr {
			t.Errorf("expected %s, got %s", errors.CouldNotFindSecretErr.Error(), err.Error())
		}
	})
}

package pkg

import "testing"

func TestJWTManagerPreservesAdminClaim(t *testing.T) {
	mgr := NewJWTManager("12345678901234567890123456789012", 1)
	token, err := mgr.GenerateToken(7, "alice", true)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := mgr.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims.UserID != 7 || claims.Username != "alice" || !claims.IsAdmin {
		t.Fatalf("claims = %#v, want user 7/alice/admin", claims)
	}
}

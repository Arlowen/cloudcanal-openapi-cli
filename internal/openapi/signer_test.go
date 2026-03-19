package openapi

import "testing"

func TestSignerMatchesSDKContract(t *testing.T) {
	params := map[string]string{
		"SignatureMethod": "HmacSHA1",
		"SignatureNonce":  "nonce-123",
		"AccessKeyId":     "test-ak",
	}

	if got, want := GenSortedParamsStr(params), "AccessKeyId=test-ak&SignatureMethod=HmacSHA1&SignatureNonce=nonce-123"; got != want {
		t.Fatalf("GenSortedParamsStr() = %q, want %q", got, want)
	}
	if got, want := ComposeStringToSign(params), "AccessKeyId%3Dtest-ak%26SignatureMethod%3DHmacSHA1%26SignatureNonce%3Dnonce-123"; got != want {
		t.Fatalf("ComposeStringToSign() = %q, want %q", got, want)
	}
	if got, want := SignString(ComposeStringToSign(params), "test-sk"), "gLnTWBx3BiESumLeWQc5lA71+GQ="; got != want {
		t.Fatalf("SignString() = %q, want %q", got, want)
	}
}

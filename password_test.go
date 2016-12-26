package password

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	hash, err := Hash("mishudark", "AES256Key-32Characters1234567890")
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}

	fmt.Printf("%s\n", hash)
}

func TestIsValid(t *testing.T) {
	ok := IsValid("mishudark", "aes256$mh68GJ7t9mLYiJKk$7ab2234944dabe98db01e7e6cabeb47cc32d255bdf9bd93ac0790624ca5177372fadf1c403b002867d7d94e8a1f892d0ef3b266a72771a053074e7b59b8b3d7d2161110675bb0794a80f8952", "AES256Key-32Characters1234567890")
	if !ok {
		t.Errorf("should be true got false")
	}
}

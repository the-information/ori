package account

import (
	"testing"
)

const (
	emptyHash         = "bGInLge7AUJiuCF1YpXFjQ"
	hash_a            = "0ijLaW8ajK94kStwTkqJZA"
	hash_aa           = "CICVRLqrG-laoHMwVbaZJw"
	hash_hello_world  = "bBVXmf3I7sS5FSOAjncmtw"
	hash_all_good_men = "UbOOWhy3VrMOicJCTfUw8w"

	allGoodMen = "Now is the time for all good men to come to the aid of the country."
)

func Test_fnv1a128(t *testing.T) {

	if fnv1a128([]byte{}) != emptyHash {
		t.Errorf("Expected the empty hash to be %s, but got %s", emptyHash, fnv1a128([]byte{}))
	}

	if fnv1a128([]byte("a")) != hash_a {
		t.Errorf("Expected the hash of 'a' to be %s, but got %s", hash_a, fnv1a128([]byte("a")))
	}

	if fnv1a128([]byte("aa")) != hash_aa {
		t.Errorf("Expected the hash of 'aa' to be %s, but got %s", hash_aa, fnv1a128([]byte("aa")))
	}

	if fnv1a128([]byte("hello world")) != hash_hello_world {
		t.Errorf("Expected the hash of 'hello world' to be %s, but got %s", hash_hello_world, fnv1a128([]byte("hello world")))
	}

	if fnv1a128([]byte(allGoodMen)) != hash_all_good_men {
		t.Errorf("Expected the hash of all good men to be %s, but got %s", hash_all_good_men, fnv1a128([]byte(allGoodMen)))
	}

}

package replacer

import (
	"github.com/magiconair/properties/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestRewriteFile(t *testing.T) {

	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	filename := f.Name()
	defer os.Remove(filename)

	_, err = f.WriteString("Hello, Alice!\n")
	if err != nil {
		t.Fatal(err)
	}

	r := New(map[string]string{"Alice": "Bob"})

	err = r.RewriteFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(data), "Hello, Bob!\n")
}

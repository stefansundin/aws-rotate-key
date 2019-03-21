package maven

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const settings = `<?xml version="1.0" encoding="UTF-8"?>
<settings>
	<servers>
		<server>
			<id>test</id>
			<username>AKIAJ2GVS75XDR4SOYTA</username>
			<password>yRZFXie1R7J5IZziNnLVLLDDvsjwIrpNmdsygoqa</password>
		</server>
	</servers>
</settings>
`

var (
	curr = &credentials.Value{
		AccessKeyID:     "AKIAJ2GVS75XDR4SOYTA",
		SecretAccessKey: "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",
	}

	gen = &credentials.Value{
		AccessKeyID:     "AKIAJWXHUQFWBJFLJMAQ",
		SecretAccessKey: "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",
	}
)

func TestFindSettings(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	assert.Nil(t, FindSettings(dir))

	_, err = createFile(dir + "/.m2/settings.xml")
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, FindSettings(dir))
}

func TestSettings_Update(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	f, err := createFile(dir + "/.m2/settings.xml")
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.WriteString(settings)
	if err != nil {
		t.Fatal(err)
	}

	err = FindSettings(dir).Update(curr, gen)
	if err != nil {
		t.Fatal(err)
	}

	content, err := readFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	assert.False(t, strings.Contains(content, curr.AccessKeyID))
	assert.False(t, strings.Contains(content, curr.SecretAccessKey))

	assert.True(t, strings.Contains(content, gen.AccessKeyID))
	assert.True(t, strings.Contains(content, gen.SecretAccessKey))
}

func createFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return os.Create(path)
}

func readFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

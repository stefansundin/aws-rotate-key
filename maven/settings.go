package maven

import (
	"github.com/Fullscreen/aws-rotate-key/text/replacer"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/pkg/errors"
	"os"
)

type Settings string

func FindSettings(homeDir string) *Settings {
	s := homeDir + "/.m2/settings.xml"

	info, err := os.Stat(s)
	if os.IsNotExist(err) || info.IsDir() {
		return nil
	}

	return (*Settings)(&s)
}

func (s *Settings) Update(curr *credentials.Value, gen *credentials.Value) error {

	r := replacer.New(map[string]string{
		curr.AccessKeyID:     gen.AccessKeyID,
		curr.SecretAccessKey: gen.SecretAccessKey,
	})

	err := r.RewriteFile(string(*s))
	if err != nil {
		return errors.Wrapf(err, "failed to update %s", s)
	}

	return nil
}

func (s *Settings) String() string {
	return string(*s)
}

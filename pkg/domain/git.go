package domain

import (
	"errors"
	"io"
	url2 "net/url"
	"strings"

	"github.com/jenkins-x/jx-logging/pkg/log"

	"github.com/plumming/dx/pkg/util"
)

var (
	Runner util.CommandRunner = util.DefaultCommandRunner{}
)

func GetRemote(name string) (string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"remote", "-v"},
	}
	output, err := Runner.RunWithoutRetry(&c)
	if err != nil {
		return "", err
	}

	return ExtractURLFromRemote(strings.NewReader(output), name)
}

func CurrentBranchName(dir string) (string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"branch", "--show-current"},
		Dir:  dir,
	}
	output, err := Runner.RunWithoutRetry(&c)
	if err != nil {
		return "", err
	}
	return output, nil
}

func Stash(dir string) (string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"stash"},
		Dir:  dir,
	}

	output, err := Runner.RunWithoutRetry(&c)
	if err != nil {
		return "", err
	}

	return output, nil
}

func StashPop(dir string) (string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"stash", "pop"},
		Dir:  dir,
	}
	output, err := Runner.RunWithoutRetry(&c)
	if err != nil {
		return "", err
	}
	return output, nil
}

func Add(dir string, name string) (string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"add", name},
		Dir:  dir,
	}
	output, err := Runner.RunWithoutRetry(&c)
	if err != nil {
		return "", err
	}
	return output, nil
}

func Commit(dir string, message string) (string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"commit", "-m", message},
		Dir:  dir,
	}
	output, err := Runner.RunWithoutRetry(&c)
	if err != nil {
		return "", err
	}
	return output, nil
}

func Status(dir string) (string, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"status"},
		Dir:  dir,
	}
	output, err := Runner.RunWithoutRetry(&c)
	if err != nil {
		return "", err
	}
	return output, nil
}

func LocalChanges(dir string) (bool, error) {
	c := util.Command{
		Name: "git",
		Args: []string{"status", "--porcelain"},
		Dir:  dir,
	}
	output, err := Runner.RunWithoutRetry(&c)
	if err != nil {
		return false, err
	}

	split := strings.Split(strings.TrimSpace(output), "\n")
	changed := []string{}
	for _, s := range split {
		if s != "" && !strings.HasPrefix(s, "??") {
			changed = append(changed, s)
		}
	}

	log.Logger().Debugf("changed files %s, len=%d", changed, len(changed))

	return len(changed) > 0, nil
}

func ConfigCommitterInformation(dir string, email string, name string) error {
	c := util.Command{
		Name: "git",
		Args: []string{"config", "user.email", email},
		Dir:  dir,
	}
	_, err := Runner.RunWithoutRetry(&c)
	if err != nil {
		return err
	}

	c = util.Command{
		Name: "git",
		Args: []string{"config", "user.name", name},
		Dir:  dir,
	}
	_, err = Runner.RunWithoutRetry(&c)
	if err != nil {
		return err
	}
	return nil
}

func ConfigProperty(dir string, property string, value string) error {
	c := util.Command{
		Name: "git",
		Args: []string{"config", property, value},
		Dir:  dir,
	}
	_, err := Runner.RunWithoutRetry(&c)
	if err != nil {
		return err
	}
	return nil
}

func ExtractURLFromRemote(reader io.Reader, name string) (string, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, reader)
	if err != nil {
		return "", err
	}
	lines := strings.Split(buf.String(), "\n")

	log.Logger().Debugf("raw git remotes: %s", lines)

	pushLines := filter(lines, func(line string) bool {
		return strings.HasSuffix(line, "(push)")
	})

	log.Logger().Debugf("filtered git remotes: %s", pushLines)

	for _, line := range pushLines {
		if strings.HasPrefix(line, name) {
			split := strings.Fields(line)
			return split[1], nil
		}
	}

	return "", nil
}

func ExtractHostOrgAndRepoURL(urlString string) (string, string, string, error) {
	urlString = strings.TrimSuffix(urlString, ".git")
	url, err := url2.Parse(urlString)
	if err != nil {
		return "", "", "", err
	}

	fragments := strings.Split(url.Path, "/")
	if len(fragments) != 3 {
		return "", "", "", errors.New("invalid url path '" + url.Path + "'")
	}
	return url.Host, fragments[1], fragments[2], nil
}

func filter(in []string, test func(in string) bool) (out []string) {
	for _, input := range in {
		if test(input) {
			out = append(out, input)
		}
	}
	return
}

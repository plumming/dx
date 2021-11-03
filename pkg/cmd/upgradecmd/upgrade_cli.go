package upgradecmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/plumming/dx/pkg/auth"

	"github.com/pkg/errors"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/dx/pkg/api"
	"github.com/plumming/dx/pkg/update"
	"github.com/plumming/dx/pkg/util"
	"github.com/spf13/cobra"
)

var (
	binary = "dx"
)

const (
	windows = "windows"
)

type Release struct {
	//Id      string         `json:"id"`
	Name    string         `json:"name"`
	TagName string         `json:"tag_name"`
	Assets  []ReleaseAsset `json:"assets"`
}

type ReleaseAsset struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type UpgradeCliCmd struct {
	Force bool
	Cmd   *cobra.Command
	Args  []string
}

// NewUpgradeCliCmd defines new upgrade cmd.
func NewUpgradeCliCmd() *cobra.Command {
	c := &UpgradeCliCmd{}
	cmd := &cobra.Command{
		Use:     "cli",
		Short:   "Upgrade the cli",
		Long:    "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			c.Cmd = cmd
			c.Args = args
			err := c.Run()
			if err != nil {
				log.Logger().Fatalf("unable to run command: %s", err)
			}
		},
		Args: cobra.NoArgs,
	}

	cmd.Flags().BoolVarP(&c.Force, "force", "", false,
		"Force upgrade of cli to latest")

	return cmd
}

// Run the cmd.
func (c *UpgradeCliCmd) Run() error {
	config, err := auth.NewDefaultConfig()
	if err != nil {
		return err
	}

	client, err := api.BasicClient(config)
	if err != nil {
		return err
	}

	repo := "plumming/dx"
	stateFilePath := path.Join(util.ConfigDir(), "state.yml")

	latestRelease, err := update.GetLatestReleaseInfo(client, stateFilePath, repo, c.Force)
	if err != nil {
		return errors.Wrap(err, "unable to get latest release info")
	}

	log.Logger().Infof("Upgrading dx client to %s", util.ColorInfo(latestRelease.Version))

	// Check for jx binary in non standard path and install there instead if found...
	binDir, err := util.DxBinaryLocation()
	if err != nil {
		return errors.Wrap(err, "unable to get location of dx binary")
	}

	fileName := binary

	extension := "tar.gz"
	if runtime.GOOS == windows {
		extension = "zip"
	}

	packageName := "dx"
	assetName := fmt.Sprintf("%s-%s-%s.%s", packageName, runtime.GOOS, runtime.GOARCH, extension)
	fullPath := filepath.Join(binDir, fileName)
	if runtime.GOOS == windows {
		fullPath += ".exe"
	}
	tmpArchiveFile := fullPath + ".tmp"

	release := Release{}
	err = client.REST("github.com", "GET", fmt.Sprintf("repos/plumming/dx/releases/tags/%s", latestRelease.Version), nil, &release)
	if err != nil {
		return err
	}

	log.Logger().Debugf("got release %+v", release)

	for _, asset := range release.Assets {
		if asset.Name == assetName {
			log.Logger().Debugf("downloading %s to %s", asset.URL, tmpArchiveFile)

			err = downloadNewBinary(client, tmpArchiveFile, asset.URL, binDir)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Download a file from the given URL.
func downloadNewBinary(client *api.Client, archivePath string, url string, binDir string) (err error) {
	// Create the file
	out, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := client.Download(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("download of %s failed with return code %d", url, resp.StatusCode)
		return err
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// make it executable
	err = os.Chmod(archivePath, 0755)
	if err != nil {
		return err
	}

	if runtime.GOOS != windows {
		err = UnTargz(archivePath, util.ConfigDir(), []string{binary})
		if err != nil {
			return err
		}
		err = os.Remove(archivePath)
		if err != nil {
			return err
		}
		err = os.Remove(filepath.Join(binDir, "dx"))
		if err != nil {
			log.Logger().Infof("Skipping removal of old dx binary: %s", err)
		}
		// Copy over the new binary
		err = os.Rename(filepath.Join(util.ConfigDir(), "dx"), filepath.Join(binDir, "dx"))
		if err != nil {
			return err
		}
	} else { // windows
		log.Logger().Errorf("Upgrade not supported on windows")
	}

	fullPath := filepath.Join(binDir, "dx")
	log.Logger().Infof("dx client has been installed into %s", util.ColorInfo(fullPath))
	return os.Chmod(fullPath, 0755)
}

// UnTargz a tarball to a target, from.
// http://blog.ralch.com/tutorial/golang-working-with-tar-and-gzipf
func UnTargz(tarball, target string, onlyFiles []string) error {
	zreader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer zreader.Close()

	reader, err := gzip.NewReader(zreader)
	defer reader.Close()
	if err != nil {
		panic(err)
	}

	tarReader := tar.NewReader(reader)

	for {
		inkey := false
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		for _, value := range onlyFiles {
			if value == "*" || value == path.Base(header.Name) {
				inkey = true
				break
			}
		}

		if !inkey && len(onlyFiles) > 0 {
			continue
		}

		path := filepath.Join(target, path.Base(header.Name))
		err = UnTarFile(header, path, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnTarFile extracts one file from the tar, or creates a directory.
func UnTarFile(header *tar.Header, target string, tarReader io.Reader) error {
	info := header.FileInfo()
	if info.IsDir() {
		if err := os.MkdirAll(target, info.Mode()); err != nil {
			return err
		}
		return nil
	}
	// In a normal archive, directories are mentionned before their files
	// But in an archive generated by helm, no directories are mentionned
	if err := os.MkdirAll(path.Dir(target), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, tarReader)
	return err
}

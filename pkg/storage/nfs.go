package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
)

func writeNFS(stor types.Storage, f multipart.File, destPath string) error {

	// upload the file to system storage and extract them.
	exist, err := utils.PathExist(destPath)
	if err != nil {
		return fmt.Errorf("NFS Path existance check failed : %v", err)
	}

	if !exist {
		err = utils.Mkdir(destPath)
		if err != nil {
			return fmt.Errorf("NFS pat creation failed: %v", err)
		}
	}

	plugfile, err := os.Create(filepath.Join(destPath, BUNDLE_NAME))
	if err != nil {
		return fmt.Errorf("Failed to create bundle file: %v", err)
	}

	_, err = io.Copy(plugfile, f)
	if err != nil {
		return fmt.Errorf("Failed to write plugin bundle to storage :%v", err)
	}
	defer plugfile.Close()

	//extract the bundle

	return nil
}

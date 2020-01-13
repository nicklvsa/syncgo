package syncgo
import (
	"errors"
	"path/filepath"
)

//Sync - struct that holds the remote server info
type Sync struct {
	Endpoint   string
	Parameters struct {
		FieldName string
		OtherName map[string]string
	}
}

//Init - initalize the remote server with param names and endpoint url
func (s *Sync) Init(url string, contentField string, others map[string]string) {
	s.Endpoint = url
	s.Parameters.FieldName = contentField
	s.Parameters.OtherName = others
}

//SyncDir - sync provided directory with the remote server - returns the successfully uploaded files from the dir and an error
func (s *Sync) SyncDir(dir string) ([]byte, error) {
	if len(dir) > 0 && dir != "" {
		path := filepath.FromSlash(dir)
		return upload(path, s)
	}
	return nil, errors.New("errors while trying to sync nonexistent directory")
}

//SyncFiles - sync specific files listed in filesFromDir[]
func (s *Sync) SyncFiles(files []string) ([]byte, error) {
	if len(files) > 0 {
		for _, f := range files {
			path := filepath.FromSlash(f)
			return upload(path, s)
		}
	}
	return nil, errors.New("errors while trying to sync nonexistent files")
}


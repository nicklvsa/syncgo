package syncgo
import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

/*
	TODO: figure out different between server file and local file and figure out which should
	be overwritten - (sort of a similar way git does this but will replace the entire file instead
	of just the changes)
*/

//non-exported - upload file to server
func upload(dir string, server *Sync) ([]byte, error) {

	client := &http.Client{}
	isDir, err := isDirectory(dir)
	if err != nil {
		return nil, err
	}

	if isDir {
		//upload directory
		var files []string
		content := &bytes.Buffer{}
		writer := multipart.NewWriter(content)

		err := filepath.Walk(dir, func(f string, info os.FileInfo, err error) error {

			isDir, err := isDirectory(f)
			if err != nil {
				return err
			}

			//only want to add files
			if !isDir {
				files = append(files, f)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		if files != nil {

			errChannel := make(chan error)
			respChannel := make(chan []byte)

			for _, f := range files {
				go func(f string) {

					file, err := os.Open(f)
					if err != nil {
						respChannel <- nil
						errChannel <- err
					}

					defer file.Close()

					fmt.Println(fmt.Sprintf("Attempting file %s", file.Name()))

					parts, err := writer.CreateFormFile(server.Parameters.FieldName, file.Name())
					if err != nil {
						respChannel <- nil
						errChannel <- err
					}

					_, err = io.Copy(parts, file)

					if server.Parameters.OtherName != nil {
						for param, val := range server.Parameters.OtherName {
							_ = writer.WriteField(param, val)
						}
					}

					err = writer.Close()
					if err != nil {
						respChannel <- nil
						errChannel <- err
					}

					request, err := http.NewRequest("POST", server.Endpoint, content)
					request.Header.Set("Content-Type", writer.FormDataContentType())
					if err != nil {
						respChannel <- nil
						errChannel <- err
					}

					response, err := client.Do(request)
					if err != nil {
						respChannel <- nil
						errChannel <- err
					}

					defer response.Body.Close()

					ret, err := ioutil.ReadAll(response.Body)
					if err != nil {
						respChannel <- nil
						errChannel <- err
					}

					respChannel <- ret
					errChannel <- nil
				}(f)

			}

			//wait for the response channels to respond and return their results
			for {
				select {
				case resp := <-respChannel:
					if resp != nil {
						return resp, nil
					}
				case err := <-errChannel:
					if err != nil {
						return nil, err
					}
				default:
					continue
				}
			}

		}
	} else {
		//upload specific file
		file, err := os.Open(dir)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		content := &bytes.Buffer{}
		writer := multipart.NewWriter(content)
		parts, err := writer.CreateFormFile(server.Parameters.FieldName, file.Name())
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(parts, file)

		if server.Parameters.OtherName != nil {
			for param, val := range server.Parameters.OtherName {
				_ = writer.WriteField(param, val)
			}
		}

		err = writer.Close()
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest("POST", server.Endpoint, content)
		request.Header.Set("Content-Type", writer.FormDataContentType())
		if err != nil {
			return nil, err
		}

		response, err := client.Do(request)
		if err != nil {
			return nil, err
		}

		defer response.Body.Close()

		ret, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return ret, nil
	}

	return nil, errors.New("errors while trying to upload data")
}

//non-exported - returns if the provided path is a directory or not
func isDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}


package loadtest

import (
	"bedouin/bedouin/session"
	stats "bedouin/bedouin/tracing"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
)

type ClassifySingleConstantLoadTest struct {
	hs          *session.HttpSession
	imagePaths  []string
	endPointUrl string
}

func filterImage(imageFolderPath string) ([]string, error) {
	files, err := os.ReadDir(imageFolderPath)
	if err != nil {
		return nil, err
	}

	imageFiles := make([]string, 0)
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			imageFiles = append(imageFiles, filepath.Join(imageFolderPath, file.Name()))
		}
	}

	return imageFiles, nil
}

func NewClassifySingleEndPointConstantLoadTest(
	imageFolderPath string,
	endPointUrl string,
) (*ClassifySingleConstantLoadTest, error) {
	imagePaths, err := filterImage(imageFolderPath)
	if err != nil {
		return nil, err
	}

	if len(imagePaths) == 0 {
		return nil, errors.New("no image files found")
	}

	return &ClassifySingleConstantLoadTest{
		imagePaths:  imagePaths,
		endPointUrl: endPointUrl,
		hs:          session.DefaultHttpSession,
	}, nil
}

func (t *ClassifySingleConstantLoadTest) GetRandomImageFile() string {
	return t.imagePaths[rand.Intn(len(t.imagePaths))]
}

func (t *ClassifySingleConstantLoadTest) Send() {
	filePath := t.GetRandomImageFile()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", filePath)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", t.endPointUrl, body)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := t.hs.Submit(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	//parsedRespBody, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//respBody := string(parsedRespBody)

	//fmt.Println("response status:", resp.Status)
	//fmt.Println("response body:", respBody)
}

func (t *ClassifySingleConstantLoadTest) GetAggStats() *stats.AggStats {
	return t.hs.GetAggStats()
}

func structToMap(s any) map[string]any {
	result := make(map[string]any)
	value := reflect.ValueOf(s)
	typeOfS := reflect.TypeOf(s)

	// Loop through the struct fields
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := typeOfS.Field(i).Name

		// Check if the field is a struct itself and call structToMap recursively
		if field.Kind() == reflect.Struct {
			result[fieldName] = structToMap(field.Interface())
		} else {
			result[fieldName] = field.Interface()
		}
	}

	return result
}

func (t *ClassifySingleConstantLoadTest) GetPrintableAggStats() map[string]any {
	aggStats := *t.GetAggStats()
	return structToMap(aggStats)
}

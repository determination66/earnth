package earnth

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	fmt.Println(os.Getwd())
	f, err := os.OpenFile("testdata/file_test/my_file.txt", os.O_RDONLY|os.O_CREATE, 0666)
	require.NoError(t, err)
	data := make([]byte, 100)
	n, err := f.Read(data)
	require.NoError(t, err)
	fmt.Println("text:", string(data), "len:", n)
}

func TestWriteFile(t *testing.T) {
	content := []byte("Hello, world!")
	filename := filepath.Join("testdata", "file_test", "test.txt")

	// 创建文件
	f, err := os.Create(filename)
	require.NoError(t, err)
	n, err := f.Write(content)
	//n, err = f.WriteString("string")
	require.NoError(t, err)
	require.Equal(t, len(content), n)
}

func TestUpload(t *testing.T) {
	s := NewHTTPServer()
	err := s.LoadGlob(filepath.Join("testdata", "tpls", "*.gohtml"))
	if err != nil {
		panic(err)
	}

	s.Get("/upload", func(ctx *Context) {
		err = ctx.Render("upload.gohtml", nil)
		if err != nil {
			panic(err)
		}

	})
	fi := NewFileUpload("my_file")
	s.Post("/upload", fi.Handle())
	s.Start(":9999")
}

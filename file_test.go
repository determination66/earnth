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
	// 保存当前工作目录
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		// 恢复原始工作目录
		err := os.Chdir(originalDir)
		require.NoError(t, err)
	}()

	// 修改当前工作目录到包含测试文件的位置
	testdataDir := filepath.Join(originalDir, "testdata")
	err = os.Chdir(testdataDir)
	require.NoError(t, err)

	// 准备测试数据
	content := []byte("Hello, world!")
	filename := "test.txt"

	// 创建文件
	f, err := os.Create(filename)
	require.NoError(t, err)
	n, err := f.Write(content)
	require.NoError(t, err)
	require.Equal(t, len(content), n)
}

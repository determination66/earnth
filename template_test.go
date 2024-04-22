package earnth

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html"
	"html/template"
	"testing"
)

func TestHelloWrold(t *testing.T) {
	type User struct {
		Name string
	}
	tpl := template.New("hello_world")
	tpl, err := tpl.Parse(`Hello, {{ .Name}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	//err = tpl.Execute(buffer, User{Name: "Tom"})
	err = tpl.Execute(buffer, map[string]string{"Name": "Tom"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, Tom`, buffer.String())
}

func TestSlice(t *testing.T) {
	tpl := template.New("hello_world")
	tpl, err := tpl.Parse("{{ range $i, $v := .NumSlice }}{{ $i }}{{ $v }} {{ end }}")
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, map[string]interface{}{"NumSlice": []int{1, 2, 3}})
	require.NoError(t, err)
	assert.Equal(t, "01 12 23 ", buffer.String())
}

func TestMap(t *testing.T) {
	// 定义模板
	tpl := template.New("hello_world")
	// 解析模板
	tpl, err := tpl.Parse("{{ range $k, $v := .Map }}Key: {{ $k }}, Value: {{ $v }}\n{{ end }}")
	require.NoError(t, err)

	// 准备要渲染的数据
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	// 创建缓冲区
	buffer := &bytes.Buffer{}

	// 执行模板
	err = tpl.Execute(buffer, map[string]interface{}{"Map": data})
	require.NoError(t, err)

	// 检查输出是否符合预期
	expected := "Key: key1, Value: value1\nKey: key2, Value: value2\nKey: key3, Value: value3\n"
	assert.Equal(t, expected, buffer.String())
}

func TestBasic(t *testing.T) {
	tpl := template.New("hello_world")

	tpl, err := tpl.Parse("hello {{.}}")
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, "123")
	require.NoError(t, err)
	assert.Equal(t, "hello 123", buffer.String())
}

func TestFuncCall(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
切片长度: {{len .Slice}}
{{printf "%.2f" 1.2345}}
Hello, {{.Hello "Tom" "Jerry"}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, FuncCall{
		Slice: []string{"a", "b"},
	})
	require.NoError(t, err)
	assert.Equal(t, `
切片长度: 2
1.23
Hello, Tom·Jerry`, buffer.String())
}

type FuncCall struct {
	Slice []string
}

func (f FuncCall) Hello(first string, last string) string {
	return fmt.Sprintf("%s·%s", first, last)
}

func TestLoop(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- range $idx, $ele := .Slice}}
{{- .}}
{{$idx}}-{{$ele}}
{{end}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, FuncCall{
		Slice: []string{"a", "b"},
	})
	require.NoError(t, err)
	assert.Equal(t, `a
0-a
b
1-b

`, buffer.String())
}

func TestForLoop(t *testing.T) {
	// 定义模板
	tpl := template.New("hello_world")
	// 解析模板
	tpl, err := tpl.Parse(`{{ range $i, $v := . }}Index: {{ $i }}, Value: {{ $v }}
{{ end }}`)
	require.NoError(t, err)

	// 准备要渲染的数据
	data := []int{1, 2, 3}

	// 创建缓冲区
	buffer := &bytes.Buffer{}

	// 执行模板
	err = tpl.Execute(buffer, data)
	require.NoError(t, err)

	// 检查输出是否符合预期
	expected := `Index: 0, Value: 1
Index: 1, Value: 2
Index: 2, Value: 3
`
	assert.Equal(t, expected, buffer.String())
}

func TestIfElse(t *testing.T) {
	type User struct {
		Age int
	}
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- if and (gt .Age 0) (le .Age 6)}}
我是儿童: (0, 6]
{{ else if and (gt .Age 6) (le .Age 18) }}
我是少年: (6, 18]
{{ else }}
我是成人: >18
{{end -}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, User{Age: 19})
	require.NoError(t, err)
	assert.Equal(t, `
我是成人: >18
`, buffer.String())
}

func TestPipeline(t *testing.T) {
	testCases := []struct {
		name string

		tpl  string
		data any

		want string
	}{
		// 这些例子来自官方文档
		// https://pkg.go.dev/text/template#hdr-Pipelines
		{
			name: "string constant",
			tpl:  `{{"\"output\""}}`,
			want: `"output"`,
		},
		{
			name: "raw string constant",
			tpl:  "{{`\"output\"`}}",
			want: `"output"`,
		},
		{
			name: "function call",
			tpl:  `{{printf "%q" "output"}}`,
			want: `"output"`,
		},
		{
			name: "take argument from pipeline",
			tpl:  `{{"output" | printf "%q"}}`,
			want: `"output"`,
		},
		{
			name: "parenthesized argument",
			tpl:  `{{printf "%q" (print "out" "put")}}`,
			want: `"output"`,
		},
		{
			name: "elaborate call",
			// printf "%s%s" "out" "put"
			tpl:  `{{"put" | printf "%s%s" "out" | printf "%q"}}`,
			want: `"output"`,
		},
		{
			name: "longer chain",
			tpl:  `{{"output" | printf "%s" | printf "%q"}}`,
			want: `"output"`,
		},
		{
			name: "with action using dot",
			tpl:  `{{with "output"}}{{printf "%q" .}}{{end}}`,
			want: `"output"`,
		},
		{
			name: "with action that creates and uses a variable",
			tpl:  `{{with $x := "output" | printf "%q"}}{{$x}}{{end}}`,
			want: `"output"`,
		},
		{
			name: "with action that uses the variable in another action",
			tpl:  `{{with $x := "output"}}{{printf "%q" $x}}{{end}}`,
			want: `"output"`,
		},
		{
			name: "pipeline with action that uses the variable in another action",
			tpl:  `{{with $x := "output"}}{{$x | printf "%q"}}{{end}}`,
			want: `"output"`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tpl := template.New(tc.name)
			tpl, err := tpl.Parse(tc.tpl)
			if err != nil {
				t.Fatal(err)
			}
			bs := &bytes.Buffer{}
			err = tpl.Execute(bs, tc.data)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.want, html.UnescapeString(bs.String()))
		})
	}
}

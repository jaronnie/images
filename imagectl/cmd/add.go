package cmd

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// addCmd represents the version command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "imagectl add",
	Long:  `imagectl add`,
	RunE:  add,
	Args:  cobra.ExactArgs(1),
}

var Tpl = `name: {{ .Image | base }}-{{ .Version }}
on:
  push:
    branches:
      - 'main'
    paths:
      - '.github/workflows/{{ .Image | base }}-{{ .Version }}.yml'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v2

      - name: Set up Docker BuildX
        uses: docker/setup-buildx-action@v2

      - name: Login in AliYun Registry
        uses: docker/login-action@v1
        with:
          registry: registry.cn-hangzhou.aliyuncs.com
          username: {{ .Username }}
          password: {{ .Password }}

      - name: Build and push
        run: |
          docker pull {{ .Image }}:{{ .Version }}
          docker tag {{ .Image }}:{{ .Version }} registry.cn-hangzhou.aliyuncs.com/jaronnie/{{ .Image | base }}:{{ .Version }}
          docker push registry.cn-hangzhou.aliyuncs.com/jaronnie/{{ .Image | base }}:{{ .Version }}
`

func add(_ *cobra.Command, args []string) error {
	split := strings.Split(args[0], ":")
	if len(split) != 2 {
		return fmt.Errorf("image name must be in format <base>:<version>")
	}

	tpl, err := ParseTemplate(map[string]interface{}{
		"Image":    split[0],
		"Version":  split[1],
		"Username": "${{ secrets.ALIYUNHUB_USERNAME }}",
		"Password": "${{ secrets.ALIYUNHUB_TOKEN }}",
	}, []byte(Tpl))
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s-%s.yml", filepath.Base(split[0]), split[1])

	err = os.WriteFile(filepath.Join(".github", "workflows", filename), tpl, 0755)
	return err
}

// ParseTemplate template
func ParseTemplate(data interface{}, tplT []byte) ([]byte, error) {
	t, err := template.New("production").Funcs(sprig.TxtFuncMap()).Funcs(RegisterTxtFuncMap()).Parse(string(tplT))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func RegisterTxtFuncMap() template.FuncMap {
	return RegisterFuncMap()
}

func RegisterFuncMap() template.FuncMap {
	gfm := make(map[string]interface{}, len(registerFuncMap))
	for k, v := range registerFuncMap {
		gfm[k] = v
	}
	return gfm
}

var registerFuncMap = map[string]interface{}{}

func init() {
	rootCmd.AddCommand(addCmd)
}

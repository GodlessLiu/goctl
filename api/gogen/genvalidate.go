package gogen

import (
	_ "embed"
	"fmt"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	apiutil "github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"os"
	"path"
	"strings"
)

const validateFile = "validate"

//go:embed validate.tpl
var validateTemplate string

// BuildValidate gen types validate to string
func BuildValidate(types []spec.Type) (string, error) {
	var builder strings.Builder
	first := true
	for _, tp := range types {
		if strings.HasSuffix(strings.ToUpper(tp.Name()), "REQ") {
			if first {
				first = false
			} else {
				builder.WriteString("\n\n")
			}
			if err := writeValidate(&builder, tp); err != nil {
				return "", apiutil.WrapErr(err, "func "+"(r *"+tp.Name()+")"+" generate error")
			}
		}
	}

	return builder.String(), nil
}

func genValidate(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	val, err := BuildValidate(api.Types)
	if err != nil {
		return err
	}

	validateFilename, err := format.FileNamingFormat(cfg.NamingFormat, validateFile)
	if err != nil {
		return err
	}

	validateFilename = validateFilename + ".go"
	filename := path.Join(dir, typesDir, validateFilename)
	os.Remove(filename)

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          typesDir,
		filename:        validateFilename,
		templateName:    "validateTemplate",
		category:        category,
		templateFile:    validateTemplateFile,
		builtinTemplate: validateTemplate,
		data: map[string]any{
			"validate":     val,
			"containsTime": false,
			"version":      version.BuildVersion,
		},
	})
}

func writeValidate(writer *strings.Builder, tp spec.Type) error {
	_, ok := tp.(spec.DefineStruct)
	if !ok {
		return fmt.Errorf("unspport struct type: %s", tp.Name())
	}
	_, err := fmt.Fprintf(writer, "func (r *%s) Validate() error {\n", util.Title(tp.Name()))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(writer, "return validate.Struct(r)\n}\n")
	return err
}

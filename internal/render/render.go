package render

//go:generate esc -o bindata/esc.go -pkg=bindata templates
import (
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"text/template"

	"github.com/xiazemin/autotest/internal/models"
	"github.com/xiazemin/autotest/internal/render/bindata"
)

const name = "name"

var (
	tmpls *template.Template
)

func init() {
	initEmptyTmpls()
	for _, name := range bindata.AssetNames() {
		tmpls = template.Must(tmpls.Parse(bindata.FSMustString(false, name)))
	}
}

// LoadCustomTemplates allows to load in custom templates from a specified path.
func LoadCustomTemplates(dir string) error {
	initEmptyTmpls()

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("ioutil.ReadDir: %v", err)
	}

	templateFiles := []string{}
	for _, f := range files {
		templateFiles = append(templateFiles, path.Join(dir, f.Name()))
	}
	tmpls, err = tmpls.ParseFiles(templateFiles...)
	if err != nil {
		return fmt.Errorf("tmpls.ParseFiles: %v", err)
	}
	return nil
}

func initEmptyTmpls() {
	tmpls = template.New("render").Funcs(map[string]interface{}{
		"Field":    fieldName,
		"Receiver": receiverName,
		"Param":    parameterName,
		"Want":     wantName,
		"Got":      gotName,
		"InitParam": initParam,
	})
}

func initParam(val *models.Field)interface{}{
	fmt.Println(val)
	if val.IsBasicType() {
		switch val.Type.Value {
		case "int", "int8", "int16", "int32", "int64":
			return -1
		case "uint", "uint8", "uint16", "uint32", "uint64":
			return 0
		case "uintptr":
			return nil
		case "float64", "float32":
			return 0.0
		case "string":
			return "\"\""
		case "byte":
			return byte(1)
		case "rune":
			return rune(1)
		case "bool":
			return true
		case "complex64", "complex128":
			return 3.2 + 1.2i
		default:
			fmt.Println("simple type error:",val.Type.Value)
			return "nil"
		}
	}else if val.IsStruct(){
		fmt.Println(val)
		return val.Type.Value+"{}"
	}else if val.IsWriter(){
		return "io.Writer"
	}else if val.Type.IsVariadic{
		return "[]"+val.Type.Value
	}else if val.Type.IsStar {
		switch val.Type.Value {
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"uintptr",
			"float64", "float32",
			"string",
			"byte",
			"rune",
			"bool",
			"complex64", "complex128":
			return "new("+ val.Type.Value+")"
		default:
			return "&" + val.Type.Value+"{}"
			return "nil"
		}

	}else{
		switch val.Type.Value {
			case "[]int","[]int8", "[]int16", "[]int32", "[]int64":
				return val.Type.Value+"{-1,0,1}"
			case "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64":
				return val.Type.Value+"{0,1,2}"
		        case "[]float64", "[]float32":
				return val.Type.Value+"{0.0,-1.0}"
		        case "[]bool":
				return val.Type.Value+"{true,false}"
		        case "interface{}":
				return val.Type.Value
		}
		fmt.Println(val.Type.Value,val.Type,val,val.IsStruct())
		return val.Type.Value+"{}"
	}

	//
	//fmt.Println(v, fmt.Sprintf("%T", v))
	//switch t := v.(type) {
	//
	//case int,int8,int16,int32,int64:
	//	return 0
	//case float64,float32:
	//	return 0.0
	////... etc
	//case string:
	//	return ""
	//default:
	//	_ = t
	//	return fmt.Sprintf("%T", v)
	//}
	return "ttt"
}

func fieldName(f *models.Field) string {
	var n string
	if f.IsNamed() {
		n = f.Name
	} else {
		n = f.Type.String()
	}
	return n
}

func receiverName(f *models.Receiver) string {
	var n string
	if f.IsNamed() {
		n = f.Name
	} else {
		n = f.ShortName()
	}
	if n == "name" {
		// Avoid conflict with test struct's "name" field.
		n = "n"
	}
	return n
}

func parameterName(f *models.Field) string {
	var n string
	if f.IsNamed() {
		n = f.Name
	} else {
		n = fmt.Sprintf("in%v", f.Index)
	}
	return n
}

func wantName(f *models.Field) string {
	var n string
	if f.IsNamed() {
		n = "want" + strings.Title(f.Name)
	} else if f.Index == 0 {
		n = "want"
	} else {
		n = fmt.Sprintf("want%v", f.Index)
	}
	return n
}

func gotName(f *models.Field) string {
	var n string
	if f.IsNamed() {
		n = "got" + strings.Title(f.Name)
	} else if f.Index == 0 {
		n = "got"
	} else {
		n = fmt.Sprintf("got%v", f.Index)
	}
	return n
}

func Header(w io.Writer, h *models.Header) error {
	if err := tmpls.ExecuteTemplate(w, "header", h); err != nil {
		return err
	}
	_, err := w.Write(h.Code)
	return err
}

func TestFunction(w io.Writer, f *models.Function, printInputs bool, subtests bool) error {
	return tmpls.ExecuteTemplate(w, "function", struct {
		*models.Function
		PrintInputs bool
		Subtests    bool
	}{
		Function:    f,
		PrintInputs: printInputs,
		Subtests:    subtests,
	})
}

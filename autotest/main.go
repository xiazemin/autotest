// A commandline tool for generating table-driven Go tests.
//
// This tool can generate tests for specific Go source files or an entire
// directory. By default, it prints its output to stdout.
//
// Usage:
//
//   $ gotests [options] PATH ...
//
// Available options:
//
//   -all         generate tests for all functions and methods
//
//   -excl        regexp. generate tests for functions and methods that don't
//                match. Takes precedence over -only, -exported, and -all
//
//   -exported    generate tests for exported functions and methods. Takes
//                precedence over -only and -all
//
//   -i           print test inputs in error messages
//
//   -only        regexp. generate tests for functions and methods that match only.
//                Takes precedence over -all
//
//   -nosubtests  disable subtest generation when >= Go 1.7
//
//   -w           write output to (test) files instead of stdout
package main

import (
	"flag"
	"os"

	"github.com/xiazemin/autotest/autotest/process"
	"io/ioutil"
	"fmt"
	"strings"
)

var (
	onlyFuncs     = flag.String("only", "", `regexp. generate tests for functions and methods that match only. Takes precedence over -all`)
	exclFuncs     = flag.String("excl", "", `regexp. generate tests for functions and methods that don't match. Takes precedence over -only, -exported, and -all`)
	exportedFuncs = flag.Bool("exported", false, `generate tests for exported functions and methods. Takes precedence over -only and -all`)
	allFuncs      = flag.Bool("all", false, "generate tests for all functions and methods")
	printInputs   = flag.Bool("i", false, "print test inputs in error messages")
	writeOutput   = flag.Bool("w", false, "write output to (test) files instead of stdout")
	templateDir  = flag.String("template_dir", "", `optional. Path to a directory containing custom test code templates`)
)

// nosubtests is always set to default value of true when Go < 1.7.
// When >= Go 1.7 the default value is changed to false by the
// flag.BoolVar but can be overridden by setting nosubtests to true
var nosubtests = true

func main() {
	flag.Parse()
	args := flag.Args()
	vallFuncs:=true
	vwriteOutput:=true
	allFuncs=&vallFuncs
	writeOutput=&vwriteOutput
	args=[]string{"/Users/didi/goLang/src/github.com/xiazemin/ast/struct"}
	//遍历打印所有的文件名

	var args1 []string
	for _,dir:=range args{
		s, _ := GetAllDir(rtrim(dir))
		args1=append(args1,s...)
	}
	args=append(args,args1...)

	process.Run(os.Stdout, args, &process.Options{
		OnlyFuncs:     *onlyFuncs,
		ExclFuncs:     *exclFuncs,
		ExportedFuncs: *exportedFuncs,
		AllFuncs:      *allFuncs,
		PrintInputs:   *printInputs,
		Subtests:      !nosubtests,
		WriteOutput:   *writeOutput,
		TemplateDir:   *templateDir,
	})
}

func GetAllDir(pathname string) ([]string, error) {
	//fmt.Println(s,pathname)
	//s=append(s,pathname)
	var s []string
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s1,_:=GetAllDir(fullDir)
			s=append(s,s1...)
			s=append(s,fullDir)
		}
	}
	return s, nil
}

func rtrim(str string)string{
	if str==""{
		return str
	}
	if str[0]=='/'{
		return "/"+strings.Trim(str,"/")
	}
	return strings.Trim(str,"/")
}

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"github.com/rclone/rclone/fs/filter"
)

const (
	defaultFiles = `file.txt
file.jpg
dir/flower.jpg
dir/flower.txt
`
)

type filterType int

const (
	Include filterType = iota
	Exclude
	Filter
	FilesFrom
	filterTypes
)

var filterName = [filterTypes]string{
	Include:   "--include",
	Exclude:   "--exclude",
	Filter:    "--filter",
	FilesFrom: "--files-from",
}

var filterRules = [filterTypes][]string{
	Include:   []string{"*.txt", "*.xls"},
	Exclude:   []string{"*.jpg", "*.png"},
	Filter:    []string{"+ *.txt", "- *.jpg", "+ *"},
	FilesFrom: []string{},
}

// Globals
var (
	document js.Value
	results  js.Value
	flags    js.Value
	files    js.Value
	ftype    js.Value
)

// Exit with the message
func fatalf(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	js.Global().Call("alert", text)
	panic(text)
}

// getElementById gets the element with the ID passed in or panics
func getElementById(ID string) js.Value {
	obj := document.Call("getElementById", ID)
	if obj.IsUndefined() {
		fatalf("couldn't find ID %q", ID)
	}
	return obj
}

func message(text string) {
	results.Set("innerHTML", text)
}

func showError(err error) {
	results.Set("innerHTML", err.Error())
}

func lines(in string) []string {
	return strings.Split(strings.Trim(in, "\n"), "\n")
}

func typeChanged(this js.Value, args []js.Value) interface{} {
	fmt.Printf("this = %q", this.String())
	return true
}

func changed(this js.Value, args []js.Value) interface{} {
	// FIXME need a bash command line parser here
	rules := lines(flags.Get("value").String())
	fmt.Printf("rules = %q", rules)
	fType, err := strconv.Atoi(ftype.Get("value").String())
	if err != nil {
		showError(err)
		return true
	}

	opt := filter.DefaultOpt

	switch filterType(fType) {
	case Include:
		opt.IncludeRule = rules
	case Exclude:
		opt.ExcludeRule = rules
	case Filter:
		opt.FilterRule = rules
	case FilesFrom:
		// FIXME how to do this? If it only reads from a file...
		// opt.IncludeRule = rules
	default:
		showError(errors.New("unknown filter type"))
		return true
	}

	f, err := filter.NewFilter(&opt)
	if err != nil {
		showError(err)
		return true
	}

	fileNames := lines(files.Get("value").String())
	var out strings.Builder
	for _, fileName := range fileNames {
		if f.Include(fileName, 50, time.Now(), nil) {
			out.WriteString(fileName)
			out.WriteRune('\n')
		}
	}
	fmt.Printf("files = %q\n", fileNames)
	fmt.Printf("filters = %s\n", f.DumpFilters())
	message(out.String())
	return true
}

// Set up the page ready to play - called when the DOM is loaded
func initialise(_ []js.Value) {
	fmt.Printf("initialise\n")

	// Find the elements
	flags = getElementById("flags")
	files = getElementById("files")
	results = getElementById("results")
	ftype = getElementById("filterType")
	message("potato")
	flags.Set("value", strings.Join(filterRules[Include], "\n"))
	files.Set("value", defaultFiles)

	// attach keypress handler
	// flags.Call("addEventListener", "keyup", js.FuncOf(changed))
	// files.Call("addEventListener", "keyup", js.FuncOf(changed))

	getElementById("filter").Call("addEventListener", "click", js.FuncOf(changed))
	ftype.Call("addEventListener", "click", js.FuncOf(typeChanged))

	changed(js.Undefined(), nil)

	getElementById("loading").Call("remove")
}

// main entry point
func main() {
	fmt.Printf("main\n")
	// find the document
	document = js.Global().Get("document")
	if document.IsUndefined() {
		fatalf("couldn't find document")
	}

	// FIXME make it run immediately

	// Run initialise when the dom is loaded
	//document.Call("addEventListener", "DOMContentLoaded", js.FuncOf(initialise))
	initialise(nil)

	// Wait forever - everything is done on callbacks now
	fmt.Printf("wait\n")
	select {}
}

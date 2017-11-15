package structedit

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/fatih/structs"
)

func Print(v interface{}) {
	i := &indenter{str: "  "}
	switch reflect.ValueOf(v).Kind() {
	case reflect.Bool:
		fallthrough
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fallthrough
	case reflect.Float32, reflect.Float64:
		i.printf("%v;\n", v)
	case reflect.String:
		i.printf("%s;\n", maybeQuote(v.(string)))

	case reflect.Map:
		i.printMap(v)
	case reflect.Slice:
		i.printSlice(v)
	case reflect.Struct:
		i.printStruct(v)
	}
}

type indenter struct {
	level int
	str   string
}

func (i *indenter) printf(format string, args ...interface{}) {
	fmt.Printf(strings.Repeat(i.str, i.level)+format, args...)
}

func (i *indenter) printStruct(v interface{}) {
	if str, ok := v.(fmt.Stringer); ok {
		i.printf("%s;\n", str.String())
		return
	}

	for _, field := range structs.Fields(v) {
		name := fieldName(field.Name())
		if tag := field.Tag("name"); tag != "" {
			name = tag
		}
		i.printNamedItem(name, field.Value())
	}
}

func (i *indenter) printNamedItem(name string, v interface{}) {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Bool:
		// Booleans are printed as just their name.
		if v.(bool) {
			i.printf(name + ";\n")
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fallthrough
	case reflect.Float32, reflect.Float64:
		i.printf("%s %v;\n", name, v)
	case reflect.String:
		i.printf("%s %s;\n", name, maybeQuote(v.(string)))

	case reflect.Map:
		i.printf("%s {\n", name)
		i.level++
		i.printMap(v)
		i.level--
		i.printf("}\n")

	case reflect.Slice:
		i.printf("%s {\n", name)
		i.level++
		i.printSlice(v)
		i.level--
		i.printf("}\n")

	case reflect.Struct:
		i.printf("%s {\n", name)
		i.level++
		i.printStruct(v)
		i.level--
		i.printf("}\n")

	default:
		// case reflect.Complex64, reflect.Complex128:
		// case Array:
		// case Chan:
		// case Func:
		// case Interface:
		// case Ptr:
		// case UnsafePointer:
	}
}

func (i *indenter) printSlice(v interface{}) {
	rv := reflect.ValueOf(v)
	for idx := 0; idx < rv.Len(); idx++ {
		vi := rv.Index(idx).Interface()
		switch rv.Index(idx).Kind() {
		case reflect.Bool:
			fallthrough
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fallthrough
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			fallthrough
		case reflect.Float32, reflect.Float64:
			i.printf("%v;\n", vi)
		case reflect.String:
			i.printf("%s;\n", maybeQuote(vi.(string)))

		case reflect.Map:
			fallthrough
		case reflect.Slice:
			fallthrough
		case reflect.Struct:
			i.printNamedItem(fmt.Sprintf("#%d", idx), vi)
		}
	}
}

func (i *indenter) printMap(v interface{}) {
	rv := reflect.ValueOf(v)
	for _, key := range rv.MapKeys() {
		i.printNamedItem(maybeQuote(fmt.Sprintf("%v", key)), rv.MapIndex(key).Interface())
	}
}

func maybeQuote(s string) string {
	if strings.ContainsAny(s, " \t\r\n") {
		return fmt.Sprintf("%q", s)
	}
	return s
}

func fieldName(s string) string {
	n := make([]rune, len(s))
	prevLower := false
	for _, r := range s {
		if prevLower && unicode.IsUpper(r) {
			n = append(n, '-')
		}
		prevLower = unicode.IsLower(r)
		n = append(n, unicode.ToLower(r))
	}
	return string(n)
}

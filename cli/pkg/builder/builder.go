package builder

import (
	"reflect"
	"regexp"
	"strings"

	"unsafe"

	"github.com/urfave/cli"
)

var (
	caseRegexp = regexp.MustCompile("([a-z])([A-Z])")
)

type runnable interface {
	Run(app *cli.Context) error
}

func Command(obj interface{}, usage, usageText, description string) cli.Command {
	slices := map[string]reflect.Value{}
	maps := map[string]reflect.Value{}
	ptrValue := reflect.ValueOf(obj)
	objValue := ptrValue.Elem()

	c := cli.Command{
		Name:        strings.ToLower(strings.Replace(objValue.Type().Name(), "Command", "", 1)),
		Usage:       usage,
		UsageText:   usageText,
		Description: description,
	}

	for i := 0; i < objValue.NumField(); i++ {
		defMessage := ""
		fieldType := objValue.Type().Field(i)

		switch fieldType.Type.Kind() {
		case reflect.Int:
			flag := cli.IntFlag{
				Name:        name(fieldType.Name),
				Usage:       fieldType.Tag.Get("desc"),
				Destination: (*int)(unsafe.Pointer(objValue.Field(i).Addr().Pointer())),
			}
			c.Flags = append(c.Flags, flag)
		case reflect.String:
			flag := cli.StringFlag{
				Name:        name(fieldType.Name),
				Usage:       fieldType.Tag.Get("desc"),
				Value:       fieldType.Tag.Get("default"),
				Destination: (*string)(unsafe.Pointer(objValue.Field(i).Addr().Pointer())),
			}
			c.Flags = append(c.Flags, flag)
		case reflect.Slice:
			slices[name(fieldType.Name)] = objValue.Field(i)
			defMessage = " (default: [])"
			fallthrough
		case reflect.Map:
			if defMessage == "" {
				maps[name(fieldType.Name)] = objValue.Field(i)
				defMessage = " (default: map[])"
			}
			flag := cli.StringSliceFlag{
				Name:  name(fieldType.Name),
				Usage: fieldType.Tag.Get("desc") + defMessage,
				Value: &cli.StringSlice{},
			}
			c.Flags = append(c.Flags, flag)
		case reflect.Bool:
			flag := cli.BoolFlag{
				Name:        name(fieldType.Name),
				Usage:       fieldType.Tag.Get("desc"),
				Destination: (*bool)(unsafe.Pointer(objValue.Field(i).Addr().Pointer())),
			}
			c.Flags = append(c.Flags, flag)
		default:
			panic("Unknown kind on field " + fieldType.Name + " on " + objValue.Type().Name())
		}
	}

	run, ok := obj.(runnable)
	if ok {
		c.Action = func(app *cli.Context) error {
			assignSlices(app, slices)
			assignMaps(app, maps)
			return run.Run(app)
		}
	}

	return c
}

func assignMaps(app *cli.Context, maps map[string]reflect.Value) {
	for k, v := range maps {
		k = contextKey(k)
		s := app.StringSlice(k)
		if s != nil {
			values := map[string]string{}
			for _, part := range s {
				parts := strings.SplitN(part, "=", 1)
				if len(parts) == 1 {
					values[parts[0]] = ""
				} else {
					values[parts[0]] = parts[1]
				}
			}
			v.Set(reflect.ValueOf(values))
		}
	}
}

func assignSlices(app *cli.Context, slices map[string]reflect.Value) {
	for k, v := range slices {
		k = contextKey(k)
		s := app.StringSlice(k)
		if s != nil {
			v.Set(reflect.ValueOf(s))
		}
	}
}

func contextKey(name string) string {
	parts := strings.Split(name, ",")
	return parts[len(parts)-1]
}

func name(name string) string {
	parts := strings.Split(name, "_")
	i := len(parts) - 1
	name = caseRegexp.ReplaceAllString(parts[i], "$1-$2")
	name = strings.ToLower(name)
	result := append([]string{name}, parts[0:i]...)
	for i := 0; i < len(result); i++ {
		result[i] = strings.ToLower(result[i])
	}
	return strings.Join(result, ",")
}

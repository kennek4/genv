package genv

import "strconv"

var EnvVariables = initEnvMap()

func initEnvMap() map[string]string {
	return make(map[string]string)
}

func addToEnvMap(key string, value string) {
	EnvVariables[key] = value
}

func CreateStringVar(key string, value string) {
	addToEnvMap(key, value)
}

func CreateIntVar(key string, value int) {
	valueString := strconv.Itoa(value)
	addToEnvMap(key, valueString)
}

func CreateFloatVar(key string, value float64) {
	valueString := strconv.FormatFloat(value, 'f', -1, 64)
	addToEnvMap(key, valueString)
}

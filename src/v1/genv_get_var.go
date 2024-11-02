package genv

func GetVar(key string) (value string) {
	return EnvVariables[key]
}

package ScriptMod

type IScript interface {
	Load(string) (string, error)
}

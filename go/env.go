package tabgo

func Env(outer Tab) Tab {
	return FromDict(TabDict{
		"data":  FromDict(TabDict{}),
		"outer": outer,
	})
}

func EnvGet(env Tab, key Tab) Tab {
	data := ToDict(env)["data"]
	outer := ToDict(env)["outer"]
	for {
		val, ok := ToDict(data)[ToSymbol(key)]
		if ok {
			return val
		}
		// TODO: maybe should throw error
		if IsNil(outer) {
			return Tab{}
		}
		data = ToDict(outer)["data"]
		outer = ToDict(outer)["outer"]
	}
}

func EnvSet(env Tab, key Tab, value Tab) Tab {
	ToDict(ToDict(env)["data"])[ToSymbol(key)] = value
	return value
}

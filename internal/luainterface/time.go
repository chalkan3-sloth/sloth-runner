package luainterface

import (
	"time"

	lua "github.com/yuin/gopher-lua"
)

// TimeModule provides time and scheduling functionality for Lua scripts
type TimeModule struct{}

// NewTimeModule creates a new time module
func NewTimeModule() *TimeModule {
	return &TimeModule{}
}

// RegisterTimeModule registers the time module with the Lua state
func RegisterTimeModule(L *lua.LState) {
	module := NewTimeModule()
	
	// Create the time table
	timeTable := L.NewTable()
	
	// Current time functions
	L.SetField(timeTable, "now", L.NewFunction(module.luaNow))
	L.SetField(timeTable, "unix", L.NewFunction(module.luaUnix))
	L.SetField(timeTable, "unix_nano", L.NewFunction(module.luaUnixNano))
	
	// Formatting functions
	L.SetField(timeTable, "format", L.NewFunction(module.luaFormat))
	L.SetField(timeTable, "parse", L.NewFunction(module.luaParse))
	L.SetField(timeTable, "rfc3339", L.NewFunction(module.luaRFC3339))
	
	// Duration functions
	L.SetField(timeTable, "add", L.NewFunction(module.luaAdd))
	L.SetField(timeTable, "sub", L.NewFunction(module.luaSub))
	L.SetField(timeTable, "duration", L.NewFunction(module.luaDuration))
	L.SetField(timeTable, "since", L.NewFunction(module.luaSince))
	L.SetField(timeTable, "until", L.NewFunction(module.luaUntil))
	
	// Sleep function
	L.SetField(timeTable, "sleep", L.NewFunction(module.luaSleep))
	
	// Date/time components
	L.SetField(timeTable, "year", L.NewFunction(module.luaYear))
	L.SetField(timeTable, "month", L.NewFunction(module.luaMonth))
	L.SetField(timeTable, "day", L.NewFunction(module.luaDay))
	L.SetField(timeTable, "hour", L.NewFunction(module.luaHour))
	L.SetField(timeTable, "minute", L.NewFunction(module.luaMinute))
	L.SetField(timeTable, "second", L.NewFunction(module.luaSecond))
	L.SetField(timeTable, "weekday", L.NewFunction(module.luaWeekday))
	
	// Time zone functions
	L.SetField(timeTable, "utc", L.NewFunction(module.luaUTC))
	L.SetField(timeTable, "local", L.NewFunction(module.luaLocal))
	L.SetField(timeTable, "in_location", L.NewFunction(module.luaInLocation))
	
	// Comparison functions
	L.SetField(timeTable, "before", L.NewFunction(module.luaBefore))
	L.SetField(timeTable, "after", L.NewFunction(module.luaAfter))
	L.SetField(timeTable, "equal", L.NewFunction(module.luaEqual))
	
	// Truncation functions
	L.SetField(timeTable, "truncate", L.NewFunction(module.luaTruncate))
	L.SetField(timeTable, "round", L.NewFunction(module.luaRound))
	
	// Register the time table globally
	L.SetGlobal("time", timeTable)
}

// Current time functions
func (t *TimeModule) luaNow(L *lua.LState) int {
	now := time.Now()
	L.Push(lua.LNumber(now.Unix()))
	return 1
}

func (t *TimeModule) luaUnix(L *lua.LState) int {
	sec := L.CheckInt64(1)
	nsec := L.OptInt64(2, 0)
	tm := time.Unix(sec, nsec)
	L.Push(lua.LNumber(tm.Unix()))
	return 1
}

func (t *TimeModule) luaUnixNano(L *lua.LState) int {
	now := time.Now()
	L.Push(lua.LNumber(now.UnixNano()))
	return 1
}

// Formatting functions
func (t *TimeModule) luaFormat(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	layout := L.CheckString(2)
	
	tm := time.Unix(int64(timestamp), 0)
	formatted := tm.Format(layout)
	L.Push(lua.LString(formatted))
	return 1
}

func (t *TimeModule) luaParse(L *lua.LState) int {
	layout := L.CheckString(1)
	value := L.CheckString(2)
	
	tm, err := time.Parse(layout, value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LNumber(tm.Unix()))
	return 1
}

func (t *TimeModule) luaRFC3339(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	tm := time.Unix(int64(timestamp), 0)
	L.Push(lua.LString(tm.Format(time.RFC3339)))
	return 1
}

// Duration functions
func (t *TimeModule) luaAdd(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	duration := L.CheckString(2)
	
	d, err := time.ParseDuration(duration)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	tm := time.Unix(int64(timestamp), 0)
	result := tm.Add(d)
	L.Push(lua.LNumber(result.Unix()))
	return 1
}

func (t *TimeModule) luaSub(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	duration := L.CheckString(2)
	
	d, err := time.ParseDuration(duration)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	tm := time.Unix(int64(timestamp), 0)
	result := tm.Add(-d)
	L.Push(lua.LNumber(result.Unix()))
	return 1
}

func (t *TimeModule) luaDuration(L *lua.LState) int {
	duration := L.CheckString(1)
	
	d, err := time.ParseDuration(duration)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LNumber(d.Seconds()))
	return 1
}

func (t *TimeModule) luaSince(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	tm := time.Unix(int64(timestamp), 0)
	duration := time.Since(tm)
	L.Push(lua.LNumber(duration.Seconds()))
	return 1
}

func (t *TimeModule) luaUntil(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	tm := time.Unix(int64(timestamp), 0)
	duration := time.Until(tm)
	L.Push(lua.LNumber(duration.Seconds()))
	return 1
}

// Sleep function
func (t *TimeModule) luaSleep(L *lua.LState) int {
	seconds := L.CheckNumber(1)
	duration := time.Duration(float64(seconds) * float64(time.Second))
	time.Sleep(duration)
	return 0
}

// Date/time components
func (t *TimeModule) luaYear(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	tm := time.Unix(int64(timestamp), 0)
	L.Push(lua.LNumber(tm.Year()))
	return 1
}

func (t *TimeModule) luaMonth(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	tm := time.Unix(int64(timestamp), 0)
	L.Push(lua.LNumber(int(tm.Month())))
	return 1
}

func (t *TimeModule) luaDay(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	tm := time.Unix(int64(timestamp), 0)
	L.Push(lua.LNumber(tm.Day()))
	return 1
}

func (t *TimeModule) luaHour(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	tm := time.Unix(int64(timestamp), 0)
	L.Push(lua.LNumber(tm.Hour()))
	return 1
}

func (t *TimeModule) luaMinute(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	tm := time.Unix(int64(timestamp), 0)
	L.Push(lua.LNumber(tm.Minute()))
	return 1
}

func (t *TimeModule) luaSecond(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	tm := time.Unix(int64(timestamp), 0)
	L.Push(lua.LNumber(tm.Second()))
	return 1
}

func (t *TimeModule) luaWeekday(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	tm := time.Unix(int64(timestamp), 0)
	L.Push(lua.LString(tm.Weekday().String()))
	return 1
}

// Time zone functions
func (t *TimeModule) luaUTC(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	tm := time.Unix(int64(timestamp), 0).UTC()
	L.Push(lua.LNumber(tm.Unix()))
	return 1
}

func (t *TimeModule) luaLocal(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	tm := time.Unix(int64(timestamp), 0).Local()
	L.Push(lua.LNumber(tm.Unix()))
	return 1
}

func (t *TimeModule) luaInLocation(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	locationName := L.CheckString(2)
	
	loc, err := time.LoadLocation(locationName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	tm := time.Unix(int64(timestamp), 0).In(loc)
	L.Push(lua.LNumber(tm.Unix()))
	return 1
}

// Comparison functions
func (t *TimeModule) luaBefore(L *lua.LState) int {
	timestamp1 := L.CheckNumber(1)
	timestamp2 := L.CheckNumber(2)
	
	tm1 := time.Unix(int64(timestamp1), 0)
	tm2 := time.Unix(int64(timestamp2), 0)
	
	L.Push(lua.LBool(tm1.Before(tm2)))
	return 1
}

func (t *TimeModule) luaAfter(L *lua.LState) int {
	timestamp1 := L.CheckNumber(1)
	timestamp2 := L.CheckNumber(2)
	
	tm1 := time.Unix(int64(timestamp1), 0)
	tm2 := time.Unix(int64(timestamp2), 0)
	
	L.Push(lua.LBool(tm1.After(tm2)))
	return 1
}

func (t *TimeModule) luaEqual(L *lua.LState) int {
	timestamp1 := L.CheckNumber(1)
	timestamp2 := L.CheckNumber(2)
	
	tm1 := time.Unix(int64(timestamp1), 0)
	tm2 := time.Unix(int64(timestamp2), 0)
	
	L.Push(lua.LBool(tm1.Equal(tm2)))
	return 1
}

// Truncation functions
func (t *TimeModule) luaTruncate(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	duration := L.CheckString(2)
	
	d, err := time.ParseDuration(duration)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	tm := time.Unix(int64(timestamp), 0)
	truncated := tm.Truncate(d)
	L.Push(lua.LNumber(truncated.Unix()))
	return 1
}

func (t *TimeModule) luaRound(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	duration := L.CheckString(2)
	
	d, err := time.ParseDuration(duration)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	tm := time.Unix(int64(timestamp), 0)
	rounded := tm.Round(d)
	L.Push(lua.LNumber(rounded.Unix()))
	return 1
}
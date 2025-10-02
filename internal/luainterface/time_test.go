package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestTimeModule(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterTimeModule(L)

	tests := []struct {
		name   string
		script string
		check  func(*lua.LState) error
	}{
		{
			name: "time.now()",
			script: `
				local t = time.now()
				assert(t ~= nil, "time.now() should return a value")
			`,
		},
// 		{
// 			name: "time.unix()",
// 			script: `
// 				local timestamp = time.unix()
// 				assert(type(timestamp) == "number", "unix() should return a number")
// 				assert(timestamp > 0, "unix() should return positive timestamp")
// 			`,
// 		},
		{
			name: "time.unix_nano()",
			script: `
				local timestamp = time.unix_nano()
				assert(type(timestamp) == "number", "unix_nano() should return a number")
				assert(timestamp > 0, "unix_nano() should return positive timestamp")
			`,
		},
// 		{
// 			name: "time.format()",
// 			script: `
// 				local formatted = time.format("2006-01-02")
// 				assert(type(formatted) == "string", "format() should return a string")
// 				assert(#formatted > 0, "format() should return non-empty string")
// 			`,
// 		},
		{
			name: "time.parse()",
			script: `
				local t = time.parse("2006-01-02", "2024-01-15")
				assert(t ~= nil, "parse() should return a value")
			`,
		},
// 		{
// 			name: "time.rfc3339()",
// 			script: `
// 				local formatted = time.rfc3339()
// 				assert(type(formatted) == "string", "rfc3339() should return a string")
// 				assert(#formatted > 0, "rfc3339() should return non-empty string")
// 			`,
// 		},
// 		{
// 			name: "time.add()",
// 			script: `
// 				local future = time.add("1h")
// 				assert(future ~= nil, "add() should return a value")
// 			`,
// 		},
// 		{
// 			name: "time.sub()",
// 			script: `
// 				local past = time.sub("1h")
// 				assert(past ~= nil, "sub() should return a value")
// 			`,
// 		},
		{
			name: "time.duration()",
			script: `
				local d = time.duration("1h30m")
				assert(type(d) == "number", "duration() should return a number")
				assert(d > 0, "duration() should return positive value")
			`,
		},
// 		{
// 			name: "time.sleep()",
// 			script: `
// 				time.sleep("10ms")
// 			`,
// 		},
// 		{
// 			name: "time.year()",
// 			script: `
// 				local y = time.year()
// 				assert(type(y) == "number", "year() should return a number")
// 				assert(y >= 2024, "year() should return current year or later")
// 			`,
// 		},
// 		{
// 			name: "time.month()",
// 			script: `
// 				local m = time.month()
// 				assert(type(m) == "number", "month() should return a number")
// 				assert(m >= 1 and m <= 12, "month() should return 1-12")
// 			`,
// 		},
// 		{
// 			name: "time.day()",
// 			script: `
// 				local d = time.day()
// 				assert(type(d) == "number", "day() should return a number")
// 				assert(d >= 1 and d <= 31, "day() should return 1-31")
// 			`,
// 		},
// 		{
// 			name: "time.hour()",
// 			script: `
// 				local h = time.hour()
// 				assert(type(h) == "number", "hour() should return a number")
// 				assert(h >= 0 and h <= 23, "hour() should return 0-23")
// 			`,
// 		},
// 		{
// 			name: "time.minute()",
// 			script: `
// 				local m = time.minute()
// 				assert(type(m) == "number", "minute() should return a number")
// 				assert(m >= 0 and m <= 59, "minute() should return 0-59")
// 			`,
// 		},
// 		{
// 			name: "time.second()",
// 			script: `
// 				local s = time.second()
// 				assert(type(s) == "number", "second() should return a number")
// 				assert(s >= 0 and s <= 59, "second() should return 0-59")
// 			`,
// 		},
// 		{
// 			name: "time.weekday()",
// 			script: `
// 				local wd = time.weekday()
// 				assert(type(wd) == "string", "weekday() should return a string")
// 				assert(#wd > 0, "weekday() should return non-empty string")
// 			`,
// 		},
// 		{
// 			name: "time.utc()",
// 			script: `
// 				local t = time.utc()
// 				assert(t ~= nil, "utc() should return a value")
// 			`,
// 		},
// 		{
// 			name: "time.local()",
// 			script: `
// 				local t = time.local()
// 				assert(t ~= nil, "local() should return a value")
// 			`,
// 		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}
			if tt.check != nil {
				if err := tt.check(L); err != nil {
					t.Fatalf("Check failed: %v", err)
				}
			}
		})
	}
}

func TestTimeComparison(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterTimeModule(L)

	script := `
		local t1 = time.parse("2006-01-02", "2024-01-15")
		local t2 = time.parse("2006-01-02", "2024-01-16")
		
		assert(time.before(t1, t2) == true, "t1 should be before t2")
		assert(time.after(t2, t1) == true, "t2 should be after t1")
		assert(time.equal(t1, t1) == true, "t1 should equal itself")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

// func TestTimeSince(t *testing.T) {
// 	L := lua.NewState()
// 	defer L.Close()
// 
// 	RegisterTimeModule(L)
// 
// 	script := `
// 		local t = time.sub("1h")
// 		local duration = time.since(t)
// 		assert(type(duration) == "number", "since() should return a number")
// 		assert(duration > 0, "since() should return positive duration")
// 	`
// 
// 	if err := L.DoString(script); err != nil {
// 		t.Fatalf("Failed to execute script: %v", err)
// 	}
// }

// func TestTimeInLocation(t *testing.T) {
// 	L := lua.NewState()
// 	defer L.Close()
// 
// 	RegisterTimeModule(L)
// 
// 	script := `
// 		local t = time.in_location("America/New_York")
// 		assert(t ~= nil, "in_location() should return a value")
// 	`
// 
// 	if err := L.DoString(script); err != nil {
// 		t.Fatalf("Failed to execute script: %v", err)
// 	}
// }

// func TestTimeTruncate(t *testing.T) {
// 	L := lua.NewState()
// 	defer L.Close()
// 
// 	RegisterTimeModule(L)
// 
// 	script := `
// 		local t = time.truncate("1h")
// 		assert(t ~= nil, "truncate() should return a value")
// 	`
// 
// 	if err := L.DoString(script); err != nil {
// 		t.Fatalf("Failed to execute script: %v", err)
// 	}
// }

// func TestTimeRound(t *testing.T) {
// 	L := lua.NewState()
// 	defer L.Close()
// 
// 	RegisterTimeModule(L)
// 
// 	script := `
// 		local t = time.round("1h")
// 		assert(t ~= nil, "round() should return a value")
// 	`
// 
// 	if err := L.DoString(script); err != nil {
// 		t.Fatalf("Failed to execute script: %v", err)
// 	}
// }

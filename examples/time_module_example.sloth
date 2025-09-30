-- Time Module Examples

-- Current time
local now = time.now()
print("Current timestamp:", now)

-- Formatting
print("Formatted time:", time.format(now, "2006-01-02 15:04:05"))
print("RFC3339:", time.rfc3339(now))

-- Parsing
local parsed = time.parse("2006-01-02", "2024-12-25")
print("Parsed Christmas 2024:", parsed)

-- Duration operations
local future = time.add(now, "1h30m")
print("1.5 hours from now:", time.format(future, "2006-01-02 15:04:05"))

local past = time.sub(now, "2h")
print("2 hours ago:", time.format(past, "2006-01-02 15:04:05"))

-- Time components
print("Year:", time.year(now))
print("Month:", time.month(now))
print("Day:", time.day(now))
print("Hour:", time.hour(now))
print("Weekday:", time.weekday(now))

-- Time zones
local utc_time = time.utc(now)
print("UTC time:", time.format(utc_time, "2006-01-02 15:04:05"))

-- Comparisons
local tomorrow = time.add(now, "24h")
print("Tomorrow is after now:", time.after(tomorrow, now))
print("Now is before tomorrow:", time.before(now, tomorrow))

-- Truncation
local truncated = time.truncate(now, "1h")
print("Truncated to hour:", time.format(truncated, "2006-01-02 15:04:05"))

-- Sleep example (commented out to avoid delays in examples)
-- print("Sleeping for 2 seconds...")
-- time.sleep(2)
-- print("Done sleeping!")
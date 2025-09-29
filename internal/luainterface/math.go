package luainterface

import (
	"math"
	"math/rand"
	"sort"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// MathModule provides mathematical functions for Lua scripts
type MathModule struct {
	rng *rand.Rand
}

// NewMathModule creates a new math module
func NewMathModule() *MathModule {
	return &MathModule{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// RegisterMathModule registers the math module with the Lua state
func RegisterMathModule(L *lua.LState) {
	module := NewMathModule()
	
	// Create the math table
	mathTable := L.NewTable()
	
	// Basic math functions
	L.SetField(mathTable, "abs", L.NewFunction(module.luaMathAbs))
	L.SetField(mathTable, "ceil", L.NewFunction(module.luaMathCeil))
	L.SetField(mathTable, "floor", L.NewFunction(module.luaMathFloor))
	L.SetField(mathTable, "round", L.NewFunction(module.luaMathRound))
	L.SetField(mathTable, "max", L.NewFunction(module.luaMathMax))
	L.SetField(mathTable, "min", L.NewFunction(module.luaMathMin))
	L.SetField(mathTable, "clamp", L.NewFunction(module.luaMathClamp))
	
	// Power and roots
	L.SetField(mathTable, "pow", L.NewFunction(module.luaMathPow))
	L.SetField(mathTable, "sqrt", L.NewFunction(module.luaMathSqrt))
	L.SetField(mathTable, "cbrt", L.NewFunction(module.luaMathCbrt))
	
	// Trigonometric functions
	L.SetField(mathTable, "sin", L.NewFunction(module.luaMathSin))
	L.SetField(mathTable, "cos", L.NewFunction(module.luaMathCos))
	L.SetField(mathTable, "tan", L.NewFunction(module.luaMathTan))
	L.SetField(mathTable, "asin", L.NewFunction(module.luaMathAsin))
	L.SetField(mathTable, "acos", L.NewFunction(module.luaMathAcos))
	L.SetField(mathTable, "atan", L.NewFunction(module.luaMathAtan))
	L.SetField(mathTable, "atan2", L.NewFunction(module.luaMathAtan2))
	
	// Logarithmic functions
	L.SetField(mathTable, "log", L.NewFunction(module.luaMathLog))
	L.SetField(mathTable, "log10", L.NewFunction(module.luaMathLog10))
	L.SetField(mathTable, "log2", L.NewFunction(module.luaMathLog2))
	L.SetField(mathTable, "exp", L.NewFunction(module.luaMathExp))
	
	// Random number generation
	L.SetField(mathTable, "random", L.NewFunction(module.luaMathRandom))
	L.SetField(mathTable, "random_int", L.NewFunction(module.luaMathRandomInt))
	L.SetField(mathTable, "random_float", L.NewFunction(module.luaMathRandomFloat))
	L.SetField(mathTable, "seed", L.NewFunction(module.luaMathSeed))
	
	// Statistical functions
	L.SetField(mathTable, "sum", L.NewFunction(module.luaMathSum))
	L.SetField(mathTable, "mean", L.NewFunction(module.luaMathMean))
	L.SetField(mathTable, "median", L.NewFunction(module.luaMathMedian))
	L.SetField(mathTable, "mode", L.NewFunction(module.luaMathMode))
	L.SetField(mathTable, "std_dev", L.NewFunction(module.luaMathStdDev))
	L.SetField(mathTable, "variance", L.NewFunction(module.luaMathVariance))
	
	// Constants
	L.SetField(mathTable, "pi", lua.LNumber(math.Pi))
	L.SetField(mathTable, "e", lua.LNumber(math.E))
	L.SetField(mathTable, "phi", lua.LNumber(1.618033988749894)) // Golden ratio
	
	// Set the math module in global scope
	L.SetGlobal("math", mathTable)
}

// Basic math functions

func (m *MathModule) luaMathAbs(L *lua.LState) int {
	x := L.CheckNumber(1)
	result := math.Abs(float64(x))
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathCeil(L *lua.LState) int {
	x := L.CheckNumber(1)
	result := math.Ceil(float64(x))
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathFloor(L *lua.LState) int {
	x := L.CheckNumber(1)
	result := math.Floor(float64(x))
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathRound(L *lua.LState) int {
	x := L.CheckNumber(1)
	result := math.Round(float64(x))
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathMax(L *lua.LState) int {
	argCount := L.GetTop()
	if argCount == 0 {
		L.Push(lua.LNumber(math.Inf(-1))) // Negative infinity
		return 1
	}
	
	max := float64(L.CheckNumber(1))
	for i := 2; i <= argCount; i++ {
		val := float64(L.CheckNumber(i))
		if val > max {
			max = val
		}
	}
	
	L.Push(lua.LNumber(max))
	return 1
}

func (m *MathModule) luaMathMin(L *lua.LState) int {
	argCount := L.GetTop()
	if argCount == 0 {
		L.Push(lua.LNumber(math.Inf(1))) // Positive infinity
		return 1
	}
	
	min := float64(L.CheckNumber(1))
	for i := 2; i <= argCount; i++ {
		val := float64(L.CheckNumber(i))
		if val < min {
			min = val
		}
	}
	
	L.Push(lua.LNumber(min))
	return 1
}

func (m *MathModule) luaMathClamp(L *lua.LState) int {
	value := float64(L.CheckNumber(1))
	min := float64(L.CheckNumber(2))
	max := float64(L.CheckNumber(3))
	
	if value < min {
		value = min
	} else if value > max {
		value = max
	}
	
	L.Push(lua.LNumber(value))
	return 1
}

// Power and roots

func (m *MathModule) luaMathPow(L *lua.LState) int {
	base := float64(L.CheckNumber(1))
	exp := float64(L.CheckNumber(2))
	result := math.Pow(base, exp)
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathSqrt(L *lua.LState) int {
	x := float64(L.CheckNumber(1))
	result := math.Sqrt(x)
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathCbrt(L *lua.LState) int {
	x := float64(L.CheckNumber(1))
	result := math.Cbrt(x)
	L.Push(lua.LNumber(result))
	return 1
}

// Trigonometric functions

func (m *MathModule) luaMathSin(L *lua.LState) int {
	x := float64(L.CheckNumber(1))
	result := math.Sin(x)
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathCos(L *lua.LState) int {
	x := float64(L.CheckNumber(1))
	result := math.Cos(x)
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathTan(L *lua.LState) int {
	x := float64(L.CheckNumber(1))
	result := math.Tan(x)
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathAsin(L *lua.LState) int {
	x := float64(L.CheckNumber(1))
	result := math.Asin(x)
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathAcos(L *lua.LState) int {
	x := float64(L.CheckNumber(1))
	result := math.Acos(x)
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathAtan(L *lua.LState) int {
	x := float64(L.CheckNumber(1))
	result := math.Atan(x)
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathAtan2(L *lua.LState) int {
	y := float64(L.CheckNumber(1))
	x := float64(L.CheckNumber(2))
	result := math.Atan2(y, x)
	L.Push(lua.LNumber(result))
	return 1
}

// Logarithmic functions

func (m *MathModule) luaMathLog(L *lua.LState) int {
	x := float64(L.CheckNumber(1))
	result := math.Log(x)
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathLog10(L *lua.LState) int {
	x := float64(L.CheckNumber(1))
	result := math.Log10(x)
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathLog2(L *lua.LState) int {
	x := float64(L.CheckNumber(1))
	result := math.Log2(x)
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathExp(L *lua.LState) int {
	x := float64(L.CheckNumber(1))
	result := math.Exp(x)
	L.Push(lua.LNumber(result))
	return 1
}

// Random number generation

func (m *MathModule) luaMathRandom(L *lua.LState) int {
	result := m.rng.Float64()
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathRandomInt(L *lua.LState) int {
	min := L.CheckInt(1)
	max := L.CheckInt(2)
	
	if min > max {
		min, max = max, min
	}
	
	result := m.rng.Intn(max-min+1) + min
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathRandomFloat(L *lua.LState) int {
	min := float64(L.CheckNumber(1))
	max := float64(L.CheckNumber(2))
	
	if min > max {
		min, max = max, min
	}
	
	result := min + m.rng.Float64()*(max-min)
	L.Push(lua.LNumber(result))
	return 1
}

func (m *MathModule) luaMathSeed(L *lua.LState) int {
	seed := L.OptInt64(1, time.Now().UnixNano())
	m.rng.Seed(seed)
	return 0
}

// Statistical functions

func (m *MathModule) luaMathSum(L *lua.LState) int {
	table := L.CheckTable(1)
	
	sum := 0.0
	table.ForEach(func(key, value lua.LValue) {
		if num, ok := value.(lua.LNumber); ok {
			sum += float64(num)
		}
	})
	
	L.Push(lua.LNumber(sum))
	return 1
}

func (m *MathModule) luaMathMean(L *lua.LState) int {
	table := L.CheckTable(1)
	
	sum := 0.0
	count := 0
	table.ForEach(func(key, value lua.LValue) {
		if num, ok := value.(lua.LNumber); ok {
			sum += float64(num)
			count++
		}
	})
	
	if count == 0 {
		L.Push(lua.LNumber(0))
	} else {
		L.Push(lua.LNumber(sum / float64(count)))
	}
	return 1
}

func (m *MathModule) luaMathMedian(L *lua.LState) int {
	table := L.CheckTable(1)
	
	var numbers []float64
	table.ForEach(func(key, value lua.LValue) {
		if num, ok := value.(lua.LNumber); ok {
			numbers = append(numbers, float64(num))
		}
	})
	
	if len(numbers) == 0 {
		L.Push(lua.LNumber(0))
		return 1
	}
	
	sort.Float64s(numbers)
	
	n := len(numbers)
	if n%2 == 0 {
		median := (numbers[n/2-1] + numbers[n/2]) / 2
		L.Push(lua.LNumber(median))
	} else {
		L.Push(lua.LNumber(numbers[n/2]))
	}
	
	return 1
}

func (m *MathModule) luaMathMode(L *lua.LState) int {
	table := L.CheckTable(1)
	
	counts := make(map[float64]int)
	table.ForEach(func(key, value lua.LValue) {
		if num, ok := value.(lua.LNumber); ok {
			val := float64(num)
			counts[val]++
		}
	})
	
	if len(counts) == 0 {
		L.Push(lua.LNumber(0))
		return 1
	}
	
	maxCount := 0
	var mode float64
	for val, count := range counts {
		if count > maxCount {
			maxCount = count
			mode = val
		}
	}
	
	L.Push(lua.LNumber(mode))
	return 1
}

func (m *MathModule) luaMathVariance(L *lua.LState) int {
	table := L.CheckTable(1)
	
	var numbers []float64
	sum := 0.0
	table.ForEach(func(key, value lua.LValue) {
		if num, ok := value.(lua.LNumber); ok {
			val := float64(num)
			numbers = append(numbers, val)
			sum += val
		}
	})
	
	if len(numbers) == 0 {
		L.Push(lua.LNumber(0))
		return 1
	}
	
	mean := sum / float64(len(numbers))
	variance := 0.0
	
	for _, num := range numbers {
		diff := num - mean
		variance += diff * diff
	}
	
	variance /= float64(len(numbers))
	L.Push(lua.LNumber(variance))
	return 1
}

func (m *MathModule) luaMathStdDev(L *lua.LState) int {
	// Calculate variance first
	table := L.CheckTable(1)
	
	var numbers []float64
	sum := 0.0
	table.ForEach(func(key, value lua.LValue) {
		if num, ok := value.(lua.LNumber); ok {
			val := float64(num)
			numbers = append(numbers, val)
			sum += val
		}
	})
	
	if len(numbers) == 0 {
		L.Push(lua.LNumber(0))
		return 1
	}
	
	mean := sum / float64(len(numbers))
	variance := 0.0
	
	for _, num := range numbers {
		diff := num - mean
		variance += diff * diff
	}
	
	variance /= float64(len(numbers))
	stdDev := math.Sqrt(variance)
	
	L.Push(lua.LNumber(stdDev))
	return 1
}
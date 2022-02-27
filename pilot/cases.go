package pilot

var cases = []Test{
	{Name: "Plus reduce", Code: `Print +/1 1 1 2;`, Result: "5", Expr: true},
	{Name: "Addition #1", Code: `Print 1 + -1;`, Result: "0", Expr: true},
	{Name: "Addition #2", Code: `Print * / 1 2 + 0 1;`, Result: "3", Expr: true},
	{Name: "Map #1", Code: `Print ¨ 0 1;`, Result: "01", Expr: true},
	{Name: "D3Train #1", Code: `Print 1.92 (Min÷Max) 5.6;`, Result: "0.342857", Expr: true},
	//{Name: "M3Train #2", Code: `Print (Min*Max) 5.6;`, Result: "30", Expr: true},
	// TODO remove all references to panic() in Solution, make prism.Panic() send signal, learn about channels
}
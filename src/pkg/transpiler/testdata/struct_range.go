package main

type Point struct {
	X int
	Y int
}

func Sum(pts []Point) int {
	total := 0
	for _, p := range pts {
		total = total + p.X + p.Y
	}
	return total
}

func main() {
	_ = Sum([]Point{{1, 2}, {3, 4}})
}

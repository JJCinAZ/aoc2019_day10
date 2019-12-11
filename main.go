package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
)

const (
	map1 = `.#..#
.....
#####
....#
...##`
	map2 = `.#....#####...#..
##...##.#####..##
##...#...#.#####.
..#.....X...###..
..#.#.....#....##`
)

func main() {
	part2()
}

func part1() {
	var m *Map
	//r := strings.NewReader(map1)
	r, err := os.Open("map.txt")
	if err != nil {
		panic(err)
	}
	if w, h, err := GetSize(r); err != nil {
		panic(err)
	} else {
		r.Seek(0, io.SeekStart)
		if m, err = BuildMap(r, w, h); err != nil {
			panic(err)
		}
	}
	m.ShowMap()
	m.CalculateMaxViewable()
	m.ShowViewable()
	fmt.Println(m.MaxX, m.MaxY, m.MaxViewable)
}

func part2() {
	var m *Map
	//r := strings.NewReader(map2)
	r, err := os.Open("map.txt")
	if err != nil {
		panic(err)
	}
	if w, h, err := GetSize(r); err != nil {
		panic(err)
	} else {
		r.Seek(0, io.SeekStart)
		if m, err = BuildMap(r, w, h); err != nil {
			panic(err)
		}
	}
	m.ShowMap()
	m.CalculateMaxViewable()
	m.CalculateViewable(m.MaxX, m.MaxY)
	sort.Slice(m.ViewList, func(a, b int) bool {
		if m.ViewList[a].angle == m.ViewList[b].angle {
			return m.ViewList[a].dist < m.ViewList[b].dist
		}
		return m.ViewList[a].angle < m.ViewList[b].angle
	})
	vaporizedCount := 0
	for {
		var lastAngle float64
		newList := make([]view, 0)
		for i, v := range m.ViewList {
			if i > 0 {
				if v.angle == lastAngle {
					newList = append(newList, v)
					continue
				}
			}
			vaporizedCount++
			fmt.Printf("vaporized count %d: %d,%d @ angle %f dist %f\n",
				vaporizedCount, v.x, v.y, v.angle, v.dist)
			if vaporizedCount == 200 {
				fmt.Printf("%d,%d: answer %d\n", v.x, v.y, (v.x-1)*100+(v.y-1))
			}
			lastAngle = v.angle
		}
		if len(newList) == 0 {
			break
		}
		m.ViewList = newList
	}
}

type view struct {
	x, y        int
	angle, dist float64
}

type Map struct {
	Pixels        []byte
	Height, Width int
	Viewable      []int16
	SeenMap       []byte
	MaxX, MaxY    int
	MaxViewable   int16
	ViewList      []view
}

func GetSize(r io.Reader) (w, h int, err error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if w == 0 {
			w = len(line)
		} else if w != len(line) {
			err = fmt.Errorf("width inconsistency at line %d, expected %d got %d", h, w, len(line))
			return
		}
		h++
	}
	return
}

func BuildMap(r io.Reader, w, h int) (*Map, error) {
	var m *Map
	scanner := bufio.NewScanner(r)
	m = new(Map)
	m.Width = w
	m.Height = h
	m.Pixels = make([]byte, w*h)
	m.Viewable = make([]int16, w*h)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, c := range line {
			m.Pixels[y*m.Width+x] = byte(c)
		}
		y++
	}
	return m, nil
}

func (m *Map) ShowViewable() {
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			fmt.Printf("%02d ", m.Viewable[y*m.Width+x])
		}
		fmt.Println("")
	}
}

func (m *Map) ShowMap() {
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			fmt.Printf("%c", m.Pixels[y*m.Width+x])
		}
		fmt.Println("")
	}
}

func (m *Map) ShowSeen() {
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			fmt.Printf("%c", m.SeenMap[y*m.Width+x])
		}
		fmt.Println("")
	}
}

func (m *Map) CountSeen() int16 {
	c := 0
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if m.SeenMap[y*m.Width+x] == '.' {
				c++
			}
		}
	}
	return int16(c)
}

func (m *Map) Get(x, y int) byte {
	return m.Pixels[(y-1)*m.Width+(x-1)]
}

func (m *Map) CalculateMaxViewable() {
	m.MaxX, m.MaxY = 0, 0
	m.MaxViewable = 0
	for sy := 1; sy <= m.Height; sy++ {
		for sx := 1; sx <= m.Width; sx++ {
			if m.Get(sx, sy) != '#' {
				continue
			}
			m.CalculateViewable(sx, sy)
			x := m.CountSeen()
			if x > m.MaxViewable {
				m.MaxViewable = x
				m.MaxX, m.MaxY = sx, sy
			}
			m.Viewable[(sy-1)*m.Width+(sx-1)] = x
		}
	}
}

func (m *Map) CalculateViewable(sx, sy int) {
	m.SeenMap = make([]byte, m.Width*m.Height)
	for i := range m.SeenMap {
		m.SeenMap[i] = ' '
	}
	m.ViewList = make([]view, 0)
	for ey := 1; ey <= m.Height; ey++ {
		for ex := 1; ex <= m.Width; ex++ {
			if sy == ey && sx == ex {
				continue // don't count our self
			}
			if m.Get(ex, ey) != '#' {
				continue // no asteroid there to check
			}
			if x, y := m.findAlongLine(sx, sy, ex, ey); x != 0 && y != 0 {
				m.SeenMap[(y-1)*m.Width+(x-1)] = '.'
				m.ViewList = append(m.ViewList, view{ex, ey, angle(sx, sy, ex, ey), distance(sx, sy, ex, ey)})
			}
		}
	}
}

func distance(sx, sy, ex, ey int) float64 {
	return math.Sqrt(math.Pow(float64(sx-ex), 2) + math.Pow(float64(sy-ey), 2))
}

func angle(sx, sy, ex, ey int) float64 {
	dx, dy := reduce(ex-sx, ey-sy)
	angle := math.Atan2(float64(dy), float64(dx))/(2*math.Pi)*360.0 + 90
	if angle < 0 {
		angle += 360
	}
	return angle
}

// Looks along line on map from sx,sy towards ex, ey and returns the x,y of the first asteroid seen
// or 0, 0 if none
func (m *Map) findAlongLine(sx, sy, ex, ey int) (int, int) {
	var dx, dy int
	dx, dy = reduce(ex-sx, ey-sy)
	x, y := sx+dx, sy+dy
	if dy > 0 {
		if dx > 0 {
			for {
				if x > ex || y > ey {
					break
				}
				if m.Get(x, y) == '#' {
					return x, y
				}
				x, y = x+dx, y+dy
			}
		} else {
			for {
				if x < ex || y > ey {
					break
				}
				if m.Get(x, y) == '#' {
					return x, y
				}
				x, y = x+dx, y+dy
			}
		}
	} else {
		if dx > 0 {
			for {
				if x > ex || y < ey {
					break
				}
				if m.Get(x, y) == '#' {
					return x, y
				}
				x, y = x+dx, y+dy
			}
		} else {
			for {
				if x < ex || y < ey {
					break
				}
				if m.Get(x, y) == '#' {
					return x, y
				}
				x, y = x+dx, y+dy
			}
		}
	}
	return 0, 0
}

func (m *Map) Vaporize() {
	count := 0
	for {
		passCount := 0
		mode := 0
		curX, curY := m.MaxX, 1
		for n := 1; n < m.Width*m.Height; n++ {
			if x, y := m.findAlongLine(m.MaxX, m.MaxY, curX, curY); x != 0 && y != 0 {
				count++
				passCount++
				fmt.Printf("Vaporized #%d at %d,%d\n", count, x, y)
				//m.Pixels[(y-1)*m.Width+(x-1)] = '.'
			}
			switch mode {
			case 0:
				if curX == m.Width {
					mode = 1
				} else {
					curX++
				}
			case 1:
				if curY == m.Height {
					mode = 2
				} else {
					curY++
				}
			case 2:
				if curX == 1 {
					mode = 3
				} else {
					curX--
				}
			case 3:
				if curY == 1 {
					mode = 0
				} else {
					curY--
				}
			}
		}
		m.ShowMap()
		break
		if passCount == 0 {
			break
		}
	}
}

func reduce(dx, dy int) (int, int) {
	for {
		absDx := dx
		if absDx < 0 {
			absDx = -absDx
		}
		absDy := dy
		if absDy < 0 {
			absDy = -absDy
		}
		if dx == 0 {
			return dx, dy / absDy
		} else if dy == 0 {
			return dx / absDx, dy
		}
		smallest := absDx
		if absDy < smallest {
			smallest = absDy
		}
		foundone := false
		for i := smallest; i > 1; i-- {
			if dx%i == 0 && dy%i == 0 {
				dx /= i
				dy /= i
				foundone = true
			}
		}
		if !foundone {
			break
		}
	}
	return dx, dy
}

func (m *Map) SetSeen(x, y int) bool {
	fmt.Printf("%d,%d  ", x, y)
	if m.Get(x, y) != '.' {
		m.SeenMap[(y-1)*m.Width+(x-1)] = '.'
		return true
	}
	return false
}

func (m *Map) buildSeenMap(x1, y1, x2, y2, w, h int) {
	var dx, dy, e, slope int

	fmt.Printf("%d, %d -> %d, %d: ", x1, y1, x2, y2)
	defer fmt.Println("")
	dx, dy = x2-x1, y2-y1
	if dy < 0 {
		dy = -dy
	}
	if dx < 0 {
		dx = -dx
	}
	cx := 1
	if x1 > x2 {
		cx = -1
	}
	cy := 1
	if y1 > y2 {
		cy = -1
	}
	switch {
	// Is line a point ?
	case x1 == x2 && y1 == y2:
		return

	// Is line an horizontal ?
	case y1 == y2:
		for ; dx != 0; dx-- {
			x1 += cx
			if m.SetSeen(x1, y1) {
				return
			}
		}

	// Is line a vertical ?
	case x1 == x2:
		for ; dy != 0; dy-- {
			y1 += cy
			if m.SetSeen(x1, y1) {
				return
			}
		}

	// Is line a diagonal ?
	case dx == dy:
		for ; dx != 0; dx-- {
			x1 += cx
			y1 += cy
			if m.SetSeen(x1, y1) {
				return
			}
		}

	// wider than high ?
	case dx > dy:
		if y1 < y2 {
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				x1 += cx
				e -= dy
				if e < 0 {
					y1++
					e += slope
				}
				if m.SetSeen(x1, y1) {
					return
				}
			}
		} else {
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				x1 += cx
				e -= dy
				if e < 0 {
					y1--
					e += slope
				}
				if m.SetSeen(x1, y1) {
					return
				}
			}
		}
		if m.SetSeen(x2, y2) {
			return
		}

	// higher than wide.
	default:
		if y1 < y2 {
			// BresenhamDyXRYD(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				y1++
				e -= dx
				if e < 0 {
					//x1++
					x1 += cx
					e += slope
				}
				if m.SetSeen(x1, y1) {
					return
				}
			}
		} else {
			// BresenhamDyXRYU(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				y1--
				e -= dx
				if e < 0 {
					//x1++
					x1 += cx
					e += slope
				}
				if m.SetSeen(x1, y1) {
					return
				}
			}
		}
		if m.SetSeen(x2, y2) {
			return
		}
	}
}

package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"os/exec"
	"sort"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/vova616/screenshot"
)

var cnt int = 0
var failedToFindBombCell int = 0
var faildedToFindFreeCell int = 0
var a [][]int
var totolBombsFound = 0

func main() {
	exec.Command("gnome-mines", "--big").Start()
	fmt.Println("lets start")

	time.Sleep(3 * time.Second)

	// click on center cell
	clickLeftXY(24+38*15+19, 93+38*8+19)

	time.Sleep(1 * time.Second)

	a = make([][]int, 16)
	for i := range a {
		a[i] = make([]int, 30)
	}

	fillAllArr()

	iteration := 0
	for iteration < 200 {

		failedToFindBombCell++
		faildedToFindFreeCell++

		markBombCells()

		clickFreeCell()

		if totolBombsFound == 99 {
			fmt.Println("Solved it :D")
			os.Exit(0)
		}
		if failedToFindBombCell >= 4 && faildedToFindFreeCell >= 4 {

			unopenedCellArr := retUnpenedCellArr(a)

			allRegions := segregate(unopenedCellArr)

			allPoints := make([]struct {
				p   point
				cnt int
			}, 0)

			numOfRegionsChecked := 0
			for _, region := range allRegions {

				if len(region) > 26 {
					continue
				}
				numOfRegionsChecked++
				points := tankSolverAlgorithm(region, 0)

				for idx := range points {
					allPoints = append(allPoints, struct {
						p   point
						cnt int
					}{region[idx], points[idx]})
				}

			}
			if numOfRegionsChecked == 0 {
				fmt.Println("sorry i can't solve it because it will take long time :(")
				os.Exit(0)
			}
			sort.Slice(allPoints, func(i, j int) bool {
				return allPoints[i].cnt < allPoints[j].cnt
			})

			if len(allPoints) > 0 && allPoints[0].cnt > 0 {
				clickLeftXY(allPoints[0].p.Y*38+24+19, allPoints[0].p.X*38+19+93)
			} else {
				for _, point := range allPoints {
					if point.cnt == 0 {
						clickLeftXY(point.p.Y*38+24+19, point.p.X*38+19+93)
					} else {
						break
					}
				}
			}
			failedToFindBombCell = 0
			faildedToFindFreeCell = 0
		}
		fillAllArr()
		iteration++
	}

}

type cellColor struct {
	R int
	G int
	B int
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func getMax(a, b, c, d, e, f, g, h, i int) int {
	return max(a, max(b, max(c, max(d, max(e, max(f, max(g, max(h, i))))))))
}

func getCellNumber(img image.Image) int {
	unopenedCell := cellColor{186, 189, 182}
	zeroCell := cellColor{222, 222, 220}
	oneCell := cellColor{221, 250, 195}
	twoCell := cellColor{236, 237, 191}
	threeCell := cellColor{237, 218, 180}
	fourCell := cellColor{237, 195, 138}
	fiveCell := cellColor{247, 161, 162}
	sixCell := cellColor{254, 167, 133}
	sevenCell := cellColor{255, 125, 96}
	bombCell := cellColor{204, 0, 0}

	rect := img.Bounds()

	cellunopendCnt := 0
	cellZeroCnt := 0
	cellOneCnt := 0
	cellTwoCnt := 0
	cellThreeCnt := 0
	cellFourCnt := 0
	cellFiveCnt := 0
	cellSixCnt := 0
	cellSevenCnt := 0

	for i := 0; i < rect.Max.Y; i++ {
		for j := 0; j < rect.Max.X; j++ {
			c := color.RGBAModel.Convert(img.At(j, i))
			r := int(c.(color.RGBA).R)
			g := int(c.(color.RGBA).G)
			b := int(c.(color.RGBA).B)
			//fmt.Printf("%d %d %d \n", r, g, b)
			if bombCell.R == r && bombCell.B == b && bombCell.G == g {
				fmt.Println("Sorry Failed to solve it :(")
				os.Exit(0)
			}
			if unopenedCell.R == r && unopenedCell.B == b && unopenedCell.G == g {
				cellunopendCnt++
			}
			if zeroCell.R == r && zeroCell.B == b && zeroCell.G == g {
				cellZeroCnt++
			}
			if oneCell.R == r && oneCell.B == b && oneCell.G == g {
				cellOneCnt++
			}
			if twoCell.R == r && twoCell.B == b && twoCell.G == g {
				cellTwoCnt++
			}
			if threeCell.R == r && threeCell.B == b && threeCell.G == g {
				cellThreeCnt++
			}
			if fourCell.R == r && fourCell.B == b && fourCell.G == g {
				cellFourCnt++
			}
			if fiveCell.R == r && fiveCell.B == b && fiveCell.G == g {
				cellFiveCnt++
			}
			if sixCell.R == r && sixCell.B == b && sixCell.G == g {
				cellSixCnt++
			}
			if sevenCell.R == r && sevenCell.B == b && sevenCell.G == g {
				cellSevenCnt++
			}
		}
	}
	ret := -100
	if cellunopendCnt == getMax(cellunopendCnt, cellZeroCnt, cellOneCnt, cellTwoCnt, cellThreeCnt, cellFourCnt, cellFiveCnt, cellSixCnt, cellSevenCnt) {
		ret = 9
	}
	if cellZeroCnt == getMax(cellunopendCnt, cellZeroCnt, cellOneCnt, cellTwoCnt, cellThreeCnt, cellFourCnt, cellFiveCnt, cellSixCnt, cellSevenCnt) {
		ret = 0
	}
	if cellOneCnt == getMax(cellunopendCnt, cellZeroCnt, cellOneCnt, cellTwoCnt, cellThreeCnt, cellFourCnt, cellFiveCnt, cellSixCnt, cellSevenCnt) {
		ret = 1
	}
	if cellTwoCnt == getMax(cellunopendCnt, cellZeroCnt, cellOneCnt, cellTwoCnt, cellThreeCnt, cellFourCnt, cellFiveCnt, cellSixCnt, cellSevenCnt) {
		ret = 2
	}
	if cellThreeCnt == getMax(cellunopendCnt, cellZeroCnt, cellOneCnt, cellTwoCnt, cellThreeCnt, cellFourCnt, cellFiveCnt, cellSixCnt, cellSevenCnt) {
		ret = 3
	}
	if cellFourCnt == getMax(cellunopendCnt, cellZeroCnt, cellOneCnt, cellTwoCnt, cellThreeCnt, cellFourCnt, cellFiveCnt, cellSixCnt, cellSevenCnt) {
		ret = 4
	}
	if cellFiveCnt == getMax(cellunopendCnt, cellZeroCnt, cellOneCnt, cellTwoCnt, cellThreeCnt, cellFourCnt, cellFiveCnt, cellSixCnt, cellSevenCnt) {
		ret = 5
	}
	if cellSixCnt == getMax(cellunopendCnt, cellZeroCnt, cellOneCnt, cellTwoCnt, cellThreeCnt, cellFourCnt, cellFiveCnt, cellSixCnt, cellSevenCnt) {
		ret = 6
	}
	if cellSevenCnt == getMax(cellunopendCnt, cellZeroCnt, cellOneCnt, cellTwoCnt, cellThreeCnt, cellFourCnt, cellFiveCnt, cellSixCnt, cellSevenCnt) {
		ret = 7
	}
	if ret == -100 {
		fmt.Println("Failed to detect the screen")
		os.Exit(1)
	}

	return ret
}

func markBombCells() {
	i := 0
	for i < 16 {
		j := 0
		for j < 30 {
			numberOfBombs := a[i][j]
			if numberOfBombs <= 0 || numberOfBombs == 9 {
				j++
				continue
			}
			numberOfNeighbours := 0
			numberOfFlagedNeighbours := 0

			h := -1
			k := -1

			for h <= 1 {
				k = -1
				for k <= 1 {
					if k == 0 && h == 0 {
						k++
						continue
					}
					newX := i + h
					newY := j + k
					if newX < 16 && newX >= 0 && newY < 30 && newY >= 0 {
						if a[newX][newY] == 9 {
							numberOfNeighbours++
						}
						if a[newX][newY] == -2 {
							numberOfFlagedNeighbours++
						}

					}
					k++
				}
				h++
			}
			if numberOfNeighbours+numberOfFlagedNeighbours == numberOfBombs {
				h = -1
				k = -1
				for h <= 1 {
					k = -1
					for k <= 1 {
						if k == 0 && h == 0 {
							k++
							continue
						}
						newX := i + h
						newY := j + k
						if newX < 16 && newX >= 0 && newY < 30 && newY >= 0 && a[newX][newY] == 9 {
							a[newX][newY] = -2
							cnt++
							failedToFindBombCell = 0
							clickRightXY(newY, newX)
						}
						k++
					}
					h++
				}
			}
			j = j + 1
		}
		i = i + 1
	}
}

func clickFreeCell() {
	i := 0
	for i < 16 {
		j := 0
		for j < 30 {
			numberOfBombs := a[i][j]
			if numberOfBombs <= 0 || numberOfBombs == 9 {
				j++
				continue
			}
			numberOfNeighbours := 0
			numberOfFlagedNeighbours := 0

			h := -1
			k := -1

			for h <= 1 {
				k = -1
				for k <= 1 {
					if k == 0 && h == 0 {
						k++
						continue
					}
					newX := i + h
					newY := j + k
					if newX < 16 && newX >= 0 && newY < 30 && newY >= 0 {
						if a[newX][newY] == 9 {
							numberOfNeighbours++
						}
						if a[newX][newY] == -2 {
							numberOfFlagedNeighbours++
						}
					}
					k++
				}
				h++
			}
			if numberOfFlagedNeighbours == numberOfBombs {
				h = -1
				k = -1
				for h <= 1 {
					k = -1
					for k <= 1 {
						if k == 0 && h == 0 {
							k++
							continue
						}
						newX := i + k
						newY := j + h
						if newX < 16 && newX >= 0 && newY < 30 && newY >= 0 && a[newX][newY] == 9 {
							faildedToFindFreeCell = 0
							clickLeftXY(24+newY*38+19, 93+newX*38+19)
						}
						k++
					}
					h++
				}
			}
			j = j + 1
		}
		i = i + 1
	}
}

func clickLeftXY(x, y int) {
	robotgo.MoveMouse(x, y)
	robotgo.MouseClick()
	robotgo.MoveMouse(-x, -y)
}

func clickRightXY(x, y int) {
	x = 24 + x*38 + 19
	y = 93 + y*38 + 19
	robotgo.MoveMouse(x, y)
	robotgo.MouseClick("right")
	totolBombsFound++
	robotgo.MoveMouse(-x, -y)
}

func fillAllArr() {

	i := 0
	currentX := 24
	startX := 24
	currentY := 93

	for i < 16 {
		j := 0
		for j < 30 {
			if a[i][j] >= 0 {
				img, err := screenshot.CaptureRect(image.Rect(currentX, currentY, currentX+38, currentY+38))
				if err != nil {
					fmt.Println("ERROR")
				}
				myImg := image.Image(img)
				a[i][j] = getCellNumber(myImg)
			}
			currentX += 38
			j = j + 1
		}
		currentX = startX
		currentY += 38
		i = i + 1
	}
}

type point struct {
	X int
	Y int
}

// return unopened cell neighbored with opened cell
func retUnpenedCellArr(a [][]int) []point {
	var arr []point
	i := 0
	for i < 16 {
		j := 0
		for j < 30 {
			if a[i][j] == 9 {
				numberOfOpenedNeighbours := 0
				h := -1
				k := -1

				for h <= 1 {
					k = -1
					for k <= 1 {
						if k == 0 && h == 0 {
							k++
							continue
						}
						newX := i + h
						newY := j + k
						if newX < 16 && newX >= 0 && newY < 30 && newY >= 0 {
							if a[newX][newY] >= 0 && a[newX][newY] <= 7 {
								numberOfOpenedNeighbours++
							}
						}
						k++
					}
					h++
				}
				if numberOfOpenedNeighbours > 0 {
					arr = append(arr, point{i, j})
				}
			}
			j = j + 1
		}
		i = i + 1
	}
	return arr
}

func validate(pointArr []point) bool {
	for _, p := range pointArr {
		i := -1
		j := -1
		for i <= 1 {
			j = -1
			for j <= 1 {
				if j == 0 && i == 0 {
					j++
					continue
				}
				x := i + p.X
				y := j + p.Y
				if x < 16 && x >= 0 && y < 30 && y >= 0 {
					if a[x][y] > 0 && a[x][y] < 9 {
						numberOfBombs := a[x][y]
						numberOfFlagedNeighbours := 0
						numberOfUntrustedBombs := 0

						h := -1
						k := -1

						for h <= 1 {
							k = -1
							for k <= 1 {
								if k == 0 && h == 0 {
									k++
									continue
								}
								newX := x + h
								newY := y + k
								if newX < 16 && newX >= 0 && newY < 30 && newY >= 0 {

									if a[newX][newY] == -2 {
										numberOfFlagedNeighbours++
									}
									if a[newX][newY] == -3 {
										numberOfUntrustedBombs++
									}

								}
								k++
							}
							h++
						}
						if numberOfFlagedNeighbours+numberOfUntrustedBombs != numberOfBombs {
							return false
						}
					}
				}
				j++
			}
			i++
		}
	}
	return true
}

func tankSolverAlgorithm(pointArr []point, idx int) []int {
	ans := make([]int, len(pointArr))
	if idx == len(pointArr) {
		if validate(pointArr) {
			i := 0
			for i < len(pointArr) {
				if a[pointArr[i].X][pointArr[i].Y] == -3 {
					ans[i]++
				}
				i++
			}
		}
		return ans
	}
	// assume this cell contains bomb
	a[pointArr[idx].X][pointArr[idx].Y] = -3
	tempAns := tankSolverAlgorithm(pointArr, idx+1)
	a[pointArr[idx].X][pointArr[idx].Y] = 9
	j := 0
	for j < len(ans) {
		ans[j] += tempAns[j]
		j++
	}

	// assume this cell doesn't contain bomb
	tempAns = tankSolverAlgorithm(pointArr, idx+1)
	j = 0
	for j < len(ans) {
		ans[j] += tempAns[j]
		j++
	}
	return ans
}

func contains(s []point, p point) bool {
	for _, a := range s {
		if a == p {
			return true
		}
	}
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func segregate(pointArr []point) [][]point {
	allRegions := make([][]point, 0)
	coverd := make([]point, 0)
	for true {
		queue := make([]point, 0)
		finishedRegion := make([]point, 0)

		for _, p := range pointArr {
			if !contains(coverd, p) {
				queue = append(queue, p)
				break
			}
		}
		if len(queue) == 0 {
			break
		}
		for len(queue) > 0 {
			var p point
			p, queue = queue[len(queue)-1], queue[:len(queue)-1] // pop from queue
			ci := p.X
			cj := p.Y

			finishedRegion = append(finishedRegion, p)
			coverd = append(coverd, p)

			for _, tile := range pointArr {
				ti := tile.X
				tj := tile.Y

				isConnected := false

				if contains(finishedRegion, tile) {
					continue
				}

				if abs(ci-ti) > 2 || abs(cj-tj) > 2 {
					isConnected = false
				} else {
					i := 0
					for i < 16 && !isConnected {
						j := 0
						for j < 30 {
							numberOfBombs := a[i][j]
							if numberOfBombs > 0 && numberOfBombs <= 8 {
								if abs(ci-i) <= 1 && abs(cj-j) <= 1 && abs(ti-i) <= 1 && abs(tj-j) <= 1 {
									isConnected = true
									break
								}
							}
							j = j + 1
						}
						i = i + 1
					}
				}

				if !isConnected {
					continue
				}

				if !contains(queue, tile) {
					queue = append(queue, tile)
				}
			}
		}
		allRegions = append(allRegions, finishedRegion)
	}
	return allRegions
}

// AILab3 project main.go
package main

import (
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"os"
	"reflect"
)

type Net struct {
	weights   [][]float64
	output    []int
	size      int
	templates [][]int
}

const (
	letterWidth  = 6
	letterHeight = 6
)

func NewNet(templates [][]int) *Net {
	n := &Net{}
	n.templates = templates
	for key := range templates {
		if n.size == 0 {
			n.size = len(templates[key])
			continue
		}
		if len(templates[key]) != n.size {
			return nil
		}
	}
	n.output = make([]int, n.size)
	n.weights = make([][]float64, n.size)
	for i := 0; i < n.size; i++ {
		n.weights[i] = make([]float64, n.size)
	}
	for i := 0; i < n.size; i++ {
		for j := 0; j < n.size; j++ {
			if i == j {
				continue
			}
			sum := 0.0
			for _, template := range templates {
				sum += float64(template[i] * template[j])
			}
			n.weights[i][j] = sum / float64(n.size)
			n.weights[j][i] = sum / float64(n.size)
		}
	}
	//fmt.Println(n.weights)
	return n
}

func (n *Net) Detect(name string, input []int) []int {
	copy(n.output, input)
	oldOutput := make([]int, n.size)
	iteration := 0
	for iteration = 0; iteration == 0 || !reflect.DeepEqual(oldOutput, n.output); iteration++ {
		copy(oldOutput, n.output)
		for i := 0; i < n.size; i++ {
			z := 0.0
			for j := 0; j < n.size; j++ {
				z += n.weights[i][j] * float64(oldOutput[j])
			}
			if z > 0 {
				n.output[i] = 1
			} else if z < 0 {
				n.output[i] = -1
			}
		}
	}
	for _, template := range n.templates {
		if reflect.DeepEqual(template, n.output) {
			fmt.Println("Match found ", iteration, input, n.output)
			saveImages(name, input, n.output)
			return n.output
		}
	}
	fmt.Println("Match not found ", iteration, input, n.output)
	saveImages(name, input, n.output)
	return nil
}

func saveImages(name string, first []int, second []int) {
	file, _ := os.Create(name + ".png")
	defer file.Close()
	img := image.NewRGBA(image.Rect(0, 0, 2*letterWidth, letterHeight))
	for i := 0; i < letterHeight; i++ {
		for j := 0; j < letterWidth; j++ {
			if first[i*letterWidth+j] > 0 {
				img.Set(j, i, color.RGBA{255, 255, 255, 255})
			} else {
				img.Set(j, i, color.RGBA{0, 0, 0, 255})
			}
			if second[i*letterWidth+j] > 0 {
				img.Set(letterWidth+j, i, color.RGBA{255, 255, 255, 255})
			} else {
				img.Set(letterWidth+j, i, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	bmp.Encode(file, img)
}

func main() {
	file, _ := os.Open("font.bmp")
	defer file.Close()
	img, _ := bmp.Decode(file)
	lettersSet := []rune{}
	lettersSet = append(lettersSet, 'a', 'b', 'c', 'd')
	/*for c := 'a'; c < 'l'; c++ {
		lettersSet = append(lettersSet, c)
	}*/
	letters := make(map[rune][]int)
	for _, c := range lettersSet {
		letterX := (int(c-'a') * letterWidth)
		buf := make([]int, 0, letterHeight*letterWidth)
		for y := 0; y < letterHeight; y++ {
			for x := 0; x < letterWidth; x++ {
				r, g, b, _ := img.At(letterX+x, y).RGBA()
				if r != 0 && g != 0 && b != 0 {
					buf = append(buf, 1)
				} else {
					buf = append(buf, -1)
				}
			}
		}
		letters[c] = buf
	}
	templates := [][]int{}
	for _, val := range letters {
		templates = append(templates, val)
	}
	net := NewNet(templates)
	for key, val := range letters {
		detectedVal := net.Detect(string(key), val)
		if detectedVal == nil || !reflect.DeepEqual(detectedVal, val) {
			fmt.Println("Fail at letter ", string(key))
		} else {
			fmt.Println("Success at letter ", string(key))
		}
	}
}

package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Neuron struct {
	Inputs []float64
	Fire   FireHandler
}

func NewNeuron(numInputs int) *Neuron {
	return &Neuron{
		Inputs: make([]float64, numInputs),
	}
}

type FireHandler interface {
	Fire(*Neuron) (out float64, right func(), wrong func())
}

type WeightingHandler struct {
	Weights  []float64
	Interval int
}

func (h *WeightingHandler) Fire(n *Neuron) (out float64, right func(), ouch func()) {
	for k, w := range h.Weights {
		out += w * n.Inputs[k]
	}
	out = sigmoid(out)
	return out, func() {
			// correct
			h.Interval++

		}, func() {
			// wrong
			h.Interval--
			if h.Interval < 0 {
				h.Interval = 0
			}
			for k := range h.Weights {
				dw := (2 * rand.Float64()) - 1
				dw *= 0.1
				h.Weights[k] += dw / math.Pow(10, float64(h.Interval))
				if h.Weights[k] < 0 {
					h.Weights[k] = 0
				} else if h.Weights[k] > 5 {
					h.Weights[k] = 5
				}
				fmt.Println(h.Weights[k])
			}

		}
}

func main() {

	const pivot = 0.3

	n := NewNeuron(1)
	// n := NewNeuron(3)
	h := &WeightingHandler{
		Weights:  []float64{0.5},
		Interval: 0,
		// Weights:         []float64{0.5, 0.5, 0.5},
	}
	n.Fire = h
	for i := 0; i < 1000; i++ {
		n.Inputs[0] = rand.Float64()
		// n.Inputs[1] = rand.Float64()
		// n.Inputs[2] = rand.Float64()

		for {
			out, yes, no := n.Fire.Fire(n)
			saysYes := out >= sigmoid(0.5)
			if saysYes {
				if n.Inputs[0] > pivot {
					yes()
					break
				} else {
					no()
				}
			} else {
				if n.Inputs[0] > pivot {
					no()
				} else {
					yes()
					break
				}
			}
			fmt.Printf("[%.4f] -> %.4f: expected: %t, guess: %t\n", n.Inputs[0], out, n.Inputs[0] > pivot, saysYes)
		}
	}

	// fmt.Printf("Final Weights: [%.10f][%.10f][%.10f]\n", h.Weights[0], h.Weights[1], h.Weights[2])
	fmt.Printf("Final Weights: [%.10f]\n", h.Weights[0])

	for {
		fmt.Println("Test a float : ")
		_, err := fmt.Scanf("%f", &n.Inputs[0])
		if err != nil {
			fmt.Println(err)
		}
		out, _, _ := n.Fire.Fire(n)

		fmt.Printf("Is above %.2f: %t (out = %.3f)\n", pivot, out > sigmoid(0.5), out)
	}
}

// Sigmoid returns the input values in the range of -1 to 1
// along the sigmoid or s-shaped curve, commonly used in
// machine learning while training neural networks as an
// activation function.
func sigmoid(input float64) float64 {
	x := 1 / (1 + math.Exp(-input))
	fmt.Println("sigm", x)
	return x
}

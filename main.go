package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Payload struct {
	TotalPlates      int `json:"n"`
	TargetPalindrome int `json:"x"`
}

type Response struct {
	IterativeTime   string    `json:"iterativeTime"`
	RecursiveTime   string    `json:"recursiveTime"`
	IterativeRaw    float64   `json:"iterativeRaw"`
	RecursiveRaw    float64   `json:"recursiveRaw"`
	GraphLabels     []int     `json:"graphLabels"`
	GraphIterative  []float64 `json:"graphIterative"`
	GraphRecursive  []float64 `json:"graphRecursive"`
	DetectedSamples []string  `json:"samples"`
}

// --- ALGORITMA ---

// Iteratif
func isPalindromeIterative(s string) bool {
	clean := strings.ReplaceAll(s, " ", "")
	n := len(clean)
	for i := 0; i < n/2; i++ {
		if clean[i] != clean[n-1-i] {
			return false
		}
	}
	return true
}

// Rekursif
func isPalindromeRecursive(s string) bool {
	clean := strings.ReplaceAll(s, " ", "")
	return recursiveHelper(clean)
}

func recursiveHelper(s string) bool {
	if len(s) <= 1 {
		return true
	}
	if s[0] != s[len(s)-1] {
		return false
	}
	return recursiveHelper(s[1 : len(s)-1])
}

// --- GENERATOR PLAT KENDARAAN ---

func generatePlates(n, x int) []string {
	plates := make([]string, n)
	rand.Seed(time.Now().UnixNano())

	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for i := 0; i < n; i++ {
		// Logika: Buat X plat pertama jadi palindrom
		isTargetPalindrome := i < x 

		if isTargetPalindrome {
			// Generate Palindrom Sengaja (Cth: A 121 A)
			l := string(letters[rand.Intn(len(letters))])
			d1 := rand.Intn(9) + 1
			d2 := rand.Intn(10)
			plates[i] = fmt.Sprintf("%s %d%d%d %s", l, d1, d2, d1, l)
		} else {
			// Generate Random (Bisa jadi palindrom tidak sengaja, tapi peluang kecil)
			l1 := string(letters[rand.Intn(len(letters))])
			l2 := string(letters[rand.Intn(len(letters))])
			nums := rand.Intn(899) + 100 // 3 digit random (100-999)
			plates[i] = fmt.Sprintf("%s %d %s", l1, nums, l2)
		}
	}

	// Acak posisi agar palindrom menyebar
	rand.Shuffle(len(plates), func(i, j int) { plates[i], plates[j] = plates[j], plates[i] })

	return plates
}

// --- HANDLER ---

func benchmarkHandler(w http.ResponseWriter, r *http.Request) {
	var p Payload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 1. Generate Data Utama
	plates := generatePlates(p.TotalPlates, p.TargetPalindrome)

	// 2. Siapkan Data untuk Grafik (Simulation Growth)
	steps := 5
	graphLabels := make([]int, steps)
	graphIter := make([]float64, steps)
	graphRec := make([]float64, steps)

	for i := 1; i <= steps; i++ {
		currentN := (p.TotalPlates / steps) * i
		subset := plates[:currentN]
		graphLabels[i-1] = currentN

		// Ukur Iteratif
		start := time.Now()
		for _, plate := range subset {
			isPalindromeIterative(plate)
		}
		graphIter[i-1] = float64(time.Since(start).Microseconds()) / 1000.0 // ke milisecond

		// Ukur Rekursif
		start = time.Now()
		for _, plate := range subset {
			isPalindromeRecursive(plate)
		}
		graphRec[i-1] = float64(time.Since(start).Microseconds()) / 1000.0
	}

	// 3. Ambil Sampel Palindrom
	var detected []string
	count := 0

	// Tentukan batas tampilan (display limit)
	// Jika user minta X=2, kita tampilkan max 2.
	// Jika user minta X=1000, kita tampilkan max 20 saja agar UI tidak lag.
	displayLimit := p.TargetPalindrome
	if displayLimit > 20 {
		displayLimit = 20
	}
	
	for _, plate := range plates {
		if isPalindromeIterative(plate) {
			detected = append(detected, plate)
			count++
			
			// Stop mencari sampel jika sudah memenuhi kuota yang diminta user
			// atau maksimal 20 (mana yang lebih kecil)
			if count >= displayLimit {
				break
			}
		}
	}

	// 4. Kirim Response
	resp := Response{
		IterativeTime:   fmt.Sprintf("%.2f ms", graphIter[steps-1]),
		RecursiveTime:   fmt.Sprintf("%.2f ms", graphRec[steps-1]),
		IterativeRaw:    graphIter[steps-1],
		RecursiveRaw:    graphRec[steps-1],
		GraphLabels:     graphLabels,
		GraphIterative:  graphIter,
		GraphRecursive:  graphRec,
		DetectedSamples: detected,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	// Serve file static (index.html, style.css, script.js) dari folder saat ini
	http.Handle("/", http.FileServer(http.Dir("./"))) 
	http.HandleFunc("/api/benchmark", benchmarkHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
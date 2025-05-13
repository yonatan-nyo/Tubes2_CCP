package models

import (
	"fmt"
	"time"
)

type QueueItem struct {
	Element *ElementsGraphNode
	From    string
}

// Fungsi utama algoritma Bidirectional Search
func BidirectionalFindTrees(
	rootRecipeTree *RecipeTreeNode,
	targetGraphNode *ElementsGraphNode,
	maxTreeCount int,
	signalTreeChange func(*RecipeTreeNode, int, int32),
	globalStartTime time.Time,
	globalNodeCounter *int32,
	delayMs int,
) ([]*RecipeTreeNode, error) {
	// Validasi awal apakah graph node target valid
	if targetGraphNode == nil {
		return nil, fmt.Errorf("targetGraphNode is nil")
	}

	// Validasi jika elemen tidak memiliki resep atau merupakan elemen dasar
	if len(targetGraphNode.RecipesToMakeThisElement) == 0 || IsBaseElement(targetGraphNode.Name) {
		node := &RecipeTreeNode{
			Name:      targetGraphNode.Name,
			ImagePath: GetImagePath(targetGraphNode.ImagePath),
		}
		return []*RecipeTreeNode{node}, nil
	}

	// Menyimpan elemen yang sudah dikunjungi dari dua arah
	visitedUpper := make(map[string]bool)
	visitedLower := make(map[string]bool)
	// Menyimpan node yang sudah pernah ditemukan
	seenMeeting := make(map[string]bool)

	// queueUpper: proses dimulai dari target
	queueUpper := []*QueueItem{{Element: targetGraphNode}}

	// queueLower: proses dimulai dari base elements yang relevan
	// Pemilihan base elements menggunakan IsThisMadeFrom
	queueLower := []*QueueItem{}
	for _, base := range GetBaseElements() {
		if targetGraphNode.IsThisMadeFrom(base) {
			if node, ok := GetElementsGraphNodeByName(base); ok {
				queueLower = append(queueLower, &QueueItem{Element: node})
			}
		}
	}

	var resultTrees []*RecipeTreeNode

	// Proses utama loop pencarian dua arah
	for len(queueUpper) > 0 && len(queueLower) > 0 && len(resultTrees) < maxTreeCount {
		// Delay setiap iterasi untuk visualisasi live update
		if delayMs > 0 {
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}

		// Proses pencarian dari arah target menuju base elements
		nQueueUpper, newUpperNames := processUpper(queueUpper, visitedUpper)
		for _, name := range newUpperNames {
			// Jika ditemukan pertemuan dan belum diproses sebelumnya
			if visitedLower[name] && name == targetGraphNode.Name {
				if _, already := seenMeeting[name]; already {
					continue
				}
				seenMeeting[name] = true
				if node, ok := GetElementsGraphNodeByName(name); ok {
					// DFS dipanggil setelah upper dan lower bertemu untuk membangun tree secara lengkap
					treesFromDFS, err := DFSFindTrees(nil, node, maxTreeCount, signalTreeChange, globalStartTime, globalNodeCounter, delayMs)
					if err == nil {
						resultTrees = appendAllValidTargetTrees(resultTrees, treesFromDFS, targetGraphNode.Name, maxTreeCount)
						if len(resultTrees) >= maxTreeCount {
							return resultTrees, nil
						}
					}
				}
			}
		}
		queueUpper = nQueueUpper

		if delayMs > 0 {
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}

		// Proses pencarian dari base menuju target
		nQueueLower, newLowerNames := processLower(queueLower, visitedLower)
		for _, name := range newLowerNames {
			if visitedUpper[name] && name == targetGraphNode.Name {
				if _, already := seenMeeting[name]; already {
					continue
				}
				seenMeeting[name] = true
				if node, ok := GetElementsGraphNodeByName(name); ok {
					treesFromDFS, err := DFSFindTrees(nil, node, maxTreeCount, signalTreeChange, globalStartTime, globalNodeCounter, delayMs)
					if err == nil {
						resultTrees = appendAllValidTargetTrees(resultTrees, treesFromDFS, targetGraphNode.Name, maxTreeCount)
						if len(resultTrees) >= maxTreeCount {
							return resultTrees, nil
						}
					}
				}
			}
		}
		queueLower = nQueueLower
	}
	return resultTrees, nil
}

// Helper function untuk memroses queueUpper (pencarian dari target ke base)
func processUpper(queue []*QueueItem, visited map[string]bool) ([]*QueueItem, []string) {
	// Queue untuk iterasi selanjutnya
	nextQueue := []*QueueItem{}
	// Menyimpan nama-nama node yang dihasilkan
	// Digunakan untuk deteksi pertemuan
	produced := []string{}

	for _, item := range queue {
		node := item.Element

		// Jika node sudah dikunjungi, maka abaikan
		if visited[node.Name] {
			continue
		}

		// Tandai viisted
		visited[node.Name] = true
		produced = append(produced, node.Name)

		// Proses seluruh resep pembentuk elemen ini
		// Lanjutkan ke child nodes yang menjadi bahan resep
		for _, recipe := range node.RecipesToMakeThisElement {
			// Tambahkan kedua elemen bahan ke antrian berikutnya
			nextQueue = append(nextQueue, &QueueItem{Element: recipe.ElementOne})
			nextQueue = append(nextQueue, &QueueItem{Element: recipe.ElementTwo})
		}
	}

	return nextQueue, produced
}

// Helper function untuk memroses queueLower (pencarian dari base ke target)
func processLower(queue []*QueueItem, visited map[string]bool) ([]*QueueItem, []string) {
	// Queue untuk iterasi selanjutnya
	nextQueue := []*QueueItem{}
	// Menyimpan nama-nama node yang dihasilkan
	// Digunakan untuk deteksi pertemuan
	produced := []string{}

	for _, item := range queue {
		node := item.Element

		// Jika node sudah dikunjungi, maka abaikan
		if visited[node.Name] {
			continue
		}

		// Tandai visited
		visited[node.Name] = true
		produced = append(produced, node.Name)

		// Proses seluruh elemen yang dapat dibuat dari elemen ini
		for _, recipe := range node.RecipesToMakeOtherElement {
			// Ambil elemen hasil dari resep
			if targetNode, ok := GetElementsGraphNodeByName(recipe.TargetElementName); ok {
				// Tambahkan elemen hasil ke antrian berikutnya
				nextQueue = append(nextQueue, &QueueItem{Element: targetNode})
			}
		}
	}

	return nextQueue, produced
}

// Menambahkan hasil pencarian ke dalam list hasil akhir (menghindari duplikat)
func appendAllValidTargetTrees(existing []*RecipeTreeNode, newTrees []*RecipeTreeNode, target string, maxTreeCount int) []*RecipeTreeNode {
	// Set untuk mendeteksi duplikasi berdasarkan key unik
	existingSet := make(map[string]bool)

	// Buat key unik untuk setiap tree yang sudah ada
	for _, t := range existing {
		if t.Name == target {
			existingSet[t.ImagePath+t.Name+treeToString(t)] = true
		}
	}

	for _, t := range newTrees {
		// Pastikan hanya pohon dengan root = target yang diterima
		if t.Name == target {
			key := t.ImagePath + t.Name + treeToString(t)
			if !existingSet[key] {
				// Clone tree sebelum dimasukkan agar tidak terjadi referensi bersama
				existing = append(existing, t.clone())
				existingSet[key] = true

				// Jika jumlah pohon sudah memenuhi batas, maka break
				if len(existing) >= maxTreeCount {
					break
				}
			}
		}
	}

	return existing
}

// Mengubah tree menjadi string
// Metode yang digunakan untuk mendeteksi duplicates
func treeToString(tree *RecipeTreeNode) string {
	if tree == nil {
		return ""
	}
	return tree.Name + "(" + treeToString(tree.Element1) + "," + treeToString(tree.Element2) + ")"
}

// Mendapatkan base elements
func GetBaseElements() []string {
	return baseElements
}

package models

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Fungsi utama algoritma BFS
func BFSFindTrees(
	targetGraphNode *ElementsGraphNode,
	maxTreeCount int,
	signalTreeChange func(*RecipeTreeNode, int, int32),
	globalStartTime time.Time,
	globalNodeCounter *int32,
	delayMs int,
) ([]*RecipeTreeNode, error) {
	// Validasi awal apakah node target valid
	if targetGraphNode == nil {
		return nil, fmt.Errorf("targetGraphNode is nil")
	}

	// Jika node merupakan base element atau tidak memiliki resep, langsung return sebagai hasil
	if len(targetGraphNode.RecipesToMakeThisElement) == 0 || IsBaseElement(targetGraphNode.Name) {
		node := &RecipeTreeNode{
			Name:      targetGraphNode.Name,
			ImagePath: GetImagePath(targetGraphNode.ImagePath),
		}
		return []*RecipeTreeNode{node}, nil
	}

	var (
		result     []*RecipeTreeNode                          // Menyimpan hasil akhir
		treesFound int                                        // Jumlah pohon yang berhasil ditemukan
		mu         sync.Mutex                                 // Mutex untuk sinkronisasi antar thread (multithreading)
		wg         sync.WaitGroup                             // WaitGroup untuk menunggu semua goroutine selesai
		resultChan = make(chan *RecipeTreeNode, maxTreeCount) // Channel untuk mengirim hasil tree
	)

	// Struktur queue BFS untuk menyimpan state saat traversal
	type QueueItem struct {
		ElementName string
		Level       int
		TreeSoFar   *RecipeTreeNode
		IsComplete  bool
	}

	// Iterasi setiap resep dari target node
	for _, recipe := range targetGraphNode.RecipesToMakeThisElement {
		wg.Add(1)
		go func(r *Recipe) {
			defer wg.Done()

			// Map untuk menyimpan semua tree parsial yang berhasil dibentuk per elemen
			elementToTrees := make(map[string][]*RecipeTreeNode)
			processedElements := make(map[string]bool) // Menandai elemen yang sudah diproses
			queue := make([]*QueueItem, 0)

			// Memasukkan dua element dari resep ke dalam queue
			queue = append(
				queue,
				&QueueItem{ElementName: r.ElementOne.Name, Level: 1, TreeSoFar: nil, IsComplete: false},
				&QueueItem{ElementName: r.ElementTwo.Name, Level: 1, TreeSoFar: nil, IsComplete: false},
			)

			// BFS loop
			for len(queue) > 0 {
				if delayMs > 0 {
					time.Sleep(time.Duration(delayMs) * time.Millisecond)
				}

				// Batasi jumlah tree yang ditemukan sesuai batas maksimum MaxTreeCount
				mu.Lock()
				if treesFound >= maxTreeCount {
					mu.Unlock()
					return
				}
				mu.Unlock()

				item := queue[0]
				queue = queue[1:]

				// Continue jika sudah pernah diproses
				if processedElements[item.ElementName] {
					continue
				}

				// Ambil node dari nama elemen
				elementNode, exists := GetElementsGraphNodeByName(item.ElementName)
				if !exists || elementNode == nil {
					continue
				}

				// Tambah counter global eksplorasi node (aman untuk goroutine)
				atomic.AddInt32(globalNodeCounter, 1)

				// Jika node adalah base element, buat node tree sederhana
				if IsBaseElement(elementNode.Name) || len(elementNode.RecipesToMakeThisElement) == 0 {
					simpleTree := &RecipeTreeNode{
						Name:      elementNode.Name,
						ImagePath: GetImagePath(elementNode.ImagePath),
					}
					elementToTrees[elementNode.Name] = []*RecipeTreeNode{simpleTree}
					processedElements[elementNode.Name] = true
					continue
				}

				// Jika belum punya tree, tunggu semua prerequisite selesai diproses
				if item.TreeSoFar == nil {
					var allPrereqsProcessed = true

					// Cek apakah semua bahan resep sudah tersedia
					for _, elementRecipe := range elementNode.RecipesToMakeThisElement {
						if !processedElements[elementRecipe.ElementOne.Name] ||
							!processedElements[elementRecipe.ElementTwo.Name] {
							allPrereqsProcessed = false

							// Tambahkan elemen yang belum tersedia ke antrian
							if !processedElements[elementRecipe.ElementOne.Name] {
								queue = append(queue, &QueueItem{
									ElementName: elementRecipe.ElementOne.Name,
									Level:       item.Level + 1,
								})
							}
							if !processedElements[elementRecipe.ElementTwo.Name] {
								queue = append(queue, &QueueItem{
									ElementName: elementRecipe.ElementTwo.Name,
									Level:       item.Level + 1,
								})
							}
						}
					}

					// Jika belum semua bahan tersedia, tunda node saat ini
					if !allPrereqsProcessed {
						queue = append(queue, &QueueItem{
							ElementName: item.ElementName,
							Level:       item.Level + 10,
						})
						continue
					}

					// Semua bahan tersedia, lalu mulai membentuk tree
					var elementTrees []*RecipeTreeNode
					for _, elementRecipe := range elementNode.RecipesToMakeThisElement {
						leftTrees := elementToTrees[elementRecipe.ElementOne.Name]
						rightTrees := elementToTrees[elementRecipe.ElementTwo.Name]

						for _, lt := range leftTrees {
							for _, rt := range rightTrees {
								newTree := &RecipeTreeNode{
									Name:      elementNode.Name,
									ImagePath: GetImagePath(elementNode.ImagePath),
									Element1:  lt,
									Element2:  rt,
								}

								// Kirim update ExploringTree ke FE Visualizer melalui WebSocket
								if signalTreeChange != nil {
									func() {
										defer func() {
											if r := recover(); r != nil {
											}
										}()
										signalTreeChange(
											newTree,
											int(time.Since(globalStartTime).Milliseconds()),
											atomic.LoadInt32(globalNodeCounter),
										)
									}()
								}

								elementTrees = append(elementTrees, newTree)
							}
						}
					}

					// Simpan hasil tree yang dibentuk dan tandai elemen sebagai telah diproses
					elementToTrees[elementNode.Name] = elementTrees
					processedElements[elementNode.Name] = true
				}
			}

			// Setelah seluruh subtree dibentuk, gabungkan kedua elemen bahan menjadi root
			leftTrees := elementToTrees[r.ElementOne.Name]
			rightTrees := elementToTrees[r.ElementTwo.Name]

			if leftTrees == nil || rightTrees == nil {
				return
			}

			for _, lt := range leftTrees {
				for _, rt := range rightTrees {
					mu.Lock()
					if treesFound >= maxTreeCount {
						mu.Unlock()
						return
					}

					// Bentuk tree lengkap dari dua subtree
					root := &RecipeTreeNode{
						Name:      targetGraphNode.Name,
						ImagePath: GetImagePath(targetGraphNode.ImagePath),
						Element1:  lt,
						Element2:  rt,
					}

					treesFound++
					mu.Unlock()

					// Kirim update ExploringTree ke FE Visualizer melalui WebSocket
					if signalTreeChange != nil {
						func() {
							defer func() {
								if r := recover(); r != nil {
								}
							}()
							signalTreeChange(
								root,
								int(time.Since(globalStartTime).Milliseconds()),
								atomic.LoadInt32(globalNodeCounter),
							)
						}()
					}

					resultChan <- root
				}
			}
		}(recipe)
	}

	// Tunggu semua goroutine selesai dan tutup channel hasil
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Ambil seluruh tree hasil dari channel
	for tree := range resultChan {
		result = append(result, tree)
		if len(result) >= maxTreeCount {
			break
		}
	}

	// Jika tidak ada tree yang ditemukan, maka return error
	if len(result) == 0 {
		return nil, fmt.Errorf("no complete tree found for target %s", targetGraphNode.Name)
	}

	return result, nil
}

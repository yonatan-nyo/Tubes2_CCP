package models

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Fungsi utama algoritma DFS
func DFSFindTrees(
	rootRecipeTree *RecipeTreeNode,
	targetGraphNode *ElementsGraphNode,
	maxTreeCount int,
	signalTreeChange func(*RecipeTreeNode, int, int32),
	globalStartTime time.Time,
	globalNodeCounter *int32,
	delayMs int,
) ([]*RecipeTreeNode, error) {
	// Validasi awal jika node target tidak tersedia
	if targetGraphNode == nil {
		return nil, fmt.Errorf("targetGraphNode is nil")
	}

	// Jika node adalah base element atau tidak memiliki resep, maka return node sederhana
	if len(targetGraphNode.RecipesToMakeThisElement) == 0 || IsBaseElement(targetGraphNode.Name) {
		node := &RecipeTreeNode{
			Name:      targetGraphNode.Name,
			ImagePath: GetImagePath(targetGraphNode.ImagePath),
		}
		return []*RecipeTreeNode{node}, nil
	}

	var (
		result []*RecipeTreeNode // Menyimpan hasil pohon recipe yang ditemukan
		mu     sync.Mutex        // Mutex untuk menghindari race condition
		wg     sync.WaitGroup    // WaitGroup untuk menunggu semua goroutine DFS selesai
		count  = 0               // Jumlah pohon yang sudah ditemukan
	)

	treeChan := make(chan *RecipeTreeNode, maxTreeCount) // Channel untuk menyimpan hasil tree secara concurrent

	// Iterasi DFS untuk setiap resep yang memungkinkan dalam menghasilkan node target
	for _, recipe := range targetGraphNode.RecipesToMakeThisElement {
		wg.Add(1)
		go func(r *Recipe) {
			defer wg.Done()

			// Konfigurasi delay untuk update ExploringTree pada FE Visualization
			if delayMs > 0 {
				time.Sleep(time.Duration(delayMs) * time.Millisecond)
			}

			// Tambah hitungan node yang dieksplorasi
			// Aman untuk goroutine
			atomic.AddInt32(globalNodeCounter, 1)

			// Recurssion DFS ke elemen kiri dan kanan dari resep
			leftTrees, err1 := DFSFindTrees(nil, r.ElementOne, maxTreeCount, signalTreeChange, globalStartTime, globalNodeCounter, delayMs)
			if err1 != nil {
				return
			}

			rightTrees, err2 := DFSFindTrees(nil, r.ElementTwo, maxTreeCount, signalTreeChange, globalStartTime, globalNodeCounter, delayMs)
			if err2 != nil {
				return
			}

			// Kombinasikan pasangan tree kiri dan tree kanan untuk membentuk root
			for _, lt := range leftTrees {
				mu.Lock()
				if count >= maxTreeCount {
					mu.Unlock()
					break
				}
				mu.Unlock()

				for _, rt := range rightTrees {
					mu.Lock()
					if count >= maxTreeCount {
						mu.Unlock()
						break
					}

					// Buat node root baru dari dua subtree
					root := &RecipeTreeNode{
						Name:      targetGraphNode.Name,
						ImagePath: GetImagePath(targetGraphNode.ImagePath),
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
								root,
								int(time.Since(globalStartTime).Milliseconds()),
								atomic.LoadInt32(globalNodeCounter),
							)
						}()
					}

					// Kirim hasil tree ke channel
					treeChan <- root
					count++
					mu.Unlock()
				}
			}
		}(recipe)

		// Jika sudah cukup banyak pohon ditemukan
		// Maka hentikan break dari loop luar
		mu.Lock()
		if count >= maxTreeCount {
			mu.Unlock()
			break
		}
		mu.Unlock()
	}

	// Menunggu seluruh goroutine selesai dan menutup channel hasil
	go func() {
		wg.Wait()
		close(treeChan)
	}()

	// Mengumpulkan semua tree dari channel ke dalam result slice
	for tree := range treeChan {
		result = append(result, tree)
		if len(result) >= maxTreeCount {
			break
		}
	}

	// Jika tidak ada tree valid ditemukan, maka return error
	if len(result) == 0 {
		return nil, fmt.Errorf("no valid trees found for %s", targetGraphNode.Name)
	}

	return result, nil
}

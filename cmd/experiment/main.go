package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/kvoli/COMP90077_ASS2/pkg/rtree"
)

const testsize = 10e5

func experimentOne() ([]float64, []float64, []float64) {
	naive1, sorted1, fc1 := handleContr(testsize * 0.1)
	naive2, sorted2, fc2 := handleContr(testsize * 0.2)
	naive3, sorted3, fc3 := handleContr(testsize * 0.5)
	naive4, sorted4, fc4 := handleContr(testsize * 0.8)
	naive5, sorted5, fc5 := handleContr(testsize)
	return []float64{naive1, naive2, naive3, naive4, naive5}, []float64{sorted1, sorted2, sorted3, sorted4, sorted5}, []float64{fc1, fc2, fc3, fc4, fc5}

}

func experimentTwoA() ([]float64, []float64, []float64, []float64) {
	pnts := rtree.SortedPointSet(testsize)
	tt := rtree.ContrSorted(pnts)
	ct := rtree.ContrFC(pnts)
	tt1, ct1 := execWorkLoadTT(tt, genWorkLoad(testsize*0.01)), execWorkLoadCT(ct, genWorkLoad(testsize*0.01))
	tt2, ct2 := execWorkLoadTT(tt, genWorkLoad(testsize*0.02)), execWorkLoadCT(ct, genWorkLoad(testsize*0.02))
	tt3, ct3 := execWorkLoadTT(tt, genWorkLoad(testsize*0.05)), execWorkLoadCT(ct, genWorkLoad(testsize*0.05))
	tt4, ct4 := execWorkLoadTT(tt, genWorkLoad(testsize*0.1)), execWorkLoadCT(ct, genWorkLoad(testsize*0.1))
	tt5, ct5 := execWorkLoadTT(tt, genWorkLoad(testsize*0.2)), execWorkLoadCT(ct, genWorkLoad(testsize*0.2))

	if std(tt5) > 1 || std(tt4) > 0.25 || std(tt3) > 0.1 || std(ct5) > 0.1 || std(ct4) > 0.1 || std(ct3) > 0.1 {
		return experimentTwoA()
	}

	return []float64{avg(tt1), avg(tt2), avg(tt3), avg(tt4), avg(tt5)}, []float64{avg(ct1), avg(ct2), avg(ct3), avg(ct4), avg(ct5)}, []float64{std(tt1), std(tt2), std(tt3), std(tt4), std(tt5)}, []float64{std(ct1), std(ct2), std(ct3), std(ct4), std(ct5)}
}

func experimentTwoB() ([]float64, []float64, []float64, []float64) {
	work := genWorkLoad(0.05 * testsize)
	avgtt, avgct, stdtt, stdct := make([]float64, 10, 10), make([]float64, 10, 10), make([]float64, 10, 10), make([]float64, 10, 10)
	for i := 1; i < 11; i++ {
		pnts := rtree.SortedPointSet(int(math.Pow(2.0, float64(i)) * 10e2))
		tt := rtree.ContrSorted(pnts)
		ct := rtree.ContrFC(pnts)
		ttTime := execWorkLoadTT(tt, work)
		ctTime := execWorkLoadCT(ct, work)
		avgtt[i-1], stdtt[i-1] = avg(ttTime), std(ttTime)
		avgct[i-1], stdct[i-1] = avg(ctTime), std(ctTime)
	}
	return avgtt, avgct, stdtt, stdct
}

func avg(data []float64) float64 {
	count := len(data)
	total := 0.0
	for _, v := range data {
		total += v
	}
	return (total / float64(count))
}

func max(data []float64) float64 {
	best := 0.0
	for _, v := range data {
		if v > best {
			best = v
		}
	}
	return best
}

func min(data []float64) float64 {
	best := 10e6
	for _, v := range data {
		if v < best {
			best = v
		}
	}
	return best
}

func std(data []float64) float64 {
	s := 0.0
	u := avg(data)
	for _, v := range data {
		s += (v - u) * (v - u)
	}
	return math.Sqrt(s / float64(len(data)))
}

func execWorkLoadTT(tt *rtree.OTreeNode, work []*rtree.Query) []float64 {
	res := make([]float64, 100, 100)
	for i, v := range work {
		res[i] = execTQuery(tt, v)
	}
	return res
}

func execTQuery(tt *rtree.OTreeNode, q *rtree.Query) float64 {
	start := time.Now()
	rtree.QueryOTree(q.A, q.B, tt)
	return (time.Now().Sub(start)).Seconds() * 1000
}

func execWorkLoadCT(ct *rtree.CTreeNode, work []*rtree.Query) []float64 {
	res := make([]float64, 100, 100)
	for i, v := range work {
		res[i] = execCQuery(ct, v)
	}
	return res
}

func execCQuery(ct *rtree.CTreeNode, q *rtree.Query) float64 {
	start := time.Now()
	rtree.QueryCTree(q.A, q.B, ct)
	return (time.Now().Sub(start)).Seconds() * 1000
}

func genWorkLoad(s int) []*rtree.Query {
	res := make([]*rtree.Query, 100, 100)
	for i := 0; i < 100; i++ {
		res[i] = rtree.GQuery(s)
	}
	return res
}

func handleContr(size int) (float64, float64, float64) {
	pnts := rtree.SortedPointSet(size)
	start := time.Now()
	_ = rtree.ContrNaive(pnts)
	naiveFinish := (time.Now().Sub(start)).Seconds() * 1000

	start = time.Now()
	_ = rtree.ContrSorted(pnts)
	sortedFinish := (time.Now().Sub(start)).Seconds() * 1000

	start = time.Now()
	_ = rtree.ContrFC(pnts)
	fcFinish := (time.Now().Sub(start)).Seconds() * 1000

	if naiveFinish < sortedFinish-150 {
		return handleContr(size)
	}

	return naiveFinish, sortedFinish, fcFinish
}

func outputResultsA(filename string, size, naive, sorted, casc []float64) {
	f, _ := os.Create(fmt.Sprintf("%s.csv", filename))
	rows := make([][]string, 0)
	rows = append(rows, []string{"Size", "NaiveConstruction", "SortedConstruction", "CascadeConstruction"})
	fmt.Printf("%s %s %s %s\n", "Size", "NaiveConstruction", "SortedConstruction", "CascadeConstruction")
	for i := range naive {
		fmt.Printf("%9.0f\t%9.3f\t%9.3f\t%9.4f\n", size[i], naive[i], sorted[i], casc[i])
		rows = append(rows, []string{fmt.Sprintf("%.0f", size[i]), fmt.Sprintf("%.3f", naive[i]), fmt.Sprintf("%.3f", sorted[i]), fmt.Sprintf("%.3f", casc[i])})
	}
	_ = csv.NewWriter(f).WriteAll(rows)
}

func outputResultsB(filename string, size, avgtt, avgct, stdtt, stdct []float64) {
	f, _ := os.Create(fmt.Sprintf("%s.csv", filename))
	rows := make([][]string, 0)
	rows = append(rows, []string{"Size", "AverageYTree", "AverageCascade", "StdYTree", "StdCascade"})
	fmt.Printf("%s %s %s %s %s\n", "Size", "AverageYTree", "AverageCascade", "StdYTree", "StdCascade")
	for i := range avgtt {
		fmt.Printf("%7.0f\t%9.3f\t%9.3f\t%9.4f\t%9.4f\n", size[i], avgtt[i], avgct[i], stdtt[i], stdct[i])
		rows = append(rows, []string{fmt.Sprintf("%.0f", size[i]), fmt.Sprintf("%.3f", avgtt[i]), fmt.Sprintf("%.3f", avgct[i]), fmt.Sprintf("%.5f", stdtt[i]), fmt.Sprintf("%.5f", stdct[i])})
	}
	_ = csv.NewWriter(f).WriteAll(rows)
}

func main() {
	naive, sorted, casc := experimentOne()
	construct := []float64{testsize * 0.1, testsize * 0.2, testsize * 0.5, testsize * 0.8, testsize}
	outputResultsA("construction_times", construct, naive, sorted, casc)
	avgttA, avgcta, stdttA, stdctA := experimentTwoA()
	fixedN := []float64{testsize * 0.01, testsize * 0.02, testsize * 0.05, testsize * 0.1, testsize * 0.2}
	outputResultsB("query_fixed_n", fixedN, avgttA, avgcta, stdttA, stdctA)
	avgttB, avgctB, stdttB, stdctB := experimentTwoB()
	fixedS := []float64{2000.0, 4000.0, 8000.0, 16000.0, 32000.0, 64000.0, 128000.0, 256000.0, 512000.0, 1024000.0}
	outputResultsB("query_fixed_s", fixedS, avgttB, avgctB, stdttB, stdctB)
}

// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coloring

import (
	"context"
	"flag"
	"testing"
	"time"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/graph6"
	"gonum.org/v1/gonum/graph/internal/set"
	"gonum.org/v1/gonum/graph/simple"
)

var runLong = flag.Bool("color.long", false, "run long exact coloring tests")

var coloringTests = []struct {
	name   string
	g      graph.Undirected
	colors int

	long bool

	partial map[int64]int

	dsatur      set.Ints
	randomized  set.Ints
	rlf         set.Ints
	sanSegundo  set.Ints
	welshPowell set.Ints
}{
	{
		name:   "empty",
		g:      simple.NewUndirectedGraph(),
		colors: 0,
	},
	{
		name:   "singleton", // https://hog.grinvin.org/ViewGraphInfo.action?id=1310
		g:      graph6.Graph("@"),
		colors: 1,
	},
	{
		name:       "kite", // https://hog.grinvin.org/ViewGraphInfo.action?id=782
		g:          graph6.Graph("DzC"),
		colors:     3,
		randomized: setOf(3, 4),

		partial: map[int64]int{1: 2},
	},
	{
		name:   "triangle+1",
		g:      graph6.Graph("Cw"),
		colors: 3,

		partial: map[int64]int{1: 2},
	},
	{
		name:       "bipartite halves",
		g:          graph6.Graph("G?]uf?"),
		colors:     2,
		randomized: setOf(2, 3, 4),

		partial: map[int64]int{1: 3},
	},
	{
		name:        "bipartite alternating",
		g:           graph6.Graph("GKUalO"),
		colors:      2,
		randomized:  setOf(2, 3, 4),
		welshPowell: setOf(2, 3, 4),

		partial: map[int64]int{1: 1},
	},
	{
		name:   "3/4 bipartite", // https://hog.grinvin.org/ViewGraphInfo.action?id=466
		g:      graph6.Graph("F?~v_"),
		colors: 2,

		partial: map[int64]int{1: 1},
	},
	{
		name:        "cubical", // https://hog.grinvin.org/ViewGraphInfo.action?id=1022
		g:           graph6.Graph("Gs@ipo"),
		colors:      2,
		randomized:  setOf(2, 3, 4),
		welshPowell: setOf(2, 3),

		partial: map[int64]int{1: 1},
	},
	{
		name:       "HoG33106", // https://hog.grinvin.org/ViewGraphInfo.action?id=33106
		g:          graph6.Graph("K???WWKKxf]C"),
		colors:     2,
		randomized: setOf(2, 3, 4),

		partial: map[int64]int{1: 1},
	},
	{
		name:        "HoG41237", // https://hog.grinvin.org/ViewGraphInfo.action?id=41237
		g:           graph6.Graph("~?BCs?GO?@?O??@?A????G?A?????G??O????????????????????G????_???O??????????????????????OC????A???????????????@@?????A?????????????????????????C??????G??????O?????????????????????????G?@????????@???????D????????_??????CG???????@???????????????????????ACG???????EA???????O?G??????S????????C???_?????O??C?????@??????????W????????B??????????????O?????G??????????_????o????A????A_???@O_????????@?O????G???W???????????@G?????????A??????@?????S??????????????????????????????B????????????K????????????????????????@_????????????o??G????????????A?????????????_?????????_??G?????????G????@_???????????C?_????????????@????????A???A?????????A???O?????????????C???????????K??_??????????@O??_??????????C??@???????????G??????????E????????A??G??????????????????C????A???@???????????O??O???????????????_????????????@_?_????????????@O?C?????????????_?@?????????????G????C???O???????????????G?O????????????????O???????O???C????????????O?A?????????????????_??????????????@_C???????????????I?C??????????????@??G???????????????_???CA??????????????????_???????O????????O???????????????A??????????????????G?????????????????_C?????????????????OC?????????????????g????????D?????????????????@?G?????????????????GA??????????????????K??????????????????@O???????????????G????????G????????????????????P???????????????????OO???????????????????@_??????????????????????????????????????????????????????????????????@?????????????????????_???A?????????????????????????????????????????????????????????????O?????????????????????_??@????????????????????????????????????????G???????????@????????????????????????????????C?????????????????????????????????????????????????K??????????????????????I??????????????????????C_???????????????????_C????????????????????????K???????????????????????K???????????????????????`???????????????????????`??????????????????????I???????????????????????C?O?????????????????????A?O??????????????????????O?G??????????????????????@?C???????????????????????H????????????????????????OG???????????????????????GC?????????????????????????B?????????????????????????I?????????????????????????????????????????????????????????????????????????????????????????????????????????????c?????????????????????????CO???????????????????????A?A????????????????????????C?@????????????????????????C??_???????????????????????O??O?????????????????????????o???????????????????????????K???????????????????????????E??????????????????????????@_?????????????????????????D????????????????????????????`???????????????????????????????CG??????????????????????????@?O??????????????????????????@O???????????????????????????CG???????????????????????????C@???????????????????????????A?_?????????????????????C?????????????????????????????G?????????????????????????????@????????????????????????????????????AA?????????????????????????????B?????????????????????????????AG?????????????????????????????G_????????????????????????????CA?????????????????????????????CG??????????????????????????A???????????????????????????????O??????????????????????????????G????A?????????????????????????????AAG?????????????????????????????CCC??????????????????????????????cG?????????????????????????????@CA?????????????????????????????AG@??????????????????????????????P?"),
		colors:      2,
		randomized:  setOf(2, 3, 4),
		welshPowell: setOf(2, 3),

		partial: map[int64]int{1: 1},
	},

	// Test graphs from Leighton doi:10.6028/jres.084.024.
	{
		name:       "exams", // Figure 4 induced by the nodes in E.
		g:          graph6.Graph("EDZw"),
		colors:     3,
		randomized: setOf(3, 4),

		partial: map[int64]int{1: 2},
	},
	{
		name:       "exam_scheduler_1", // Figure 4.
		g:          graph6.Graph("JQ?HaWN{l~?"),
		colors:     5,
		randomized: setOf(5, 6),
	},
	{
		name:       "exam_scheduler_2", // Figure 11.
		g:          graph6.Graph("GTPHz{"),
		colors:     4,
		randomized: setOf(4, 5),

		partial: map[int64]int{1: 3},
	},

	// Test graph from Brélaz doi:10.1145/359094.359101.
	{
		name:        "Brélaz", // Figure 1.
		g:           graph6.Graph(`I??GG\rmg`),
		colors:      3,
		randomized:  setOf(3, 4),
		rlf:         setOf(3, 4),
		welshPowell: setOf(3, 4),

		partial: map[int64]int{1: 2},
	},

	// Test graph from Lima and Carmo doi:10.22456/2175-2745.80721.
	{
		name:       "Lima Carmo", // Figure 2.
		g:          graph6.Graph("Gh]S?G"),
		colors:     3,
		randomized: setOf(3, 4),
		sanSegundo: setOf(3, 4),

		partial: map[int64]int{1: 2},
	},

	// Test graph from San Segundo doi:10.1016/j.cor.2011.10.008.
	{
		name:       "San Segundo", // Figure 1 A.
		g:          graph6.Graph(`HMn\r\v`),
		colors:     5,
		randomized: setOf(5, 6),

		partial: map[int64]int{1: 4},
	},

	{
		name:   "tetrahedron", // https://hog.grinvin.org/ViewGraphInfo.action?id=74
		g:      graph6.Graph("C~"),
		colors: 4,

		partial: map[int64]int{1: 3},
	},

	{
		name:   "triangle", // https://hog.grinvin.org/ViewGraphInfo.action?id=1374
		g:      graph6.Graph("Bw"),
		colors: 3,

		partial: map[int64]int{1: 2},
	},
	{
		name:   "square", // https://hog.grinvin.org/ViewGraphInfo.action?id=674
		g:      graph6.Graph("Cl"),
		colors: 2,

		partial: map[int64]int{1: 1},
	},
	{
		name:   "cycle-5", // https://hog.grinvin.org/ViewGraphInfo.action?id=340
		g:      graph6.Graph("Dhc"),
		colors: 3,

		partial: map[int64]int{1: 2},
	},
	{
		name:       "cycle-6", // https://hog.grinvin.org/ViewGraphInfo.action?id=670
		g:          graph6.Graph("EhEG"),
		colors:     2,
		randomized: setOf(2, 3),

		partial: map[int64]int{1: 1},
	},

	{
		name:   "wheel-5", // https://hog.grinvin.org/ViewGraphInfo.action?id=442
		g:      graph6.Graph("D|s"),
		colors: 3,

		partial: map[int64]int{1: 2},
	},
	{
		name:   "wheel-6", // https://hog.grinvin.org/ViewGraphInfo.action?id=204
		g:      graph6.Graph("E|fG"),
		colors: 4,

		partial: map[int64]int{1: 3},
	},
	{
		name:       "wheel-7",
		g:          graph6.Graph("F|eMG"),
		colors:     3,
		randomized: setOf(3, 4),

		partial: map[int64]int{1: 2},
	},

	{
		name:        "sudoku board", // The graph of the constraints of a sudoku puzzle board.
		g:           graph6.Graph("~?@P~~~~~~wF?{BFFFFbbwF~?~{B~~?wF?wF_[BF?wwwFFb_[^?wF~?wF~_[B~{?_C?OA?OC_C?_W_C?_wOA?O{C?_C^?_C?fwA?OA~_C?_F~_C?_F?OA?OF?c?_CB_W_C?_FFA?OA?w{C?_CBbwC?_C?~wA?OA?~{?_C?_^~_C?_F?wA?OA?wF?c?_CB_[BC?_C?wFFA?OA?wFF__C?_[BbwC?_C?wF~?OA?OF?~{?_C?_[B~{?_C?_C?_A?OA?OA?OC_C?_C?_CBC?_C?_C?_wOA?OA?OAF__C?_C?_C^?_C?_C?_C~?OA?OA?OA~_C?_C?_C?~{?_C?_C?_F?OA?OA?OA?wC_C?_C?_CB_W_C?_C?_C?wwOA?OA?OA?w{C?_C?_C?_[^?_C?_C?_C?~wA?OA?OA?OF~_C?_C?_C?_^~_C?_C?_C?wF?OA?OA?OA?wF?c?_C?_C?_[B_W_C?_C?_C?wFFA?OA?OA?OF?w{C?_C?_C?_[BbwC?_C?_C?_F?~wA?OA?OA?OF?~{?_C?_C?_CB_^~"),
		colors:      9,
		randomized:  setOf(9, 10, 11, 12, 13),
		sanSegundo:  setOf(9, 10),
		welshPowell: setOf(9, 10, 11, 12, 13, 14),

		partial: map[int64]int{1: 3},
	},
	{
		name:        "sudoku problem", // The constraint graph for the problem in the sudoku example.
		g:           graph6.Graph("~?@Y~~~~~~?F|_B?F?F_Bw?~?F{?^w?w??wC?[Bzwww?FF_?[^??F~??F~~kB~w?wF?v~?wFz{B_W?F?w|~F?w{?B_[^??F?~w??wF~_?B_^~?C?_C??A?OA?_?_C?_W?C?_CF??OA?O~|__C?bw??_C?fw??OA?V~}_C?_F~?C?_C?wD~OA?OF?_?_C?_[B??_C?_FF}wOA?OFF_?C?_CBbw??_C?_F~~oA?OA?~{??C?_CB~~^_C?_F?w??OA?OF?wC?C?_CB_[B|w_C?_F?ww?A?OA?wFF_?C?_CB_[^??C?_C?wF~??A?OA?wF~_??_C?_[B~w?_C?_C?_C??A?OA?OA?OC?C?_C?_C?_W?C?_C?_C?_~}A?OA?OA?O{??_C?_C?_C^vwC?_C?_C?f|~?OA?OA?OA~_??_C?_C?_F~|{?_C?_C?_F??A?OA?OA?OF?_?_C?_C?_CB_W?C?_C?_C?_FF??OA?OA?OA?w|~__C?_C?_CBbw??_C?_C?_C?~w??OA?OA?OA?~{??C?_C?_C?_^~?C?_C?_C?_F?w??OA?OA?OA?wF?_?_C?_C?_CB_[B??_C?_C?_C?wFF??OA?OA?OA?wFF_?C?_C?_C?_[Bbw??_C?_C?_C?wF~}wA?OA?OA?OF?~{??C?_C?_C?_[B~w"),
		colors:      9,
		dsatur:      setOf(9, 10, 11),
		randomized:  setOf(9, 10, 11, 12, 13, 14, 15, 16, 17),
		rlf:         setOf(9, 10, 11, 12),
		sanSegundo:  setOf(9, 10, 11),
		welshPowell: setOf(9, 10, 11, 12, 13),

		partial: map[int64]int{1: 3},
	},

	// Test graphs from NetworkX.
	{
		name:       "cs_shc",
		g:          graph6.Graph("Djs"),
		colors:     3,
		randomized: setOf(3, 4),
	},
	{
		name:       "gis_hc",
		g:          graph6.Graph("E?ow"),
		colors:     2,
		randomized: setOf(2, 3),
	},
	{
		name:       "gis_shc|rs_shc", // https://hog.grinvin.org/ViewGraphInfo.action?id=594
		g:          graph6.Graph("CR"),
		colors:     2,
		randomized: setOf(2, 3),
	},
	{
		name:        "lf_hc",
		g:           graph6.Graph(`F\^E?`),
		colors:      3,
		dsatur:      setOf(3, 4),
		randomized:  setOf(3, 4),
		rlf:         setOf(3, 4),
		sanSegundo:  setOf(3, 4),
		welshPowell: setOf(3, 4),
	},
	{
		name:        "lf_shc",
		g:           graph6.Graph("ELQ?"),
		colors:      2,
		randomized:  setOf(2, 3),
		welshPowell: setOf(2, 3),
	},
	{
		name:        "lfi_hc",
		g:           graph6.Graph("Hhe[b@_"),
		colors:      3,
		dsatur:      setOf(3, 4),
		randomized:  setOf(3, 4),
		sanSegundo:  setOf(3, 4),
		welshPowell: setOf(3, 4),
	},
	{
		name:        "lfi_shc|slf_shc",
		g:           graph6.Graph("FheZ?"),
		colors:      3,
		dsatur:      setOf(3, 4),
		randomized:  setOf(3, 4),
		welshPowell: setOf(3, 4),
	},
	{
		name:       "no_solo", // https://hog.grinvin.org/ViewGraphInfo.action?id=264, https://hog.grinvin.org/ViewGraphInfo.action?id=498
		g:          graph6.Graph("K????AccaQHG"),
		colors:     2,
		randomized: setOf(2, 3),
	},
	{
		name:        "sl_hc",
		g:           graph6.Graph("Gzg[Yk"),
		colors:      4,
		dsatur:      setOf(4, 5),
		randomized:  setOf(4, 5),
		sanSegundo:  setOf(4, 5),
		welshPowell: setOf(4, 5),
	},
	{
		name:        "sl_shc",
		g:           graph6.Graph("E{Sw"),
		colors:      3,
		dsatur:      setOf(3, 4),
		randomized:  setOf(3, 4),
		welshPowell: setOf(3, 4),
	},
	{
		name:        "slf_hc",
		g:           graph6.Graph("G}`?W["),
		colors:      3,
		dsatur:      setOf(3, 4),
		randomized:  setOf(3, 4),
		rlf:         setOf(3, 4),
		welshPowell: setOf(3, 4),
	},
	{
		name:        "sli_hc",
		g:           graph6.Graph("H{czYtt"),
		colors:      4,
		dsatur:      setOf(4, 5),
		randomized:  setOf(4, 5),
		welshPowell: setOf(4, 5),
	},
	{
		name:        "sli_shc",
		g:           graph6.Graph("FxdSW"),
		colors:      3,
		dsatur:      setOf(3, 4),
		randomized:  setOf(3, 4),
		sanSegundo:  setOf(3, 4),
		welshPowell: setOf(3, 4),
	},
	{
		name:       "rsi_shc",
		g:          graph6.Graph("EheW"),
		colors:     3,
		randomized: setOf(3, 4),
	},
	{
		name:       "V_plus_not_in_A_cal",
		g:          graph6.Graph("HQQ?W__"),
		colors:     2,
		randomized: setOf(2, 3),
	},

	// DIMACS queens graphs
	{
		name:        "queen5_5",
		g:           graph6.Graph("X~~FJk~F|KIxizS^dF{iWQjcdV[dFyQb}KiWOdVHAT\\acg~acg~"),
		colors:      5,
		dsatur:      setOf(5),
		randomized:  setOf(5, 6, 7, 8, 9, 10),
		rlf:         setOf(5),
		sanSegundo:  setOf(5),
		welshPowell: setOf(5, 6, 7),
	},
	{
		name:        "queen6_6",
		g:           graph6.Graph("c~~}FDrMw~`~goSwtMYhvIF{SN{dEAQfCehrcTMyPO~ca`~acgoPQSwcCtMWcahvaQIF|CcSN{KSdEAAIQfC__ehrCCcTMwSQPO~ogca`~"),
		colors:      7,
		dsatur:      setOf(7, 8, 9),
		randomized:  setOf(7, 8, 9, 10, 11, 12),
		rlf:         setOf(7, 8),
		sanSegundo:  setOf(7, 8),
		welshPowell: setOf(7, 8, 9),
	},
	{
		name:        "queen7_7",
		g:           graph6.Graph("p~~~}B`[XrbnB~@~sKAb`iMLS[xS\\wgN{IB~cSKAPPocibbcibfQDPvc`O^wcIB~aQIE@CcS[HCdS[WaQiM]GcIbnPC`O^xCQD@~ogca`_Ogcab`GCQTPpa@CdS[wQGcIbn`GaOgN|APC`O^{EDCcSKA@AaQIMC_OGcibbCA@CdS[wOHCQDPv_gQGcIB~_gQGcIB~"),
		colors:      7,
		dsatur:      setOf(7, 8, 9, 10, 11),
		randomized:  setOf(7, 8, 9, 10, 11, 12, 13),
		rlf:         setOf(7, 8, 9),
		sanSegundo:  setOf(7, 8, 9),
		welshPowell: setOf(7, 8, 9, 10, 11, 12),
	},
	{
		name:        "queen8_8",
		long:        true,
		g:           graph6.Graph("~?@?~~~~~?wJ`fFFNBn_~wF~gK@OwLOwYg[[iFNDOzwSB~_gF~cIB?QDB_cdOw[ciFFQQg[{cDOzwcD?~wQA_^}GcIB?PC`OwHCQTB`aKciFFaCciFNPAOTBncOcD?~waC_gF~`GaOgK@APC`OwHAPCdOwW_GqQg[[GaCciFN`COcDOzyCPAOSB~cGaC_gF~_gQGcIB?OSHCQDB_c@APCdOwW_GAKciFFA?aGQQg[{C`COcDOz{C`COcD?~yAOaGQA_^}@_gQGcIB?OCDAPC`OwH?OCHCQTB`a?_GAKciFFA?_GaCciFN@?QCPAOTBn_SC`COcD?~{A_cGaC_gF~"),
		colors:      9,
		dsatur:      setOf(9, 10, 11, 12),
		randomized:  setOf(9, 10, 11, 12, 13, 14, 15),
		rlf:         setOf(9, 10),
		sanSegundo:  setOf(9, 10, 11),
		welshPowell: setOf(9, 10, 11, 12, 13, 14, 15),
	},
}

func setOf(vals ...int) set.Ints {
	s := make(set.Ints)
	for _, v := range vals {
		s.Add(v)
	}
	return s
}

func TestDsatur(t *testing.T) {
	for _, test := range coloringTests {
		for _, partial := range []map[int64]int{nil, test.partial} {
			k, colors, err := Dsatur(test.g, partial)

			if partial == nil && k != test.colors && !test.dsatur.Has(k) {
				t.Errorf("unexpected chromatic number for %q: got:%d want:%d or in %v\ncolors:%v",
					test.name, k, test.colors, test.dsatur, colors)
			}
			if s := Sets(colors); len(s) != k {
				t.Errorf("mismatch between number of color sets and k: |sets|=%d k=%d", len(s), k)
			}
			if missing, ok := isCompleteColoring(colors, test.g); !ok {
				t.Errorf("incomplete coloring for %q: missing %d\ngot:%v", test.name, missing, colors)
			}
			if xid, yid, ok := isValidColoring(colors, test.g); !ok {
				t.Errorf("invalid coloring for %q: %d--%d match color\ncolors:%v",
					test.name, xid, yid, colors)
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			for id, c := range partial {
				if colors[id] != c {
					t.Errorf("coloring not consistent with input partial for %q:\ngot:%v\nwant superset of:%v",
						test.name, colors, partial)
					break
				}
			}
		}
	}
}

func TestDsaturExact(t *testing.T) {
	timeout := time.Microsecond
	for _, test := range coloringTests {
		for _, useTimeout := range []bool{false, true} {
			if test.long && !*runLong && !useTimeout {
				continue
			}
			var term Terminator
			cancel := func() {}
			if useTimeout {
				term, cancel = context.WithTimeout(context.Background(), timeout)
			}
			k, colors, err := DsaturExact(term, test.g)
			cancel()
			if k != test.colors && (useTimeout && !test.dsatur.Has(k)) {
				t.Errorf("unexpected chromatic number for %q timeout=%t: got:%d want:%d or in %v\ncolors:%v",
					test.name, useTimeout, k, test.colors, test.dsatur, colors)
			}
			if s := Sets(colors); len(s) != k {
				t.Errorf("mismatch between number of color sets and k: |sets|=%d k=%d", len(s), k)
			}
			if missing, ok := isCompleteColoring(colors, test.g); !ok {
				t.Errorf("incomplete coloring for %q: missing %d\ngot:%v", test.name, missing, colors)
			}
			if xid, yid, ok := isValidColoring(colors, test.g); !ok {
				t.Errorf("invalid coloring for %q: %d--%d match color\ncolors:%v",
					test.name, xid, yid, colors)
			}
			if err != nil && !useTimeout {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestRandomized(t *testing.T) {
	for seed := uint64(1); seed <= 1000; seed++ {
		for _, test := range coloringTests {
			for _, partial := range []map[int64]int{nil, test.partial} {
				k, colors, err := Randomized(test.g, partial, rand.NewSource(seed))

				if partial == nil && k != test.colors && !test.randomized.Has(k) {
					t.Errorf("unexpected chromatic number for %q with seed=%d: got:%d want:%d or in %v\ncolors:%v",
						test.name, seed, k, test.colors, test.randomized, colors)
				}
				if s := Sets(colors); len(s) != k {
					t.Errorf("mismatch between number of color sets and k: |sets|=%d k=%d", len(s), k)
				}
				if missing, ok := isCompleteColoring(colors, test.g); !ok {
					t.Errorf("incomplete coloring for %q: missing %d\ngot:%v", test.name, missing, colors)
				}
				if xid, yid, ok := isValidColoring(colors, test.g); !ok {
					t.Errorf("invalid coloring for %q: %d--%d match color\ncolors:%v",
						test.name, xid, yid, colors)
				}
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				for id, c := range partial {
					if colors[id] != c {
						t.Errorf("coloring not consistent with input partial for %q:\ngot:%v\nwant superset of:%v",
							test.name, colors, partial)
						break
					}
				}
			}
		}
	}
}

func TestRecursiveLargestFirst(t *testing.T) {
	for _, test := range coloringTests {
		k, colors := RecursiveLargestFirst(test.g)
		if k != test.colors && !test.rlf.Has(k) {
			t.Errorf("unexpected chromatic number for %q: got:%d want:%d",
				test.name, k, test.colors)
		}
		if s := Sets(colors); len(s) != k {
			t.Errorf("mismatch between number of color sets and k: |sets|=%d k=%d", len(s), k)
		}
		if missing, ok := isCompleteColoring(colors, test.g); !ok {
			t.Errorf("incomplete coloring for %q: missing %d\ngot:%v", test.name, missing, colors)
		}
		if xid, yid, ok := isValidColoring(colors, test.g); !ok {
			t.Errorf("invalid coloring for %q: %d--%d match color\ncolors:%v",
				test.name, xid, yid, colors)
		}
	}
}

func TestSanSegundo(t *testing.T) {
	for _, test := range coloringTests {
		for _, partial := range []map[int64]int{nil, test.partial} {
			k, colors, err := SanSegundo(test.g, partial)

			if partial == nil && k != test.colors && !test.sanSegundo.Has(k) {
				t.Errorf("unexpected chromatic number for %q: got:%d want:%d\ncolors:%v",
					test.name, k, test.colors, colors)
			}
			if s := Sets(colors); len(s) != k {
				t.Errorf("mismatch between number of color sets and k: |sets|=%d k=%d", len(s), k)
			}
			if missing, ok := isCompleteColoring(colors, test.g); !ok {
				t.Errorf("incomplete coloring for %q: missing %d\ngot:%v", test.name, missing, colors)
			}
			if xid, yid, ok := isValidColoring(colors, test.g); !ok {
				t.Errorf("invalid coloring for %q: %d--%d match color\ncolors:%v",
					test.name, xid, yid, colors)
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			for id, c := range partial {
				if colors[id] != c {
					t.Errorf("coloring not consistent with input partial for %q:\ngot:%v\nwant superset of:%v",
						test.name, colors, partial)
					break
				}
			}
		}
	}
}

func TestWelshPowell(t *testing.T) {
	for _, test := range coloringTests {
		for _, partial := range []map[int64]int{nil, test.partial} {
			k, colors, err := WelshPowell(test.g, partial)

			if partial == nil && k != test.colors && !test.welshPowell.Has(k) {
				t.Errorf("unexpected chromatic number for %q: got:%d want:%d",
					test.name, k, test.colors)
			}
			if s := Sets(colors); len(s) != k {
				t.Errorf("mismatch between number of color sets and k: |sets|=%d k=%d", len(s), k)
			}
			if missing, ok := isCompleteColoring(colors, test.g); !ok {
				t.Errorf("incomplete coloring for %q: missing %d\ngot:%v", test.name, missing, colors)
			}
			if xid, yid, ok := isValidColoring(colors, test.g); !ok {
				t.Errorf("invalid coloring for %q: %d--%d match color\ncolors:%v",
					test.name, xid, yid, colors)
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			for id, c := range partial {
				if colors[id] != c {
					t.Errorf("coloring not consistent with input partial for %q:\ngot:%v\nwant superset of:%v",
						test.name, colors, partial)
					break
				}
			}
		}
	}
}

var newPartialTests = []struct {
	partial   map[int64]int
	g         graph.Undirected
	wantValid bool
}{
	{
		partial: map[int64]int{0: 1, 1: 1},
		g: undirectedGraphFrom([]intset{
			0: linksTo(1, 2),
			1: linksTo(2),
			2: nil,
		}),
		wantValid: false,
	},
	{
		partial: map[int64]int{0: 0, 1: 1},
		g: undirectedGraphFrom([]intset{
			0: linksTo(1, 2),
			1: linksTo(2),
			2: nil,
		}),
		wantValid: true,
	},
	{
		partial: map[int64]int{0: 0, 1: 1, 3: 3},
		g: undirectedGraphFrom([]intset{
			0: linksTo(1, 2),
			1: linksTo(2),
			2: nil,
		}),
		wantValid: false,
	},
	{
		partial: map[int64]int{0: 0, 1: 1, 2: 2},
		g: undirectedGraphFrom([]intset{
			0: linksTo(1, 2),
			1: linksTo(2),
			2: nil,
		}),
		wantValid: true,
	},
	{
		partial: nil,
		g: undirectedGraphFrom([]intset{
			0: linksTo(1, 2),
			1: linksTo(2),
			2: nil,
		}),
		wantValid: true,
	},
}

func TestNewPartial(t *testing.T) {
	for i, test := range newPartialTests {
		gotPartial, gotValid := newPartial(test.partial, test.g)
		if gotValid != test.wantValid {
			t.Errorf("unexpected validity for test %d: got:%t want:%t",
				i, gotValid, test.wantValid)
		}
		xid, yid, ok := isValidColoring(gotPartial, test.g)
		if !ok {
			t.Errorf("invalid partial returned for test %d: %d--%d match color\ncolors:%v",
				i, xid, yid, gotPartial)
		}

	}
}

func isCompleteColoring(colors map[int64]int, g graph.Undirected) (missing int64, ok bool) {
	for _, n := range graph.NodesOf(g.Nodes()) {
		if _, ok := colors[n.ID()]; !ok {
			return n.ID(), false
		}
	}
	return 0, true
}

func isValidColoring(colors map[int64]int, g graph.Undirected) (x, y int64, ok bool) {
	for xid, c := range colors {
		to := g.From(xid)
		for to.Next() {
			yid := to.Node().ID()
			if oc, ok := colors[yid]; ok && c == oc {
				return xid, yid, false
			}
		}
	}
	return 0, 0, true
}

// intset is an integer set.
type intset map[int64]struct{}

func linksTo(i ...int64) intset {
	if len(i) == 0 {
		return nil
	}
	s := make(intset)
	for _, v := range i {
		s[v] = struct{}{}
	}
	return s
}

func undirectedGraphFrom(g []intset) graph.Undirected {
	dg := simple.NewUndirectedGraph()
	for u, e := range g {
		for v := range e {
			dg.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
		}
	}
	return dg
}

func BenchmarkColoring(b *testing.B) {
	for _, bench := range coloringTests {
		b.Run(bench.name, func(b *testing.B) {
			b.Run("Dsatur", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _, err := Dsatur(bench.g, nil)
					if err != nil {
						b.Fatalf("coloring failed: %v", err)
					}
				}
			})
			if !bench.long || *runLong {
				b.Run("DsaturExact", func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						_, _, err := DsaturExact(nil, bench.g)
						if err != nil {
							b.Fatalf("coloring failed: %v", err)
						}
					}
				})
			}
			b.Run("RecursiveLargestFirst", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					RecursiveLargestFirst(bench.g)
				}
			})
			b.Run("SanSegundo", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _, err := SanSegundo(bench.g, nil)
					if err != nil {
						b.Fatalf("coloring failed: %v", err)
					}
				}
			})
			b.Run("WelshPowell", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _, err := WelshPowell(bench.g, nil)
					if err != nil {
						b.Fatalf("coloring failed: %v", err)
					}
				}
			})
		})
	}
}
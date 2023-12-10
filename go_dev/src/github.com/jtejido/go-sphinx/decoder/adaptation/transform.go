package adaptation

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jtejido/linear"
)

type Transform struct {
	a            [][][][]float32
	b            [][][]float32
	loader       *tiedstate.Sphinx3Loader
	nrOfClusters int
}

func NewTransform(loader *tiedstate.Sphinx3Loader, nrOfClusters int) *Transform {
	return &Transform{
		loader:       loader,
		nrOfClusters: nrOfClusters,
	}
}

/**
 * Used for access to A matrix.
 *
 * @return A matrix (representing A from A*x + B = C)
 */
func (t *Transform) As() [][][][]float32 {
	return t.a
}

/**
 * Used for access to B matrix.
 *
 * @return B matrix (representing B from A*x + B = C)
 */
func (t *Transform) Bs() [][][]float32 {
	return t.b
}

/**
 * Writes the transformation to file in a format that could further be used
 * in Sphinx3 and Sphinx4.
 *
 * @param filePath
 *            path to store transform matrix
 * @param index
 *            index of transform to store
 * @throws Exception
 *             if something went wrong
 */
func (t *Transform) Store(filePath string, index int) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// nMllrClass
	fmt.Fprintln(writer, "1")
	fmt.Fprintln(writer, t.loader.NumStreams())

	for i := 0; i < t.loader.NumStreams(); i++ {
		fmt.Fprintln(writer, t.loader.VectorLength()[i])

		for j := 0; j < t.loader.VectorLength()[i]; j++ {
			for k := 0; k < t.loader.VectorLength()[i]; k++ {
				fmt.Fprint(writer, t.a[index][i][j][k])
				fmt.Fprint(writer, " ")
			}
			fmt.Fprintln(writer)
		}

		for j := 0; j < t.loader.VectorLength()[i]; j++ {
			fmt.Fprint(writer, t.b[index][i][j])
			fmt.Fprint(writer, " ")

		}
		fmt.Fprintln(writer)

		for j := 0; j < t.loader.VectorLength()[i]; j++ {
			fmt.Fprint(writer, "1.0 ")

		}
		fmt.Fprintln(writer)
	}

	return nil
}

/**
 * Used for computing the actual transformations (A and B matrices). These
 * are stored in As and Bs.
 */
func (t *Transform) computeMllrTransforms(regLs [][][][][]float64, regRs [][][][]float64) (err error) {
	var len int
	var solver linear.DecompositionSolver
	var coef linear.RealMatrix
	var vect, ABloc linear.RealVector
	var lud *linear.LUDecomposition

	for c := 0; c < t.nrOfClusters; c++ {
		t.a[c] = make([][][]float32, t.loader.NumStreams())
		t.b[c] = make([][]float32, t.loader.NumStreams())

		for i := 0; i < t.loader.NumStreams(); i++ {
			len = t.loader.VectorLength()[i]
			t.a[c][i] = make([][]float32, len)
			for ii := 0; ii < len; ii++ {
				t.a[c][i][ii] = make([]float32, len)
			}

			t.b[c][i] = make([]float32, len)

			for j := 0; j < len; j++ {
				coef, err = linear.NewArray2DRowRealMatrixFromSlices(regLs[c][i][j], false)
				if err != nil {
					return
				}

				lud, err = linear.NewLUDecomposition(coef)
				if err != nil {
					return
				}
				solver = lud.Solver()
				vect, err = linear.NewArrayRealVector(regRs[c][i][j], false)
				if err != nil {
					return
				}
				ABloc = solver.SolveVector(vect)

				for k := 0; k < len; k++ {
					t.a[c][i][j][k] = float32(ABloc.At(k))
				}

				t.b[c][i][j] = float32(ABloc.At(len))
			}
		}
	}
	return
}

func (t *Transform) Load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var numStreams, nMllrClass int
	if scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d", &nMllrClass)
	}

	if nMllrClass != 1 {
		return fmt.Errorf("Unexpected nMllrClass value: %d", nMllrClass)
	}

	if scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d", &numStreams)
	}

	t.a = make([][][][]float32, nMllrClass)
	t.b = make([][][]float32, nMllrClass)

	for i := 0; i < nMllrClass; i++ {
		t.a[i] = make([][][]float32, numStreams)
		t.b[i] = make([][]float32, numStreams)

		for j := 0; j < numStreams; j++ {
			var length int
			if scanner.Scan() {
				fmt.Sscanf(scanner.Text(), "%d", &length)
			}

			t.a[i][j] = make([][]float32, length)
			t.b[i][j] = make([]float32, length)

			for k := 0; k < length; k++ {
				t.a[i][j][k] = make([]float32, length)

				for l := 0; l < length; l++ {
					if scanner.Scan() {
						fmt.Sscanf(scanner.Text(), "%f", &t.a[i][j][k][l])
					}
				}
			}

			for k := 0; k < length; k++ {
				if scanner.Scan() {
					fmt.Sscanf(scanner.Text(), "%f", &t.b[i][j][k])
				}
			}

			for k := 0; k < length; k++ {
				// Skip MLLR variance scale
				if scanner.Scan() {
					fmt.Sscanf(scanner.Text(), "%f", new(float32))
				}
			}
		}
	}

	return scanner.Err()
}

/**
 * Stores in current object a transform generated on the provided stats.
 *
 * @param stats
 *            provided stats that were previously collected from Result
 *            objects.
 */
func (t *Transform) Update(stats *Stats) error {
	stats.FillRegLowerPart()
	t.a = make([][][][]float32, t.nrOfClusters)
	t.b = make([][][]float32, t.nrOfClusters)
	return t.computeMllrTransforms(stats.RegLs(), stats.RegRs())
}

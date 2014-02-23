package numerics

import (
	"math"
	"reflect"
)

const (
	FractionalDigitsToRoundTo = 10 // Anything smaller than this will be assumed to be rounding error in terms of equality matching
)

var (
	ErrorTolerance = math.Pow(0.1, FractionalDigitsToRoundTo) // e.g. 1/10^10
)

// Matrix represents an MxN matrix with double precision values.
type Matrix struct {
	Rows            int
	Columns         int
	matrixRowValues [][]float64
}

// Note: some properties like Determinant, Inverse, etc are properties instead
// of methods to make the syntax look nicer even though this sort of goes against
// Framework Design Guidelines that properties should be "cheap" since it could take
// a long time to compute these properties if the matrices are "big."

func NewMatrixFromValues(rows, columns int, allRowValues ...float64) Matrix {
	m := Matrix{
		Rows:            rows,
		Columns:         columns,
		matrixRowValues: make([][]float64, rows),
	}

	currentIndex := 0
	for currentRow := 0; currentRow < m.Rows; currentRow++ {
		m.matrixRowValues[currentRow] = make([]float64, m.Columns)

		for currentColumn := 0; currentColumn < m.Columns; currentColumn++ {
			if currentIndex < len(allRowValues) {
				m.matrixRowValues[currentRow][currentColumn] = allRowValues[currentIndex]
				currentIndex++
			}
		}
	}

	return m
}

func NewMatrix(rows, columns int, matrixRowValues [][]float64) Matrix {
	return Matrix{
		Rows:            rows,
		Columns:         columns,
		matrixRowValues: matrixRowValues,
	}
}

//         public double this[int row, int column]
//         {
//             get { return _MatrixRowValues[row][column]; }
//         }

func (m Matrix) Transpose() Matrix {

	// Just flip everything
	transposeMatrix := make([][]float64, m.Columns)
	for currentRowTransposeMatrix := 0; currentRowTransposeMatrix < m.Columns; currentRowTransposeMatrix++ {
		transposeMatrixCurrentRowColumnValues := make([]float64, m.Rows)
		transposeMatrix[currentRowTransposeMatrix] = transposeMatrixCurrentRowColumnValues

		for currentColumnTransposeMatrix := 0; currentColumnTransposeMatrix < m.Rows; currentColumnTransposeMatrix++ {
			transposeMatrixCurrentRowColumnValues[currentColumnTransposeMatrix] =
				m.matrixRowValues[currentColumnTransposeMatrix][currentRowTransposeMatrix]
		}
	}

	return Matrix{
		Rows:            m.Columns,
		Columns:         m.Rows,
		matrixRowValues: transposeMatrix,
	}
}

func (m Matrix) IsSquare() bool {
	return (m.Rows == m.Columns) && m.Rows > 0
}

func (m Matrix) Determinant() float64 {

	// Basic argument checking
	if !m.IsSquare() {
		panic("Matrix must be square!")
	}

	if m.Rows == 1 {
		// Really happy path :)
		return m.matrixRowValues[0][0]
	}

	if m.Rows == 2 {
		// Happy path!
		// Given:
		// | a b |
		// | c d |
		// The determinant is ad - bc
		a := m.matrixRowValues[0][0]
		b := m.matrixRowValues[0][1]
		c := m.matrixRowValues[1][0]
		d := m.matrixRowValues[1][1]
		return a*d - b*c
	}

	// I use the Laplace expansion here since it's straightforward to implement.
	// It's O(n^2) and my implementation is especially poor performing, but the
	// core idea is there. Perhaps I should replace it with a better algorithm
	// later.
	// See http://en.wikipedia.org/wiki/Laplace_expansion for details

	result := 0.0

	// I expand along the first row
	for currentColumn := 0; currentColumn < m.Columns; currentColumn++ {
		firstRowColValue := m.matrixRowValues[0][currentColumn]
		cofactor := m.GetCofactor(0, currentColumn)
		itemToAdd := firstRowColValue * cofactor
		result += itemToAdd
	}

	return result
}

func (m Matrix) Adjugate() Matrix {

	if !m.IsSquare() {
		panic("Matrix must be square!")
	}

	// See http://en.wikipedia.org/wiki/Adjugate_matrix
	if m.Rows == 2 {
		// Happy path!
		// Adjugate of:
		// | a b |
		// | c d |
		// is
		// | d -b |
		// | -c a |

		a := m.matrixRowValues[0][0]
		b := m.matrixRowValues[0][1]
		c := m.matrixRowValues[1][0]
		d := m.matrixRowValues[1][1]

		return NewSquareMatrix(d, -b, -c, a)
	}

	// The idea is that it's the transpose of the cofactors
	result := make([][]float64, m.Columns)

	for currentColumn := 0; currentColumn < m.Columns; currentColumn++ {
		result[currentColumn] = make([]float64, m.Rows)

		for currentRow := 0; currentRow < m.Rows; currentRow++ {
			result[currentColumn][currentRow] = m.GetCofactor(currentRow, currentColumn)
		}
	}

	return Matrix{m.Columns, m.Rows, result}
}

func (m Matrix) Inverse() Matrix {
	if (m.Rows == 1) && (m.Columns == 1) {
		return NewSquareMatrix(1.0 / m.matrixRowValues[0][0])
	}

	// Take the simple approach:
	// http://en.wikipedia.org/wiki/Cramer%27s_rule#Finding_inverse_matrix
	return MultiplyBy((1.0 / m.Determinant()), m.Adjugate())
}

func MultiplyBy(scalarValue float64, matrix Matrix) Matrix {
	rows := matrix.Rows
	columns := matrix.Columns
	newValues := make([][]float64, rows)

	for currentRow := 0; currentRow < rows; currentRow++ {
		newRowColumnValues := make([]float64, columns)
		newValues[currentRow] = newRowColumnValues

		for currentColumn := 0; currentColumn < columns; currentColumn++ {
			newRowColumnValues[currentColumn] = scalarValue * matrix.matrixRowValues[currentRow][currentColumn]
		}
	}

	return Matrix{rows, columns, newValues}
}

func Add(left, right Matrix) Matrix {
	if (left.Rows != right.Rows) || (left.Columns != right.Columns) {
		panic("Matrices must be of the same size")
	}

	// simple addition of each item

	resultMatrix := make([][]float64, left.Rows)

	for currentRow := 0; currentRow < left.Rows; currentRow++ {
		rowColumnValues := make([]float64, right.Columns)
		resultMatrix[currentRow] = rowColumnValues
		for currentColumn := 0; currentColumn < right.Columns; currentColumn++ {
			rowColumnValues[currentColumn] =
				left.matrixRowValues[currentRow][currentColumn] + right.matrixRowValues[currentRow][currentColumn]
		}
	}

	return Matrix{left.Rows, right.Columns, resultMatrix}
}

func Multiply(left Matrix, others ...Matrix) Matrix {
	for _, right := range others {
		left = multiply(left, right)
	}

	return left
}

func multiply(left, right Matrix) Matrix {
	// Just your standard matrix multiplication.
	// See http://en.wikipedia.org/wiki/Matrix_multiplication for details

	if left.Columns != right.Rows {
		panic("The width of the left matrix must match the height of the right matrix")
	}

	resultRows := left.Rows
	resultColumns := right.Columns

	resultMatrix := make([][]float64, resultRows)

	for currentRow := 0; currentRow < resultRows; currentRow++ {
		resultMatrix[currentRow] = make([]float64, resultColumns)

		var productValue, leftValue, rightValue, vectorIndexProduct float64

		for currentColumn := 0; currentColumn < resultColumns; currentColumn++ {
			productValue = 0

			for vectorIndex := 0; vectorIndex < left.Columns; vectorIndex++ {
				leftValue = left.matrixRowValues[currentRow][vectorIndex]
				rightValue = right.matrixRowValues[vectorIndex][currentColumn]
				vectorIndexProduct = leftValue * rightValue
				productValue += vectorIndexProduct
			}

			resultMatrix[currentRow][currentColumn] = productValue
		}
	}

	return Matrix{resultRows, resultColumns, resultMatrix}
}

func (m Matrix) GetMinorMatrix(rowToRemove, columnToRemove int) Matrix {

	// See http://en.wikipedia.org/wiki/Minor_(linear_algebra)

	// I'm going to use a horribly naïve algorithm... because I can :)
	result := make([][]float64, (m.Rows - 1))
	resultRow := 0

	for currentRow := 0; currentRow < m.Rows; currentRow++ {
		if currentRow == rowToRemove {
			continue
		}

		result[resultRow] = make([]float64, (m.Columns - 1))

		resultColumn := 0

		for currentColumn := 0; currentColumn < m.Columns; currentColumn++ {
			if currentColumn == columnToRemove {
				continue
			}

			result[resultRow][resultColumn] = m.matrixRowValues[currentRow][currentColumn]
			resultColumn++
		}

		resultRow++
	}

	return Matrix{m.Rows - 1, m.Columns - 1, result}
}

func (m Matrix) GetCofactor(rowToRemove, columnToRemove int) float64 {
	// See http://en.wikipedia.org/wiki/Cofactor_(linear_algebra) for details
	// REVIEW: should things be reversed since I'm 0 indexed?
	sum := rowToRemove + columnToRemove
	isEven := (sum%2 == 0)

	if isEven {
		return m.GetMinorMatrix(rowToRemove, columnToRemove).Determinant()
	} else {
		return -1.0 * m.GetMinorMatrix(rowToRemove, columnToRemove).Determinant()
	}
}

func (a Matrix) Equals(b Matrix) bool {
	return reflect.DeepEqual(a, b)
}

func (a Matrix) NotEquals(b Matrix) bool {
	return !a.Equals(b)
}

func (m Matrix) GetHashCode() int32 {
	result := float64(m.Rows)
	result += 2 * float64(m.Columns)

	for currentRow := 0; currentRow < m.Rows; currentRow++ {
		eventRow := (currentRow % 2) == 0
		var multiplier float64
		if eventRow {
			multiplier = 1.0
		} else {
			multiplier = 2.0
		}

		for currentColumn := 0; currentColumn < m.Columns; currentColumn++ {
			cellValue := m.matrixRowValues[currentRow][currentColumn]
			roundedValue := Round(cellValue, FractionalDigitsToRoundTo)
			result += multiplier * roundedValue
		}
	}

	// Ok, now convert that double to an int
	resultBytes := Float64ToBytes(result)

	finalBytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		finalBytes[i] = byte(resultBytes[i] ^ resultBytes[i+4])
	}

	hashCode := BytesToInt32(finalBytes)

	return hashCode
}

func NewDiagonalMatrix(diagonalValues []float64) Matrix {
	dm := Matrix{
		Rows:            len(diagonalValues),
		Columns:         len(diagonalValues),
		matrixRowValues: make([][]float64, len(diagonalValues)),
	}

	for i := 0; i < len(diagonalValues); i++ {
		dm.matrixRowValues[i] = make([]float64, len(diagonalValues))
		dm.matrixRowValues[i][i] = diagonalValues[i]
	}

	return dm
}

type Vector struct {
	Matrix
}

func NewVector(vectorValues []float64) Vector {
	v := Vector{
		Matrix: NewMatrixFromValues(len(vectorValues), 1, vectorValues...),
	}

	return v
}

func NewSquareMatrix(allValues ...float64) Matrix {
	n := int(math.Sqrt(float64(len(allValues))))
	m := Matrix{
		Rows:    n,
		Columns: n,
	}
	m.matrixRowValues = make([][]float64, m.Rows)
	allValuesIndex := 0

	for currentRow := 0; currentRow < m.Rows; currentRow++ {
		currentRowValues := make([]float64, m.Columns)
		m.matrixRowValues[currentRow] = currentRowValues

		for currentColumn := 0; currentColumn < m.Columns; currentColumn++ {
			currentRowValues[currentColumn] = allValues[allValuesIndex]
			allValuesIndex++
		}
	}
	return m
}

func NewIdentityMatrix(rows int) Matrix {
	return NewDiagonalMatrix(createDiagonal(rows))
}

func createDiagonal(rows int) []float64 {
	result := make([]float64, rows)
	for i := 0; i < rows; i++ {
		result[i] = 1.0
	}

	return result
}

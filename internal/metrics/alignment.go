package metrics

type OperationType string

const (
	OperationEqual      OperationType = "equal"
	OperationSubstitute OperationType = "substitute"
	OperationDelete     OperationType = "delete"
	OperationInsert     OperationType = "insert"
)

// Operation описывает одно действие при преобразовании
// эталона в гипотезу.
type Operation struct {
	Type       OperationType `json:"type"`
	Reference  string        `json:"reference,omitempty"`
	Hypothesis string        `json:"hypothesis,omitempty"`
}

// Alignment содержит расстояние, операции и счётчики ошибок.
type Alignment struct {
	Distance      int
	Hits          int
	Substitutions int
	Deletions     int
	Insertions    int
	Operations    []Operation
}

// Align выравнивает две последовательности строк.
//
// Для CER элементами будут Unicode-символы,
// для WER — слова.
func Align(reference, hypothesis []string) Alignment {
	rows := len(reference) + 1
	columns := len(hypothesis) + 1

	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, columns)
	}

	for i := 0; i < rows; i++ {
		matrix[i][0] = i
	}

	for j := 0; j < columns; j++ {
		matrix[0][j] = j
	}

	for i := 1; i < rows; i++ {
		for j := 1; j < columns; j++ {
			substitutionCost := 0
			if reference[i-1] != hypothesis[j-1] {
				substitutionCost = 1
			}

			deletion := matrix[i-1][j] + 1
			insertion := matrix[i][j-1] + 1
			substitution := matrix[i-1][j-1] + substitutionCost

			matrix[i][j] = min3(
				deletion,
				insertion,
				substitution,
			)
		}
	}

	result := Alignment{
		Distance: matrix[len(reference)][len(hypothesis)],
	}

	i := len(reference)
	j := len(hypothesis)

	for i > 0 || j > 0 {
		//совпадение
		if i > 0 &&
			j > 0 &&
			reference[i-1] == hypothesis[j-1] &&
			matrix[i][j] == matrix[i-1][j-1] {

			result.Operations = append(result.Operations, Operation{
				Type:       OperationEqual,
				Reference:  reference[i-1],
				Hypothesis: hypothesis[j-1],
			})
			result.Hits++

			i--
			j--
			continue
		}

		//  замена
		if i > 0 &&
			j > 0 &&
			matrix[i][j] == matrix[i-1][j-1]+1 {

			result.Operations = append(result.Operations, Operation{
				Type:       OperationSubstitute,
				Reference:  reference[i-1],
				Hypothesis: hypothesis[j-1],
			})
			result.Substitutions++

			i--
			j--
			continue
		}

		// удаление
		if i > 0 &&
			matrix[i][j] == matrix[i-1][j]+1 {

			result.Operations = append(result.Operations, Operation{
				Type:      OperationDelete,
				Reference: reference[i-1],
			})
			result.Deletions++

			i--
			continue
		}

		// вставка
		result.Operations = append(result.Operations, Operation{
			Type:       OperationInsert,
			Hypothesis: hypothesis[j-1],
		})
		result.Insertions++
		j--
	}

	reverseOperations(result.Operations)

	return result
}

func reverseOperations(operations []Operation) {
	for left, right := 0, len(operations)-1; left < right; left, right = left+1, right-1 {
		operations[left], operations[right] = operations[right], operations[left]
	}
}

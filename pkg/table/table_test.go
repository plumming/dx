package table

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	simpleTable = `A B C D E
1 2 3 4 5
6 7 8 9 10
`

	variableTable = `AAA   B     CC  DDDD    EEEEE
1     22222 333 4444444 555555555
66666 7777  8   99      10000000000000
`

	variableTableWithHeader = `AAA   B     CC  DDDD    EEEEE
1     22222 333 4444444 555555555
This is a really long header row
66666 7777  8   99      10000000000000
`

	variableTableWithMultipleHeaders = `This is a really long header row
AAA   B     CC  DDDD    EEEEE
1     22222 333 4444444 555555555
This is a really long header row
66666 7777  8   99      10000000000000
`
)

func TestNewTable_Simple(t *testing.T) {
	var b bytes.Buffer
	table := NewTable(&b)
	table.AddRow("A", "B", "C", "D", "E")
	table.AddRow("1", "2", "3", "4", "5")
	table.AddRow("6", "7", "8", "9", "10")
	table.Render()

	assert.Equal(t, b.String(), simpleTable)
}

func TestNewTable_VariableLength(t *testing.T) {
	var b bytes.Buffer
	table := NewTable(&b)
	table.AddRow("AAA", "B", "CC", "DDDD", "EEEEE")
	table.AddRow("1", "22222", "333", "4444444", "555555555")
	table.AddRow("66666", "7777", "8", "99", "10000000000000")
	table.Render()

	assert.Equal(t, b.String(), variableTable)
}

func TestNewTable_VariableLength_WithHeader(t *testing.T) {
	var b bytes.Buffer
	table := NewTable(&b)
	table.AddRow("AAA", "B", "CC", "DDDD", "EEEEE")
	table.AddRow("1", "22222", "333", "4444444", "555555555")
	table.AddRow("# This is a really long header row")
	table.AddRow("66666", "7777", "8", "99", "10000000000000")
	table.Render()

	assert.Equal(t, b.String(), variableTableWithHeader)
}

func TestNewTable_VariableLength_WithMultipleHeader(t *testing.T) {
	var b bytes.Buffer
	table := NewTable(&b)
	table.AddRow("# This is a really long header row")
	table.AddRow("AAA", "B", "CC", "DDDD", "EEEEE")
	table.AddRow("1", "22222", "333", "4444444", "555555555")
	table.AddRow("# This is a really long header row")
	table.AddRow("66666", "7777", "8", "99", "10000000000000")
	table.Render()

	assert.Equal(t, b.String(), variableTableWithMultipleHeaders)
}

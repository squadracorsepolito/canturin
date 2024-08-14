package main

// type dummyOperation struct{}

// func (d *dummyOperation) before() any { return nil }
// func (d *dummyOperation) after() any  { return nil }

// func Test_history(t *testing.T) {
// 	assert := assert.New(t)

// 	history := newHistory()

// 	dummyOp := &dummyOperation{}

// 	history.push(dummyOp)
// 	history.push(dummyOp)
// 	assert.Equal(1, history.currIdx)
// 	assert.Len(history.operations, 2)

// 	history.undo()
// 	history.undo()
// 	assert.Equal(0, history.currIdx)

// 	history.redo()
// 	history.redo()
// 	assert.Equal(1, history.currIdx)

// 	history.push(dummyOp)
// 	assert.Equal(2, history.currIdx)
// 	assert.Len(history.operations, 3)

// 	history.undo()
// 	history.undo()
// 	assert.Equal(0, history.currIdx)

// 	history.push(dummyOp)
// 	assert.Equal(1, history.currIdx)
// 	assert.Len(history.operations, 2)

// 	history.redo()
// 	assert.Equal(1, history.currIdx)
// }

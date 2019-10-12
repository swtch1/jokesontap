package jokesontap

//func TestPutAndGetFromQueue(t *testing.T) {
//	t.Parallel()
//	assert := assert.New(t)
//
//	tests := []struct {
//		name string
//		items []int
//	}{
//		{"basic", []int{1, 5, 7, 9}},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			nq := NewNameQueue(int64(len(tt.items)))
//			for _, item := range tt.items {
//				// we do not expect any errors adding to the Queue
//				assert.Nil(nq.Queue.Put(item))
//			}
//			for _, item := range tt.items {
//				v, err := nq.Queue.Get(1)
//				// we do not expect to get any errors dequeueing
//				assert.Nil(err)
//				// ensure items come out the way they went in
//				assert.Equal(item, v[0])
//			}
//		})
//	}
//
//}

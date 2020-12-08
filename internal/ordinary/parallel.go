package ordinary

import (
	"context"
	"sync"
)

// GroupFunc 调用函数
type GroupFunc func() error

// NewParallel 创建
func NewParallel(ctx context.Context, handles ...GroupFunc) error {
	errChan := make(chan error)
	doneChan := make(chan *struct{})

	for _, handle := range handles {
		currentHandle := handle

		go func() {
			if err := currentHandle(); err != nil {
				errChan <- err
			}

			doneChan <- nil
		}()
	}

	count := len(handles)

	for {
		select {
		case <-ctx.Done():
			{
				return ctx.Err()
			}

		case err := <-errChan:
			{
				return err
			}

		case <-doneChan:
			{
				count--
				if count <= 0 {
					return nil
				}
			}
		}
	}
}

// Merge 扇入函数（组件），把多个 channel 中的数据发送到一个 channel 中
func Merge(ins ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	p := func(in <-chan string) {
		defer wg.Done()
		for c := range in {
			out <- c
		}

	}

	wg.Add(len(ins))

	for _, cs := range ins {
		go p(cs)
	}

	// 等待所有输入的数据ins处理完，再关闭输出out
	go func() {
		wg.Wait()
		close(out)

	}()

	return out

}

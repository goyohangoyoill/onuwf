package util

import "context"

// ChkWithCtx 함수는 각 함수가 종료된 상태인지를 확인하는 함수입니다.
func ChkWithCtx(ctx context.Context, isDone chan bool) error {
	select {
	//	case result := <-isDone:
	case <-isDone:
		// 해당 함수 정상종료시.
		return nil
	case <-ctx.Done():
		// 강제종료 입력시.
		return ctx.Err()
	}
}

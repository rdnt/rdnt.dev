package router

import (
	"fmt"
	"testing"
)

func TestPagination(t *testing.T) {
	t.Skip()

	total := 9
	for i := 1; i <= total; i++ {
		curr := i
		p := Pagination{
			Page:     curr,
			PrevPage: curr - 1,
			NextPage: curr + 1,
			Total:    total,
		}

		p.UpdatePagination()

		fmt.Print("=== ", curr, " -> ")
		fmt.Println(p.Pages)
	}

}

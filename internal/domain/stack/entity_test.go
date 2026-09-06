package stack_test

import (
	"testing"

	"dockermanager/internal/domain/container"
	"dockermanager/internal/domain/stack"
)

func TestComposeStackAggregate(t *testing.T) {
	st := stack.ComposeStack{
		Name: "web_stack",
		Containers: []container.Container{
			{ID: "c1", Name: "web", State: container.StateRunning},
			{ID: "c2", Name: "db", State: container.StateRunning},
			{ID: "c3", Name: "redis", State: container.StateExited},
		},
	}

	st.RecalculateCounts()

	if st.TotalCount != 3 {
		t.Errorf("expected TotalCount = 3, got %d", st.TotalCount)
	}
	if st.RunningCount != 2 {
		t.Errorf("expected RunningCount = 2, got %d", st.RunningCount)
	}

	// Stop c1
	st.Containers[0].State = container.StateExited
	st.RecalculateCounts()

	if st.RunningCount != 1 {
		t.Errorf("expected RunningCount = 1 after stopping container, got %d", st.RunningCount)
	}
}

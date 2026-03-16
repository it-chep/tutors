package permissions

type Request struct {
	CanViewContracts      bool    `json:"can_view_contracts"`
	CanPenalizeAssistants []int64 `json:"can_penalize_assistants"`
}

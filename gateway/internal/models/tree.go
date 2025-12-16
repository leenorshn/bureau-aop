package models

// TreeNode représente un nœud dans l'arbre client
type TreeNode struct {
	ID                string  `json:"id"`
	ClientID          string  `json:"clientId"`
	Name              string  `json:"name"`
	Phone             *string `json:"phone,omitempty"`
	ParentID          *string `json:"parentId,omitempty"`
	Level             int     `json:"level"`
	Position          *string `json:"position,omitempty"`
	NetworkVolumeLeft float64 `json:"networkVolumeLeft"`
	NetworkVolumeRight float64 `json:"networkVolumeRight"`
	BinaryPairs       int     `json:"binaryPairs"`
	TotalEarnings     float64 `json:"totalEarnings"`
	WalletBalance     float64 `json:"walletBalance"`
	IsActive          bool    `json:"isActive"`
	LeftActives       int     `json:"leftActives"`
	RightActives      int     `json:"rightActives"`
	IsQualified       bool    `json:"isQualified"`
	CyclesAvailable   *int    `json:"cyclesAvailable,omitempty"`
	CyclesPaidToday   *int    `json:"cyclesPaidToday,omitempty"`
}

// ClientTreeResponse représente la réponse complète de l'arbre client
type ClientTreeResponse struct {
	Root       *TreeNode  `json:"root"`
	Nodes      []*TreeNode `json:"nodes"`
	TotalNodes int        `json:"totalNodes"`
	MaxLevel   int        `json:"maxLevel"`
}




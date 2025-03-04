package cmd

import (
	"artifacts/internal"
	"os"

	"github.com/0xN0x/go-artifactsmmo"
	"github.com/0xN0x/go-artifactsmmo/models"
)

func GetBankItemsAction(client *artifactsmmo.ArtifactsMMO, code string, page int, size int) *[]models.SimpleItem {
	items, err := client.GetBankItems(code, page, size)
	if err != nil {
		internal.Logger.Error("Error getting bank items", "err", err)
		os.Exit(1)
	}

	return items
}

func DepositItemAction(client *artifactsmmo.ArtifactsMMO, code string, quantity int) *models.BankItemTransaction {
	transaction, err := client.DepositBank(code, quantity)
	if err != nil {
		internal.Logger.Error("Error depositing items", "err", err)
		os.Exit(1)
	}

	return transaction
}

func WithdrawItemAction(client *artifactsmmo.ArtifactsMMO, code string, quantity int) *models.BankItemTransaction {
	transaction, err := client.WithdrawBank(code, quantity)
	if err != nil {
		internal.Logger.Error("Error withdrawing items", "err", err)
		os.Exit(1)
	}

	return transaction
}

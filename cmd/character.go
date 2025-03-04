package cmd

import (
	"artifacts/internal"
	"os"

	"github.com/0xN0x/go-artifactsmmo"
	"github.com/0xN0x/go-artifactsmmo/models"
)

func CharacterInfoAction(client *artifactsmmo.ArtifactsMMO, charName string) *models.Character {
	info, err := client.GetCharacterInfo(charName)
	if err != nil {
		internal.Logger.Error("Error getting character info", "err", err)
		os.Exit(1)
	}

	return info
}

func CharacterInventoryAction(client *artifactsmmo.ArtifactsMMO, charName string) (*[]models.InventorySlots, int) {
	info := CharacterInfoAction(client, charName)
	return &info.Inventory, info.InventoryMaxItems
}

func MoveAction(client *artifactsmmo.ArtifactsMMO, charName string, x, y int) *models.CharacterMovementData {
	move, err := client.Move(x, y)
	if err != nil {
		internal.Logger.Error("Error moving character", "err", err)
		os.Exit(1)
	}

	internal.Cooldown(move.Cooldown.RemainingSeconds)
	return move
}

func RestAction(client *artifactsmmo.ArtifactsMMO) *models.Rest {
	rest, err := client.Rest()
	if err != nil {
		internal.Logger.Error("Error resting", "err", err)
		os.Exit(1)
	}

	internal.Cooldown(rest.Cooldown.RemainingSeconds)
	return rest
}

func GatherAction(client *artifactsmmo.ArtifactsMMO) *models.SkillData {
	gather, err := client.Gather()
	if err != nil {
		internal.Logger.Error("Error gathering", "err", err)
		os.Exit(1)
	}

	internal.Cooldown(gather.Cooldown.RemainingSeconds)
	return gather
}

func FightAction(client *artifactsmmo.ArtifactsMMO) *models.CharacterFight {
	fight, err := client.Fight()
	if err != nil {
		internal.Logger.Error("Error fighting", "err", err)
		os.Exit(1)
	}

	internal.Cooldown(fight.Cooldown.RemainingSeconds)
	return fight
}

func CraftAction(client *artifactsmmo.ArtifactsMMO, item string, quantity int) *models.SkillData {
	craft, err := client.Craft(item, quantity)
	if err != nil {
		internal.Logger.Error("Error crafting", "err", err)
	}

	internal.Cooldown(craft.Cooldown.RemainingSeconds)
	return craft
}

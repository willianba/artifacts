package cmd

import (
	"artifacts/internal"

	"github.com/0xN0x/go-artifactsmmo"
)

type FlowName string

const (
	CopperFlow FlowName = "copper"
)

type FlowAction func(client *artifactsmmo.ArtifactsMMO, charName string, goal int)

type Flow struct {
	Name   FlowName
	Action FlowAction
}

var Flows = []Flow{
	{
		Name:   CopperFlow,
		Action: copperFlowAction,
	},
}

func GetFlow(name FlowName) FlowAction {
	for _, flow := range Flows {
		if flow.Name == name {
			return flow.Action
		}
	}

	return nil
}

func copperFlowAction(client *artifactsmmo.ArtifactsMMO, charName string, goal int) {
	var goalFulfilled bool
	for {
		if goalFulfilled {
			break
		}
		info := CharacterInfoAction(client, charName)

		// move to copper rocks if not there yet
		if info.X != 2 && info.Y != 0 {
			internal.Logger.Info("Moving to copper rocks")
			MoveAction(client, charName, 2, 0)
		}

		var isInventoryFull bool
		for {
			if isInventoryFull {
				break
			}

			internal.Logger.Info("Gathering copper ores")
			GatherAction(client)
			inventory, maxItems := CharacterInventoryAction(client, charName)

			itemsQty := 0
			for _, item := range *inventory {
				itemsQty += item.Quantity
			}

			isInventoryFull = itemsQty == maxItems
		}

		internal.Logger.Info("Moving to the furnace")
		MoveAction(client, charName, 1, 5)

		internal.Logger.Info("Crafting copper bars")
		copperOres := countCodeInInventory(client, charName, "copper_ore")
		CraftAction(client, "copper", copperOres/10)

		goalFulfilled = canGoalBeFullfiled(client, charName, "copper_ore", "copper", goal)
	}
}

func canGoalBeFullfiled(client *artifactsmmo.ArtifactsMMO, charName string, rawCode string, goalCode string, goal int) bool {
	inventory, _ := CharacterInventoryAction(client, charName)

	totalGoalItem := 0
	totalRawItem := 0

	for _, item := range *inventory {
		if item.Code == goalCode {
			totalGoalItem = item.Quantity
		}

		if item.Code == rawCode {
			totalRawItem = item.Quantity
		}
	}

	return totalGoalItem+totalRawItem/10 >= goal
}

func countCodeInInventory(client *artifactsmmo.ArtifactsMMO, charName string, code string) int {
	inventory, _ := CharacterInventoryAction(client, charName)

	var count int
	for _, item := range *inventory {
		if item.Code == code {
			count = item.Quantity
			break
		}
	}

	return count
}

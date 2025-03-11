package cmd

import (
	"artifacts/internal"
	"os"

	"github.com/0xN0x/go-artifactsmmo"
	"github.com/0xN0x/go-artifactsmmo/models"
)

type FlowName string

const (
	CopperFlow FlowName = "copper"
	IronFlow   FlowName = "iron"
)

var (
	flowMap = map[string]FlowName{
		"copper": CopperFlow,
		"iron":   IronFlow,
	}
)

func ParseToFlowName(str string) FlowName {
	c := flowMap[str]
	return c
}

type FlowAction func(client *artifactsmmo.ArtifactsMMO, charName string, goal int)

type Flow struct {
	Name   FlowName
	Action FlowAction
}

var Flows = []Flow{
	{
		Name: CopperFlow,
		Action: miningFlowAction(MiningConfig{
			OreCode:     "copper_ore",
			BarCode:     "copper",
			OreLocation: models.Movement{X: 2, Y: 0},
			OresPerBar:  10,
		}),
	},
	{
		Name: IronFlow,
		Action: miningFlowAction(MiningConfig{
			OreCode:     "iron_ore",
			BarCode:     "iron",
			OreLocation: models.Movement{X: 1, Y: 7},
			OresPerBar:  10,
		}),
	},
}

func GetFlow(name FlowName) FlowAction {
	for _, flow := range Flows {
		if flow.Name == name {
			return flow.Action
		}
	}

	internal.Logger.Error("Flow not found", "flow", name)
	os.Exit(1)
	return nil
}

type MiningConfig struct {
	OreCode     string
	BarCode     string
	OreLocation models.Movement
	OresPerBar  int
}

func miningFlowAction(config MiningConfig) FlowAction {
	return func(client *artifactsmmo.ArtifactsMMO, charName string, goal int) {
		var goalFulfilled bool

		for {
			if goalFulfilled {
				break
			}
			info := CharacterInfoAction(client, charName)

			// Move to ore location if not there yet
			if info.X != config.OreLocation.X || info.Y != config.OreLocation.Y {
				internal.Logger.Infof("Moving to %s rocks", config.OreCode)
				MoveAction(client, charName, config.OreLocation.X, config.OreLocation.Y)
			}

			var isInventoryFull bool
			for {
				if isInventoryFull {
					break
				}

				internal.Logger.Infof("Gathering %s ores", config.OreCode)
				GatherAction(client)
				inventory, maxItems := CharacterInventoryAction(client, charName)

				itemsQty := 0
				for _, item := range *inventory {
					itemsQty += item.Quantity
				}

				isInventoryFull = itemsQty == maxItems
			}

			internal.Logger.Infof("Moving to the furnace")
			MoveAction(client, charName, 1, 5)

			internal.Logger.Infof("Crafting %s bars", config.BarCode)
			oreCount := countCodeInInventory(client, charName, config.OreCode)
			CraftAction(client, config.BarCode, oreCount/config.OresPerBar)

			goalFulfilled = canGoalBeFullfiled(client, charName, config.OreCode, config.BarCode, goal)
		}
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

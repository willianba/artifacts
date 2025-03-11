package main

import (
	"artifacts/cmd"
	"artifacts/internal"
	"os"
	"strconv"

	"github.com/0xN0x/go-artifactsmmo"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		internal.Logger.Error("Error loading .env file", "err", err)
		os.Exit(1)
	}

	app := createApp()
	err = app.Run(os.Args)
	if err != nil {
		internal.Logger.Error("Error starting app", "err", err)
		os.Exit(1)
	}
}

func createApp() *cli.App {
	var charName string

	apiToken := os.Getenv("TOKEN")
	if apiToken == "" {
		internal.Logger.Error("TOKEN environment variable is not set")
		os.Exit(1)
	}

	app := &cli.App{
		Name:  "Artifacts CLI",
		Usage: "Interact with the ArtifactsMMO API",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Aliases:     []string{"n"},
				Usage:       "Name of the character",
				Required:    true,
				Destination: &charName,
			},
			&cli.IntFlag{
				Name:    "loop",
				Aliases: []string{"l"},
				Usage:   "The amount of times to loop the action",
				Value:   1,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "character",
				Aliases: []string{"ch"},
				Usage:   "Get information about your character",
				Action: func(c *cli.Context) error {
					client := artifactsmmo.NewClient(apiToken, charName)
					info := cmd.CharacterInfoAction(client, charName)
					internal.Logger.Infof("Name: %s\nXP: %d/%d\n", info.Name, info.Xp, info.MaxXp)
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:    "inventory",
						Aliases: []string{"i"},
						Usage:   "Get information about your character's inventory",
						Action: func(c *cli.Context) error {
							client := artifactsmmo.NewClient(apiToken, charName)
							inventory, maxItems := cmd.CharacterInventoryAction(client, charName)

							itemsQty := 0
							for _, item := range *inventory {
								itemsQty += item.Quantity
							}

							internal.Logger.Infof("Items in inventory %v\n", inventory)
							internal.Logger.Infof("Slots used: %d/%d\n", itemsQty, maxItems)
							return nil
						},
					},
					{
						Name:            "move",
						Aliases:         []string{"m"},
						Usage:           "Move your character",
						SkipFlagParsing: true,
						Action: func(c *cli.Context) error {
							client := artifactsmmo.NewClient(apiToken, charName)
							info := cmd.CharacterInfoAction(client, charName)
							x, _ := strconv.Atoi(c.Args().Get(0))
							y, _ := strconv.Atoi(c.Args().Get(1))
							move := cmd.MoveAction(client, charName, x, y)

							internal.Logger.Infof("Moved %s from (%d, %d) to (%d, %d)\n", info.Name, info.X, info.Y, move.Destination.X, move.Destination.Y)
							return nil
						},
					},
					{
						Name:    "gather",
						Aliases: []string{"g"},
						Usage:   "Gather resources in the location",
						Action: func(c *cli.Context) error {
							client := artifactsmmo.NewClient(apiToken, charName)
							loops := c.Int("l")

							internal.LoopAction(func() {
								gather := cmd.GatherAction(client)
								internal.Logger.Info("Done gathering", "items", gather.Details.Items, "XP", gather.Details.Xp)
							}, loops)

							return nil
						},
					},
					{
						Name:    "craft",
						Aliases: []string{"c"},
						Usage:   "Craft an item",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "item",
								Aliases:  []string{"i"},
								Usage:    "Code of the item to craft",
								Required: true,
							},
							&cli.IntFlag{
								Name:    "quantity",
								Aliases: []string{"q"},
								Usage:   "Quantity of the item to craft",
								Value:   1,
							},
						},
						Action: func(c *cli.Context) error {
							client := artifactsmmo.NewClient(apiToken, charName)
							code := c.String("item")
							quantity := c.Int("quantity")
							craft := cmd.CraftAction(client, code, quantity)

							internal.Logger.Infof("Crafted %v. Got %d XP\n", craft.Details.Items, craft.Details.Xp)
							return nil
						},
					},
					{
						Name:    "fight",
						Aliases: []string{"f"},
						Usage:   "Fight a monster",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "rest",
								Aliases: []string{"r"},
								Usage:   "Rest after the fight if HP is below 50%",
								Value:   true,
							},
						},
						Action: func(c *cli.Context) error {
							client := artifactsmmo.NewClient(apiToken, charName)
							loops := c.Int("l")
							shouldRest := c.Bool("r")

							internal.LoopAction(func() {
								fight := cmd.FightAction(client)
								internal.Logger.Info("Fight done", "result", fight.Fight.Result, "XP", fight.Fight.Xp, "gold", fight.Fight.Gold)

								if shouldRest {
									if fight.Character.Hp <= fight.Character.MaxHp/2 {
										internal.Logger.Info("Character HP is below 50%, resting...")
										rest := cmd.RestAction(client)
										internal.Logger.Info("HP restored", rest.HpRestored)
									}
								}

							}, loops)

							return nil
						},
					},
					{
						Name:    "rest",
						Aliases: []string{"r"},
						Usage:   "Rest to recover HP",
						Action: func(c *cli.Context) error {
							client := artifactsmmo.NewClient(apiToken, charName)
							rest := cmd.RestAction(client)
							internal.Logger.Infof("Restored %d HP\n", rest.HpRestored)
							return nil
						},
					},
					{
						Name:    "use",
						Aliases: []string{"u"},
						Usage:   "Use a consumable",
						Action: func(c *cli.Context) error {
							client := artifactsmmo.NewClient(apiToken, charName)
							code := c.Args().Get(0)
							quantity, _ := strconv.Atoi(c.Args().Get(1))
							use := cmd.UseConsumableAction(client, code, quantity)
							internal.Logger.Infof("Used %d %s\n", quantity, use.Item.Name)
							return nil
						},
					},
					{
						Name:    "recycle",
						Aliases: []string{"rc"},
						Usage:   "Recycle an item",
						Action: func(c *cli.Context) error {
							client := artifactsmmo.NewClient(apiToken, charName)
							code := c.Args().Get(0)
							quantity, _ := strconv.Atoi(c.Args().Get(1))
							recycle := cmd.RecycleAction(client, code, quantity)
							internal.Logger.Infof("Recycled %d %s. Got %d %s", quantity, code, recycle.Details.Items[0].Quantity, recycle.Details.Items[0].Code)
							return nil
						},
					},
				},
			},
			{
				Name:    "bank",
				Aliases: []string{"b"},
				Usage:   "Interact with the bank",
				Subcommands: []*cli.Command{
					{
						Name:    "deposit",
						Aliases: []string{"d"},
						Usage:   "Deposit an item in the bank",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "item",
								Aliases:  []string{"i"},
								Usage:    "Code of the item to deposit",
								Required: true,
							},
							&cli.IntFlag{
								Name:    "quantity",
								Aliases: []string{"q"},
								Usage:   "Quantity of the item to deposit",
								Value:   1,
							},
						},
						Action: func(c *cli.Context) error {
							client := artifactsmmo.NewClient(apiToken, charName)
							code := c.String("item")
							quantity := c.Int("quantity")
							transaction := cmd.DepositItemAction(client, code, quantity)
							internal.Logger.Info("Item deposited successfully", "transaction", transaction.Item)
							return nil
						},
					},
					{
						Name:    "withdraw",
						Aliases: []string{"w"},
						Usage:   "Withdraw an item from the bank",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "item",
								Aliases:  []string{"i"},
								Usage:    "Code of the item to deposit",
								Required: true,
							},
							&cli.IntFlag{
								Name:    "quantity",
								Aliases: []string{"q"},
								Usage:   "Quantity of the item to deposit",
								Value:   1,
							},
						},
						Action: func(c *cli.Context) error {
							client := artifactsmmo.NewClient(apiToken, charName)
							code := c.String("item")
							quantity := c.Int("quantity")
							transaction := cmd.WithdrawItemAction(client, code, quantity)
							internal.Logger.Info("Item withdrawn successfully", "transaction", transaction.Item)
							return nil
						},
					},
					{
						Name:    "items",
						Aliases: []string{"i"},
						Usage:   "Get items in the bank",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "code",
								Aliases: []string{"c"},
								Usage:   "Code of the item to filter by",
							},
							&cli.IntFlag{
								Name:    "page",
								Aliases: []string{"p"},
								Usage:   "Page number",
								Value:   1,
							},
							&cli.IntFlag{
								Name:    "size",
								Aliases: []string{"s"},
								Usage:   "Number of items per page",
								Value:   50,
								Action: func(c *cli.Context, flag int) error {
									if flag < 1 {
										internal.Logger.Error("Size can't be less than 1")
										os.Exit(1)
									}

									if flag > 100 {
										internal.Logger.Error("Size can't be greater than 100")
										os.Exit(1)
									}

									return nil
								},
							},
						},
						Action: func(c *cli.Context) error {
							client := artifactsmmo.NewClient(apiToken, charName)
							items := cmd.GetBankItemsAction(client, c.String("code"), c.Int("page"), c.Int("size"))
							internal.Logger.Info("Items in bank", "items", items)
							return nil
						},
					},
				},
			},
			{
				Name:    "tasks",
				Aliases: []string{"t"},
				Usage:   "Interact with the task master",
				Subcommands: []*cli.Command{
					{
						Name:    "accept",
						Aliases: []string{"a"},
						Usage:   "Accept a new task",
						Action: func(c *cli.Context) error {
							client := artifactsmmo.NewClient(apiToken, charName)
							task := cmd.AcceptNewTaskAction(client)
							internal.Logger.Info("New task accepted", "task", task.Task)
							return nil
						},
					},
				},
			},
			{
				Name:    "flow",
				Aliases: []string{"f"},
				Usage:   "Perform a series of actions",
				Action: func(c *cli.Context) error {
					client := artifactsmmo.NewClient(apiToken, charName)
					flow := cmd.ParseToFlowName(c.Args().Get(0))
					goal, err := strconv.Atoi(c.Args().Get(1))
					if err != nil {
						internal.Logger.Error("Error parsing goal", "err", err)
						os.Exit(1)
					}

					cmd.GetFlow(flow)(client, charName, goal)

					return nil
				},
			},
		},
	}

	return app
}

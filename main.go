package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Gopher struct {
	name              string
	hitpoints         int
	weapon            Weapon
	inventory         map[string]int
	activeConsumables []Consumable
	strength          int
	agility           int
	intellect         int
	coins             int
}

type Weapon struct {
	name            string
	damage          [2]int
	strengthReq     int
	agilityReq      int
	intelligenceReq int
	cost            int
}

type Consumable struct {
	name            string
	duration        int
	hitpointsEffect int
	strengthEffect  int
	agilityEffect   int
	intellectEffect int
	cost            int
}

type action interface {
	attack()
	buy()
	use()
	work()
	train()
	exit()
}

func main() {
	rand.NewSource(time.Now().UnixNano())

	g1 := Gopher{"Gopher 1", 30, weapons[barehand], make(map[string]int), []Consumable{}, 0, 0, 0, 20}
	g2 := Gopher{"Gopher 2", 30, weapons[barehand], make(map[string]int), []Consumable{}, 0, 0, 0, 20}

	for true {
		fmt.Printf("\n\n")
		printGopherStatus(&g1)
		printGopherStatus(&g2)
		play(&g1, &g2)
		fmt.Printf("\n\n")
		printGopherStatus(&g1)
		printGopherStatus(&g2)
		play(&g2, &g1)
	}
}

func play(player *Gopher, opponent *Gopher) {
	removeInactiveConsumables(player)
	played := false
	for !played {
		fmt.Printf("\n%s", player.name)
		showActions()
		var actionChoice int
		fmt.Scan(&actionChoice)
		played = performAction(actionChoice, player, opponent)
	}
}

func printGopherStatus(g *Gopher) {
	fmt.Printf("%s:\n", g.name)
	fmt.Printf("Hitpoints: %d\n", g.hitpoints)
	fmt.Printf("Strength: %d\n", g.strength)
	fmt.Printf("Agility: %d\n", g.agility)
	fmt.Printf("Intellect: %d\n", g.intellect)
	fmt.Printf("Coins: %d\n", g.coins)
	fmt.Printf("Equipped Weapon: %+v\n", g.weapon.name)
	fmt.Printf("Inventory:\n")
	for consumable, quantity := range g.inventory {
		fmt.Printf("%s: %d\n", consumable, quantity)
	}
	fmt.Printf("\n")
}

func showActions() {
	fmt.Printf("\nActions:\n")
	fmt.Printf("1. Attack\n")
	fmt.Printf("2. Buy\n")
	fmt.Printf("3. Use\n")
	fmt.Printf("4. Work\n")
	fmt.Printf("5. Train\n")
	fmt.Printf("6. Exit\n")

	fmt.Print("Enter your choice (1-6): ")
}

func performAction(actionChoice int, player *Gopher, opponent *Gopher) bool {
	switch actionChoice {
	case 1:
		return attack(player, opponent)
	case 2:
		return buyMenu(player)
	case 3:
		return useMenu(player)
	case 4:
		return work(player)
	case 5:
		return trainMenu(player)
	case 6:
		fmt.Printf("%s forfeits the game. %s wins!\n", player.name, opponent.name)
		os.Exit(0)
	default:
		fmt.Printf("Invalid choice. Try again.\n")
	}
	return false
}

func buyMenu(g *Gopher) bool {
	fmt.Printf("\nBuy Items:\n")
	fmt.Printf("1. Consumables\n")
	fmt.Printf("2. Weapons\n")
	fmt.Print("Enter your choice (1-2): ")
	var buyChoice int
	fmt.Scanln(&buyChoice)

	switch buyChoice {
	case 1:
		return buyConsumablesMenu(g)
	case 2:
		return buyWeaponsMenu(g)
	default:
		fmt.Printf("Invalid choice. Try again.\n")
	}
	return false
}

func buyConsumablesMenu(g *Gopher) bool {
	fmt.Printf("\nBuy Consumables:\n")
	fmt.Printf("1. Health Potion (5 gold)\n")
	fmt.Printf("2. Strength Potion (10 gold)\n")
	fmt.Printf("3. Agility Potion (10 gold)\n")
	fmt.Printf("4. Intellect Potion (10 gold)\n")
	fmt.Print("Enter the number of the consumable you want to buy (1-4): ")
	var consumableChoice int
	fmt.Scanln(&consumableChoice)

	if consumableChoice < 1 && consumableChoice > 4 {
		fmt.Printf("Invalid consumable choice. Try again.\n")
		return false
	}

	return buy(g, consumableOptions[consumableChoice], true)
}

func buyWeaponsMenu(g *Gopher) bool {
	fmt.Printf("\nBuy Weapons:\n")
	fmt.Printf("1. Knife (10 gold)\n")
	fmt.Printf("2. Sword (35 gold)\n")
	fmt.Printf("3. Ninjaku (25 gold)\n")
	fmt.Printf("4. Wand (30 gold)\n")
	fmt.Printf("5. Gophermourne (65 gold)\n")
	fmt.Print("Enter the number of the weapon you want to buy (1-5): ")
	var weaponChoice int
	fmt.Scanln(&weaponChoice)

	if weaponChoice < 1 && weaponChoice > 5 {
		fmt.Printf("Invalid weapon choice. Try again.\n")
		return false
	}

	return buy(g, weaponOptions[weaponChoice], false)
}

func useMenu(g *Gopher) bool {
	inventoryMenu := make([]string, 0)
	for itemName, itemCount := range g.inventory {
		inventoryMenu = append(inventoryMenu, fmt.Sprintf("%s (x%d)", itemName, itemCount))
	}

	fmt.Printf("Select an item to use:\n")
	for i, item := range inventoryMenu {
		fmt.Printf("%d. %s\n", i+1, item)
	}

	fmt.Print("Enter the number of the item you want to use: ")
	var itemChoice int
	fmt.Scanln(&itemChoice)

	if itemChoice < 1 || itemChoice > len(inventoryMenu) {
		fmt.Printf("Invalid choice. Try again.\n")
		return false
	}

	selectedItem := inventoryMenu[itemChoice-1]
	itemName := strings.Split(selectedItem, " ")[0]
	return use(g, itemName)
}

func trainMenu(g *Gopher) bool {
	fmt.Println("\nTrain Skill:")
	fmt.Println("1. Strength")
	fmt.Println("2. Agility")
	fmt.Println("3. Intellect")
	fmt.Print("Enter the number of the skill you want to train (1-3): ")

	var skillChoice int
	fmt.Scanln(&skillChoice)

	if skillChoice < 1 && skillChoice > 3 {
		fmt.Println("Invalid choice. Try again.")
		return false
	}

	return train(g, skillOptions[skillChoice])
}

func removeInactiveConsumables(g *Gopher) {
	newActiveConsumables := []Consumable{}
	for _, item := range g.activeConsumables {
		item.duration--
		if item.duration != 0 {
			newActiveConsumables = append(newActiveConsumables, item)
		} else {
			g.hitpoints -= item.hitpointsEffect
			g.strength -= item.strengthEffect
			g.agility -= item.agilityEffect
			g.intellect -= item.intellectEffect
			fmt.Printf("%s duration expired and %s lost the following effects:\n", item.name, g.name)
			fmt.Printf("Hitpoints: -%d\n", item.hitpointsEffect)
			fmt.Printf("Strength: -%d\n", item.strengthEffect)
			fmt.Printf("Agility: -%d\n", item.agilityEffect)
			fmt.Printf("Intellect: -%d\n", item.intellectEffect)
		}
	}
	g.activeConsumables = newActiveConsumables
}

func attack(attacker *Gopher, target *Gopher) bool {
	damage := rand.Intn(attacker.weapon.damage[1]-attacker.weapon.damage[0]+1) + attacker.weapon.damage[0]
	target.hitpoints -= damage
	fmt.Printf("%s attacks %s for %d damage!\n", attacker.name, target.name, damage)
	if target.hitpoints <= 0 {
		fmt.Printf("%s is dead. %s wins!\n", target.name, attacker.name)
		os.Exit(0)
	}
	return true
}

func buy(g *Gopher, itemToBuy string, isConsumable bool) bool {
	var exists bool

	if isConsumable {
		_, exists = consumables[itemToBuy]
	} else {
		_, exists = weapons[itemToBuy]
	}

	if !exists {
		fmt.Printf("Item not found in the shop. Try again\n")
		return false
	}

	if isConsumable {
		c := consumables[itemToBuy]
		if g.coins < c.cost {
			fmt.Printf("Not enough coins to buy the item. Try again\n")
			return false
		}
		g.coins -= c.cost
		g.inventory[itemToBuy]++
		fmt.Printf("%s bought %s for %d gold\n", g.name, itemToBuy, c.cost)
	} else {
		w := weapons[itemToBuy]
		if g.coins < w.cost {
			fmt.Printf("Not enough coins to buy the weapon. Try again\n")
			return false
		}
		if g.strength < w.strengthReq && g.agility < w.agilityReq && g.intellect < w.intelligenceReq {
			fmt.Printf("Insufficient skills to buy and equip the weapon. Try again\n")
			return false
		}
		g.coins -= w.cost
		g.weapon = w
		fmt.Printf("%s bought %s for %d gold and equipped it\n", g.name, itemToBuy, w.cost)
	}
	return true
}

func use(g *Gopher, itemToUse string) bool {
	if g.inventory[itemToUse] == 0 {
		fmt.Printf("Item not found in inventory. Try again\n")
		return false
	}

	c := consumables[itemToUse]
	g.hitpoints += c.hitpointsEffect
	g.strength += c.strengthEffect
	g.agility += c.agilityEffect
	g.intellect += c.intellectEffect
	fmt.Printf("%s used %s and gained the following effects:\n", g.name, itemToUse)
	fmt.Printf("Hitpoints: +%d\n", c.hitpointsEffect)
	fmt.Printf("Strength: +%d\n", c.strengthEffect)
	fmt.Printf("Agility: +%d\n", c.agilityEffect)
	fmt.Printf("Intellect: +%d\n", c.intellectEffect)
	g.inventory[itemToUse]--
	if g.inventory[itemToUse] == 0 {
		delete(g.inventory, itemToUse)
	}
	g.activeConsumables = append(g.activeConsumables, c)
	return true
}

func work(g *Gopher) bool {
	coinsEarned := rand.Intn(maxCoinsEarnedWorking-minCoinsEarnedWorking+1) + 5
	g.coins += coinsEarned
	fmt.Printf("%s worked and earned %d coins.\n", g.name, coinsEarned)
	return true
}

func train(g *Gopher, skillToTrain string) bool {
	if g.coins < trainingCost {
		fmt.Printf("Not enough coins to train. Try again\n")
		return false
	}

	switch skillToTrain {
	case strength, agility, intellect:
		switch skillToTrain {
		case strength:
			g.strength += trainingSkillIncrement
		case agility:
			g.agility += trainingSkillIncrement
		case intellect:
			g.intellect += trainingSkillIncrement
		}
		g.coins -= trainingCost
		fmt.Printf("%s trained %s and gained %d points.\n", g.name, skillToTrain, trainingSkillIncrement)
	default:
		fmt.Printf("Invalid skill to train. Try again\n")
		return false
	}
	return true
}

// weapons
var weapons = map[string]Weapon{
	barehand:     Weapon{barehand, [2]int{1, 1}, 0, 0, 0, 0},
	knife:        Weapon{knife, [2]int{2, 3}, 0, 0, 0, 10},
	sword:        Weapon{sword, [2]int{3, 5}, 2, 0, 0, 35},
	ninjaku:      Weapon{ninjaku, [2]int{1, 7}, 0, 2, 0, 25},
	wand:         Weapon{wand, [2]int{3, 3}, 0, 0, 2, 30},
	gophermourne: Weapon{gophermourne, [2]int{6, 7}, 3, 0, 2, 65},
}

// consumables
var consumables = map[string]Consumable{
	health_potion:    Consumable{health_potion, -1, 5, 0, 0, 0, 5},
	strength_potion:  Consumable{strength_potion, 3, 0, 3, 0, 0, 10},
	agility_potion:   Consumable{agility_potion, 3, 0, 0, 3, 0, 10},
	intellect_potion: Consumable{intellect_potion, 3, 0, 0, 0, 3, 10},
}

// weapon constants
const (
	barehand     = "barehand"
	knife        = "knife"
	sword        = "sword"
	ninjaku      = "ninjaku"
	wand         = "wand"
	gophermourne = "gophermourne"
)

// consumable constants
const (
	health_potion    = "health_potion"
	strength_potion  = "strength_potion"
	agility_potion   = "agility_potion"
	intellect_potion = "intellect_potion"
)

// skills
const (
	strength  = "strength"
	agility   = "agility"
	intellect = "intellect"
)

// consumable options
var consumableOptions = [5]string{"", health_potion, strength_potion, agility_potion, intellect_potion}

// weapon options
var weaponOptions = [6]string{"", knife, sword, ninjaku, wand, gophermourne}

// skill options
var skillOptions = [4]string{"", strength, agility, intellect}

// magic numbers
const (
	minCoinsEarnedWorking  = 5
	maxCoinsEarnedWorking  = 15
	trainingCost           = 5
	trainingSkillIncrement = 2
)

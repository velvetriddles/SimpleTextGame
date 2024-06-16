package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Player struct {
	Location  *Room
	Inventory map[string]bool
}

type Room struct {
	Name        string
	Description string
	Exits       map[string]*Room
	Items       map[string]string
	Actions     map[string]func(*Player, []string) string
	DoorOpen    bool
}

var player *Player
var world map[string]*Room

const kitchenN = "kitchen"

func initGame() {
	kitchen := &Room{
		Name:        "kitchen",
		Description: "ты находишься на кухне, ",
		Exits:       make(map[string]*Room),
		Items:       map[string]string{"чай": "на столе"},
		Actions:     make(map[string]func(*Player, []string) string),
		DoorOpen:    true,
	}
	corridor := &Room{
		Name:        "corridor",
		Description: "ничего интересного",
		Exits:       make(map[string]*Room),
		Items:       make(map[string]string),
		Actions:     make(map[string]func(*Player, []string) string),
		DoorOpen:    false,
	}
	bedroom := &Room{
		Name:        "bedroom",
		Description: "ты в своей комнате",
		Exits:       make(map[string]*Room),
		Items:       map[string]string{"ключи": "на столе", "конспекты": "на столе", "рюкзак": "на стуле"},
		Actions:     make(map[string]func(*Player, []string) string),
	}
	outside := &Room{
		Name:        "outside",
		Description: "на улице весна",
		Exits:       make(map[string]*Room),
		Items:       make(map[string]string),
		Actions:     make(map[string]func(*Player, []string) string),
	}

	kitchen.Exits["коридор"] = corridor
	corridor.Exits["кухня"] = kitchen
	corridor.Exits["комната"] = bedroom
	corridor.Exits["улица"] = outside
	bedroom.Exits["коридор"] = corridor
	outside.Exits["домой"] = corridor

	for _, room := range []*Room{kitchen, corridor, bedroom, outside} {
		room.Actions["осмотреться"] = lookAround
	}

	bedroom.Actions["взять"] = pickUp
	bedroom.Actions["надеть"] = wear

	for _, room := range []*Room{corridor, bedroom, kitchen, outside} {
		room.Actions["идти"] = moveTo
	}

	corridor.Actions["применить"] = apply
	outside.Actions["применить"] = apply

	world = map[string]*Room{
		"кухня":   kitchen,
		"коридор": corridor,
		"комната": bedroom,
		"улица":   outside,
	}

	player = &Player{
		Location:  kitchen,
		Inventory: make(map[string]bool),
	}
}

func lookAround(p *Player, params []string) string {
	room := p.Location
	var desc string

	// Определяем, является ли текущее действие осмотром комнаты
	isLookingAround := len(params) == 0
	// Использование карты для группирования предметов по расположениям - костыль более или менее
	locationItems := make(map[string][]string)
	for item, location := range room.Items {
		locationItems[location] = append(locationItems[location], item)
	}
	// Из [ключи: на столе, чай: на столе] преобразуем В [ на столе: item []string {"ключи", "чай"} ]
	// Делаем массив строк по расположению
	itemsDescription := []string{}
	for location, items := range locationItems {
		sort.Strings(items)
		itemsDescription = append(itemsDescription, location+": "+strings.Join(items, ", "))
	}
	sort.Strings(itemsDescription)

	// Делаем полноценную строку из описания
	itemsDesc := strings.Join(itemsDescription, ", ")
	if room.Name == "bedroom" && isLookingAround && len(itemsDescription) == 0 {
		desc = "пустая комната"
	} else {
		if isLookingAround && len(itemsDescription) > 0 && room.Name == "bedroom" {
			desc = itemsDesc
		} else if room.Name != kitchenN {
			desc = room.Description
		}
	}

	// Специально для кухни, осмотреться/идти кухня
	if room.Name == kitchenN && isLookingAround {
		if len(world["комната"].Items) > 0 {
			desc += room.Description + itemsDesc + ", надо собрать рюкзак и идти в универ"
		} else {
			desc += room.Description + itemsDesc + ", надо идти в универ"
		}
	} else if room.Name == kitchenN && !isLookingAround {
		desc += "кухня, ничего интересного"
	}

	// Формируем уже конечную строку с добавлением выходов
	desc += ". можно пройти - " + strings.Join(getRoomExits(room), ", ")

	return desc
}

// Были непонятки с порядком выходов, поэтому добавил вручную порядок

var exitOrder = map[string]int{
	"кухня":   1,
	"комната": 2,
	"улица":   3,
	"коридор": 4,
	"домой":   5,
}

func getRoomExits(room *Room) []string {
	exits := make([]string, 0, len(room.Exits))
	for exit := range room.Exits {
		exits = append(exits, exit)
	}
	// Сортировка
	sort.Slice(exits, func(i, j int) bool {
		return exitOrder[exits[i]] < exitOrder[exits[j]]
	})
	return exits
}

func pickUp(p *Player, params []string) string {
	if len(params) < 1 {
		return "не указан предмет для взятия"
	}
	item := params[0]
	location, exists := p.Location.Items[item]
	if exists && location != "" {
		if p.Inventory["рюкзак"] {
			p.Inventory[item] = true
			delete(p.Location.Items, item)
			return "предмет добавлен в инвентарь: " + item
		}
		return "некуда класть"
	}
	return "нет такого"
}

func wear(p *Player, params []string) string {
	if len(params) < 1 {
		return "не указан предмет для ношения"
	}

	item := params[0]

	if item != "рюкзак" {
		return "нельзя надеть: " + item
	}

	location, exists := p.Location.Items[item]
	if exists && location != "" {
		p.Inventory[item] = true
		delete(p.Location.Items, item)
		return "вы надели: " + item
	}
	return "нет такого"
}

func moveTo(p *Player, params []string) string {
	if len(params) < 1 {
		return "не указано направление"
	}
	direction := params[0]
	nextRoom, exists := p.Location.Exits[direction]
	if !exists {
		return "нет пути в " + direction
	}
	if nextRoom == nil {
		return "ошибка: следующая комната не найдена"
	}
	if direction == "улица" && !world["коридор"].DoorOpen {
		return "дверь закрыта"
	}
	p.Location = nextRoom
	return lookAround(p, []string{"move"}) // Передаём параметр, указывающий на перемещение
}

func apply(p *Player, params []string) string {
	if len(params) < 2 {
		return "необходимо указать предмет и объект применения"
	}
	item := params[0]
	object := params[1]

	if !p.Inventory[item] {
		return "нет предмета в инвентаре - " + item
	}

	if item == "ключи" && object == "дверь" && p.Location == world["коридор"] && !world["коридор"].DoorOpen {
		world["коридор"].DoorOpen = true
		return "дверь открыта"
	} else if item == "ключи" && object == "дверь" && p.Location == world["коридор"] && world["коридор"].DoorOpen {
		return "дверь уже открыта"
	}
	return "не к чему применить"
}

func handleCommand(command string) string {
	command = strings.TrimSpace(command) // Удаляем пробелы в начале и конце строки
	parts := strings.Split(command, " ") // Разделяем команду на части
	if len(parts) == 0 {
		return "неизвестная команда"
	}

	action := parts[0]
	params := parts[1:]
	actionFunc, exists := player.Location.Actions[action]
	if !exists {
		return "неизвестная команда"
	}

	if action == "идти" && len(params) < 1 {
		return "не указано направление"
	}

	return actionFunc(player, params)
}

func main() {
	initGame()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите команду:")
	for scanner.Scan() {
		command := scanner.Text()
		command = strings.TrimSpace(command)
		if command == "выход" {
			break
		}
		fmt.Println(handleCommand(command))
		fmt.Println("Введите команду:")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения команды:", err)
	}
}


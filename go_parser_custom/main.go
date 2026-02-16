package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	demoinfocs "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

type aimSnapshot struct {
	tick  int
	viewX float32
	viewY float32
}

func angleDelta(oldX, oldY, newX, newY float32) float64 {
	dx := float64(newX - oldX)
	dy := float64(newY - oldY)
	return math.Sqrt(dx*dx + dy*dy)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a .dem file as an argument")
	}

	demoPath := os.Args[1]

	f, err := os.Open(demoPath)
	checkError(err)
	defer f.Close()

	// Correct parser initialization
	p := demoinfocs.NewParser(f)
	defer p.Close()

	// Extract map name from the header
	var mapName string
	header, err := p.ParseHeader()
	if err != nil {
		log.Println("Error parsing header:", err)
		mapName = "unknown"
	} else {
		mapName = header.MapName
		fmt.Println("Map:", mapName)
	}

	// CSV setup
	fileExists := true
	if _, err := os.Stat("enhanced_kills.csv"); os.IsNotExist(err) {
		fileExists = false
	}

	outFile, err := os.OpenFile("enhanced_kills.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	checkError(err)
	defer outFile.Close()
	writer := csv.NewWriter(outFile)
	writer.Comma = ';'
	defer writer.Flush()

	if !fileExists {
		writer.Write([]string{
			"map", "tick", "killer", "victim", "weapon",
			"posX", "posY", "posZ",
			"viewX", "viewY", "aimDelta", "enemyVisible",
		})
	}

	aimHistory := make(map[string][]aimSnapshot)

	p.RegisterEventHandler(func(e events.Kill) {
		if e.Killer == nil || e.Victim == nil {
			return
		}

		killer := e.Killer
		victim := e.Victim

		snapshots := aimHistory[killer.Name]
		var aimDelta float64 = -1
		if len(snapshots) >= 6 {
			prev := snapshots[len(snapshots)-6]
			curr := snapshots[len(snapshots)-1]
			aimDelta = angleDelta(prev.viewX, prev.viewY, curr.viewX, curr.viewY)
		}

		pos := killer.Position()
		viewX := killer.ViewDirectionX()
		viewY := killer.ViewDirectionY()

		record := []string{
			mapName,
			strconv.Itoa(p.GameState().IngameTick()),
			killer.Name,
			victim.Name,
			e.Weapon.String(),
			fmt.Sprintf("%.4f", pos.X),
			fmt.Sprintf("%.4f", pos.Y),
			fmt.Sprintf("%.4f", pos.Z),
			fmt.Sprintf("%.4f", viewX),
			fmt.Sprintf("%.4f", viewY),
			fmt.Sprintf("%.4f", aimDelta),
			"unknown",
		}

		writer.Write(record)
	})

	for ok, err := p.ParseNextFrame(); ok; ok, err = p.ParseNextFrame() {
		if err != nil {
			log.Fatal(err)
		}

		gs := p.GameState()
		for _, pl := range gs.Participants().Playing() {
			snap := aimSnapshot{
				tick:  gs.IngameTick(),
				viewX: pl.ViewDirectionX(),
				viewY: pl.ViewDirectionY(),
			}
			aimHistory[pl.Name] = append(aimHistory[pl.Name], snap)
			if len(aimHistory[pl.Name]) > 20 {
				aimHistory[pl.Name] = aimHistory[pl.Name][1:]
			}
		}
	}

	fmt.Printf("Parsed %s on map %s and appended data to enhanced_kills.csv", demoPath, mapName)
}

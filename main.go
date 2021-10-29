package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Spells struct {
	XMLName xml.Name `xml:"spells"`
	Spells  []Spell  `xml:"instant"`
}

type Spell struct {
	XMLName                 xml.Name `xml:"instant"`
	Group                   string   `xml:"group,attr"`
	SpellID                 string   `xml:"spellid,attr"`
	Level                   string   `xml:"lvl,attr"`
	Mana                    string   `xml:"mana,attr"`
	GroupCooldown           string   `xml:"groupcooldown,attr"`
	Premium                 string   `xml:"prem,attr"`
	CasterTargetOrDirection string   `xml:"castertargetordirection,attr"`
	Name                    string   `xml:"name,attr"`
	Words                   string   `xml:"words,attr"`
	Aggressive              string   `xml:"aggressive,attr"`
	BlockWalls              string   `xml:"blockwalls,attr"`
	NeedTarget              string   `xml:"needtarget,attr"`
	NeedLearn               string   `xml:"needlearn,attr"`
	Direction               string   `xml:"direction,attr"`
	Exhaustion              string   `xml:"exhaustion,attr"`
	SelfTarget              string   `xml:"selftarget,attr"`
	Range                   string   `xml:"range,attr"`
	Script                  string   `xml:"script,attr"`
}

type Vocations struct {
	XMLName   xml.Name   `xml:"vocations"`
	Vocations []Vocation `xml:"vocation"`
}

type Vocation struct {
	XMLName      xml.Name `xml:"vocation"`
	ID           string   `xml:"id,attr"`
	Name         string   `xml:"name,attr"`
	FromVocation string   `xml:"fromvoc,attr"`
}

func openFile(path string) ([]byte, error) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return input, err
}

func main() {
	xmlSpells, err := os.Open("spells.xml")
	if err != nil {
		log.Fatal(err)
	}
	xmlVoc, err := os.Open("vocations.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer xmlSpells.Close()
	defer xmlVoc.Close()
	var spells *Spells
	var voc *Vocations
	byteValue, _ := ioutil.ReadAll(xmlSpells)
	xml.Unmarshal(byteValue, &spells)
	byteValue, _ = ioutil.ReadAll(xmlVoc)
	xml.Unmarshal(byteValue, &voc)
	for i := 0; i < len(spells.Spells); i++ {
		spell := spells.Spells[i]
		if !strings.Contains("monster/", "") {
			log.Fatal("Arquivo não está na pasta monster.")
		}
		input, _ := openFile(spells.Spells[i].Script)
		if input == nil {
			replace := strings.ReplaceAll(spells.Spells[i].Script, " ", "_")
			log.Println("O sistema tentará abrir o arquivo de um novo jeito: " + replace)
			input, err = openFile(replace)
			if err != nil {
				log.Fatal(err)
			}
		}

		output := bytes.Replace(input, []byte("function onCastSpell("), []byte("local spell = Spell(\"instant\")\n\nfunction spell.onCastSpell("), -1)
		if err = ioutil.WriteFile("new/"+spell.Script, output, 0666); err != nil {
			log.Fatal(err)
		}

		f, err := os.OpenFile("new/"+spell.Script, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		var attr []string

		if spell.Group != "" {
			attr = append(attr, "spell:group(\""+spell.Group+"\")")
		}
		if spell.SpellID != "" {
			attr = append(attr, "spell:spellid(\""+spell.SpellID+"\")")
		}
		if spell.Name != "" {
			attr = append(attr, "spell:name(\""+spell.Name+"\")")
		}
		if spell.Words != "" {
			attr = append(attr, "spell:words(\""+spell.Words+"\")")
		}
		if spell.Level != "" {
			attr = append(attr, "spell:lvl(\""+spell.Level+"\")")
		}
		if spell.Mana != "" {
			attr = append(attr, "spell:mana(\""+spell.Mana+"\")")
		}
		if spell.Range != "" {
			attr = append(attr, "spell:range(\""+spell.Range+"\")")
		}
		if spell.Exhaustion != "" {
			attr = append(attr, "spell:cooldown(\""+spell.Exhaustion+"\")")
		}
		if spell.GroupCooldown != "" {
			attr = append(attr, "spell:groupcooldown(\""+spell.GroupCooldown+"\")")
		}
		if spell.Premium == "1" {
			attr = append(attr, "spell:isPremium(true)")
		}
		if spell.CasterTargetOrDirection == "1" {
			attr = append(attr, "spell:needCasterTargetOrDirection(true)")
		}
		if spell.Aggressive == "1" {
			attr = append(attr, "spell:isAggressive(true)")
		}
		if spell.BlockWalls == "1" {
			attr = append(attr, "spell:blockWalls(true)")
		}
		if spell.NeedTarget == "1" {
			attr = append(attr, "spell:needTarget(true)")
		}
		if spell.NeedLearn == "1" {
			attr = append(attr, "spell:needLearn(true)")
		}
		if spell.Direction == "1" {
			attr = append(attr, "spell:needDirection(true)")
		}
		if spell.SelfTarget == "1" {
			attr = append(attr, "spell:isSelfTarget(true)")
		}

		for u := 0; u < len(attr); u++ {
			_, err := f.WriteString("\n" + attr[u])
			if err != nil {
				log.Fatal(err)
			}
		}

		_, err = f.WriteString("\n" + "spell:register()")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Arquivos processados: %v:%v\r", i, spell.Script)
	}
	fmt.Println("Concluido!")
}

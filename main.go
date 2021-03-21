package main

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

//go:embed check_values.grl
var grls []byte

type MyFact struct {
	IntAttribute     int64
	StringAttribute  string
	BooleanAttribute bool
	FloatAttribute   float64
	TimeAttribute    time.Time
	WhatToSay        string
}

func (mf *MyFact) GetWhatToSay(sentence string) string {
	return fmt.Sprintf("Let's say %s", sentence)
}

//go:embed check_values.grl
var checkValuesRuleGrl []byte

func main() {
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	err := ruleBuilder.BuildRuleFromResource("TutorialRules", "0.0.1", pkg.NewBytesResource(grls))
	if err != nil {
		panic(err)
	}

	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("TutorialRules", "0.0.1")

	myFact := &MyFact{
		IntAttribute:     123,
		StringAttribute:  "Some string value",
		BooleanAttribute: true,
		FloatAttribute:   1.234,
		TimeAttribute:    time.Now(),
	}
	if err := checkValues(knowledgeBase, myFact); err != nil {
		panic(err)
	}
	fmt.Println(myFact.WhatToSay)

	myFact.IntAttribute = 1234
	myFact.WhatToSay = ""

	if err := checkValues(knowledgeBase, myFact); err != nil {
		panic(err)
	}

	fmt.Println(myFact.WhatToSay)
}

func checkValues(kb *ast.KnowledgeBase, fact *MyFact) error {
	dataCtx := ast.NewDataContext()
	err := dataCtx.Add("MF", fact)
	if err != nil {
		return err
	}
	eng := engine.NewGruleEngine()
	err = eng.Execute(dataCtx, kb)
	if err != nil {
		return err
	}
	return nil
}

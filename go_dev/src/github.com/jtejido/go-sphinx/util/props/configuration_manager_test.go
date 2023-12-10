package props

import (
	"fmt"
	"testing"
)

func TestSetComponentProperty(t *testing.T) {
	cm, err := NewConfigurationManager("test_config.sxl")
	if err != nil {
		t.Errorf("config manager error: %s", err.Error())
	}
	fmt.Println(cm.symbolTable)
	//  newBeamWidth := 4711;
	// SetProperty(cm, "beamWidth", String.valueOf(newBeamWidth));

	// DummyComp dummyComp = (DummyComp) cm.lookup("duco");
	// Assert.assertEquals(newBeamWidth, dummyComp.getBeamWidth());
}

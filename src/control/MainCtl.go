
package control

// import (
// 	"fmt"

// 	. "model"
// )

type MainCtl struct {
	DomainGroupStatic *DomainGroup;
	DomainGroupDynamic *DomainGroup;
	DomainMacCtl *MacCtl;
}

var insMainCtl *MainCtl = nil;

func GetMainCtl() *MainCtl {
	if(insMainCtl == nil) {
		ins := &MainCtl{};
		ins.DomainGroupStatic = CreateDomainGroup();
		ins.DomainGroupDynamic = CreateDomainGroup();
		ins.DomainMacCtl = CreateMacCtl();

		insMainCtl = ins;
	}
	return insMainCtl
}
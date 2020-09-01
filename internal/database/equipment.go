package database

// Equipment is the model
type Equipment struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Status   string `json:"status"`
}

// Inventory is the main resource of true data.
// Inventory["VesselCode"] = vesselInventory
// data = vesselInventory["EquipmentCode"]
var Inventory = map[string]map[string]Equipment{}

// EquipmentSet is a index of equipments to find the vessel owner
// vesselOwnerCode = EquipmentSet["EquipmentCode"]
var EquipmentSet = make(map[string]string)

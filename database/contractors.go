package database

import (
	utils "github.com/YaleOpenLab/opensolar/utils"
)

// When contractors are proposing a contract towards something,
// we need to be sure that they are not following the price and are instead giving
// their best quote possible. In this case, a blind auction method is the best
// and that's what we have right now. If we want this to be an auction as well, we
// need to have a specific date of sorts where all the contractors can propose
// contracts immmediately, without latency.

// TODO: Consider some kind of security deposit for Contractors (eg. 5% ) so that they
// don't withdraw on their bid once they are in stage 3 (i.e. chosen by recipient); funds could be given to recipients or other involved stakeholders.
// An alternative is to have a reputation system for contractors.

//TODO: A given Contractor right now is allowed only for one final bid for blind
// auction advantages (i.e. no price disovcery, etc). If we want to change this, we must
// have an auction handler that will take care of this.

// Contractors are created here inheriting properties from Users/Entities
func NewContractor(uname string, pwd string, seedpwd string, Name string, Address string, Description string) (Entity, error) {
	// Create a new entity with the boolean of 'contractor' set to 'true.' This is done just by passing the string "contractor"
	// TODO: Consider other specific information needed for contractors, other than the ones set for users and entities. It can go here, or set as a separate struct.
	return NewEntity(uname, pwd, seedpwd, Name, Address, Description, "contractor")
}

func (contractor *Entity) ProposeContract(panelSize string, totalValue int, location string, years int, metadata string, recIndex int, projectIndex int) (Project, error) {
	var pc Project
	var err error

	// for this, create a new contract and store in the contracts db. We are sorting
	// by stage, so it shouldn't matter a whole lot
	indexCheck, err := RetrieveAllProjects()
	if err != nil {
		return pc, err
	}
	pc.Params.Index = len(indexCheck) + 1
	pc.Params.PanelSize = panelSize
	pc.Params.TotalValue = totalValue
	pc.Params.Location = location
	pc.Params.Years = years
	pc.Params.Metadata = metadata
	pc.Params.DateInitiated = utils.Timestamp()
	iRecipient, err := RetrieveRecipient(recIndex)
	if err != nil {
		return pc, err
	}
	pc.Params.ProjectRecipient = iRecipient
	pc.Stage = 2 // 2 since we need to filter this out while retrieving the propsoed contracts
	pc.Contractor = *contractor
	// instead of storing in this proposedcontracts slice, store it as a project, but not a contract and retrieve by stage
	err = pc.Save()
	// don't insert the project since the contractor's projects are not final
	return pc, err
}

package statusdb

import (
	"github.com/alexeyqian/gochain/utils"
	"fmt"
	"github.com/alexeyqian/gochain/entity"
)

type UndoableSet struct{
	storage *Storage
	dataBucket string
	stateBucket string	
	latestRevision uint32
	count uint32
}

type UndoState struct{	
	revision uint32
	newIDs []string
	removedValues map[string][]byte
	oldValues map[string][]byte	
}

func NewUndoableSet(s Storage, name, elementType string) (*UndoableSet, error){
	if s.HasBucket(name){
		return nil, fmt.Errorf("cannot create same set in storage")
	}

	us := UndoableSet{
		storage: s,
		dataBucket: name + "_data",
		stateBucket: name + "_state",
		latestRevision: 0,
		count: 0,
	}

	// create data bucket and state bucket
	err := s.CreateBucket(us.dataBucket)
	if err != nil{
		return nil, err
	}
	err = s.CreateBucket(us.stateBucket)
	if err != nil{
		return nil, err
	}

	return &us, nil
}

func (us *UndoableSet) Create(key string, data []byte) error{	
	if key == "" {
		return fmt.Errorf("create: entity must have ID.")
	}
	
	if us.storage.HasKey(us.dataBucket, key) {
		return fmt.Errorf("create: entity already exist.")
	}

	// 1. update state bucket
	err := us.onCreate(key)
	if err != nil{
		return err
	}

	// 2. save to data bucket
	err = us.storage.Put(us.dataBucket, key, data)
	if err != nil{ 
		return err
	}

	us.count++
	return nil
}

func (us *UndoableSet) onCreate(key string) error{
	if !us.hasSession() { 
		return nil
	}

	state := us.latestState()
	state.newIDs = append(state.newIDs, key)
	us.storage.Put(us.stateBucket, state.revision, utils.Serialize(state))

	return nil
}

func (us *UndoableSet) Update(key string, data []byte) error{
	if key == "" {
		return fmt.Errorf("update: must have ID.")
	}
	if !us.storage.HasKey(us.dataBucket, key) {
		return fmt.Errorf("update: key not exist.")
	}

	existing, err := us.storage.Get(us.dataBucket, key)
	if err != nil{
		return err
	}
	// 1. update state bucket
	err = us.onUpdate(key, existing)
	if err != nil{
		return err
	}

	// 2. udpate data bucket
	return us.storage.Put(us.dataBucket, key, data)	
}

func (us *UndoableSet) onUpdate(key string, existing []byte) error{
	if !us.hasSession() { 
		return
	}

	state := us.latestState()

	_, found := utils.FindString(state.newIDs, key)
	if found {
		return nil// do nothing if it's new added
	}

	_, ok := state.oldValues[key]
	if ok{
		return nil // do nothing if old value already exists
	}

	_, ok = state.removedValues[key]
	if ok{
		panic("cannot modify deleted entity")
	}

	state.oldValues[key] = existing
	us.storage.Put(us.stateBucket, state.revision, utils.Serialize(state))
	
	return nil
}


func (us *UndoableSet) Delete(key string) error{
	if key == ""{
		return fmt.Errorf("delete: id cannot be empty")
	}
	if !us.storage.HasKey(us.dataBucket, key) {
		return fmt.Errorf("delete: key not exist.")
	}

	existing, err := us.storage.Get(us.dataBucket, id)
	if err != nil{
		return err
	}
	
	// 1. udpate state bucket
	err = us.onDelete(key, existing)
	if err != nil{
		return err
	}

	// 2. update data bucket
	err = us.storage.Delete(us.dataBucket, id)	
	if err != nil{
		return err
	}

	us.count--
}


func onDelete(key string, existing []byte) error {
	if !us.hasSession() { 
		return nil
	}

	state := us.latestState()

	// if the removed on is new added, just remove it from new id list
	i, found := utils.FindString(state.newIDs, key)
	if found{
		// remove key from newIDs
		state.newIDs = append(state.newIDs[:i], state.newIDs[i+1:]...)
		return nil
	}

	// if the removed on is updated
	// add it to removed list
	// remove it from old value list
	_, ok := state.oldValues[key]
	if ok {
		delete(state.oldValues, key)
		state.removedValues[key] = existing
	}
		
	// if the removed one is already removed
	_, ok = state.removedValues[key]
	if ok {
		return nil
	}

	state.RemoveValues[key] = existing
	
	us.storage.Put(us.stateBucket, state.revision, utils.Serialize(state))
}

func (us *UndoableSet) Get(key string) ([]byte, error){
	return us.storage.Get(us.dataBucket, key)
}

func (us *UndoableSet) Size() int{
	return us.count
}

func (us *UndoableSet) StartUndoSession(){
	var state UndoState
	state.revision = us.latestRevision + 1
	us.storage.Put(us.stateBucket, state.revision, utils.Serialize(state))
}

func (us *UndoableSet) CommitFromLastSession(){
	if !us.hasSession() {
		return
	}

	lastState := us.latestState()
	// remove the latest element
	us.storage.Delete(name+"state", lastState.revision)
}

func (us *UndoableSet) UndoChangesFromLastSession(){
	if !us.hasSession(){
		return
	}

	head := us.latestState()

	// undo modifications
	for _, item := range head.oldValues{
		// restore old values to database
	}

	// undo creations
	for _, item := range head.newIDs{
		// remove new created entity from database
	}

	// undo deletions
	for  _, item := range head.removedValues{
		// restore removed values to database
	}

	us.entityStateTable.Pop() // remove latest
	--us.revision
}

// helpers
func (us *UndoableSet) hasSession() bool{
	return us.latestRevision > 0
}

func (us *UndoableSet) latestState() *UndoState{
	data, err := us.storage.Get(us.stateBucket, us.latestRevision)
	if err != nil{
		panic("undoable set: cannot get latest state")
	}

	var state UndoState
	entity.Deserialize(&state, data)
	return &state
}
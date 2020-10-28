package statusdb

import (
	"github.com/alexeyqian/gochain/entity"
)

type UndoState struct{	
	Revision uint32
	NewIDs []string
	RemovedValues map[string][]byte
	OldValues map[string][]byte	
}

type UndoableSet struct{
	storage *Storage
	name string
	elementType string
	revision uint32	
}

func NewUndoableSet(s Storage, name, elementType string) (*UndoableSet, error){
	if s.HasBucket(name){
		return nil, fmt.Errorf("cannot create same set in storage")
	}

	us := UndoableSet{
		storage: s,
		name: name,
		elementType: elementType,
		revision: 0,
	}

	// create data set and data state set
	s.CreateSet(name, entityType)
	s.CreateSet(name + "_state", "UndoState")

	return &us, nil
}

func (us *UndoableSet) Create(key string, data []byte) error{	
	if key == "" {
		return fmt.Errorf("entity must have ID.")
	}

	err := storage.Put(entity.GetID(e), e)
	if err != nil{ // sth like duplicated id
		return err
	}

	us.onCreate(entity)

	return nil
}

func (us *UndoableSet) Update(e Entity) error{
	if !entity.HasID(e) {
		return fmt.Errorf("entity must have ID.")
	}

	existing, err := entityTable.Find(e.ID)
	if err != nil{
		return err
	}

	us.onModify(existing)

	//e.UpdatedAt = time.Now().Unix()
	storage.Put(entity.GetID(e), e)

	return nil
}

func (us *UndoableSet) Remove(id string) error{
	if id == ""{
		return fmt.Errorf("id cannot be empty")
	}

	existing, err := entityTable.Find(id)
	if err != nil{
		return err
	}
	
	us.onRemove(existing)

	Storage.Remove(addPrefix(id))

	return nil
}

func (us *UndoableSet) Find(id string) (Entity, error){
	return entityTable.Find(id)
}

func (us *UndoableSet) Latest() (Entity, error){
	return entityTable.Latest()
}

func (us *UndoableSet) Size() int{
	return entityTable.Size()
}

func (us *UndoableSet) StartUndoSession(){
	var state UndoState
	state.Revision = ++us.Revision
	us.StateList = append(u.StateList, state)
}

func (us *UndoableSet) CommitFromLastSession(){
	if len(us.StateList) <= 0{
		return
	}
	// remove the first element
	us.StateList = us.StateList[1:]
}

func (us *UndoableSet) UndoChangesFromLastSession(){
	if !us.hasSession(){
		return
	}

	head := us.latestState()

	// undo modifications
	for _, item := range head.OldValues{
		// restore old values to database
	}

	// undo creations
	for _, item := range head.NewIDs{
		// remove new created entity from database
	}

	// undo deletions
	for  _, item := range head.RemovedValues{
		// restore removed values to database
	}

	us.entityStateTable.Pop() // remove latest
	--us.Revision
}


func (us *UndoableSet) onCreate(e *Entity){
	if !us.hasSession() { 
		return
	}

	state := us.latestState()
	state.NewIDs := append(state.NewIDs, e.ID)

}

func onRemove(e *Entity) {
	if !us.hasSession() { 
		return
	}

	state := us.latestState()

	// if the removed on is new added since last session
	// just remove it from new id list
	if state.NewIDs.Contains(e.ID){
		state.NewIDs.Remove(e.ID)
		return
	}

	// 2. if the removed on is updated one
	// add it to removed list
	// remove it from old value list
	if oldValues.Contains(e.ID) {
		oldValues.Remove(e.ID)
		removedValues.Add(e)
	}
		
	// if the removed one is already removed
	if removedValues.Contains(e){
		return
	}

	state.RemoveValues.Add(e)	
}

func (us *UndoableSet) onModify(e Entity){
	if !us.hasSession() { 
		return
	}

	head := us.latestState()
	if head.NewIDs.contains(e.ID) {
		return // do nothing if it's new added
	}
	if head.OldValues.contains(e.ID){
		return  // do nothing if old value already exists
	}

	if head.RemovedValues.contains(e.ID){
		panic("cannot modify deleted entity")
	}

	head.OldValues.add(e)
}

// helpers
func (us *UndoableSet) hasSession() bool{
	return len(us.entityStateTable) > 0
}

// TODO: rename to head()
func (us *UndoableSet) latestState() *UndoState{
	return &us.entityStateTable[len(us.entityStateTable) - 1]
}
package statusdb


type UndoState struct{
	Revision uint32
	NewIDs []string
	RemovedValues map[string]Entity
	OldValues map[string]Entity
}

type UndoableSet struct{
	revision uint32
	container []entity
	stateList []UndoState
}

func (us *UndoableSet) Create(e Entity) (error){
	entity.CreatedAt = time.Now().Unix()
	entity.UpdatedAt = entity.CreatedAt
	container.Insert(entity)
	OnCreate(entity)
	retur nil
}

func (us *UndoableSet) Update(e Entity) error{
	existing, err := container.Find(e.ID)
	if err != nil{
		return err
	}

	OnModify(existing)

	e.UpdatedAt = time.Now().Unix()
	container.Replace(e.ID, e)

	return nil
}

func (us *UndoableSet) Remove(id string) error{
	existing, err := container.Find(id)
	if err != nil{
		return err
	}
	
	OnRemove(existing)

	container.Remove(e.ID)
	return nil
}

func (us *UndoableSet) Find(id string) (Entity, error){
	return container.Find(id)
}

func (us *UndoableSet) Latest() (Entity, error){
	if container.Size() <= 0{
		return nil, fmt.Errorf("set is empty")
	}

	return container.Latest()
}

func (us *UndoableSet) Size() int{
	return container.Size()
}

func NewUndoableSet(name string, type interface{}){

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

	us.stateList.Pop() // remove latest
	--us.Revision
}


func (us *UndoableSet) onCreate(e *Entity){
	if !us.hasSession() { 
		return
	}

	state := us.latestState()
	state.NewIDs := append(state.NewIDs, e.Id)

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
	return len(us.stateList) > 0
}

// TODO: rename to head()
func (us *UndoableSet) latestState() *UndoState{
	return &us.stateList[len(us.stateList) - 1]
}
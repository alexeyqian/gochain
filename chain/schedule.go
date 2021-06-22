package chain

import "time"

func (c *Chain) GetNextWitness() string {
	witnesses := c.GetAllWitnesses()
	runningTime := time.Now().Unix() - c.GetGpo().GenesisTime
	index := runningTime % 21
	if index == 0 {
		// reschedule
	} else {
		witness = witnesses[index]
	}

	return witnesses[index]
}

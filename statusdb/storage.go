package statusdb

// used as DI for easy testing
// memory data provider for tesging
// file data provider for production
type Storage interface {
	Open()
	Close()
	RemoveAll()

	// New GP style DB interface
	GetAll(bucket string) [][]byte

	/*
	 * return error if buket or key is not found
	 */
	Get(bucket, key string) ([]byte, error)

	/*
	 * if buket is not exist, create a new bucket
	 * if key not exist, create new key and value pair
	 * if key is already exist, replace the old value
	 */
	Put(bucket, key string, data []byte) error

	/*
	 * delete key/value pairs, if buket or key not exist, just do nothing
	 */
	Delete(bucket, key string) error

	/* NOT USED
	 * delete a bucket, if not exist, do nothing
	 */
	//DeleteBucket(bucket string) error
}

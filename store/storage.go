package store

// Simple key/value storage interface
type Storage interface {
	Open() error
	Close() error
	Remove() error

	// Get all items from bucket
	GetAll(bucket string) ([][]byte, error)

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

	CreateBucket(bucket string) error
	HasBucket(bucket string) bool
	HasKey(bucket, key string) bool

	/* NOT USED
	 * delete a bucket, if not exist, do nothing
	 */
	//DeleteBucket(bucket string) error
}

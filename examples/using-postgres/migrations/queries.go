package migrations

const (
	CreateTable         = "CREATE TABLE IF NOT EXISTS customers (id int PRIMARY KEY , name varchar(5) , email varchar(30) , phone bigint);"
	DroopTable          = "Drop table If EXISTS customers"
	AddCountry          = "ALTER TABLE customers ADD COLUMN country varchar(20);"
	DropCountry         = "ALTER TABLE customers DROP COLUMN country;"
	DropPhone           = "ALTER TABLE customers DROP COLUMN phone"
	AddPhone            = "ALTER TABLE customers ADD COLUMN phone bigint;"
	AlterType           = "ALTER TABLE customers ALTER COLUMN name TYPE TEXT;"
	ResetType           = "ALTER TABLE customers ALTER COLUMN name TYPE varchar(5);"
	AddNotNullColumn    = "ALTER TABLE customers ADD COLUMN date_of_joining DATE NOT NULL DEFAULT CURRENT_DATE;"
	DeleteNotNullColumn = "ALTER TABLE customers DROP COLUMN date_of_joining;"
	AlterPrimaryKey     = "ALTER TABLE customers ALTER COLUMN id SET DATA TYPE bigint USING id::bigint;"
	ResetPrimaryKey     = "ALTER TABLE customers ALTER COLUMN id SET DATA TYPE int USING id::int;"
)

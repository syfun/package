package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_package "github.com/syfun/package/pkg/package"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sqlx.DB
}


func New(dsn string) (*Database, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) InsertPackage(p *_package.Package) (*_package.Package, error) {
	q := `INSERT INTO packages (name) VALUES ($1) RETURNING id;`
	if err := d.db.QueryRow(q, p.Name).Scan(&p.ID); err != nil {
		return nil, err
	}
	return p, nil
}

func (d *Database) ListPackages(fuzzyName string) ([]*_package.Package, error) {
	q := `SELECT * FROM packages WHERE name LIKE $1`
	name := fmt.Sprintf("%%%v%%", fuzzyName)
	rows, err := d.db.Queryx(q, name)
	if err != nil {
		return nil, err
	}

	packages := make([]*_package.Package, 0)
	for rows.Next() {
		var p _package.Package
		if err := rows.StructScan(&p); err != nil {
			return nil, err
		}
		packages = append(packages, &p)
	}
	return packages, nil
}

func (d *Database) GetPackage(name string) (*_package.Package, error) {
	p := _package.Package{}
	q := `SELECT * FROM packages WHERE name=$1`
	if err := d.db.Get(&p, q, name); err != nil {
		return nil, err
	}
	return &p, nil
}

func (d *Database) InsertVersion(v *_package.Version) (*_package.Version, error) {
	q := `INSERT INTO versions (name, file_name, size, checksum, package_id) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	if err := d.db.QueryRow(q, v.Name, v.FileName, v.Size, v.Checksum, v.PackageID).Scan(&v.ID); err != nil {
		return nil, err
	}
	return v, nil
}

func (d *Database) GetVersion(packageID int64, name string) (*_package.Version, error) {
	q := `SELECT * FROM versions WHERE package_id=$1 and name=$2`
	v := new(_package.Version)
	if err := d.db.Get(&v, q, packageID, name); err != nil {
		return nil, err
	}
	return v, nil
}

func (d *Database) DeletePackage(name string) error {
	q := `DELETE FROM packages WHERE name=$1`
	if _, err := d.db.Exec(q, name); err != nil {
		return err
	}
	return nil
}

func (d *Database) DeleteVersion(packageID int64, name string) error {
	q := `DELETE FROM versions WHERE package_id=$1 and name=$2`
	if _, err := d.db.Exec(q, packageID, name); err != nil {
		return err
	}
	return nil
}


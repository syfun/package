package _package

import (
	"fmt"
)

type Package struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Version struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	FileName  string `json:"file_name"`
	Size      int64  `json:"size"`
	Checksum  string `json:"checksum"`
	PackageID int64  `json:"package_id"`
}

type PackageIn struct {
	Name string `json:"name"`
}

type VersionIn struct {
	PackageName string `json:"package_name"`
	Name        string `json:"name"`
	FileName    string `json:"file_name"`
	Content     []byte `json:"content"`
	PackageID   string `json:"package_id"`
}

type Service interface {
	AddPackage(p *PackageIn) (*Package, error)
	ListPackages(fuzzyName string) ([]*Package, error)
	GetPackage(name string) (*Package, error)
	AddVersion(v *VersionIn) (*Version, error)
	DownloadPackage(packageName string, versionName string) ([]byte, error)
	DeletePackage(name string) error
	DeleteVersion(name string) error
}

type Repo interface {
	InsertPackage(p *Package) (*Package, error)
	ListPackages(fuzzyName string) ([]*Package, error)
	GetPackage(name string) (*Package, error)
	InsertVersion(v *Version) (*Version, error)
	GetVersion(packageID int64, versionName string) ([]byte, error)
	DeletePackage(name string) error
	DeleteVersion(packageID int64, versionName string) error
}

type Storage interface {
	Upload(name string, d []byte) (size int64, err error)
	Download(name string) ([]byte, error)
}

type service struct {
	repo    Repo
	storage Storage
}

func (s *service) AddPackage(p *PackageIn) (*Package, error) {
	pg, err := s.repo.InsertPackage(&Package{Name: p.Name})
	if err != nil {
		return nil, fmt.Errorf("cannot add package: %w", err)
	}
	return pg, nil
}

func (s *service) ListPackages(fuzzyName string) ([]*Package, error) {
	pgs, err := s.repo.ListPackages(fuzzyName)
	if err != nil {
		return nil, fmt.Errorf("cannot list packages: %w", err)
	}
	return pgs, nil
}

func (s *service) GetPackage(name string) (*Package, error) {
	pg, err := s.repo.GetPackage(name)
	if err != nil {
		return nil, fmt.Errorf("cannot get package: %w", err)
	}
	return pg, nil
}

func (s *service) AddVersion(v *VersionIn) (*Version, error) {
	size, err := s.storage.Upload(fmt.Sprintf("%v/%v/%v", v.PackageName, v.Name, v.FileName), v.Content)
	if err != nil {
		return nil, fmt.Errorf("cannot add package version: %w", err)
	}
	p, err := s.repo.GetPackage(v.PackageName)
	if err != nil {
		return nil, fmt.Errorf("cannot add package version: %w", err)
	}

	version := &Version{
		Name:      v.Name,
		FileName:  v.FileName,
		Size:      size,
		Checksum:  "",
		PackageID: p.ID,
	}
	version, err = s.repo.InsertVersion(version)
	if err != nil {
		return nil, fmt.Errorf("cannot add package version: %w", err)
	}
	return version, nil
}

func (s *service) DownloadPackage(packageName, versionName string) ([]byte, error) {
	p, err := s.repo.GetPackage(packageName)
	if err != nil {
		return nil, fmt.Errorf("cannot download package version: %w", err)
	}
	d, err := s.repo.GetVersion(p.ID, versionName)
	if err != nil {
		return nil, fmt.Errorf("cannot download package: %w", err)
	}
	return d, nil
}

func (s *service) DeletePackage(name string) error {
	if err := s.repo.DeletePackage(name); err != nil {
		return fmt.Errorf("cannot delete package: %w", err)
	}
	return nil
}
func (s *service) DeleteVersion(packageName, versionName string) error {
	p, err := s.repo.GetPackage(packageName)
	if err != nil {
		return fmt.Errorf("cannot delete version: %w", err)
	}
	if err := s.repo.DeleteVersion(p.ID, versionName); err != nil {
		return fmt.Errorf("cannot delete version: %w", err)
	}
	return nil
}

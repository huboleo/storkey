package internal

type Keyring interface {
	Pull(repo string) error
	Add(repo string) error
	ListForRepo(repo string) error
	ListAll() error
}

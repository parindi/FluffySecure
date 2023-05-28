package authentication

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/go-crypt/crypt"
	"github.com/go-crypt/crypt/algorithm"
	"gopkg.in/yaml.v3"
)

type FileUserDatabase interface {
	Save() (err error)
	Load() (err error)
	GetUserDetails(username string) (user DatabaseUserDetails, err error)
	SetUserDetails(username string, details *DatabaseUserDetails)
}

// NewYAMLUserDatabase creates a new YAMLUserDatabase.
func NewYAMLUserDatabase(filePath string, searchEmail, searchCI bool) (database *YAMLUserDatabase) {
	return &YAMLUserDatabase{
		RWMutex:     &sync.RWMutex{},
		Path:        filePath,
		Users:       map[string]DatabaseUserDetails{},
		Emails:      map[string]string{},
		Aliases:     map[string]string{},
		SearchEmail: searchEmail,
		SearchCI:    searchCI,
	}
}

// YAMLUserDatabase is a user details database that is concurrency safe database and can be reloaded.
type YAMLUserDatabase struct {
	*sync.RWMutex

	Path    string
	Users   map[string]DatabaseUserDetails
	Emails  map[string]string
	Aliases map[string]string

	SearchEmail bool
	SearchCI    bool
}

// Save the database to disk.
func (m *YAMLUserDatabase) Save() (err error) {
	m.RLock()

	defer m.RUnlock()

	if err = m.ToDatabaseModel().Write(m.Path); err != nil {
		return err
	}

	return nil
}

// Load the database from disk.
func (m *YAMLUserDatabase) Load() (err error) {
	yml := &DatabaseModel{Users: map[string]UserDetailsModel{}}

	if err = yml.Read(m.Path); err != nil {
		return fmt.Errorf("error reading the authentication database: %w", err)
	}

	m.Lock()

	defer m.Unlock()

	if err = yml.ReadToFileUserDatabase(m); err != nil {
		return fmt.Errorf("error decoding the authentication database: %w", err)
	}

	return m.LoadAliases()
}

// LoadAliases performs the loading of alias information from the database.
func (m *YAMLUserDatabase) LoadAliases() (err error) {
	if m.SearchEmail || m.SearchCI {
		for k, user := range m.Users {
			if m.SearchEmail && user.Email != "" {
				if err = m.loadAliasEmail(k, user); err != nil {
					return err
				}
			}

			if m.SearchCI {
				if err = m.loadAlias(k); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (m *YAMLUserDatabase) loadAlias(k string) (err error) {
	u := strings.ToLower(k)

	if u != k {
		return fmt.Errorf("error loading authentication database: username '%s' is not lowercase but this is required when case-insensitive search is enabled", k)
	}

	for username, details := range m.Users {
		if k == username {
			continue
		}

		if strings.EqualFold(u, details.Email) {
			return fmt.Errorf("error loading authentication database: username '%s' is configured as an email for user with username '%s' which isn't allowed when case-insensitive search is enabled", u, username)
		}
	}

	m.Aliases[u] = k

	return nil
}

func (m *YAMLUserDatabase) loadAliasEmail(k string, user DatabaseUserDetails) (err error) {
	e := strings.ToLower(user.Email)

	var duplicates []string

	for username, details := range m.Users {
		if k == username {
			continue
		}

		if strings.EqualFold(e, details.Email) {
			duplicates = append(duplicates, username)
		}
	}

	if len(duplicates) != 0 {
		duplicates = append(duplicates, k)

		return fmt.Errorf("error loading authentication database: email '%s' is configured for for more than one user (users are '%s') which isn't allowed when email search is enabled", e, strings.Join(duplicates, "', '"))
	}

	if _, ok := m.Users[e]; ok && k != e {
		return fmt.Errorf("error loading authentication database: email '%s' is also a username which isn't allowed when email search is enabled", e)
	}

	m.Emails[e] = k

	return nil
}

// GetUserDetails get a DatabaseUserDetails given a username as a value type where the username must be the users actual
// username.
func (m *YAMLUserDatabase) GetUserDetails(username string) (user DatabaseUserDetails, err error) {
	m.RLock()

	defer m.RUnlock()

	u := strings.ToLower(username)

	if m.SearchEmail {
		if key, ok := m.Emails[u]; ok {
			return m.Users[key], nil
		}
	}

	if m.SearchCI {
		if key, ok := m.Aliases[u]; ok {
			return m.Users[key], nil
		}
	}

	if details, ok := m.Users[username]; ok {
		return details, nil
	}

	return user, ErrUserNotFound
}

// SetUserDetails sets the DatabaseUserDetails for a given user.
func (m *YAMLUserDatabase) SetUserDetails(username string, details *DatabaseUserDetails) {
	if details == nil {
		return
	}

	m.Lock()

	m.Users[username] = *details

	m.Unlock()
}

// ToDatabaseModel converts the YAMLUserDatabase into the DatabaseModel for saving.
func (m *YAMLUserDatabase) ToDatabaseModel() (model *DatabaseModel) {
	model = &DatabaseModel{
		Users: map[string]UserDetailsModel{},
	}

	m.RLock()

	for user, details := range m.Users {
		model.Users[user] = details.ToUserDetailsModel()
	}

	m.RUnlock()

	return model
}

// DatabaseUserDetails is the model of user details in the file database.
type DatabaseUserDetails struct {
	Username    string
	Digest      algorithm.Digest
	Disabled    bool
	DisplayName string
	Email       string
	Groups      []string
}

// ToUserDetails converts DatabaseUserDetails into a *UserDetails given a username.
func (m DatabaseUserDetails) ToUserDetails() (details *UserDetails) {
	return &UserDetails{
		Username:    m.Username,
		DisplayName: m.DisplayName,
		Emails:      []string{m.Email},
		Groups:      m.Groups,
	}
}

// ToUserDetailsModel converts DatabaseUserDetails into a UserDetailsModel.
func (m DatabaseUserDetails) ToUserDetailsModel() (model UserDetailsModel) {
	return UserDetailsModel{
		HashedPassword: m.Digest.Encode(),
		DisplayName:    m.DisplayName,
		Email:          m.Email,
		Groups:         m.Groups,
	}
}

// DatabaseModel is the model of users file database.
type DatabaseModel struct {
	Users map[string]UserDetailsModel `yaml:"users" valid:"required"`
}

// ReadToFileUserDatabase reads the DatabaseModel into a YAMLUserDatabase.
func (m *DatabaseModel) ReadToFileUserDatabase(db *YAMLUserDatabase) (err error) {
	users := map[string]DatabaseUserDetails{}

	var udm *DatabaseUserDetails

	for user, details := range m.Users {
		if udm, err = details.ToDatabaseUserDetailsModel(user); err != nil {
			return fmt.Errorf("failed to parse hash for user '%s': %w", user, err)
		}

		users[user] = *udm
	}

	db.Users = users

	return nil
}

// Read a DatabaseModel from disk.
func (m *DatabaseModel) Read(filePath string) (err error) {
	var (
		content []byte
		ok      bool
	)

	if content, err = os.ReadFile(filePath); err != nil {
		return fmt.Errorf("failed to read the '%s' file: %w", filePath, err)
	}

	if len(content) == 0 {
		return ErrNoContent
	}

	if err = yaml.Unmarshal(content, m); err != nil {
		return fmt.Errorf("could not parse the YAML database: %w", err)
	}

	if ok, err = govalidator.ValidateStruct(m); err != nil {
		return fmt.Errorf("could not validate the schema: %w", err)
	}

	if !ok {
		return fmt.Errorf("the schema is invalid")
	}

	return nil
}

// Write a DatabaseModel to disk.
func (m *DatabaseModel) Write(fileName string) (err error) {
	var (
		data []byte
	)

	if data, err = yaml.Marshal(m); err != nil {
		return err
	}

	return os.WriteFile(fileName, data, fileAuthenticationMode)
}

// UserDetailsModel is the model of user details in the file database.
type UserDetailsModel struct {
	HashedPassword string   `yaml:"password" valid:"required"`
	DisplayName    string   `yaml:"displayname" valid:"required"`
	Email          string   `yaml:"email"`
	Groups         []string `yaml:"groups"`
	Disabled       bool     `yaml:"disabled"`
}

// ToDatabaseUserDetailsModel converts a UserDetailsModel into a *DatabaseUserDetails.
func (m UserDetailsModel) ToDatabaseUserDetailsModel(username string) (model *DatabaseUserDetails, err error) {
	var d algorithm.Digest

	if d, err = crypt.Decode(m.HashedPassword); err != nil {
		return nil, err
	}

	return &DatabaseUserDetails{
		Username:    username,
		Digest:      d,
		Disabled:    m.Disabled,
		DisplayName: m.DisplayName,
		Email:       m.Email,
		Groups:      m.Groups,
	}, nil
}

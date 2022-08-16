# EnvManager

## Concept

The envManager knows **profiles** and **storages**. A **storage** is a secure location where secrets are stored (e.g.
keepass). A **profile** is a set of environment variables which should be set / unset together. The environment variables
can either have a constant value or get their value from a storage. A profile can depend on other profiles, which means
that loading the first profile will load all dependencies as well. Even circular dependencies are allowed.
**Directory mappings** are supported since version 1.1. With directory mappings, envManager will remember which profiles
are used in which directory. You can either remember the currently loaded profiles or select from all configured profiles.

## Setup

Add the function in `wrapper.sh` to your shell (e.g. `. ./wrapper.sh`). The file contains the function `envManager`
which in turn calls `envManager-bin` (assuming it is in your PATH). Should your `envManager-bin` not be in your PATH,
replace `envManager-bin` with an absolute path to the binary (e.g. `/home/john.doe/envManager/envManager-bin`).

## Usage

Create an initial config with `envManager config init`. By default, the application creates a `.envManager.yml` in your
home directory.
You can now add storages and profiles by hand or let the application do it for you. Run `envManager config add storage`
to add your first storage and `envManager config add profile` to add your first profile. After that, you can easily
copy&paste the profile and storage configurations.
To add a directory mapping, navigate to the directory you want to map and either load the profiles you will need here and
call `envManager config add mapping`. Or navigate to the directory and call `envManager config add mapping --select` to
get a list of all your profiles and check the ones you want to map to this directory.

## Available storage adapters

### Keepass / KeepassX / KeepassXC

This adapter can read keepass2 files (.kdbx). Its config contains one key, `path` which contains an absolute path to the
kdbx file.

**Example config**
```yaml
storages:
  myStorageName:
    type: keepass
    config:
      path: /home/john.doe/myKeepassFile.kdbx
```

### Pass

This adapter supports gpg encrypted secrets, as created by the [pass](https://www.passwordstore.org/) or
[gopass](https://github.com/gopasspw/gopass) password manager. Password stores outside of `~/.password-store` are
currently not supported. The `prefix` config option specifies a directory inside the `~/.password-store` directory, if
the value is not specified (like in a config from version 1.1.0 and earlier) or set to an empty string, the path of an
entry is used as absolute within the password store.

**Example**
```yaml
storages:
  myStorageName:
    type: pass
    config:
      prefix: "prefix"
```

#### Using the prefix

Assume you have the entry `~/.password-store/shared/admin-account` in your pass storage. You can now configure your
storage adapter and a profile like this:
```yaml
storages:
  old:
    type: pass
    config: {}
  shared:
    type: pass
    config:
      prefix: "shared"
profiles:
  adminAcc:
    storage: shared
    path: admin-account
    # The keys constEnv, env and dependsOn are not relevant for this example and are omitted
```
Note how the path is no longer `shared/admin-account` but only `admin-account`. When you ask envManager to load the
profile `adminAcc`, it will automatically prefix the path with the `prefix` value of the storage adapter.

## FAQ 

### Can I use multiple storages for one profile?

No, but you can create one profile which depends on multiple profiles. If you load the "main" profile, the dependencies
will be loaded automatically.

### Can I have dependencies across multiple storages?

Yes.

### I want to have two profiles depending on each other, will I run into an endless loop?

No, circular dependencies are resolved without looping endlessly.

### Can I use multiple configuration files?

Yes, since version 1.4.0 you can create `.envManager.yml` files in any directory. When running an envManager command,
all config files from the current working directory up to the root are loaded with decreasing priority (meaning the
config file in the current working directory overrides one closer to the file system root and the one in your home
directory). You can view the discovered config files and their order by running `envManager debug files`.

## Extending envManager

The envManager can be easily extended by programming other storage adapters. Each storage adapter must implement the
`StorageAdapter` interface and define a constant type identifier (like `const KeepassTypeIdentifier = "keepass"`).
The type identifier is used in the config file to select the storage type. Additionally, the storage provider must be
registered in `secretsStorage/StorageAdapter.go` in the following methods:

- `CreateStorageAdapter()` This method is a factory for storage adapters. Add your storage adapter as new `case` and
  assign a new instance of your adapter to the `storage` variable.
- `GetStorageAdapterTypes()` This method returns all available storage adapter types. Just add your type identifier in
  the slice.
- `GetStorageAdapterDefaultConfig()` This method returns the default config of a storage adapter. Add your storage
  adapter as new `case` and assign a new, empty instance to the `storage` variable.

## Test data

In the `/testData` directory is a dummy `keepass.kdbx` containing the following entries. The password for this database is `1234`.

- `entry1` with `user1` and `pass1` as well as the additional attribute `advanced1` with `advanced1-value`
- `group1/g1e1` with `g1e1-user` and `g1e1-pass` and the additional attribute `g1e1-advanced1` with `g1e1-advanced1-value`

## Used libraries

- [github.com/gopasspw/gopass](https://github.com/gopasspw/gopass)
- [github.com/josa42/go-prompt](https://github.com/josa42/go-prompt)
- [github.com/manifoldco/promptui](https://github.com/manifoldco/promptui)
- [github.com/spf13/cobra](https://github.com/spf13/cobra)
- [github.com/tobischo/gokeepasslib/v3](https://github.com/tobischo/gokeepasslib)
- [gopkg.in/errgo.v2](https://gopkg.in/errgo.v2)
- [gopkg.in/yaml.v2 v2.4.0](https://gopkg.in/yaml.v2)

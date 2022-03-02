# op-keyring

Simple wrapper for [1password CLI][op_getting_started] to integrate with [gnome-keyring][gnome_keyring_doc].

The usual way of running `op` is to first `eval $(op signin my)` within the current shell so that an environment variable containing the session token is populated. However, this results in the following friction:

1. The user is prompted for their password for every new shell, each of which will get their own session token
2. Session tokens are valid for 30 minutes by default, so the user will need to re-enter their password for every shell, every 30 minutes

`op-keyring` lessens this friction by storing the session token in the keyring, and wrapping the `op` executable to utilise that session token. Any shell with access to the keyring will utilise this shared session token; so re-authentication will only be required at most once every 30 minutes.

At a high level, it:

1. Fetches a session token for `op` from the keyring, if it exists
2. Runs the `op` command, passing the session token along with any user provided arguments.
3. If the command resulted in an "Invalid session token" error (or the token did not exist in the keyring), it runs `op signin my`, capturing the session token, writing it to the keyring, and then re-running the command

## Security considerations

[gnome-keyring][gnome_keyring_doc] uses dbus for communication. Where the keyring is unlocked, any process running as the logged-in user can access objects within the keyring, unless restricted by AppArmor or similar. By default, most distributions provide little; or more commonly; no restriction to keyring access. This is somewhat negated by the session token only being valid for 30 minutes, however the risk is valid, depending on your threat model. Note that this same risk exists for other keyring implementations such as gnupg-keyring etc.

## TODO

* write example apparmor policy for restricting access to session token object in keyring

[op_getting_started]: https://support.1password.com/command-line-getting-started/
[gnome_keyring_doc]: https://wiki.gnome.org/Projects/GnomeKeyring

## Requirements

* [gnome-keyring][gnome_keyring_doc]
* A configured [op][op_getting_started] installation (e.g you have previously run `op signin my.1password.com $username`, and `op` exists within your `$PATH`)

## Usage

```shell
# compile binary
go build

# [optional] alias `op`
echo "'op'='/home/user/path/to/op-keyring/op-keyring'" >> ~/.zsh_aliases
```

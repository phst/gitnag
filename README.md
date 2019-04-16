# Gitnag

Gitnag is a primitive tool that reminds you to keep your Git repositories up to
date.  Specifically, it checks whether the worktree is clean and whether there
are any commits that should be pushed upstream.  It also automatically
fast-forwards clean repositories from upstream if possible.  It displays issues
using the default notification method of your graphical environment.

Run the program as follows

```shell
gitnag -config gitnag.json
```

where the configuration file `gitnag.json` specifies the directories that
gitnag should search as a JSON file of the form

```json
{
    "Directories": {
        "/directory/one": {},
        "/directory/two": {},
    }
}
```

To be most effective, run gitnag as a cron job or similar.

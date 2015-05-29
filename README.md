NCAA Build Utility
==================

Command line utility to fire off NCAA Bamboo5 build plans. Usage:

```bash
$ ncaabuild <dev|qa|prod>
```

In order to run a build, the utility must authenticate with Bamboo. You must provide your username and password either
as environment variables or explicitly each time you run the command. To set environment variables, use the following
example:

```bash
export NCAA_BARCA_BAMBOO_USER=mstills
export NCAA_BARCA_BAMBOO_PASS=mattsgreatpassword
```

If you do not want to use environment variables, or if you have env vars set but want to override them, you may use
the utility's flags to pass them in from the command line:

```bash
$ ncaabuild dev --user=mstills --password=mattsgreatpassword
```


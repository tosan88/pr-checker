# pr-checker
Checks FT UPP's PRs which are too long time open

## Install

```
go get -u github.com/tosan88/pr-checker
```

## Usage

```
Usage: pr-checker [OPTIONS]

Checks FT UPP's PRs which are too long time open

Options:
  --token=""          The GitHub Api's OAuth Token. ($TOKEN)
  --contributors=""   The list of contributors. Only those repos will be considered where these people contributed. Optional ($CONTRIBUTORS)
  --min-days="14"     The number of minimum days which an open PR could stay open. Only PRs which are opened more than that number of days are retrieved. Optional ($MIN_DAYS)
```

```
#Linux
export TOKEN=<your-GH-token>
./pr-checker

#Windows
set TOKEN=<your-GH-token>
pr-checker.exe
```

### Sample output

```
2016/10/06 15:33:12 Application starting with args [pr-checker.exe]
2016/10/06 15:33:38 PR https://github.com/Financial-Times/o-ads/pull/83 (Remove unused code) open by adgad(Arjun Gadhia) since 2016-10-06T11:42:41Z, updated at 2016-10-06T12:15:44Z
2016/10/06 15:33:48 PR https://github.com/Financial-Times/up-service-files/pull/824 (Added JAVA_OPTS memory constraints to Binary Writer.) open by carlosroman(Carlos) since 2016-10-
06T10:21:22Z, updated at 2016-10-06T10:25:50Z
...
2016/10/06 15:35:36 Application finished. It was active 143.566 seconds

```
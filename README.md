# Xray CLI scan

## About this plugin
This plugin provides an easy way for getting security issue and licenses about your project dependencies.

## Installation with JFrog CLI
Installing the latest version:

`$ jfrog plugin install xray-scan`

Installing a specific version:

`$ jfrog plugin install xray-scan@version`

Uninstalling a plugin

`$ jfrog plugin uninstall xray-scan`

## Usage
### Commands
There is 2 way for using xray scan :
* Standard input : you redirect the output of "mvnw dependency:list" or "go list -m" to the scan like this
```bash
mvn dependency:list | jfrog xray scan
go list -m all | jfrog xray scan
```
This will display a summary of the vulnerabilities (high/medium/low) and license for all the dependencies found.

* "--component" flag : you will search for vulnerabilities and license for a single component
```bash
jfrog xray-scan scan --component "golang.org/x/net v1.8.2"
```


### Environment variables

## Additional info
None.

## Release Notes
The release notes are available [here](RELEASE.md).

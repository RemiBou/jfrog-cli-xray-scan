# Xray CLI scan

## About this plugin
This plugin provides an easy way for getting security issue and licenses about your project dependencies.


## Installation with JFrog CLI
Since this plugin is currently not included in [JFrog CLI Plugins Registry](https://github.com/jfrog/jfrog-cli-plugins-reg), it needs to be built and installed manually. Follow these steps to install and use this plugin with JFrog CLI.
1. Make sure JFrog CLI is installed on you machine by running ```jfrog```. If it is not installed, [install](https://jfrog.com/getcli/) it.
2. Create a directory named ```plugins``` under ```~/.jfrog/``` if it does not exist already.
3. Clone this repository.
4. CD into the root directory of the cloned project.
5. Run ```go build``` to create the binary in the current directory.
6. Copy the binary into the ```~/.jfrog/plugins``` directory.

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

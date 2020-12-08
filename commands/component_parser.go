package commands

import "regexp"

type component string
type matcher struct {
	pattern   *regexp.Regexp
	prefix    string
	extractor func(match [][]string) string
}

// Trying to extract an xray component id from a single line of free text.
//
//  see https://www.jfrog.com/confluence/display/JFROG/Xray+REST+API#XrayRESTAPI-ComponentIdentifiers
//
//	Maven			gav://group:artifact:version		gav://ant:ant:1.6.5
//	Docker			docker://Namespace/name:tag			docker://jfrog/artifactory-oss:latest
//	RPM				rpm://os-version:package:version	rpm://7:zsh:5.0.2-14.e17_2.2
//	Debian  		deb://vendor:dist:package:version	deb://ubuntu:trustee:acl:2.2.49-2
//	NuGet			nuget://module:version				nuget://log4net:9.0.1
//	Generic file	generic://sha256:<Checksum>/name	generic://sha256:244fd47e07d1004f0aed9c156aa09083c82bf8944eceb67c946ff7430510a77b/foo.jar
//	NPM				npm://package:version				npm://mocha:2.4.5
//	Python			pip://package:version				pip://raven:5.13.0
//	Composer		composer://package:version			composer://nunomaduro/collision:1.1
//	Golang			go://package:version				go://github.com/ethereum/go-ethereum:1.8.2
//	Alpine			alpine://branch:package:version		alpine://3.7:htop:2.0.2-r0
var matchers = []matcher{
	//maven
	{
		pattern: regexp.MustCompile("([a-z0-9\\-\\._]+:[a-z][a-z0-9\\-_\\.]+)(?::[a-z]+)?(?::[a-z][a-z_0-9\\-]+)?:([a-zA-Z0-9\\-\\.]+)(?::[a-z]+)?"),
		prefix:  "gav",
		extractor: func(match [][]string) string {
			return match[0][1] + ":" + match[0][2]
		},
	},
	//go
	{
		pattern: regexp.MustCompile("((?:[a-zA-Z-._,~]+\\/)+(?:[a-zA-Z-._,~]+)) v([0-9]\\.[0-9]\\.[0-9](?:-[a-zA-Z0-9\\-\\.]*)*)"),
		prefix:  "go",
		extractor: func(match [][]string) string {
			return match[0][1] + ":" + match[0][2]
		},
	},
}

func (c component) toString() string {
	return string(c)
}

func parse(line string) (component, bool) {
	for _, matcher := range matchers {
		match := matcher.pattern.FindAllStringSubmatch(line, -1)
		if len(match) > 0 {
			return component(matcher.prefix + "://" + matcher.extractor(match)), true
		}
	}
	return "", false
}

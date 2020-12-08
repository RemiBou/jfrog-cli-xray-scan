package commands

import "testing"

func Test_parse(t *testing.T) {
	tests := []struct {
		name   string
		line   string
		want   component
		wantOk bool
	}{
		{
			name:   "Mvn pkg from mvn dep list cmd",
			line:   "com.googlecode.json-simple:json-simple:jar:1.1.1:compile",
			want:   "gav://com.googlecode.json-simple:json-simple:1.1.1",
			wantOk: true,
		},
		{
			name:   "Mvn pkg from mvn dep list cmd (full line)",
			line:   "[INFO]    com.googlecode.json-simple:json-simple:jar:1.1.1:compile -- module json.simple (auto)",
			want:   "gav://com.googlecode.json-simple:json-simple:1.1.1",
			wantOk: true,
		},
		{
			name:   "Simple mvn pkg parsed without scope",
			line:   "json-simple:json-simple:1.1.1",
			want:   "gav://json-simple:json-simple:1.1.1",
			wantOk: true,
		},
		{
			name:   "Simple mvn pkg parsed with scope",
			line:   "json-simple:json-simple:1.1.1:compile",
			want:   "gav://json-simple:json-simple:1.1.1",
			wantOk: true,
		},
		{
			name:   "Too simple mvn pkg not parsed",
			line:   "json-simple:1.1.1",
			want:   "",
			wantOk: false,
		},
		{
			name:   "Spring RELEASE",
			line:   "[INFO]    org.springframework:spring-aop:jar:5.2.9.RELEASE:compile -- module spring.aop [auto]",
			want:   "gav://org.springframework:spring-aop:5.2.9.RELEASE",
			wantOk: true,
		},
		{
			name:   "Router SNAPSHOT",
			line:   " [INFO]    org.jfrog.router:router-grpc-api:jar:7.x-SNAPSHOT:compile -- module router.grpc.api (auto)",
			wantOk: true,
			want:   "gav://org.jfrog.router:router-grpc-api:7.x-SNAPSHOT",
		},
		{
			name:   "Underscore in artifact",
			line:   "[INFO]    com.google.errorprone:error_prone_annotations:jar:2.0.19:compile -- module error.prone.annotations (auto)",
			wantOk: true,
			want:   "gav://com.google.errorprone:error_prone_annotations:2.0.19",
		},
		{
			name:   "Netty",
			line:   "[INFO]    io.netty:netty-transport-native-epoll:jar:linux-x86_64:4.1.53.Final:compile -- module io.netty.transport.epoll [auto]",
			wantOk: true,
			want:   "gav://io.netty:netty-transport-native-epoll:4.1.53.Final",
		},
		{
			name:   "ignore maven step",
			line:   "[INFO] --- maven-dependency-plugin:3.1.2:list (default-cli) @ access-common-rest ---",
			wantOk: false,
		},
		{
			name:   "Simple Go pkg",
			line:   "golang.org/x/text v0.3.3",
			want:   "go://golang.org/x/text:0.3.3",
			wantOk: true,
		},
		{
			name:   "Simple Go pkg with pseudo version",
			line:   "github.com/anmitsu/go-shlex v0.0.0-20161002113705-648efa622239",
			want:   "go://github.com/anmitsu/go-shlex:0.0.0-20161002113705-648efa622239",
			wantOk: true,
		},
		{
			name:   "Simple Go pkg not parsed",
			line:   "golang.org/x/text 0.3.3",
			want:   "",
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOk := parse(tt.line)
			if got != tt.want {
				t.Errorf("parse() got = %v, want %v", got, tt.want)
			}
			if gotOk != tt.wantOk {
				t.Errorf("parse() got1 = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

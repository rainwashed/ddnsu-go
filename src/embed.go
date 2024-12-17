package main

import _ "embed"

//go:embed example.toml
var ExampleConfigEmbed []byte

//go:embed example.service
var ExampleServiceEmbed []byte

//go:embed example.sh
var ExampleShellEmbed []byte

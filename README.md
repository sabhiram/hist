# hist

A tool to convert history commands to language specific scripts.

## Idea

1. Grab commands from `history` and annotate / prune / edit them into scripts of various languages.
2. Each language supported is a plugin that registers with the tool - this should make for a clean API and allow arbitrary output types (think: ansible, go, python, shell ...)

## Install

```
go get github.com/sabhiram/hist
```

## Usage

```
$ hist [<cmd> [options]]

If "cmd" is empty, it runs the "tag" command (see below).  Valid commands 
include:

	tag <val>	-	Set the start of a tag block with the string "<val>".  If 
					"<val>" is not specified, the previous tag is closed.

	version		- 	Print the version of the "hist" tool.
```	


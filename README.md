# git-org

A Git subcommand to do github orgaization.

## Install

    go get -u github.com/yhinoz/git-org

## Usage

	git-org [flags] [command]
	
	Available Commands:
	  grep        run grep the specified github organization repository
	  help        Help about any command
	  repos       Show the specified github organization repository
	  version     Print the version
	
	Flags:
	  -h, --help         help for git-org
	  -o, --org string   organization name
	
	Use "git-org [command] --help" for more information about a command.

## Options

To access organization's private repositories, generate a GitHub Personal Access Token, make sure you have the following environment variable.

    export GITHUB_ACCESS_TOKEN=your_personal_access_token

or your configuration to your `.env` file in the home directory:

    GITHUB_ACCESS_TOKEN=your_personal_access_token


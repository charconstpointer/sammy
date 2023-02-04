# README.md for Sammy

## Overview
Sammy is a program that generates a report from a GitHub user's activity using OpenAI's GPT-3 language model. 

## Usage
`OPEN_AI_TOKEN=token GITHUB_TOKEN=token sammy -user=<username> -max_tokens=<number> -public=<true|false>`

## Flags
`-user`: Specifies the GitHub username of the user for whom the report is generated.

`-max_tokens`: Specifies the maximum number of tokens to be used by the GPT-3 model in the report.

`-public`: Specifies whether the report should be publicly accessible or not. Set to `true` for public access and `false` for private access.

## Example
`OPEN_AI_TOKEN=token GITHUB_TOKEN=token sammy -user=foo -max_tokens=50 -public=true`

This command generates a report for the GitHub user `foo`, using a maximum of 50 tokens from the GPT-3 model, and sets the report to be publicly accessible.

## Note
Please make sure to have the necessary permissions and access to generate the report, as well as sufficient tokens for the GPT-3 model and GitHub API, before running the program.

You might need to tweak `-max_tokens` if your report happens to be a lengthy one and gets only partially generated


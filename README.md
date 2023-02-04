# README.md for Sammy

## Overview
Sammy is a program that generates a report from a GitHub user's activity using OpenAI's GPT-3 language model. 

## Usage
`OPEN_AI_TOKEN=token GITHUB_TOKEN=token sammy -user=<username> -max_tokens=<number> -public=<true|false> -from=<date> -to=<date>`

## Flags
`-user`: Specifies the GitHub username of the user for whom the report is generated.

`-max_tokens`: Specifies the maximum number of tokens to be used by the GPT-3 model in the report.

`-public`: Specifies whether the report should be publicly accessible or not. Set to `true` for public access and `false` for private access.

`-from`: Specifies the start date for the report, in the format of RFC3339.

`-to`: Specifies the end date for the report, in the format of RFC3339.

## Example
`OPEN_AI_TOKEN=token GITHUB_TOKEN=token sammy -user=foo -max_tokens=100 -public=true -from="2022-01-01T00:00:00Z" -to="2022-12-31T23:59:59Z"`

This command generates a summary of the GitHub user foo activity between January 1, 2022 and December 31, 2022, using a maximum of 100 tokens from the GPT-3 model and only including public events.

## Note
Please make sure to have the necessary permissions and access to generate the report, as well as sufficient tokens for the GPT-3 model and GitHub API, before running the program.

You might need to tweak `-max_tokens` if your report happens to be a lengthy one and gets only partially generated

If time range is not provided it will default to last 24 hours
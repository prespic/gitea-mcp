# Gitea MCP Server

[繁體中文](README.zh-tw.md) | [简体中文](README.zh-cn.md)

**Gitea MCP Server** is an integration plugin designed to connect Gitea with Model Context Protocol (MCP) systems. This allows for seamless command execution and repository management through an MCP-compatible chat interface.

[![Install with Docker in VS Code](https://img.shields.io/badge/VS_Code-Install_Server-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=gitea&inputs=[{%22id%22:%22gitea_token%22,%22type%22:%22promptString%22,%22description%22:%22Gitea%20Personal%20Access%20Token%22,%22password%22:true}]&config={%22command%22:%22docker%22,%22args%22:[%22run%22,%22-i%22,%22--rm%22,%22-e%22,%22GITEA_ACCESS_TOKEN%22,%22docker.gitea.com/gitea-mcp-server%22],%22env%22:{%22GITEA_ACCESS_TOKEN%22:%22${input:gitea_token}%22}}) [![Install with Docker in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install_Server-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=gitea&inputs=[{%22id%22:%22gitea_token%22,%22type%22:%22promptString%22,%22description%22:%22Gitea%20Personal%20Access%20Token%22,%22password%22:true}]&config={%22command%22:%22docker%22,%22args%22:[%22run%22,%22-i%22,%22--rm%22,%22-e%22,%22GITEA_ACCESS_TOKEN%22,%22docker.gitea.com/gitea-mcp-server%22],%22env%22:{%22GITEA_ACCESS_TOKEN%22:%22${input:gitea_token}%22}}&quality=insiders)

## Table of Contents

- [Gitea MCP Server](#gitea-mcp-server)
  - [Table of Contents](#table-of-contents)
  - [What is Gitea?](#what-is-gitea)
  - [What is MCP?](#what-is-mcp)
  - [🚧 Installation](#-installation)
    - [Usage with VS Code](#usage-with-vs-code)
    - [📥 Download the official binary release](#-download-the-official-binary-release)
    - [🔧 Build from Source](#-build-from-source)
    - [📁 Add to PATH](#-add-to-path)
  - [🚀 Usage](#-usage)
  - [✅ Available Tools](#-available-tools)
  - [🐛 Debugging](#-debugging)
  - [🛠 Troubleshooting](#-troubleshooting)

## What is Gitea?

Gitea is a community-managed lightweight code hosting solution written in Go. It is published under the MIT license. Gitea provides Git hosting including a repository viewer, issue tracking, pull requests, and more.

## What is MCP?

Model Context Protocol (MCP) is a protocol that allows for the integration of various tools and systems through a chat interface. It enables seamless command execution and management of repositories, users, and other resources.

## 🚧 Installation

### Usage with VS Code

For quick installation, use one of the one-click install buttons at the top of this README.

For manual installation, add the following JSON block to your User Settings (JSON) file in VS Code. You can do this by pressing `Ctrl + Shift + P` and typing `Preferences: Open User Settings (JSON)`.

Optionally, you can add it to a file called `.vscode/mcp.json` in your workspace. This will allow you to share the configuration with others.

> Note that the `mcp` key is not needed in the `.vscode/mcp.json` file.

```json
{
  "mcp": {
    "inputs": [
      {
        "type": "promptString",
        "id": "gitea_token",
        "description": "Gitea Personal Access Token",
        "password": true
      }
    ],
    "servers": {
      "gitea-mcp": {
        "command": "docker",
        "args": [
          "run",
          "-i",
          "--rm",
          "-e",
          "GITEA_ACCESS_TOKEN",
          "docker.gitea.com/gitea-mcp-server"
        ],
        "env": {
          "GITEA_ACCESS_TOKEN": "${input:gitea_token}"
        }
      }
    }
  }
}
```

### 📥 Download the official binary release

You can download the official release from [official Gitea MCP binary releases](https://gitea.com/gitea/gitea-mcp/releases).

### 🔧 Build from Source

You can download the source code by cloning the repository using Git:

```bash
git clone https://gitea.com/gitea/gitea-mcp.git
```

Before building, make sure you have the following installed:

- make
- Golang (Go 1.24 or later recommended)

Then run:

```bash
make install
```

### 📁 Add to PATH

After installing, copy the binary gitea-mcp to a directory included in your system's PATH. For example:

```bash
cp gitea-mcp /usr/local/bin/
```

## 🚀 Usage

This example is for Cursor, you can also use plugins in VSCode.
To configure the MCP server for Gitea, add the following to your MCP configuration file:

- **stdio mode**

```json
{
  "mcpServers": {
    "gitea": {
      "command": "gitea-mcp",
      "args": [
        "-t",
        "stdio",
        "--host",
        "https://gitea.com"
        // "--token", "<your personal access token>"
      ],
      "env": {
        // "GITEA_HOST": "https://gitea.com",
        // "GITEA_INSECURE": "true",
        "GITEA_ACCESS_TOKEN": "<your personal access token>"
      }
    }
  }
}
```

- **sse mode**

```json
{
  "mcpServers": {
    "gitea": {
      "url": "http://localhost:8080/sse",
      "headers": {
        "Authorization": "Bearer <your personal access token>"
      }
    }
  }
}
```

- **http mode**

```json
{
  "mcpServers": {
    "gitea": {
      "url": "http://localhost:8080/mcp",
      "headers": {
        "Authorization": "Bearer <your personal access token>"
      }
    }
  }
}
```

**Default log path**: `$HOME/.gitea-mcp/gitea-mcp.log`

> [!NOTE]
> You can provide your Gitea host and access token either as command-line arguments or environment variables.
> Command-line arguments have the highest priority

Once everything is set up, try typing the following in your MCP-compatible chatbox:

```text
list all my repositories
```

## ✅ Available Tools

The Gitea MCP Server supports the following tools:

|             Tool             |    Scope     |                       Description                        |
| :--------------------------: | :----------: | :------------------------------------------------------: |
|       get_my_user_info       |     User     |      Get the information of the authenticated user       |
|        get_user_orgs         |     User     | Get organizations associated with the authenticated user |
|         create_repo          |  Repository  |                 Create a new repository                  |
|          fork_repo           |  Repository  |                    Fork a repository                     |
|        list_my_repos         |  Repository  |  List all repositories owned by the authenticated user   |
|        create_branch         |    Branch    |                   Create a new branch                    |
|        delete_branch         |    Branch    |                     Delete a branch                      |
|        list_branches         |    Branch    |            List all branches in a repository             |
|        create_release        |   Release    |           Create a new release in a repository           |
|        delete_release        |   Release    |            Delete a release from a repository            |
|         get_release          |   Release    |                      Get a release                       |
|      get_latest_release      |   Release    |          Get the latest release in a repository          |
|        list_releases         |   Release    |            List all releases in a repository             |
|          create_tag          |     Tag      |                     Create a new tag                     |
|          delete_tag          |     Tag      |                       Delete a tag                       |
|           get_tag            |     Tag      |                        Get a tag                         |
|          list_tags           |     Tag      |              List all tags in a repository               |
|      list_repo_commits       |    Commit    |             List all commits in a repository             |
|       get_file_content       |     File     |          Get the content and metadata of a file          |
|        get_dir_content       |     File     |           Get a list of entries in a directory           |
|         create_file          |     File     |                    Create a new file                     |
|         update_file          |     File     |                 Update an existing file                  |
|         delete_file          |     File     |                      Delete a file                       |
|      get_issue_by_index      |    Issue     |                Get an issue by its index                 |
|       list_repo_issues       |    Issue     |             List all issues in a repository              |
|         create_issue         |    Issue     |                    Create a new issue                    |
|     create_issue_comment     |    Issue     |               Create a comment on an issue               |
|          edit_issue          |    Issue     |                       Edit a issue                       |
|      edit_issue_comment      |    Issue     |                Edit a comment on an issue                |
| get_issue_comments_by_index  |    Issue     |          Get comments of an issue by its index           |
|  get_pull_request_by_index   | Pull Request |             Get a pull request by its index              |
|   list_repo_pull_requests    | Pull Request |          List all pull requests in a repository          |
|     create_pull_request      | Pull Request |                Create a new pull request                 |
|         search_users         |     User     |                     Search for users                     |
|       search_org_teams       | Organization |           Search for teams in an organization            |
|         search_repos         |  Repository  |                 Search for repositories                  |
| get_gitea_mcp_server_version |    Server    |         Get the version of the Gitea MCP Server          |

## 🐛 Debugging

To enable debug mode, add the `-d` flag when running the Gitea MCP Server with sse mode:

```sh
./gitea-mcp -t sse [--port 8080] --token <your personal access token> -d
```

## 🛠 Troubleshooting

If you encounter any issues, here are some common troubleshooting steps:

1. **Check your PATH**: Ensure that the `gitea-mcp` binary is in a directory included in your system's PATH.
2. **Verify dependencies**: Make sure you have all the required dependencies installed, such as `make` and `Golang`.
3. **Review configuration**: Double-check your MCP configuration file for any errors or missing information.
4. **Consult logs**: Check the logs for any error messages or warnings that can provide more information about the issue.

Enjoy exploring and managing your Gitea repositories via chat!

# GitHub Configuration

This directory contains GitHub-specific configuration files and Copilot agents.

## Copilot Agents

### CodeGuardian (`@CodeGuardian`)

A senior code reviewer agent with three personalities:
- **The Perfectionist** - Enforces highest standards
- **The Pragmatist** - Focuses on production concerns
- **The Mentor** - Teaches and guides improvement

**Usage:**
```
@CodeGuardian please review my changes
@CodeGuardian what could be improved in this file?
@CodeGuardian check the security of this implementation
```

**Location:** `~/.copilot/agents/CodeGuardian.md` (globally available)

**Note:** CodeGuardian reviews code but does not edit it. Use it for critical feedback and learning.

## Other Configurations

- **workflows/** - GitHub Actions CI/CD workflows
- **ISSUE_TEMPLATE/** - Issue templates for bugs and features
- **pull_request_template.md** - PR template
- **dependabot.yml** - Dependency update automation

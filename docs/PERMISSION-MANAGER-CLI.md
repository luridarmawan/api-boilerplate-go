# Permission Manager CLI Tool

## Overview

A command line tool to add permissions to access IDs by managing their associated group permissions.

## Architecture

```mermaid
graph TD
    A[CLI Tool] --> B[Parse Arguments]
    B --> C[Validate Access ID]
    C --> D[Get User's Group]
    D --> E[Find Permission]
    E --> F[Check if Permission Exists in Group]
    F --> G[Add Permission to Group]
    G --> H[Success Response]
    
    C --> I[Access ID Not Found]
    E --> J[Permission Not Found]
    F --> K[Permission Already Exists]
    
    I --> L[Error Exit]
    J --> L
    K --> M[Warning Message]
```

## Database Relationships

```mermaid
erDiagram
    ACCESS ||--o{ GROUP : belongs_to
    GROUP ||--o{ GROUP_PERMISSIONS : has_many
    PERMISSION ||--o{ GROUP_PERMISSIONS : has_many
    
    ACCESS {
        string id PK
        string name
        string email
        string api_key
        uint group_id FK
    }
    
    GROUP {
        uint id PK
        string name
        string description
    }
    
    PERMISSION {
        uint id PK
        string name
        string resource
        string action
    }
    
    GROUP_PERMISSIONS {
        uint group_id FK
        uint permission_id FK
    }
```

## Usage

```bash
# Basic usage
./permission-manager [access_id] [resource] [action]

# Example
./permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 configurations create
```

## Implementation Plan

### Phase 1: Core Structure
- [x] Analyze database schema and relationships
- [x] Design CLI tool architecture
- [ ] Create main CLI structure with argument parsing
- [ ] Implement database connection and configuration loading

### Phase 2: Business Logic
- [ ] Create functions to validate access ID exists
- [ ] Create functions to find and validate permissions by resource and action
- [ ] Create functions to check if permission already exists in group
- [ ] Implement core logic to add permission to group

### Phase 3: Polish & Testing
- [ ] Add proper error handling and user feedback
- [ ] Create build scripts for the CLI tool
- [ ] Test the tool with various scenarios
- [ ] Create documentation and usage examples

## Expected Behavior

### Success Cases
```bash
$ ./permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 configurations create
✓ Permission 'configurations:create' successfully added to group 'Editor' for access ID '019847a9-4efb-72c1-92fb-2c5eab3335d1'
```

### Error Cases
```bash
# Invalid access ID
$ ./permission-manager invalid-id configurations create
✗ Error: Access ID 'invalid-id' not found

# Permission not found
$ ./permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 invalid-resource create
✗ Error: Permission 'invalid-resource:create' not found in database

# Permission already exists
$ ./permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 configurations create
⚠ Warning: Permission 'configurations:create' already exists in group 'Editor'
```

## File Structure

```
cmd/
└── permission-manager/
    └── main.go
scripts/
├── build-permission-manager.sh
└── build-permission-manager.bat
docs/
└── PERMISSION-MANAGER-CLI.md
```

## Integration

The tool integrates with existing codebase:
- Uses same configuration system (`configs/config.go`)
- Leverages existing models (`access`, `permission`, `group`)
- Follows same database connection pattern (`internal/database`)
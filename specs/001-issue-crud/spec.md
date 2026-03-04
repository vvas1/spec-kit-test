# Feature Specification: Create and Edit Issues

**Feature Branch**: `001-issue-crud`  
**Created**: 2025-03-04  
**Status**: Draft  
**Input**: User description: "Create and edit issues with title, description, status, and assignee"

## Clarifications

### Session 2025-03-04

- Q: Which status set and transition rule should we use? → A: Option C — Different set (e.g. To Do, In Progress, Review, Done); exact list specified in implementation plan; transitions between statuses are free unless otherwise specified.
- Q: Where does the list of assignable users come from? → A: Option A — In-app list only; the app maintains its own list of users/assignees; no external directory or authentication required for this feature.
- Q: What are the maximum lengths for title and description? → A: Option B — Title 200 characters, description 10,000 characters.
- Q: When two users edit the same issue at once, how should the system behave? → A: Option A — Last-write-wins; the last save overwrites; no conflict detection or warning.
- Q: How should the issue list behave when there are many issues? → A: Option B — Paginated list (e.g. 20–50 issues per page); exact page size defined in implementation plan.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create a new issue (Priority: P1)

As a user, I can create a new issue by entering a title, description, optional status, and optional assignee so that work items are captured and can be tracked.

**Why this priority**: Creating issues is the foundation; without it there is nothing to edit or view.

**Independent Test**: Can be fully tested by submitting a new issue form and verifying the issue appears in the system with the entered data.

**Acceptance Scenarios**:

1. **Given** I am on the create-issue flow, **When** I enter a title and description and submit, **Then** the system creates the issue and I see confirmation; the issue has a title, description, and default status.
2. **Given** I am creating an issue, **When** I optionally set status and assignee and submit, **Then** the issue is stored with those values.
3. **Given** I submit an issue without a title, **When** validation runs, **Then** the system rejects the submission and indicates that title is required.

---

### User Story 2 - Edit an existing issue (Priority: P2)

As a user, I can open an existing issue and change its title, description, status, or assignee so that I can correct or update the issue as work evolves.

**Why this priority**: Editing is the natural follow-on to creation; users must be able to update issues after creation.

**Independent Test**: Can be fully tested by opening an existing issue, changing one or more fields, saving, and verifying the stored issue reflects the changes.

**Acceptance Scenarios**:

1. **Given** an issue exists, **When** I open it for editing and change the title and save, **Then** the issue’s title is updated and I see confirmation.
2. **Given** an issue exists, **When** I change description, status, or assignee and save, **Then** all changed fields are persisted and visible when I reopen the issue.
3. **Given** I am editing an issue, **When** I clear the title and try to save, **Then** the system rejects the change and indicates that title is required.

---

### User Story 3 - View issue list and detail (Priority: P3)

As a user, I can view a list of issues and open a single issue to see its title, description, status, and assignee so that I can find and inspect issues before creating or editing.

**Why this priority**: Viewing supports create and edit; users need to see where to create and which issue to edit.

**Independent Test**: Can be fully tested by opening the issue list and then opening one issue and verifying all displayed fields match the stored data.

**Acceptance Scenarios**:

1. **Given** at least one issue exists, **When** I open the issue list, **Then** I see a paginated set of issues (e.g. 20–50 per page) with at least title and status (or equivalent summary), and a way to move between pages.
2. **Given** I am on the list, **When** I open one issue, **Then** I see its full title, description, status, and assignee.
3. **Given** no issues exist, **When** I open the issue list, **Then** I see an empty state and a way to create the first issue.

---

### Edge Cases

- What happens when the user enters an extremely long title or description? System MUST enforce a maximum of 200 characters for title and 10,000 for description and show clear validation messages when limits are exceeded.
- How does the system handle an assignee that is removed or invalid? System MUST display assignee only when valid; allow unassigning (empty assignee); invalid assignee MUST not block saving other fields.
- What happens when two users edit the same issue at the same time? System MUST use last-write-wins: the last save overwrites the issue; no conflict detection or warning; no silent data loss beyond overwrite.
- What happens when the user leaves required fields empty on create or edit? System MUST prevent submit and show which fields are required.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow users to create an issue with at least title and description; title is required.
- **FR-002**: System MUST allow users to set optional status and assignee on create and on edit.
- **FR-003**: System MUST allow users to edit an existing issue’s title, description, status, and assignee; title remains required after edit.
- **FR-004**: System MUST persist issues so that created and updated data is retained and visible on subsequent views.
- **FR-005**: System MUST validate required fields (title) and enforce maximum length limits (title: 200 characters, description: 10,000 characters); validation errors MUST be shown to the user when limits are exceeded.
- **FR-006**: System MUST support a paginated list of issues (page size e.g. 20–50, defined in implementation plan) and a way to open a single issue to view or edit its details.
- **FR-007**: System MUST display for each issue: title, description, status, and assignee (or “unassigned” when none).

### Key Entities

- **Issue**: A single work item. Attributes: title (required), description, status, assignee (optional). Identified by a stable id so it can be opened and edited. Status is one of a defined set (e.g. To Do, In Progress, Review, Done); exact values defined in implementation plan; any status may be changed to any other (free transitions) unless restricted later.
- **Assignee**: A reference to a user that can be assigned to an issue; may be empty (unassigned). Selected from the app’s own in-app list of users/assignees; no external directory or authentication required for this feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can create a new issue with title and description in under one minute.
- **SC-002**: Users can edit an existing issue and see their changes reflected immediately after save (or after refresh where applicable).
- **SC-003**: Required-field and length validation prevents invalid data from being saved and shows clear messages.
- **SC-004**: Issue list and issue detail views show consistent data with what was created or last edited.

## Assumptions

- Status values are a defined set (e.g. To Do, In Progress, Review, Done); exact list is defined in the implementation plan; transitions between statuses are unrestricted unless specified otherwise.
- Assignee is chosen from an in-app list of users/assignees maintained by the app; user management (e.g. how users are added to the list) and authentication are out of scope for this feature unless otherwise specified.
- Concurrent edits are resolved by last-write-wins; no conflict detection or warning is required for this feature.
- The issue list is paginated; page size (e.g. 20–50) is defined in the implementation plan.

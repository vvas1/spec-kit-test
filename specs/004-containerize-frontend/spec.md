# Feature Specification: Containerize Frontend (Host Port 3000, Container Port 5137)

**Feature Branch**: `004-containerize-frontend`  
**Created**: 2026-03-05  
**Status**: Draft  
**Input**: User description: "want to containerize frontend. outer port should be 3000. inner 5137"

**Scope**: This feature is for **local development only**. The container runs the frontend dev server (e.g. with hot reload) and is not intended for production deployment.

## Clarifications

### Session 2026-03-05

- Q: When the frontend runs in a container and needs to call the backend API, how should the API base URL be configured? → A: Default to backend via internal network (e.g. http://backend:8080) when running in the same composition; document this and allow override via env.
- Q: Is this frontend containerization for local development only, or should the same setup also be suitable for production use? → A: Development only: container runs the frontend dev server (e.g. hot reload); not intended for production deployment.
- Q: Should the spec explicitly list what is out of scope (e.g. production deployment, HTTPS in container, authentication) to avoid scope creep? → A: Yes: add a short "Out of scope" subsection listing production deployment, HTTPS/TLS in container, and authentication (unless needed for dev).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Run Frontend in a Container (Priority: P1)

As a developer or operator, I want to run the frontend as a containerized service so that the runtime environment is consistent and reproducible, and I do not need to install the frontend toolchain directly on my machine.

**Why this priority**: Containerization is the foundation; without it, the specified port mapping and consistent behavior cannot be achieved.

**Independent Test**: Can be fully tested by starting the frontend via the containerized workflow and verifying the application is reachable on the specified host port. Delivers a single, consistent way to run the frontend.

**Acceptance Scenarios**:

1. **Given** the container image and run configuration exist, **When** the user starts the frontend using the containerized workflow, **Then** the frontend service runs and is ready to serve the application.
2. **Given** the frontend is running in a container, **When** a user opens the application in a browser using the specified host port, **Then** the frontend responds and the application is usable.
3. **Given** the user has not installed the frontend runtime on the host, **When** the user starts the frontend via the containerized workflow, **Then** the frontend runs successfully without host-side dependencies.

---

### User Story 2 - Access Frontend on Host Port 3000 (Priority: P2)

As a developer or user, I want to access the frontend application on the host at port 3000 so that I can use a predictable, documented URL (e.g. http://localhost:3000) from my browser or other tools on the host.

**Why this priority**: The user explicitly requested outer port 3000; this is the primary way users will reach the application when it is containerized.

**Independent Test**: Can be tested by starting the frontend container and confirming that connecting to the host on port 3000 serves the frontend application.

**Acceptance Scenarios**:

1. **Given** the frontend is running in a container with the intended configuration, **When** a client on the host machine connects to port 3000 (e.g. http://localhost:3000), **Then** the frontend application is served and responds correctly.
2. **Given** the frontend container is running, **When** the user checks the port mapping, **Then** the host port 3000 is mapped to the container’s internal port 5137 (or equivalent so that the app listening on 5137 inside the container is reachable at 3000 on the host).

---

### User Story 3 - Start and Stop Frontend via Container Lifecycle (Priority: P3)

As a developer or operator, I want to start and stop the frontend using standard container lifecycle commands so that I can integrate it with other containerized services and automation (e.g. orchestration or compose).

**Why this priority**: Enables predictable operations and composition with backend and other services.

**Independent Test**: Can be tested by starting the frontend container, verifying it is reachable on port 3000, then stopping it and verifying it is no longer reachable.

**Acceptance Scenarios**:

1. **Given** the frontend container definition exists, **When** the user starts the frontend container, **Then** the frontend is running and reachable on the host at port 3000 within a reasonable time (e.g. under 60 seconds).
2. **Given** the frontend container is running, **When** the user stops the frontend container, **Then** the frontend is no longer running and port 3000 on the host is no longer serving the application.
3. **Given** multiple services (e.g. frontend, backend, data store) are defined in the same composition, **When** the user starts the composition, **Then** the frontend starts and is reachable at the specified host port.

---

### Edge Cases

- What happens when the frontend container fails to start (e.g. build failure or port 3000 already in use)? The system should surface a clear failure (e.g. container exit or orchestration error) so the user can resolve the conflict or fix the build.
- How does the system behave when the host has no container runtime? The user should have documentation or feedback that the containerized workflow requires a container runtime to be installed.
- What happens when the frontend needs to call the backend or other services? When the frontend runs in the same composition as the backend, the default API base URL MUST be the backend via the internal network (e.g. http://backend:8080). This default MUST be overridable via an environment variable so that different environments or host-run backends can be used. Documentation MUST state the variable name and the default.

## Out of scope

- **Production deployment**: This feature does not cover deploying the frontend container to production environments.
- **HTTPS/TLS in container**: Configuring HTTPS or TLS inside the frontend container is out of scope; use plain HTTP for local dev.
- **Authentication**: Adding or configuring authentication in the frontend or container is out of scope unless required for local development (e.g. to call a protected backend).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The frontend MUST run as a containerized service that can be started and stopped via standard container lifecycle (e.g. run/start/stop).
- **FR-002**: The frontend MUST be reachable from the host at port **3000** (outer/host port). Connections to the host on port 3000 MUST be forwarded to the frontend service inside the container.
- **FR-003**: The frontend application inside the container MUST listen on port **5137** (inner/container port), or the container MUST be configured so that host port 3000 maps to the application listening on 5137 inside the container.
- **FR-004**: The frontend MUST start and become ready to serve the application within a defined, reasonable time after the container starts.
- **FR-005**: The frontend container definition MUST be documented or discoverable so that users know how to start the frontend and how to access it (e.g. http://localhost:3000).
- **FR-006**: When the frontend container runs in the same composition as the backend, the default API base URL MUST be the backend service on the internal network (e.g. http://backend:8080). The API base URL MUST be overridable via a documented environment variable.
- **FR-007**: The frontend container MUST run the frontend in development mode (e.g. dev server with hot reload). Production builds and production deployment are out of scope for this feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can start the frontend using only the containerized workflow (no frontend runtime installed on the host) and have it ready to serve the application within 60 seconds.
- **SC-002**: When the frontend is running with the intended configuration, a user opening http://localhost:3000 (or the host machine on port 3000) receives the frontend application and can use it.
- **SC-003**: The host port used to access the frontend is 3000, and the application inside the container listens on port 5137 (or is exposed to the host via 3000→5137 mapping).
- **SC-004**: A user can stop the frontend via the container lifecycle and confirm that the application is no longer reachable on port 3000.

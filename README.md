# Open Source Data API

It follows a similar approch like Contentful or Storyblok but in an enhanced way and for all types of data, not only content oriented.

The Data API is planed as a client-server application.

## Get started

Open the project in the devcontainer and run `make watch` or `make run` to start the application.

## Server

The server component of the Data API is designed for high performance and scalability, leveraging Go for its concurrency and efficiency. The key features and functionalities are grouped as follows:

### 1. **Data Model and Schema**

- **YAML-based Schema Definition**: Content types and fields are defined using YAML for flexibility and ease of configuration.
- **Global and Local Fields**: Support for global field definitions with the ability to override them locally for specific use cases.
- **Versioning and Storage**: Schemas and content snapshots are stored in Git to ensure traceability and version control.

### 2. **Architecture and Data Flow**

- **Event Sourcing & CQRS**: Implements an event-driven architecture using NATS Jetstream and UUIDv7 for event IDs, ensuring a clear separation of read and write operations.
- **Hierarchical Structure**: The internal hierarchy is organized as **Account -> Space -> Environment**.
- **Dependency Tree**: Each data node maintains knowledge of its ancestors and descendants, enabling a tree-like structure of dependencies.

### 3. **API Features**

- **Data Validation**: Endpoints provide options to validate data without performing insert or update operations.
- **Compact Responses**: Option to return compacted content without excessive metadata (e.g., content types used).
- **Search and Sorting**: Supports field-level control for searching and sorting content.
- **CRN (Content Resource Number)**: Each content node is uniquely identified by a CRN, similar to AWS ARN (`crn:account:space:env:type:*`).

### 4. **Security and Access Control**

- **ACL (Access Control List)**: Permissions can be defined at the field level with support for read, create, update, and delete operations.
- **Rate Limiting**: API usage can be limited based on IP, user, or API key to prevent abuse.
- **Multi-Tenant Support**: Designed to support multiple tenants, enabling isolated management of different user groups.

### 5. **Performance and Scalability**

- **Caching**: Integration with Redis for caching to improve performance and reduce latency.
- **Aggregated Data Storage**: Aggregated data nodes are stored in MongoDB or Jetstream for efficient querying and retrieval.

### 6. **Real-Time and Automation Features**

- **Webhooks**: Webhooks can be defined using glob patterns to select relevant data nodes based on the data node tree.
- **WebSocket Support**: Enables real-time updates for user interfaces.
- **Schedulable Data Nodes**: Supports scheduling for data nodes and datasets.

### 7. **Maintenance and Monitoring**

- **Backup and Restore**: Provides functionality for automated backups and data restoration to ensure data integrity.
- **Monitoring**: Integrated with OpenTelemetry for monitoring and logging.

### 8. **Internationalization**

- **i18n Support**: All fields and data nodes can be defined with multilingual support to cater to global audiences.

## Client

The UI makes the context what type of data you want to manage. The Data API should be fully built with open source software.

- Developed with Nuxt/Vue as an SPA. For Admin UIs a good solution because of the complexity of the functionality.
- Responsive Design: Optimized for mobile and tablets too.
- Extendable: Support for modules or plugins for specific requirements.

## Further ideas

- Watching schema files and reload schema manager: [fsnotify](https://github.com/fsnotify/fsnotify)
